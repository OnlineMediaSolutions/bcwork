package rtb_house_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/bccron"
	"time"
)

const (
	CompassRequestEndpoint = "/report-dashboard/report-new-bidder"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	EmailCreds  map[string]string
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	skipInitRun bool
	ReportName  string
}

type RequestData struct {
	Data RequestDetails `json:"data"`
}

type RequestDetails struct {
	Date       Date     `json:"date"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
	Filters    Filters  `json:"filters"`
}

type Filters struct {
	DemandPartner []string `json:"DemandPartner"`
}

type Date struct {
	Range    []string `json:"range"`
	Interval string   `json:"interval"`
}

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

type Report struct {
	Time          string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	Revenue       float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	SoldImps      int     `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	DemandPartner string  `boil:"demand_partner" json:"demand_partner" toml:"demand_partner" yaml:"demand_partner"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	worker.Cron, _ = conf.GetStringValue("cron")
	worker.ReportName = "rtb_house_report"

	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	emailCredsMap, err := config.FetchConfigValues([]string{worker.ReportName})
	worker.EmailCreds = emailCredsMap

	if err != nil {
		return fmt.Errorf("failed to get email credentials %w", err)
	}

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return fmt.Errorf("failed initialize DB for real time report in environment %s,%w", worker.DatabaseEnv, err)
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	//if worker.skipInitRun {
	//	worker.skipInitRun = false
	//	return nil
	//}

	var emailCreds EmailCreds
	credsRaw := worker.EmailCreds[worker.ReportName]

	report, err := worker.getReportFromCompass()

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		return err
	}

	worker.prepareEmail(report, err, emailCreds)

	return nil
}

func (worker *Worker) getReportFromCompass() (map[string]interface{}, error) {
	compassClient := compass.NewCompass()
	requestData := PrepareRequestData()

	reportData, err := compassClient.Request(CompassRequestEndpoint, "POST", requestData, true, true)
	if err != nil {
		return nil, fmt.Errorf("request to Compass failed: %w", err)
	}

	return reportData, nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
