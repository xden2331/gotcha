package logger

import (
	"context"
	"log"
	"os"
)

type Logger interface {
	CtxInfo(ctx context.Context, format string, v ...interface{})
	CtxWarn(ctx context.Context, format string, v ...interface{})
	CtxError(ctx context.Context, format string, v ...interface{})
	CtxFatal(ctx context.Context, format string, v ...interface{})
}

type stdLogger struct {
	logger *log.Logger
}

func NewStdLogger() Logger {
	return &stdLogger{logger: log.New(os.Stdout, "", log.LstdFlags)}
}

func (s *stdLogger) CtxInfo(ctx context.Context, format string, v ...interface{}) {
	s.logger.Printf(format, v)
}

func (s *stdLogger) CtxWarn(ctx context.Context, format string, v ...interface{}) {
	s.logger.Printf(format, v)
}

func (s *stdLogger) CtxError(ctx context.Context, format string, v ...interface{}) {
	s.logger.Printf(format, v)
}

func (s *stdLogger) CtxFatal(ctx context.Context, format string, v ...interface{}) {
	s.logger.Printf(format, v)
}
