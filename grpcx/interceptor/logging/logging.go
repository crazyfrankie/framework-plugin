package logging

import (
	"context"
)

type LoggerFunc func(ctx context.Context, level Level, msg string, fields ...any)

type Logger struct {
	loggerFunc LoggerFunc
}

func NewLogger(fn func(ctx context.Context, level Level, msg string, fields ...any)) *Logger {
	return &Logger{loggerFunc: fn}
}
