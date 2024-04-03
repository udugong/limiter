package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	queuemocks "github.com/udugong/limiter/internal/ratelimit/mocks"
)

func TestLocalSlideWindowLimiter_Limit(t *testing.T) {
	now := time.UnixMilli(1695571200000)
	l := &LocalSlideWindowLimiter{
		Window: 10 * time.Second,
	}
	WithTimeFunc(func() time.Time {
		return now
	}).apply(l)
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) BoundedQueue
		want    bool
		wantErr error
	}{
		{
			// 窗口已满,清理后可以通过
			name: "queue_is_full_but_passed_after_cleaning",
			mock: func(ctrl *gomock.Controller) BoundedQueue {
				beforeTime := time.UnixMilli(1685571200000)
				q := queuemocks.NewMockQueue(ctrl)
				q.EXPECT().IsFull().Return(true)
				q.EXPECT().Peek().Return(beforeTime, nil)
				q.EXPECT().Dequeue().Return(beforeTime, nil)
				q.EXPECT().Peek().Return(time.Time{}, errors.New("队列为空"))
				q.EXPECT().IsFull().Return(false)
				q.EXPECT().Enqueue(gomock.Any()).Return(nil)
				return q
			},
			want:    false,
			wantErr: nil,
		},
		{
			// 窗口已满
			name: "queue_is_full",
			mock: func(ctrl *gomock.Controller) BoundedQueue {
				beforeTime := time.UnixMilli(1685571200000)
				q := queuemocks.NewMockQueue(ctrl)
				q.EXPECT().IsFull().Return(true)
				q.EXPECT().Peek().Return(beforeTime, nil)
				q.EXPECT().Dequeue().Return(beforeTime, nil)
				q.EXPECT().Peek().Return(time.Time{}, errors.New("队列为空"))
				q.EXPECT().IsFull().Return(true)
				return q
			},
			want:    true,
			wantErr: nil,
		},
		{
			// 窗口已满,但第一次时间晚于窗口开始时间
			name: "queue_is_full_but_first_time_is_later_than_window_start_time",
			mock: func(ctrl *gomock.Controller) BoundedQueue {
				q := queuemocks.NewMockQueue(ctrl)
				q.EXPECT().IsFull().Return(true)
				q.EXPECT().Peek().Return(now, nil)
				q.EXPECT().IsFull().Return(true)
				return q
			},
			want:    true,
			wantErr: nil,
		},
		{
			// 窗口未满
			name: "queue_is_full_but_first_time_is_later_than_window_start_time",
			mock: func(ctrl *gomock.Controller) BoundedQueue {
				q := queuemocks.NewMockQueue(ctrl)
				q.EXPECT().IsFull().Return(false)
				q.EXPECT().Enqueue(now).Return(nil)
				return q
			},
			want:    false,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			l.Queue = tt.mock(ctrl)
			got, err := l.Limit(context.Background(), "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
