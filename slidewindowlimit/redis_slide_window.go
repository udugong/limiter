package slidewindowlimit

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/udugong/limiter"
	"github.com/udugong/limiter/internal/slidewindowlimit"
)

// NewRedisSlidingWindowLimiter 创建一个基于 redis 的滑动窗口限流器.
// cmd: 可传入 redis 的客户端
// interval: 窗口大小
// rate: 阈值
// 表示: 在 interval 内允许 rate 个请求
// 示例: 1s 内允许 3000 个请求 NewRedisSlidingWindowLimiter(redis.Client, time.Second, 3000)
func NewRedisSlidingWindowLimiter(cmd redis.Cmdable,
	interval time.Duration, rate int) limiter.Limiter {
	return &slidewindowlimit.RedisSlidingWindowLimiter{
		Cmd:      cmd,
		Interval: interval,
		Rate:     rate,
	}
}
