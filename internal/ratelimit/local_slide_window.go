package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/udugong/limiter/internal/queue"
)

type LocalSlideWindowLimiter struct {
	// 窗口大小
	Window time.Duration

	// 有界队列
	Queue    queue.BoundedQueue
	lock     sync.Mutex
	timeFunc func() time.Time
}

// NewLocalSlideWindowLimiter 本地的滑动窗口算法限流器实现
func NewLocalSlideWindowLimiter(window time.Duration, queue queue.BoundedQueue, opts ...Option) *LocalSlideWindowLimiter {
	l := &LocalSlideWindowLimiter{
		Window:   window,
		Queue:    queue,
		timeFunc: func() time.Time { return time.Now() },
	}
	for _, opt := range opts {
		opt.apply(l)
	}
	return l
}

type Option interface {
	apply(*LocalSlideWindowLimiter)
}

type optionFunc func(*LocalSlideWindowLimiter)

func (f optionFunc) apply(limiter *LocalSlideWindowLimiter) {
	f(limiter)
}

// WithTimeFunc 控制生成当前时间
func WithTimeFunc(fn func() time.Time) Option {
	return optionFunc(func(limiter *LocalSlideWindowLimiter) {
		limiter.timeFunc = fn
	})
}

func (l *LocalSlideWindowLimiter) Limit(_ context.Context, _ string) (bool, error) {
	l.lock.Lock()
	now := l.timeFunc()
	if !l.Queue.IsFull() {
		_ = l.Queue.Enqueue(now)
		l.lock.Unlock()
		return false, nil
	}
	windowStart := now.Add(-l.Window)
	for {
		first, err := l.Queue.Peek()
		if err != nil {
			break
		}
		if first.Before(windowStart) {
			_, _ = l.Queue.Dequeue()
		} else {
			break
		}
	}
	if !l.Queue.IsFull() {
		_ = l.Queue.Enqueue(now)
		l.lock.Unlock()
		return false, nil
	}
	l.lock.Unlock()
	return true, nil
}
