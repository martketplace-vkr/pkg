package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/martketplace-vkr/pkg/inbox/dto"
)

// pgAdvisoryLockerSQLX реализует session-level advisory locks через sqlx.
type pgAdvisoryLockerSQLX struct {
	db *sqlx.DB

	mu    sync.Mutex
	conns map[string]*sql.Conn // key -> pinned connection (session)
}

func NewPgSessionDistributedLocker(db *sqlx.DB) *pgAdvisoryLockerSQLX {
	return &pgAdvisoryLockerSQLX{
		db:    db,
		conns: make(map[string]*sql.Conn),
	}
}

func (l *pgAdvisoryLockerSQLX) Lock(
	ctx context.Context,
	events dto.Events,
) (clean dto.Events, dirty []int64, _ error) {

	seen := make(map[string]struct{})

	for _, e := range events {
		key := e.Key
		if key == "" {
			dirty = append(dirty, e.ID)
			continue
		}

		if l.isLockedLocally(key) {
			clean = append(clean, e)
			continue
		}

		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			if err := l.tryLockKey(ctx, key); err != nil {
				switch {
				case errors.Is(err, errAlreadyLockedLocally):
				case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
					dirty = append(dirty, e.ID)
					continue
				case errors.Is(err, errKeyBusy):
					dirty = append(dirty, e.ID)
					continue
				default:
					dirty = append(dirty, e.ID)
					continue
				}
			}
		}

		// сюда попадаем, если ключ у нас
		clean = append(clean, e)
	}

	return clean, dirty, nil
}

func (l *pgAdvisoryLockerSQLX) Unlock(ctx context.Context, events dto.Events) error {
	seen := make(map[string]struct{})

	for _, e := range events {
		key := e.Key
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if err := l.unlockKey(ctx, key); err != nil {
			// можно аккумулировать multierror; для простоты — первый же err наружу
			return err
		}
	}
	return nil
}

// ============== internals ==============

func (l *pgAdvisoryLockerSQLX) isLockedLocally(key string) bool {
	l.mu.Lock()
	_, ok := l.conns[key]
	l.mu.Unlock()
	return ok
}

func (l *pgAdvisoryLockerSQLX) tryLockKey(ctx context.Context, key string) error {
	lockID := hashKeyToInt64(key)

	// быстрый путь: вдруг уже есть
	l.mu.Lock()
	if _, ok := l.conns[key]; ok {
		l.mu.Unlock()
		return errAlreadyLockedLocally
	}
	l.mu.Unlock()

	// закрепляем конкретное соединение (session)
	conn, err := l.db.DB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("db.Conn: %w", err)
	}

	locked := false
	// нельзя использовать sqlx на уровне *sql.Conn, работаем через database/sql API
	if err = conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock($1)`, lockID).Scan(&locked); err != nil {
		_ = conn.Close() // вернуть коннект пулу
		return fmt.Errorf("pg_try_advisory_lock: %w", err)
	}
	if !locked {
		_ = conn.Close()
		return errKeyBusy
	}

	// регистрируем коннект за ключом
	l.mu.Lock()
	if _, exists := l.conns[key]; exists {
		// редкая гонка: кто-то успел положить
		_, _ = conn.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, lockID)
		l.mu.Unlock()
		_ = conn.Close()
		return errAlreadyLockedLocally
	}
	l.conns[key] = conn
	l.mu.Unlock()

	return nil
}

func (l *pgAdvisoryLockerSQLX) unlockKey(ctx context.Context, key string) error {
	lockID := hashKeyToInt64(key)

	l.mu.Lock()
	conn := l.conns[key]
	if conn != nil {
		delete(l.conns, key)
	}
	l.mu.Unlock()

	if conn == nil {
		// лок уже не у нас или снят — ок
		return nil
	}

	// снять лок и освободить соединение
	_, _ = conn.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, lockID)
	return conn.Close() // вернёт соединение пулу
}
