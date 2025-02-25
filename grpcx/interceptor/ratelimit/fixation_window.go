package ratelimit

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FixationWindowLimiter 固定窗口算法
type FixationWindowLimiter struct {
	// 窗口大小
	window time.Duration
	// 上一个窗口的起始时间
	lastStart time.Time
	// 当前窗口的请求数量
	cnt int
	// 窗口允许的最大请求数量
	threshold int

	mux sync.Mutex
}

func BuildFixWindow(window time.Duration, threshold int) *FixationWindowLimiter {
	return &FixationWindowLimiter{
		window:    window,
		threshold: threshold,
	}
}

func (f *FixationWindowLimiter) NewSeverInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		f.mux.Lock()

		now := time.Now()
		// 要换窗口了
		if now.After(f.lastStart.Add(f.window)) {
			f.lastStart = now
			f.cnt = 0
		}
		f.cnt++
		if f.cnt <= f.threshold {
			f.mux.Unlock()
			return handler(ctx, req)
		}
		f.mux.Unlock()
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}
