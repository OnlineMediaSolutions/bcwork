package hello

import (
	"context"
	"fmt"

	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bccron"
)

type Worker struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	w.Name = conf.GetStringValueWithDefault("name", "Stranger")
	w.Cron, _ = conf.GetStringValue("cron")

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	fmt.Println("Hello", w.Name, w.Cron)

	return nil
}

func (w *Worker) GetSleep() int {
	if w.Cron != "" {
		return bccron.Next(w.Cron)
	}

	return 0
}
