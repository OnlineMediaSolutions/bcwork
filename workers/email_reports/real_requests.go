package email_reports

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	Cron  string                `json:"cron"`
	Quest []string              `json:"quest_instances"`
	Start time.Time             `json:"start"`
	End   time.Time             `json:"end"`
	Slack *messager.SlackModule `json:"slack_instances"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	err := worker.InitializeValues(conf)
	if err != nil {
		message := fmt.Sprintf("failed to initialize values. Error: %s", err.Error())
		//log.Error().Msg(message)
		worker.Alert(message)
		return errors.New(message)
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	worker.End = time.Now().UTC()

	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	realTimeRecordsMap, err := worker.FetchFromQuest(ctx, worker.Start, worker.End)

	if err != nil {
		fmt.Println("err", err)
		//log.Err(err)
	}

	for key, report := range realTimeRecordsMap {
		fmt.Printf("Key: %s, BidRequests: %.2f\n", key, report.BidRequests)
	}
	return nil
}

func (worker *Worker) GetSleep() int {
	//if worker.Cron != "" {
	//	return bccron.Next(worker.Cron)
	//}
	return 0
}

func (worker *Worker) InitializeValues(conf config.StringMap) error {
	stringErrors := make([]string, 0)
	var err error
	var questExist bool

	worker.Slack, err = messager.NewSlackModule()
	if err != nil {
		message := fmt.Sprintf("failed to initalize Slack module, err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	if len(stringErrors) != 0 {
		return errors.New(strings.Join(stringErrors, "\n"))
	}
	return nil

}

func (worker *Worker) Alert(message string) {
	err := worker.Slack.SendMessage(message)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error sending slack alert: %s", err))
	}
}
