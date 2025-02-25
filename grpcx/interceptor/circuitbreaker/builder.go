package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func NewInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{breaker: sre.NewBreaker()}
}

func (b *InterceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if err := b.breaker.Allow(); err != nil {
			return nil, status.Errorf(codes.Unavailable, "触发熔断")
		}

		// 尝试处理请求
		resp, err = handler(ctx, req)
		if err != nil {
			b.breaker.MarkFailed()
			return nil, err
		}

		b.breaker.MarkSuccess()
		return resp, nil
	}
}
