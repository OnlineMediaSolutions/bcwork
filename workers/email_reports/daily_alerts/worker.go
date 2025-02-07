package daily_alerts

import (
	"context"
	"fmt"
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
	Date         string  `json:"date"`
	DataStamp    int64   `json:"DateStamp"`
	Publisher    string  `json:"publisher"`
	Domain       string  `json:"domain"`
	PaymentType  string  `json:"PaymentType"`
	AM           string  `json:"am"`
	PubImps      int64   `json:"pub_imps"`
	LoopingRatio float64 `json:"looping_ratio"`
	Ratio        float64 `json:"estRatio"`
	CPM          float64 `json:"est_cpm"`
	Cost         float64 `json:"est_cost"`
	RPM          float64 `json:"mergedEstRpm"`
	DPRPM        float64 `json:"estDpRpm"`
	Revenue      float64 `json:"EstRevenue"`
	GP           float64 `json:"mergedEstGp"`
	GPP          float64 `json:"mergedEstGpp"`
}

type AlertsEmailRepo struct {
	AM           string             `json:"AM"`
	Email        string             `json:"Email"`
	FirstReport  AggregatedReport   `json:"FirstReport"`
	SecondReport []AggregatedReport `json:"SecondReport"`
}

type Report struct {
	Data struct {
		Result []Result `json:"result"`
	} `json:"data"`
}

type Result struct {
	Date         string  `json:"date"`
	DataStamp    int64   `json:"DateStamp"`
	Publisher    string  `json:"publisher"`
	Domain       string  `json:"domain"`
	PaymentType  string  `json:"PaymentType"`
	AM           string  `json:"am"`
	PubImps      int64   `json:"pub_imps"`
	LoopingRatio float64 `json:"looping_ratio"`
	Ratio        float64 `json:"estRatio"`
	CPM          float64 `json:"est_cpm"`
	Cost         float64 `json:"est_cost"`
	RPM          float64 `json:"mergedEstRpm"`
	DPRPM        float64 `json:"estDpRpm"`
	Revenue      float64 `json:"EstRevenue"`
	GP           float64 `json:"mergedEstGp"`
	GPP          float64 `json:"mergedEstGpp"`
}

type Worker struct {
	Cron          string                `json:"cron"`
	Slack         *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv   string                `json:"dbenv"`
	Start         string                `json:"start"`
	End           string                `json:"end"`
	Test          string                `json:"test"`
	ThreeHoursAgo int64                 `json:"three_hours_ago"`
	Report        Report
	AlertTypes    []string
	userService   *core.UserService
	CurrentTime   time.Time
	UserData      map[string]string
	CompassClient *compass.Compass
}

type Alert struct {
	Name             string `yaml:"name"`
	SlackMessageName string `yaml:"slack_message_name"`
}

type Config struct {
	Alerts []Alert `yaml:"alerts"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {

	loc, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(loc)
	worker.CurrentTime = now

	startString := now.Truncate(time.Hour).Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	endString := now.Truncate(time.Hour).Format("2006-01-02 15:04:05")

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

	worker.AlertTypes = alertTypes
	worker.Start = startString
	worker.End = endString
	worker.ThreeHoursAgo = now.Add(-3*time.Hour).Unix() / 100

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	userData, _ := worker.GetUsers()
	//worker.UserData = userData
	fmt.Println(userData, "userData")

	for _, alertType := range worker.AlertTypes {
		alert := GetAlerts(alertType)
		if alert != nil {
			report, err := alert.Request(worker)
			fmt.Println(report, "report")
			aggData := alert.Aggregate(report)
			avgData := alert.ComputeAverage(aggData, worker)
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
		Types: filter.String2DArrayFilter(filter.StringArrayFilter{"account_manager"}),
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
		userMap[user.FirstName+user.LastName] = user.Email
	}

	return userMap, nil

}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
