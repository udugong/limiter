package bucketlimit

import (
	"time"

	"github.com/udugong/limiter/internal/bucketlimit"
)

// NewLeakyBucketLimiter 创建一个漏桶限流器.
// interval 每个请求之间的间隔
func NewLeakyBucketLimiter(interval time.Duration) *bucketlimit.Bucket {
	return bucketlimit.NewLeakyBucket(interval)
}
