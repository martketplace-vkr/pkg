package rediscomponent

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/redisconnector"
	"github.com/redis/go-redis/v9"
)

const (
	cmpName = "redis"
)

type RedisConnector struct {
	cfg Config
	*redis.Client
}

func New(cfg Config) *RedisConnector {
	client, err := redisconnector.NewRedisClient(cfg.Config)
	if err != nil {
		log.Fatalf("failed to create redis client: %v", err)
	}

	return &RedisConnector{
		cfg:    cfg,
		Client: client,
	}
}

func (r *RedisConnector) Start(ctx context.Context) (err error) {
	_, err = r.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisConnector) Stop(_ context.Context) (err error) {
	return r.Close()
}

func (r *RedisConnector) GetStartTimeout() time.Duration {
	return r.cfg.StartTimeout.Duration
}

func (r *RedisConnector) GetStopTimeout() time.Duration {
	return r.cfg.StopTimeout.Duration
}

func (r *RedisConnector) GetShutdownDelay() time.Duration {
	return r.cfg.ShutdownDelay.Duration
}

func (r *RedisConnector) GetName() string {
	return cmpName
}
