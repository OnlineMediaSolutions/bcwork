package alerts

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/utils/bccron"
	"os/exec"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rotisserie/eris"
)

type Worker struct {
	DatabaseEnv string            `json:"dbenv"`
	TaskCrons   map[string]string `json:"task_crons"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	w.TaskCrons = make(map[string]string)

	// Get the tasks and cron schedules from the configuration
	taskCronPairs, _ := conf.GetStringValue("tasks")
	taskCronPairs = strings.Trim(taskCronPairs, "{}") // Remove curly braces
	pairs := strings.Split(taskCronPairs, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			task := strings.TrimSpace(parts[0])
			cron := strings.TrimSpace(parts[1])
			w.TaskCrons[task] = cron
		}
	}

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	// Iterate over tasks and execute them according to their cron schedules
	for task, cron := range w.TaskCrons {
		fmt.Printf("Executing task: %s with cron: %s\n", task, cron)

		// Construct the file path based on the task name
		filePath := fmt.Sprintf("./%s.sh", task) // assuming the file is a shell script

		// Execute the file as a separate process
		cmd := exec.CommandContext(ctx, "/bin/sh", filePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to execute task %s: %v\n", task, err)
		} else {
			fmt.Printf("Output of task %s: %s\n", task, string(output))
		}
	}

	fmt.Println("All tasks processed")

	return nil
}

func (w *Worker) GetSleep() int {
	for _, cron := range w.TaskCrons {
		return bccron.Next(cron)
	}
	return 0
}
