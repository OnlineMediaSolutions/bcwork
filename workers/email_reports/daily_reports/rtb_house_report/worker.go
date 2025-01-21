package rtb_house_report

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/bccron"
	"log"
	"time"
)

const (
	CompassRequestEndpoint = "/report-dashboard/report-new-bidder"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	skipInitRun bool
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

type RequestData struct {
	Data RequestDetails `json:"data"`
}

type RequestDetails struct {
	Date       Date     `json:"date"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}

type Date struct {
	Range    []string `json:"range"`
	Interval string   `json:"interval"`
}

type EmailCredentials struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	worker.Cron, _ = conf.GetStringValue("cron")

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	if err := bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return fmt.Errorf("failed to initialize DB for sellers: %w", err)
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	//if worker.skipInitRun {
	//	worker.skipInitRun = false
	//	return nil
	//}

	compassClient := compass.NewCompass()
	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-7 * 24 * time.Hour)

	requestData := PrepareRequestData()

	data, err := compassClient.Request(CompassRequestEndpoint, "POST", requestData, true, true)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Printf("Response: %+v\n", data)

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
