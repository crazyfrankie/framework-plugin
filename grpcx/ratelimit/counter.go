package ratelimit

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CounterLimiter 固定计数器算法
type CounterLimiter struct {
	cnt       *atomic.Int32
	threshold int32
}

func BuildCounter(threshold int32) *CounterLimiter {
	return &CounterLimiter{
		cnt:       nil,
		threshold: threshold,
	}
}

func (c *CounterLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		c.cnt.Add(1)
		defer func() {
			c.cnt.Add(-1)
		}()
		if c.cnt.Load() > c.threshold {
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}

		return handler(ctx, req)
	}

}
