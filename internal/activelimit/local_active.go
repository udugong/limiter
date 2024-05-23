package activelimit

import (
	"context"
	"errors"
	"sync/atomic"
)

type LocalActiveLimiter struct {
	maxActive int64
	count     atomic.Int64
}

func NewLocalActiveLimiter(maxActive int64) *LocalActiveLimiter {
	return &LocalActiveLimiter{
		maxActive: maxActive,
	}
}

func (l *LocalActiveLimiter) Limit(_ context.Context, _ string) (bool, error) {
	count := l.count.Add(1)
	return count > l.maxActive, nil
}

func (l *LocalActiveLimiter) Decr(_ context.Context, _ string) error {
	v := l.count.Add(-1)
	if v < 0 {
		return errors.New("错误使用 LocalActiveLimiter.Decr")
	}
	return nil
}
