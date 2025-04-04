package logging

import "time"

type GRPCType string

const (
	Unary        GRPCType = "unary"
	ClientStream GRPCType = "client_stream"
	ServerStream GRPCType = "server_stream"
	BidiStream   GRPCType = "bidi_stream"
)

type Reporter interface {
	MsgReceive(replyProto any, err error, rpcDuration time.Duration)
	MsgSend(reqProto any, err error, rpcDuration time.Duration)
	MsgCall(err error, rpcDuration time.Duration)
}

type NoopReporter struct{}

func (NoopReporter) MsgCall(error, time.Duration)         {}
func (NoopReporter) MsgSend(any, error, time.Duration)    {}
func (NoopReporter) MsgReceive(any, error, time.Duration) {}

type report struct {
	callMeta  CallMeta
	startTime time.Time
}

func newReport(callMeta CallMeta) report {
	return report{
		callMeta:  callMeta,
		startTime: time.Now(),
	}
}
