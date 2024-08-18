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

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	pairs := worker.getParams(conf)
	err2, done := worker.runAlerts(pairs, conf)
	if done {
		return err2
	}

	return nil
}

func (worker *Worker) runAlerts(pairs []string, conf config.StringMap) (error, bool) {
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			task := strings.TrimSpace(parts[0])
			cron := strings.TrimSpace(parts[1])
			worker.TaskCrons[task] = cron

			if task == "factor" {
				worker.Factor = &factor.Factor{
					DatabaseEnv: worker.DatabaseEnv,
					Cron:        cron,
				}
				err := worker.Factor.Init(conf)
				if err != nil {
					return eris.Wrapf(err, "Failed to initialize Factor task"), true
				}
			}
		}
	}
	return nil, false
}

func (worker *Worker) getParams(conf config.StringMap) []string {
	worker.TaskCrons = make(map[string]string)
	taskCronPairs, _ := conf.GetStringValue("tasks")
	taskCronPairs = strings.Trim(taskCronPairs, "{}")
	pairs := strings.Split(taskCronPairs, ",")
	return pairs
}

func (worker *Worker) Do(ctx context.Context) error {
	for task, cron := range worker.TaskCrons {
		fmt.Printf("Executing task: %s with cron: %s\n", task, cron)

		switch task {
		case "factor":
			err := worker.Factor.Do(ctx)
			if err != nil {
				fmt.Printf("Failed to execute Factor task: %v\n", err)
			}
			worker.Factor.Cron = cron

		default:
			fmt.Printf("No handler found for task: %s\n", task)
		}
	}

	fmt.Println("All tasks processed")

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Factor != nil && worker.Factor.Cron != "" {
		return bccron.Next(worker.Factor.Cron)
	}
	return 0
}
