package alerts

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/utils/bccron"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/workers/alerts/factor"
	"github.com/rotisserie/eris"
)

type Worker struct {
	DatabaseEnv string            `json:"dbenv"`
	TaskCrons   map[string]string `json:"task_crons"`
	Factor      *factor.Factor
}

func (w *Worker) Init(conf config.StringMap) error {
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	pairs := w.getParams(conf)
	err2, done := w.runAlerts(pairs, conf)
	if done {
		return err2
	}

	return nil
}

func (w *Worker) runAlerts(pairs []string, conf config.StringMap) (error, bool) {
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			task := strings.TrimSpace(parts[0])
			cron := strings.TrimSpace(parts[1])
			w.TaskCrons[task] = cron

			if task == "factor" {
				w.Factor = &factor.Factor{
					DatabaseEnv: w.DatabaseEnv,
					Cron:        cron,
				}
				err := w.Factor.Init(conf)
				if err != nil {
					return eris.Wrapf(err, "Failed to initialize Factor task"), true
				}
			}
		}
	}
	return nil, false
}

func (w *Worker) getParams(conf config.StringMap) []string {
	w.TaskCrons = make(map[string]string)
	taskCronPairs, _ := conf.GetStringValue("tasks")
	taskCronPairs = strings.Trim(taskCronPairs, "{}")
	pairs := strings.Split(taskCronPairs, ",")
	return pairs
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
			w.Factor.Cron = cron

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
