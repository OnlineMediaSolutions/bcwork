package dpo

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/core"
	"github.com/rs/zerolog"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	Sleep                     time.Duration           `json:"sleep"`
	DatabaseEnv               string                  `json:"dbenv"`
	Cron                      string                  `json:"cron"`
	Demands                   map[string]*DemandSetup `json:"domains"`
	Start                     time.Time               `json:"start"`
	End                       time.Time               `json:"end"`
	RevenueThreshold          float64                 `json:"revenue_threshold"`
	DpRevenueThreshold        float64                 `json:"dp_revenue_threshold"`
	PlacementRevenueThreshold float64                 `json:"placement_revenue_threshold"`
	Slack                     *messager.SlackModule   `json:"slack_instances"`
	LogSeverity               int                     `json:"logsev"`
	httpClient                httpclient.Doer
	skipInitRun               bool
	bulkService               *bulk.BulkService
	dpoService                *core.DPOService
}

// Worker functions
func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	err := worker.InitializeValues(conf)
	if err != nil {
		message := fmt.Sprintf("failed to initialize values. Error: %s", err.Error())
		log.Error().Msg(message)
		worker.Alert(message)
		return errors.New(message)
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	if worker.skipInitRun {
		fmt.Println("Skipping work as per the skip_init_run flag.")
		worker.skipInitRun = false
		return nil
	}

	var data DpoData
	var ruleUpdate map[string]*DpoChanges
	var ruleDelete map[string]*DpoChanges
	var err error

	worker.GenerateTimes()

	data = worker.FetchData(ctx)
	if data.Error != nil {
		message := fmt.Sprintf("failed to fetch data at %s: %s", worker.End.Format(constant.PostgresTimestampLayout), data.Error.Error())
		worker.Alert(message)
		return errors.Wrap(data.Error, message)
	}

	ruleUpdate, ruleDelete, err = worker.calculateRules(data)
	if err != nil {
		message := fmt.Sprintf("failed to calculate rules. Error: %s", err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
	}

	err = worker.UpdateAndLogChanges(ctx, ruleUpdate, ruleDelete)
	if err != nil {
		message := fmt.Sprintf("Error updating and logging changes. Error: %s", err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds. V1.1", bccron.Next(worker.Cron)))
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

func (worker *Worker) calculateRules(data DpoData) (map[string]*DpoChanges, map[string]*DpoChanges, error) {
	var dpoUpdates = make(map[string]*DpoChanges)
	var dpoDeletes = make(map[string]*DpoChanges)

	for _, record := range data.DpoReport {
		if worker.CheckDemand(record.DP) && record.Country != "other" && record.Os != "-" {
			oldFactor := 0.0
			key := record.Key()
			apiKey := record.ApiKey()

			revenueFlag := record.Revenue < worker.RevenueThreshold
			demandFlag := record.Revenue < (worker.DpRevenueThreshold * data.DpReport[record.DP].Revenue)
			placementFlag := record.Revenue < (worker.PlacementRevenueThreshold * data.PlacementReport[record.PlacementKey()].Revenue)
			erpmFlag := record.Erpm < worker.Demands[record.DP].Threshold

			item, exists := data.DpoApi[apiKey]
			if exists {
				oldFactor = item.Factor
			}

			if revenueFlag && demandFlag && placementFlag && erpmFlag && oldFactor != 90 {
				dpoUpdates[key] = &DpoChanges{
					Time:       record.Time,
					EvalTime:   record.EvalTime,
					Publisher:  record.Publisher,
					DP:         worker.Demands[record.DP].ApiName,
					Domain:     record.Domain,
					Country:    record.Country,
					Os:         record.Os,
					Revenue:    record.Revenue,
					BidRequest: record.BidRequest,
					Erpm:       record.Erpm,
					OldFactor:  oldFactor,
					NewFactor:  90,
				}
			} else if exists && !erpmFlag {
				dpoDeletes[key] = &DpoChanges{
					Time:       record.Time,
					EvalTime:   record.EvalTime,
					Publisher:  record.Publisher,
					DP:         worker.Demands[record.DP].ApiName,
					Domain:     record.Domain,
					Country:    record.Country,
					Os:         record.Os,
					Revenue:    record.Revenue,
					BidRequest: record.BidRequest,
					Erpm:       record.Erpm,
					OldFactor:  oldFactor,
					NewFactor:  0,
					RuleId:     item.RuleId,
				}
			}
		}
	}

	return dpoUpdates, dpoDeletes, nil
}

// Columns variable to check conflict on the price_factor_log table
var Columns = []string{
	models.DpoAutomationLogColumns.Time,
	models.DpoAutomationLogColumns.Publisher,
	models.DpoAutomationLogColumns.Domain,
	models.DpoAutomationLogColumns.Country,
	models.DpoAutomationLogColumns.Os,
	models.DpoAutomationLogColumns.DP,
}

// Update the Dpo Rules via API and push logs
func (worker *Worker) UpdateAndLogChanges(ctx context.Context, dpoUpdate map[string]*DpoChanges, dpoDelete map[string]*DpoChanges) error {
	stringErrors := make([]string, 0)

	err, dpoUpdate := worker.updateFactors(ctx, dpoUpdate, dpoDelete)
	if err != nil {
		message := fmt.Sprintf("Error bulk Updating dpo rules. err: %s", err.Error())
		stringErrors = append(stringErrors, message)
		log.Error().Msg(message)
	}

	err = UpsertLogs(ctx, dpoUpdate)
	if err != nil {
		message := fmt.Sprintf("Error Upserting logs into db. err: %s", err)
		stringErrors = append(stringErrors, message)
		log.Error().Msg(message)
	}

	if len(stringErrors) != 0 {
		return errors.New(strings.Join(stringErrors, "\n"))
	}
	return nil
}

// Utils
func (worker *Worker) InitializeValues(conf config.StringMap) error {
	stringErrors := make([]string, 0)
	var err error
	var cronExists bool

	worker.LogSeverity, err = conf.GetIntValueWithDefault(config.LogSeverityKey, int(zerolog.InfoLevel))
	if err != nil {
		message := fmt.Sprintf("failed to set Log severity, err: %s", err)
		stringErrors = append(stringErrors, message)
	}
	zerolog.SetGlobalLevel(zerolog.Level(worker.LogSeverity))

	worker.httpClient = httpclient.New(true)

	historyModule := history.NewHistoryClient()
	worker.bulkService = bulk.NewBulkService(historyModule)
	worker.dpoService = core.NewDPOService(historyModule)

	worker.Slack, err = messager.NewSlackModule()
	if err != nil {
		message := fmt.Sprintf("failed to initalize Slack module, err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		message := fmt.Sprintf("failed to initalize Postgres DB. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.Cron, cronExists = conf.GetStringValue("cron")
	if !cronExists {
		message := fmt.Sprintf("failed to get Cron. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.RevenueThreshold, err = conf.GetFloat64ValueWithDefault("revenue_threshold", 5)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.DpRevenueThreshold, err = conf.GetFloat64ValueWithDefault("dp_revenue_threshold", 0.05)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.PlacementRevenueThreshold, err = conf.GetFloat64ValueWithDefault("revenue_threshold", 0.015)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.Demands = make(map[string]*DemandSetup)
	worker.Demands["onetag-bcm"] = &DemandSetup{
		Name:      "onetag-bcm",
		ApiName:   "onetagbcm",
		Threshold: 0.001,
	}
	worker.Demands["pubmatic-pbs"] = &DemandSetup{
		Name:      "pubmatic-pbs",
		ApiName:   "pubmaticbcm",
		Threshold: 0.001,
	}
	worker.Demands["index-pbs"] = &DemandSetup{
		Name:      "index-pbs",
		ApiName:   "indexs2s",
		Threshold: 0.001,
	}
	//worker.Demands["appnexusbcm"] = &DemandSetup{
	//	Name:      "appnexusbcm",
	//	ApiName:   "appnexusbcm",
	//	Threshold: 0.002,
	//}
	worker.Demands["yieldmo-audienciad"] = &DemandSetup{
		Name:      "yieldmo-audienciad",
		ApiName:   "yieldmo",
		Threshold: 0.001,
	}
	worker.Demands["sovrn"] = &DemandSetup{
		Name:      "sovrn",
		ApiName:   "sovrnbcm",
		Threshold: 0.001,
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

func (worker *Worker) GenerateTimes() {
	worker.End = time.Now().UTC().Truncate(time.Hour)
	worker.Start = worker.End.Add(-time.Duration(1) * time.Hour)
}

func (worker *Worker) DemandPartners() []string {
	var demands []string
	for _, item := range worker.Demands {
		demands = append(demands, item.Name)
	}
	return demands
}
