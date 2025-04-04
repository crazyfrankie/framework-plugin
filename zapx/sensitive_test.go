package zapx

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSensitiveLog(t *testing.T) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zapcore.DebugLevel)
	customCore := NewCustomCore(core)
	l := zap.New(customCore)

	l.Info("info msg", zap.String("phone", "13117127078")) // print {"level":"info","msg":"info msg","phone":"131****7078"}
}
