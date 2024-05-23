package activelimit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalActiveLimiter_Limit(t *testing.T) {
	l := NewLocalActiveLimiter(1)
	tests := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		want    bool
		wantErr error
	}{
		{
			name: "normal",
			before: func(t *testing.T) {
				l.count.Store(0)
			},
			after: func(t *testing.T) {
				l.count.Store(0)
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "limited",
			before: func(t *testing.T) {
				l.count.Store(1)
			},
			after: func(t *testing.T) {
				l.count.Store(0)
			},
			want:    true,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			defer tt.after(t)
			got, err := l.Limit(context.Background(), "")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLocalActiveLimiter_Decr(t *testing.T) {
	l := NewLocalActiveLimiter(1)
	l.count.Store(1)
	tests := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantErr error
	}{
		{
			name: "normal",
			before: func(t *testing.T) {
				l.count.Store(1)
			},
			after: func(t *testing.T) {
				l.count.Store(0)
			},
			wantErr: nil,
		},
		{
			name: "count_less_than_0",
			before: func(t *testing.T) {
				l.count.Store(0)
			},
			after: func(t *testing.T) {
				l.count.Store(0)
			},
			wantErr: errors.New("错误使用 LocalActiveLimiter.Decr"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			defer tt.after(t)
			err := l.Decr(context.Background(), "")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestLocalActiveLimiter_Lifecycle(t *testing.T) {
	l := NewLocalActiveLimiter(1)
	tests := []struct {
		name    string
		op      func() (bool, error)
		want    bool
		wantErr error
	}{
		{
			name: "add",
			op: func() (bool, error) {
				return l.Limit(context.Background(), "")
			},
			want: false,
		},
		{
			name: "another_add",
			op: func() (bool, error) {
				return l.Limit(context.Background(), "")
			},
			want: true,
		},
		{
			name: "decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), "")
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "another_decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), "")
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "bad_decr",
			op: func() (bool, error) {
				return false, l.Decr(context.Background(), "")
			},
			want:    false,
			wantErr: errors.New("错误使用 LocalActiveLimiter.Decr"),
		},
	}
	for _, tt := range tests {
		got, err := tt.op()
		assert.Equalf(t, tt.wantErr, err, "%s: failed", tt.name)
		assert.Equalf(t, tt.want, got, "%s: failed", tt.name)
	}
}
