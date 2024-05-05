package hello

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/config"
)

type Worker struct {
	Name string `json:"name"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.Name = conf.GetStringValueWithDefault("name", "Stranger")

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	fmt.Println("Hello ", w.Name)

	return nil
}

func (w *Worker) GetSleep() int {
	return 0
}
