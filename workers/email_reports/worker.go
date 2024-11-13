package email_reports

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"time"

	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/rs/zerolog/log"
)

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

type Worker struct {
	Cron        string                `json:"cron"`
	Quest       []string              `json:"quest"`
	Start       time.Time             `json:"start"`
	End         time.Time             `json:"end"`
	Slack       *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv string                `json:"dbenv"`
	EmailCreds  map[string]string     `json:"email_creads"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var questExist bool

	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	emailCredsMap, err := config.FetchConfigValues([]string{"real_time_report"})

	if err != nil {
		return eris.Wrapf(err, "failed to get email credentials from  DB ", worker.DatabaseEnv)
	}

	worker.EmailCreds = emailCredsMap

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return eris.Wrapf(err, "failed to initialize DB for real time report in environment: %s", worker.DatabaseEnv)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	fmt.Println("Starting real time reports worker task")

	var emailCreds EmailCreds
	credsRaw := worker.EmailCreds["real_time_report"]

	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	report, err := worker.FetchFromQuest(ctx, worker.Start, worker.End)
	if err != nil {
		fmt.Println("Error fetching records:", err)
		log.Error().Err(err).Msg("Failed to fetch records from Quest")
		return err
	}

	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		//message := fmt.Println("Error unmarshalling email credentials for real time report")
		return err
	}

	var reports []RealTimeReport
	for _, r := range report {
		reports = append(reports, *r)
	}

	subject := fmt.Sprintf("Real time reports for %s - %s\n", worker.End, worker.Start)
	//body := "<h1>Hello</h1>"

	err = SendCustomHTMLEmail(
		emailCreds.TO,
		emailCreds.BCC,
		subject,
		subject,
		reports)

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

func (worker *Worker) Alert(message string) {
	err := worker.Slack.SendMessage(message)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error sending slack alert: %s", err))
	}
}
