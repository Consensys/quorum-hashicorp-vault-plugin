package utils

import (
	"context"
	"github.com/hashicorp/go-hclog"
)

type loggerCtxKey string

const loggerKey loggerCtxKey = "logger"

func WithLogger(ctx context.Context, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) hclog.Logger {
	logger, _ := ctx.Value(loggerKey).(hclog.Logger)
	return logger
}
