package real_time_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/rotisserie/eris"
	"sort"
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
	Cron           string                `json:"cron"`
	Quest          []string              `json:"quest"`
	Start          time.Time             `json:"start"`
	End            time.Time             `json:"end"`
	Slack          *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv    string                `json:"dbenv"`
	EmailCreds     map[string]string     `json:"email_creads"`
	Fees           map[string]float64    `json:"fees"`
	ConsultantFees map[string]float64    `json:"consultant_fees"`
	HttpClient     httpclient.Doer
	Publishers     map[string]string
	skipInitRun    bool
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var questExist bool

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")

	emailCredsMap, err := config.FetchConfigValues([]string{"real_time_report"})
	worker.EmailCreds = emailCredsMap

	if err != nil {
		return eris.Wrapf(err, fmt.Sprintf("failed to get email credentials %s", worker.DatabaseEnv))
	}

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return eris.Wrap(err, fmt.Sprintf("failed initialize DB for real time report in environment: %s", worker.DatabaseEnv))
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	if worker.skipInitRun {
		fmt.Println("Skipping work as per the skip_init_run flag real time report.")
		worker.skipInitRun = false
		return nil
	}

	fmt.Println("Starting real time reports worker task")

	var emailCreds EmailCreds
	credsRaw := worker.EmailCreds["real_time_report"]

	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	report, err := worker.FetchFromQuest(ctx, worker.Start, worker.End)
	if err != nil {
		fmt.Println("Error fetching records for real time report:", err)
		log.Error().Err(err).Msg("Failed to fetch records from Quest for real time report")
		return err
	}

	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		fmt.Println("Error unmarshalling email credentials for real time report:", err)
		return err
	}

	worker.prepareEmail(report, err, emailCreds)

	return nil
}

func (worker *Worker) prepareEmail(report map[string]*RealTimeReport, err error, emailCreds EmailCreds) {
	var reports []RealTimeReport
	for _, r := range report {
		reports = append(reports, *r)
	}

	sort.Slice(reports, func(i, j int) bool {
		dateI := helpers.FormatDate(reports[i].Time)
		dateJ := helpers.FormatDate(reports[j].Time)
		return dateI < dateJ
	})

	body, subject, reportName := GenerateReportDetails(worker)

	err = SendCustomHTMLEmail(
		emailCreds.TO,
		emailCreds.BCC,
		subject,
		body,
		reports,
		reportName)
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
