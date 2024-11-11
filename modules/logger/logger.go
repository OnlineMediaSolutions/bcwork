package logger

import (
	"context"

	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(ctx context.Context) *zerolog.Logger {
	logger, ok := ctx.Value(constant.LoggerContextKey).(*zerolog.Logger)
	if !ok {
		return &log.Logger
	}

	return logger
}
