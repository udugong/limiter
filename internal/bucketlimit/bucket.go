package bucketlimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Bucket struct {
	// 每隔多久一个令牌
	interval time.Duration
	buckets  chan struct{}
	closeCh  chan struct{}
	once     sync.Once
}

// NewTokenBucket 令牌桶算法
func NewTokenBucket(interval time.Duration, capacity int) *Bucket {
	return &Bucket{
		interval: interval,
		buckets:  make(chan struct{}, capacity),
		closeCh:  make(chan struct{}),
		once:     sync.Once{},
	}
}

// NewLeakyBucket 漏桶算法
func NewLeakyBucket(interval time.Duration) *Bucket {
	return &Bucket{
		interval: interval,
		buckets:  make(chan struct{}),
		closeCh:  make(chan struct{}),
		once:     sync.Once{},
	}
}

func (b *Bucket) Put() {
	b.buckets <- struct{}{}
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-b.closeCh:
			return
		case <-ticker.C:
			b.buckets <- struct{}{}
		}
	}
}

func (b *Bucket) Close() {
	b.once.Do(func() {
		close(b.closeCh)
	})
}

func (b *Bucket) BlockLimit(ctx context.Context, _ string) (bool, error) {
	select {
	case <-b.buckets:
		return false, nil
	case <-ctx.Done():
		return true, ctx.Err()
	case <-b.closeCh:
		return false, errors.New("限流器被关闭了")
	}
}

func (b *Bucket) Limit(ctx context.Context, _ string) (bool, error) {
	select {
	case <-b.buckets:
		return false, nil
	case <-ctx.Done():
		return true, ctx.Err()
	case <-b.closeCh:
		return false, errors.New("限流器被关闭了")
	default:
		return true, nil
	}
}
