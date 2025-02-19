package rpm_decrease

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/workers/email_reports"
	"golang.org/x/net/context"
)

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}
type EmailData struct {
	Body   string
	Report []AlertsEmails
}

type AggregatedReport struct {
	Date                 string  `json:"date"`
	DataStamp            int64   `json:"DateStamp"`
	Publisher            string  `json:"publisher"`
	Domain               string  `json:"domain"`
	PaymentType          string  `json:"PaymentType"`
	AM                   string  `json:"am"`
	PubImps              string  `json:"PubImps"`
	LoopingRatio         float64 `json:"looping_ratio"`
	Ratio                float64 `json:"ratio"`
	CPM                  float64 `json:"cpm"`
	Cost                 float64 `json:"cost"`
	RPM                  float64 `json:"rpm"`
	DpRPM                float64 `json:"dpRpm"`
	Revenue              float64 `json:"Revenue"`
	GP                   float64 `json:"Gp"`
	GPP                  float64 `json:"Gpp"`
	PublisherBidRequests string  `json:"PublisherBidRequests"`
}

type AlertsEmails struct {
	AM           string             `json:"AM"`
	Email        string             `json:"Email"`
	FirstReport  AggregatedReport   `json:"FirstReport"`
	SecondReport []AggregatedReport `json:"SecondReport"`
}

type Worker struct {
	Cron          string                `json:"cron"`
	Slack         *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv   string                `json:"dbenv"`
	UserData      map[string]string
	CompassClient *compass.Compass
	skipInitRun   bool
	BCC           string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return err
	}

	credentialsMap, err := config.FetchConfigValues([]string{"rpm_decrease_alert"})
	if err != nil {
		return fmt.Errorf("error fetching config values: %w", err)
	}

	creds := credentialsMap["rpm_decrease_alert"]

	var emailConfig EmailCreds
	err = json.Unmarshal([]byte(creds), &emailConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.BCC = emailConfig.BCC

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	if worker.skipInitRun {
		worker.skipInitRun = false
		return nil
	}

	userData, err := email_reports.GetUsers(dto.UserTypeAccountManager)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	worker.UserData = userData
	report, err := getReport()
	if err != nil {
		return err
	}
	aggData := aggregate(report)
	avgData := computeAverage(aggData, worker)
	err = prepareAndSendEmail(avgData, worker)
	if err != nil {
		return fmt.Errorf("error sending rpm decrease email alerts %w", err)
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}

	return 0
}
