package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/eapache/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SlideWindowLimiter 带快慢路径的滑动窗口算法
type SlideWindowLimiter struct {
	// 窗口大小
	window time.Duration
	// 请求的时间戳
	queue     *queue.Queue
	mux       sync.Mutex
	threshold int
}

func BuildSlideWindow(window time.Duration, threshold int) *SlideWindowLimiter {
	return &SlideWindowLimiter{
		window:    window,
		queue:     queue.New(),
		threshold: threshold,
	}
}

func (s *SlideWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		s.mux.Lock()
		now := time.Now()

		// 快路径
		if s.queue.Length() < s.threshold {
			s.queue.Add(now)
			s.mux.Unlock()
			return handler(ctx, req)
		}

		// 慢路径
		// 当前窗口的初始时间
		windowStart := now.Add(-s.window)
		for {
			first := s.queue.Peek().(time.Time)
			if first.Before(windowStart) {
				// 就是删了 first
				s.queue.Remove()
			} else {
				break
			}
		}

		if s.queue.Length() < s.threshold {
			s.queue.Add(now)
			s.mux.Unlock()
			return handler(ctx, req)
		}
		s.mux.Unlock()

		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}
