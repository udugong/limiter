package activelimit

import "github.com/udugong/limiter/internal/activelimit"

// NewLocalActiveLimiter 创建一个本地活跃请求数限流器.
// maxActive 最大请求数
func NewLocalActiveLimiter(maxActive int64) *activelimit.LocalActiveLimiter {
	return activelimit.NewLocalActiveLimiter(maxActive)
}
