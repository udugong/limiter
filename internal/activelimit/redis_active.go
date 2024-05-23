package activelimit

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisActiveLimiter struct {
	maxActive int64
	cli       redis.Cmdable
}

func NewRedisActiveLimiter(maxActive int64, cli redis.Cmdable) *RedisActiveLimiter {
	return &RedisActiveLimiter{
		maxActive: maxActive,
		cli:       cli,
	}
}

func (r *RedisActiveLimiter) Limit(ctx context.Context, key string) (bool, error) {
	count, err := r.cli.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > r.maxActive, nil
}

func (r *RedisActiveLimiter) Decr(ctx context.Context, key string) error {
	count, err := r.cli.Decr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count < 0 {
		return errors.New("错误使用 RedisActiveLimiter.Decr")
	}
	return nil
}
