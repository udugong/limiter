package bucketlimit

import (
	"time"

	"github.com/udugong/limiter/internal/bucketlimit"
)

// NewTokenBucketLimiter 创建一个令牌桶算法限流器.
// interval 每 interval 的时间放置一个令牌
// capacity 存放的令牌数
func NewTokenBucketLimiter(interval time.Duration, capacity int) *bucketlimit.Bucket {
	return bucketlimit.NewTokenBucket(interval, capacity)
}
