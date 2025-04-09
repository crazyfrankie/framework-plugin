package zapx

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSensitiveLog(t *testing.T) {
	cfg := zap.NewProductionConfig()
	l, err := cfg.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return NewCustomCore(core)
	}))
	if err != nil {
		panic(err)
	}

	l.Info("info msg", zap.String("phone", "13117127078")) // print {"level":"info","ts":1744177043.0410442,"caller":"zapx/sensitive_test.go:19","msg":"info msg","phone":"131****7078"}
}
