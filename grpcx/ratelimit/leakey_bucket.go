package ratelimit

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LeakyBucketLimiter 漏桶算法
type LeakyBucketLimiter struct {
	interval time.Duration
	closeCh  chan struct{}
	once     sync.Once
}

func BuildLeakyBucket(interval time.Duration) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		interval: interval,
		closeCh:  make(chan struct{}),
	}
}

func (l *LeakyBucketLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			return handler(ctx, req)
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-l.closeCh:
			return nil, status.Errorf(codes.Unavailable, "限流器已关闭")
		}
	}
}

func (l *LeakyBucketLimiter) Close() error {
	l.once.Do(func() {
		close(l.closeCh)
	})
	return nil
}
