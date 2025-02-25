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
		if b.breaker.Allow() == nil {
			resp, err = handler(ctx, req)
			if err != nil {
				// 标记处理失败
				// 可以进一步判断是否为业务错误, 根据实际情况来选
				b.breaker.MarkFailed()
			} else {
				b.breaker.MarkSuccess()
			}
		}
		// 触发了熔断器
		b.breaker.MarkFailed()
		return nil, status.Errorf(codes.Unavailable, "触发熔断")
	}
}
