package bucketlimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_BlockLimit(t *testing.T) {
	b := NewTokenBucket(10*time.Millisecond, 2)
	defer b.Close()
	go b.Put()
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	tests := []struct {
		name     string
		ctx      context.Context
		interval time.Duration
		want     bool
		wantErr  error
	}{
		{
			name:     "normal",
			ctx:      context.Background(),
			interval: 2,
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "another_normal",
			ctx:      context.Background(),
			interval: 5,
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "timeout",
			ctx:      ctx1,
			interval: 10,
			want:     true,
			wantErr:  context.DeadlineExceeded,
		},
		{
			name:     "wait_for_token",
			ctx:      context.Background(),
			interval: 15,
			want:     false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			<-time.After(tt.interval)
			got, err := b.BlockLimit(tt.ctx, "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
	t.Run("after_close", func(t *testing.T) {
		b.Close()
		<-time.After(10)
		got, err := b.BlockLimit(context.Background(), "")
		assert.Equal(t, false, got)
		assert.Equal(t, errors.New("限流器被关闭了"), err)
	})
}

func TestLeakyBucket_BlockLimit(t *testing.T) {
	b := NewLeakyBucket(20 * time.Millisecond)
	defer b.Close()
	go b.Put()
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	tests := []struct {
		name     string
		ctx      context.Context
		interval time.Duration
		want     bool
		wantErr  error
	}{
		{
			name:     "normal",
			ctx:      context.Background(),
			interval: 2,
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "timeout",
			ctx:      ctx1,
			interval: 5,
			want:     true,
			wantErr:  context.DeadlineExceeded,
		},
		{
			name:     "wait_for_token",
			ctx:      context.Background(),
			interval: 10,
			want:     false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			<-time.After(tt.interval)
			got, err := b.BlockLimit(tt.ctx, "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
	t.Run("after_close", func(t *testing.T) {
		b.Close()
		<-time.After(10)
		got, err := b.BlockLimit(context.Background(), "")
		assert.Equal(t, false, got)
		assert.Equal(t, errors.New("限流器被关闭了"), err)
	})
}

func TestTokenBucket_Limit(t *testing.T) {
	b := NewTokenBucket(10*time.Millisecond, 2)
	defer b.Close()
	ctx1, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	go b.Put()
	tests := []struct {
		name     string
		ctx      context.Context
		interval time.Duration
		want     bool
		wantErr  error
	}{
		{
			name:     "normal",
			ctx:      context.Background(),
			interval: 2,
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "limited",
			ctx:      context.Background(),
			interval: 5,
			want:     true,
			wantErr:  nil,
		},
		{
			name:     "timeout",
			ctx:      ctx1,
			interval: 10,
			want:     true,
			wantErr:  context.DeadlineExceeded,
		},
		{
			name:     "wait_5_mill_limited",
			ctx:      context.Background(),
			interval: 5 * time.Millisecond,
			want:     true,
			wantErr:  nil,
		},
		{
			name:     "wait_for_token",
			ctx:      context.Background(),
			interval: 5 * time.Millisecond,
			want:     false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			<-time.After(tt.interval)
			got, err := b.Limit(tt.ctx, "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
	t.Run("after_close", func(t *testing.T) {
		b.Close()
		<-time.After(10)
		got, err := b.Limit(context.Background(), "")
		assert.Equal(t, false, got)
		assert.Equal(t, errors.New("限流器被关闭了"), err)
	})
}

func TestLeakyBucket_Limit(t *testing.T) {
	b := NewLeakyBucket(10 * time.Millisecond)
	defer b.Close()
	go b.Put()
	ctx1, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	tests := []struct {
		name     string
		ctx      context.Context
		interval time.Duration
		want     bool
		wantErr  error
	}{
		{
			name:     "normal",
			ctx:      context.Background(),
			interval: 2,
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "limited",
			ctx:      context.Background(),
			interval: 5,
			want:     true,
			wantErr:  nil,
		},
		{
			name:     "timeout",
			ctx:      ctx1,
			interval: 10,
			want:     true,
			wantErr:  context.DeadlineExceeded,
		},
		{
			name:     "wait_5_mill_limited",
			ctx:      context.Background(),
			interval: 5 * time.Millisecond,
			want:     true,
			wantErr:  nil,
		},
		{
			name:     "wait_for_token",
			ctx:      context.Background(),
			interval: 5 * time.Millisecond,
			want:     false,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			<-time.After(tt.interval)
			got, err := b.Limit(tt.ctx, "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
	t.Run("after_close", func(t *testing.T) {
		b.Close()
		<-time.After(10)
		got, err := b.Limit(context.Background(), "")
		assert.Equal(t, false, got)
		assert.Equal(t, errors.New("限流器被关闭了"), err)
	})
}
