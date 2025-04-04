package logging

import (
	"context"

	"google.golang.org/grpc"
)

type reporter struct {
}

func (l Logger) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		_ = newReport(NewServerCallMeta(info.FullMethod, nil, req))
		// 调用用户实现的日志接口
		// 构建 reporter 进行日志信息的构建
		// 1. 请求进入时
		// 2. 返回响应时
		return nil, nil
	}
}
