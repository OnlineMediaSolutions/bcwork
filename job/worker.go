package job

import (
	"context"

	"github.com/m6yf/bcwork/config"
)

type Worker interface {
	Do(context.Context) error
	Init(context.Context, config.StringMap) error
	GetSleep() int
}
