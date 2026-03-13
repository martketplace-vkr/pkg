package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/inbox/internal/locker"
	"github.com/redis/go-redis/v9"
)

const distributedLockPatternf = "inbox:lock:%s:%s"

type (
	redisDistributedLocker struct {
		rd  *redis.Client
		cfg locker.Config
	}
)

func NewRedisDistributedLocker(
	cfg locker.Config,
	client *redis.Client,
) *redisDistributedLocker {
	return &redisDistributedLocker{
		cfg: cfg,
		rd:  client,
	}
}

func (l redisDistributedLocker) Lock(
	ctx context.Context,
	events dto.Events,
) (cleanEvents dto.Events, dirtyEvents []int64, err error) {
	if l.cfg.Disabled {
		return events, dirtyEvents, err
	}

	seen := make(map[string]struct{})

	for _, e := range events {
		// если один и тот же ключ уже залочили, то скипаем.
		if _, ok := seen[e.Key]; ok {
			dirtyEvents = append(dirtyEvents, e.ID)

			continue
		}

		lockKey := fmt.Sprintf(distributedLockPatternf, e.Key, l.cfg.InstanceID)

		ok, redisErr := l.rd.SetNX(ctx, lockKey, l.cfg.InstanceID, l.cfg.EventLockTTL).Result()
		if redisErr != nil {
			err = redisErr
			return
		}

		// ключ уже залочен другим инстансом/воркером – просто помечаем как dirty
		// лочим в setnx запись по key + instance
		if !ok {
			dirtyEvents = append(dirtyEvents, e.ID)

			continue
		}

		seen[e.Key] = struct{}{}
		cleanEvents = append(cleanEvents, e)
	}

	return
}

func (l redisDistributedLocker) Unlock(
	ctx context.Context,
	events dto.Events,
) error {
	if l.cfg.Disabled {
		return nil
	}

	seen := make(map[string]struct{})

	for _, e := range events {
		if _, ok := seen[e.Key]; ok {
			continue
		}
		seen[e.Key] = struct{}{}

		if err := l.unlockKey(ctx, e.Key); err != nil {
			return err
		}
	}

	return nil
}

func (l redisDistributedLocker) unlockKey(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf(distributedLockPatternf, key, l.cfg.InstanceID)

	return l.rd.Watch(ctx, func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, lockKey).Result()
		if errors.Is(err, redis.Nil) {
			return nil
		}
		if err != nil {
			return err
		}

		if val != l.cfg.InstanceID {
			return nil
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Del(ctx, lockKey)
			return nil
		})

		return err
	}, lockKey)
}
