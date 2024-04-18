package ratelimit

import (
	"time"

	"github.com/udugong/limiter/internal/queue"
	"github.com/udugong/limiter/internal/ratelimit"
)

// NewLocalSlideWindowLimiter 创建一个本地滑动窗口限流器.
// window 窗口大小
// boundedQueue 有界队列
// 表示: 在 window 内允许有界队列大小的请求
func NewLocalSlideWindowLimiter(window time.Duration, boundedQueue queue.BoundedQueue,
	opts ...ratelimit.Option) *ratelimit.LocalSlideWindowLimiter {
	return ratelimit.NewLocalSlideWindowLimiter(window, boundedQueue)
}

// WithTimeFunc 控制时间.
func WithTimeFunc(fn func() time.Time) ratelimit.Option {
	return ratelimit.WithTimeFunc(fn)
}
