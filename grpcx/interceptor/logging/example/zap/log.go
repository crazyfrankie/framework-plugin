package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/crazyfrankie/framework-plugin/grpcx/interceptor/logging"
)

// zapLog example Zap
func zapLog(l *zap.Logger) logging.LoggerFunc {
	return func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch val := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), val))
			case int:
				f = append(f, zap.Int(key.(string), val))
			case bool:
				f = append(f, zap.Bool(key.(string), val))
			default:
				f = append(f, zap.Any(key.(string), val))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch level {
		case logging.LevelDebug:
			logger.Debug(msg, f...)
		case logging.LevelInfo:
			logger.Info(msg, f...)
		case logging.LevelWarn:
			logger.Warn(msg, f...)
		case logging.LevelError:
			logger.Error(msg, f...)
		}
	}
}
