package logging

import (
	"fmt"
	"strings"

	"google.golang.org/grpc"
)

type CallMeta struct {
	Req      any
	Typ      GRPCType
	Service  string
	Method   string
	IsClient bool
}

func NewClientCallMeta(fullMethod string, streamDesc *grpc.StreamDesc, req any) CallMeta {
	c := CallMeta{Req: req, IsClient: true, Typ: Unary}
	if streamDesc != nil {
		c.Typ = clientStreamType(streamDesc)
	}
	c.Service, c.Method = splitFullMethodName(fullMethod)
	return c
}

func NewServerCallMeta(fullMethod string, streamInfo *grpc.StreamServerInfo, req any) CallMeta {
	c := CallMeta{Req: req, IsClient: false, Typ: Unary}
	if streamInfo != nil {
		c.Typ = serverStreamType(streamInfo)
	}
	c.Service, c.Method = splitFullMethodName(fullMethod)
	return c
}

func (c CallMeta) FullMethod() string {
	return fmt.Sprintf("/%s/%s", c.Service, c.Method)
}

func clientStreamType(desc *grpc.StreamDesc) GRPCType {
	if desc.ClientStreams && !desc.ServerStreams {
		return ClientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return ServerStream
	}
	return BidiStream
}

func serverStreamType(info *grpc.StreamServerInfo) GRPCType {
	if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return ServerStream
	}
	return BidiStream
}

func splitFullMethodName(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/") // remove leading slash
	if i := strings.Index(fullMethod, "/"); i >= 0 {
		return fullMethod[:i], fullMethod[i+1:]
	}
	return "unknown", "unknown"
}
