package redis

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/martketplace-vkr/pkg/inbox/internal/locker"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/martketplace-vkr/pkg/inbox/dto"
)

const (
	ttl = time.Minute
)

func newTestRedisClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err, "failed to start miniredis")

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
		DB:   0,
	})

	// sanity ping
	require.NoError(t, client.Ping(context.Background()).Err(), "failed to ping redis")

	return mr, client
}

func TestNewRedisDistributedLocker(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	locker := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)
	require.NotNil(t, locker, "expected non-nil locker")
}

func TestDistributedLocker_Lock_AllClean(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	locker := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)

	events := dto.Events{
		{Key: "resource:1"},
		{Key: "resource:2"},
		{Key: "resource:3"},
	}

	ctx := context.Background()
	clean, dirty, err := locker.Lock(ctx, events)
	require.NoError(t, err)

	require.Len(t, clean, len(events))
	require.Empty(t, dirty)
}

func TestDistributedLocker_Lock_Twice_SecondIsDirty(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	locker1 := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)
	locker2 := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)

	events := dto.Events{
		{Key: "resource:A"},
		{Key: "resource:B"},
	}

	ctx := context.Background()

	// First lock should cleanly lock everything
	clean1, dirty1, err := locker1.Lock(ctx, events)
	require.NoError(t, err, "unexpected error on first lock")
	require.Len(t, clean1, len(events))
	require.Empty(t, dirty1)

	// Second lock (from another locker instance) should mark them all as dirty
	clean2, dirty2, err := locker2.Lock(ctx, events)
	require.NoError(t, err, "unexpected error on second lock")
	require.Empty(t, clean2)
	require.Len(t, dirty2, len(events))
}

func TestDistributedLocker_Unlock_AllowsRelock(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	locker := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)

	events := dto.Events{
		{Key: "user:42"},
		{Key: "user:43"},
	}

	ctx := context.Background()

	// Lock
	clean, dirty, err := locker.Lock(ctx, events)
	require.NoError(t, err, "lock error")
	require.Len(t, clean, len(events))
	require.Empty(t, dirty)

	// Unlock
	require.NoError(t, locker.Unlock(ctx, events), "unexpected unlock error")

	// Lock again should be clean
	clean2, dirty2, err := locker.Lock(ctx, events)
	require.NoError(t, err, "second lock error")
	require.Len(t, clean2, len(events))
	require.Empty(t, dirty2)
}

func TestDistributedLocker_Lock_PartialDirty(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	locker1 := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)
	locker2 := NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client)

	lockedByFirst := dto.Events{
		{Key: "order:1"},
		{Key: "order:2"},
	}
	allEvents := dto.Events{
		{Key: "order:1"},
		{Key: "order:2"},
		{Key: "order:3"},
	}

	ctx := context.Background()

	// First locker grabs two keys
	_, _, err := locker1.Lock(ctx, lockedByFirst)
	require.NoError(t, err, "first lock error")

	// Second locker tries to grab all three, two should be dirty
	clean, dirty, err := locker2.Lock(ctx, allEvents)
	require.NoError(t, err, "second lock error")

	require.Len(t, clean, 1)
	require.Len(t, dirty, 2)
}

func TestDistributedLocker_ConcurrentLocking_SingleWinner(t *testing.T) {
	t.Parallel()

	mr, client := newTestRedisClient(t)
	defer mr.Close()
	defer client.Close()

	const goroutines = 10
	lockers := make([]*redisDistributedLocker, 0, goroutines)
	for i := 0; i < goroutines; i++ {
		lockers = append(lockers, NewRedisDistributedLocker(locker.Config{EventLockTTL: ttl, InstanceID: "inbox-instance-id"}, client))
	}

	events := dto.Events{
		{Key: "shared:1"},
	}
	ctx := context.Background()

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(goroutines)

	type result struct {
		clean int
		dirty int
		err   error
	}
	results := make(chan result, goroutines)

	for i := 0; i < goroutines; i++ {
		l := lockers[i]
		go func() {
			defer wg.Done()
			<-start
			clean, dirty, err := l.Lock(ctx, events)
			results <- result{
				clean: len(clean),
				dirty: len(dirty),
				err:   err,
			}
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	var winners, losers, errors int
	for res := range results {
		if res.err != nil {
			errors++
			continue
		}
		if res.clean == 1 && res.dirty == 0 {
			winners++
		} else if res.clean == 0 && res.dirty == 1 {
			losers++
		} else {
			require.FailNowf(t, "unexpected result", "got: %+v", res)
		}
	}

	require.Zero(t, errors, "unexpected errors during concurrent locking")
	require.Equal(t, 1, winners, "expected exactly 1 winner")
	require.Equal(t, goroutines-1, losers, "expected losers to be goroutines-1")
}
