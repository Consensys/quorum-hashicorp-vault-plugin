package log

import (
	"context"
	"os"
	"sync"

	"github.com/hashicorp/go-hclog"
)

type loggerCtxKey string

const loggerKey loggerCtxKey = "logger"

var protect sync.Once

type Logger interface {
	hclog.Logger
}

func Context(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}

	return hclog.Default()
}

func Default() Logger {
	return hclog.Default()
}

func InitLogger() {
	protect.Do(func() {
		logger := hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: false,
		})

		hclog.SetDefault(logger)
	})
}
