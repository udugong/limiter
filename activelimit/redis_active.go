package activelimit

import (
	"github.com/redis/go-redis/v9"

	"github.com/udugong/limiter/internal/activelimit"
)

// NewRedisActiveLimiter 创建一个基于 redis 的活跃请求数限流器.
func NewRedisActiveLimiter(cli redis.Cmdable, maxActive int64) *activelimit.RedisActiveLimiter {
	return activelimit.NewRedisActiveLimiter(maxActive, cli)
}
