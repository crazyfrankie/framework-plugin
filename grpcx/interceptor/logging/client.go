package logging

import (
	"context"

	"google.golang.org/grpc"
)

func (l Logger) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_ = newReport(NewClientCallMeta(method, nil, req))

		return nil
	}
}
