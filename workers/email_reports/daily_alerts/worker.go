package daily_alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type EmailData struct {
	Body   string
	Report []AlertsEmailRepo
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

type AlertsEmailRepo struct {
	AM           string             `json:"AM"`
	Email        string             `json:"Email"`
	FirstReport  AggregatedReport   `json:"FirstReport"`
	SecondReport []AggregatedReport `json:"SecondReport"`
}

type Email struct {
	Bcc []string `json:"bcc"`
}

type Worker struct {
	Cron                string                `json:"cron"`
	Slack               *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv         string                `json:"dbenv"`
	Start               string                `json:"start"`
	End                 string                `json:"end"`
	StartOfLastWeek     int64                 `json:"start_of_last_week"`
	EndOfLastWeek       int64                 `json:"end_of_last_week"`
	StartOfLastWeekUnix int64                 `json:"start_of_last_week_unix"`
	StartOfLastWeekStr  string                `json:"start_of_last_week_str"`
	EndOfLastWeekStr    string                `json:"end_of_last_week_str"`
	Yesterday           int64                 `json:"yesterday"`
	Today               int64                 `json:"today"`
	Test                string                `json:"test"`
	ThreeHoursAgo       int64                 `json:"three_hours_ago"`
	AlertTypes          []string
	userService         *core.UserService
	CurrentTime         time.Time
	UserData            map[string]string
	CompassClient       *compass.Compass
	skipInitRun         bool
	BCC                 string
}

type Alert struct {
	Name         string `yaml:"name"`
	EmailSubject string `yaml:"email_subject"`
}

type Config struct {
	Alerts []Alert `yaml:"alerts"`
}

const (
	TimeZone = "America/New_York"
)

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	loc, _ := time.LoadLocation(TimeZone)
	now := time.Now().In(loc)

	data, err := os.ReadFile("workers/email_reports/daily_alerts/daily_alerts.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var configuration Config
	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	var alertTypes []string
	for _, alert := range configuration.Alerts {
		alertTypes = append(alertTypes, alert.Name)
	}

	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err = bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return err
	}

	worker.AlertTypes = alertTypes
	worker.CurrentTime = now
	worker.Cron, _ = conf.GetStringValue("cron")
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	if worker.skipInitRun {
		worker.skipInitRun = false
		return nil
	}

	userData, _ := worker.GetUsers()
	worker.UserData = userData

	for _, alertType := range worker.AlertTypes {
		alert := GetAlerts(alertType)
		if alert != nil {
			report, err := alert.Request(worker)
			aggData := alert.Aggregate(report)
			avgData := alert.ComputeAverage(aggData, worker)
			err = WriteAlertsToJSON(avgData, "alerts.json")
			if err != nil {
				return err
			}
			err = alert.PrepareAndSendEmail(avgData, worker)
			if err != nil {
				fmt.Println("Error sending email alerts:", err)
			}
		} else {
			fmt.Println("Alert type not found.")
		}
	}
	return nil
}

func (worker *Worker) GetUsers() (map[string]string, error) {
	filters := core.UserFilter{
		Types: filter.String2DArrayFilter(filter.StringArrayFilter{"Account Manager"}),
	}

	options := core.UserOptions{
		Filter:     filters,
		Pagination: nil,
		Order:      nil,
		Selector:   "",
	}

	users, err := worker.userService.GetUsers(context.Background(), &options)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)

	for _, user := range users {
		key := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		userMap[key] = user.Email
	}

	return userMap, nil

}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

func WriteAlertsToJSON(alerts map[string][]AlertsEmailRepo, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(alerts)
	if err != nil {
		return fmt.Errorf("error writing JSON: %v", err)
	}

	fmt.Println("Alerts written to", filename)
	return nil
}
