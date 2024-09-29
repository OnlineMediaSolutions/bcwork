package templates

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/config"
)

type Worker struct {
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	fmt.Println("Migrating Templates")

	return nil
}

func (w *Worker) GetSleep() int {

	return 0
}
