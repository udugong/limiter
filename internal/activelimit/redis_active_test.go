package activelimit

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/udugong/limiter/internal/mocks/redismocks"
)

const testKey = "active_limiter_test"

func TestRedisActiveLimiter_Limit(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		want    bool
		wantErr error
	}{
		{
			name: "normal",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetVal(1)
				cmd.EXPECT().Incr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "limited",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetVal(2)
				cmd.EXPECT().Incr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "redis_error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetErr(errors.New("mock redis error"))
				cmd.EXPECT().Incr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			want:    false,
			wantErr: errors.New("mock redis error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			l := NewRedisActiveLimiter(1, tt.mock(ctrl))
			got, err := l.Limit(context.Background(), testKey)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRedisActiveLimiter_Decr(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		wantErr error
	}{
		{
			name: "normal",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetVal(0)
				cmd.EXPECT().Decr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			wantErr: nil,
		},
		{
			name: "count_less_than_0",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetVal(-1)
				cmd.EXPECT().Decr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			wantErr: errors.New("错误使用 RedisActiveLimiter.Decr"),
		},
		{
			name: "redis_error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewIntCmd(context.Background())
				res.SetErr(errors.New("mock redis error"))
				cmd.EXPECT().Decr(gomock.Any(), testKey).Return(res)
				return cmd
			},
			wantErr: errors.New("mock redis error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			l := NewRedisActiveLimiter(1, tt.mock(ctrl))
			err := l.Decr(context.Background(), testKey)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestRedisActiveLimiter_Lifecycle(t *testing.T) {
	cli := initRedis()
	l := NewRedisActiveLimiter(1, cli)
	tests := []struct {
		name    string
		op      func() (bool, error)
		want    bool
		wantErr error
	}{
		{
			name: "add",
			op: func() (bool, error) {
				return l.Limit(context.Background(), testKey)
			},
			want: false,
		},
		{
			name: "another_add",
			op: func() (bool, error) {
				return l.Limit(context.Background(), testKey)
			},
			want: true,
		},
		{
			name: "decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), testKey)
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "another_decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), testKey)
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "bad_decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), testKey)
			},
			want:    false,
			wantErr: errors.New("错误使用 RedisActiveLimiter.Decr"),
		},
	}
	err := cli.Del(context.Background(), testKey).Err()
	require.NoError(t, err)
	defer cli.Del(context.Background(), testKey)
	for _, tt := range tests {
		got, err := tt.op()
		assert.Equalf(t, tt.wantErr, err, "%s: failed", tt.name)
		assert.Equalf(t, tt.want, got, "%s: failed", tt.name)
	}
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	return redisClient
}
