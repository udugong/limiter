package slidewindowlimit

import (
	"time"

	"github.com/udugong/limiter/internal/queue"
	"github.com/udugong/limiter/internal/slidewindowlimit"
)

// NewLocalSlideWindowLimiter 创建一个本地滑动窗口限流器.
// window 窗口大小
// boundedQueue 有界队列
// 表示: 在 window 内允许有界队列大小的请求
func NewLocalSlideWindowLimiter(window time.Duration, boundedQueue queue.BoundedQueue,
	opts ...slidewindowlimit.Option) *slidewindowlimit.LocalSlideWindowLimiter {
	return slidewindowlimit.NewLocalSlideWindowLimiter(window, boundedQueue, opts...)
}

// WithTimeFunc 控制时间.
func WithTimeFunc(fn func() time.Time) slidewindowlimit.Option {
	return slidewindowlimit.WithTimeFunc(fn)
}
