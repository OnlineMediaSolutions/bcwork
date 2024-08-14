package alerts

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/utils/bccron"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rotisserie/eris"
)

type Worker struct {
	DatabaseEnv string            `json:"dbenv"`
	TaskCrons   map[string]string `json:"task_crons"`
	Factor      *Factor
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	w.TaskCrons = make(map[string]string)

	taskCronPairs, _ := conf.GetStringValue("tasks")
	taskCronPairs = strings.Trim(taskCronPairs, "{}")
	pairs := strings.Split(taskCronPairs, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			task := strings.TrimSpace(parts[0])
			cron := strings.TrimSpace(parts[1])
			w.TaskCrons[task] = cron

			if task == "factor" {
				w.Factor = &Factor{
					DatabaseEnv: w.DatabaseEnv,
					Cron:        cron,
				}
				err := w.Factor.Init(ctx, conf)
				if err != nil {
					return eris.Wrapf(err, "Failed to initialize Factor task")
				}
			}
		}
	}

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	for task, cron := range w.TaskCrons {
		fmt.Printf("Executing task: %s with cron: %s\n", task, cron)

		switch task {
		case "factor":
			err := w.Factor.Do(ctx)
			if err != nil {
				fmt.Printf("Failed to execute Factor task: %v\n", err)
			}
		default:
			fmt.Printf("No handler found for task: %s\n", task)
		}
	}

	fmt.Println("All tasks processed")

	return nil
}

func (w *Worker) GetSleep() int {
	if w.Factor != nil && w.Factor.Cron != "" {
		return bccron.Next(w.Factor.Cron)
	}
	return 0
}
