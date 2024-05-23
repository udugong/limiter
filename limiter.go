package limiter

import "context"

type Limiter interface {
	// Limit 有没有触发限流。key 就是限流对象
	// bool 代表是否限流, true 就是要限流
	// err 限流器本身有没有错误
	Limit(ctx context.Context, key string) (bool, error)
}

// ActiveLimiter 活跃请求数限流
type ActiveLimiter interface {
	// Limit 有没有触发限流。key 就是限流对象
	// 活跃请求数增加1
	// bool 代表是否限流, true 就是要限流
	// error 限流器本身有没有错误
	Limit(ctx context.Context, key string) (bool, error)

	// Decr 活跃请求数减少1
	Decr(ctx context.Context, key string) error
}

// BucketLimiter 桶限流
type BucketLimiter interface {
	// Put 往桶里放置。该方法需要异步执行 go Put()
	Put()

	// Close 关闭 Put() 方法
	Close()

	// Limit 有没有触发限流。
	// bool 代表是否限流, true 就是要限流, 若 Context.Err() == nil 也会返回 error
	// error 当调用了 Close() 时返回错误
	Limit(ctx context.Context, _ string) (bool, error)

	// BlockLimit 限流时阻塞直到超时。
	// bool 代表是否限流, true 就是要限流, 同时会返回 Context.Err()
	// error 当调用了 Close() 时返回错误
	BlockLimit(ctx context.Context, _ string) (bool, error)
}
