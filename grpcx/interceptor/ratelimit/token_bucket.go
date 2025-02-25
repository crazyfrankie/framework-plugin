package ratelimit

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// TokenBucketLimiter 令牌桶算法实现
type TokenBucketLimiter struct {
	buckets chan struct{}
	// 每隔多久一个令牌
	interval time.Duration

	closeCh chan struct{}
	// 确保只关闭一次
	once sync.Once
}

func BuildTokenBucket(interval time.Duration, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:  make(chan struct{}, capacity),
		interval: interval,
	}
}

func (t *TokenBucketLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	ticker := time.NewTicker(t.interval)
	go func() {
		for {
			select {
			case <-t.closeCh:
				return
				// 发令牌
			case <-ticker.C:
				select {
				case t.buckets <- struct{}{}:
					// 发到桶里面
				default:
					// 桶满了
				}
			}
		}
	}()
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-t.buckets:
			// 拿到了令牌
			return handler(ctx, req)
		case <-ctx.Done():
			// 没有令牌就等, 直到超时
			return nil, ctx.Err()
			// 超高并发选择下面的写法, 上面的写法是阻塞等待获取令牌, 直到超时
			//default:
			//	return nil, status.Errorf(codes.ResourceExhausted, "限流")
		}
	}
}

func (t *TokenBucketLimiter) Close() error {
	t.once.Do(func() {
		close(t.closeCh)
	})
	return nil
}
