package postgres

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	// замени на свой модуль
	"github.com/martketplace-vkr/pkg/inbox/dto"
)

var (
	errAlreadyLockedLocally = errors.New("key is already locked by this instance")
)

// pgAdvisoryLocker реализует DistributedLocker через session-level advisory locks.
type pgAdvisoryLocker struct {
	pool *pgxpool.Pool

	mu    sync.Mutex
	conns map[string]*pgxpool.Conn
}

// NewPgxSessionDistributedLocker создаёт локер.
func NewPgxSessionDistributedLocker(pool *pgxpool.Pool) *pgAdvisoryLocker {
	return &pgAdvisoryLocker{
		pool:  pool,
		conns: make(map[string]*pgxpool.Conn),
	}
}

func (l *pgAdvisoryLocker) Lock(
	ctx context.Context,
	events dto.Events,
) (clean dto.Events, dirty []int64, _ error) {

	seen := make(map[string]bool)

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

		if !seen[key] {
			seen[key] = true

			if err := l.tryLockKey(ctx, key); err != nil {
				if errors.Is(err, errAlreadyLockedLocally) {
					clean = append(clean, e)
				} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					dirty = append(dirty, e.ID)
				} else if errors.Is(err, errKeyBusy) {
					dirty = append(dirty, e.ID)
				} else {
					dirty = append(dirty, e.ID)
				}
				continue
			}
		}

		clean = append(clean, e)
	}

	return clean, dirty, nil
}

func (l *pgAdvisoryLocker) Unlock(ctx context.Context, events dto.Events) error {
	seen := make(map[string]bool)
	for _, e := range events {
		if e.Key == "" || seen[e.Key] {
			continue
		}
		seen[e.Key] = true

		if err := l.unlockKey(ctx, e.Key); err != nil {
			// не прерываем массовый unlock; собираем первую ошибку
			// при желании можно вернуть multierror
			return err
		}
	}
	return nil
}

var errKeyBusy = errors.New("key busy")

func (l *pgAdvisoryLocker) isLockedLocally(key string) bool {
	l.mu.Lock()
	_, ok := l.conns[key]
	l.mu.Unlock()
	return ok
}

func (l *pgAdvisoryLocker) tryLockKey(ctx context.Context, key string) error {
	lockID := hashKeyToInt64(key)

	l.mu.Lock()
	if _, ok := l.conns[key]; ok {
		l.mu.Unlock()
		return errAlreadyLockedLocally
	}
	l.mu.Unlock()

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pool acquire: %w", err)
	}

	defer func() {
		if err != nil {
			conn.Release()
		}
	}()

	var locked bool
	if err = conn.QueryRow(ctx, `select pg_try_advisory_lock($1)`, lockID).Scan(&locked); err != nil {
		return fmt.Errorf("pg_try_advisory_lock scan: %w", err)
	}
	if !locked {
		// кто-то уже держит лок (любой инстанс/воркер)
		return errKeyBusy
	}

	l.mu.Lock()
	if _, exists := l.conns[key]; exists {
		_, _ = conn.Exec(ctx, `select pg_advisory_unlock($1)`, lockID)
		l.mu.Unlock()
		return errAlreadyLockedLocally
	}
	l.conns[key] = conn
	l.mu.Unlock()
	return nil
}

func (l *pgAdvisoryLocker) unlockKey(ctx context.Context, key string) error {
	lockID := hashKeyToInt64(key)

	l.mu.Lock()
	conn := l.conns[key]
	if conn != nil {
		delete(l.conns, key)
	}
	l.mu.Unlock()

	if conn == nil {
		// нечего делать — лок не у нас (или уже снят)
		return nil
	}

	// Снимаем лок и отдаём соединение пулу.
	_, _ = conn.Exec(ctx, `select pg_advisory_unlock($1)`, lockID)
	conn.Release()
	return nil
}

func hashKeyToInt64(s string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	// advisory lock принимает signed BIGINT — приводим без переполнений
	return int64(h.Sum64())
}
