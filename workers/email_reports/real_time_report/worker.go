package real_time_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/helpers"
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
	ReportName     string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var questExist bool

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	worker.HttpClient = httpclient.New(true)
	worker.ReportName = "real_time_report"

	emailCredsMap, err := config.FetchConfigValues([]string{worker.ReportName})
	worker.EmailCreds = emailCredsMap

	if err != nil {
		return fmt.Errorf("failed to get email credentials %w", err)
	}

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return fmt.Errorf("failed initialize DB for real time report in environment %s,%w", worker.DatabaseEnv, err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2", "sfoquest2"}
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting real time reports worker task")

	if worker.skipInitRun {
		log.Info().Msg("Skipping work as per the skip_init_run flag real time report.")
		worker.skipInitRun = false
		return nil
	}

	var emailCreds EmailCreds
	credsRaw := worker.EmailCreds[worker.ReportName]

	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	report, err := worker.FetchAndMergeQuestReports(ctx, worker.Start, worker.End)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
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
		firstDate := helpers.FormatDate(reports[i].Time)
		secondDate := helpers.FormatDate(reports[j].Time)
		return firstDate < secondDate
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
