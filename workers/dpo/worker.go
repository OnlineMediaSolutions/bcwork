package dpo

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/rotisserie/eris"
	"strings"
	"time"

	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/rs/zerolog"

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
	demandPartnerService      *core.DemandPartnerService
}

type DemandItem struct {
	ApiName        string  `json:"demand_partner_id"`
	Threshold      float64 `json:"threshold"`
	AutomationName string  `json:"automation_name"`
	Automation     bool    `json:"automation"`
}

// Worker functions
func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	err := worker.InitializeValues(ctx, conf)
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

	jsonData, err := worker.getDpFromDB(ctx, err)
	if err != nil {
		return err
	}

	if len(jsonData) == 0 {
		return nil
	}

	worker.Demands, err = worker.getDemandPartners(jsonData)

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
	errSlice := make([]string, 0)

	if len(dpoUpdate) == 0 && len(dpoDelete) == 0 {
		err := fmt.Errorf("No rules to update or delete")
		return eris.Wrapf(err, "No rules to update or delete")
	}

	err, dpoUpdate := worker.updateFactors(ctx, dpoUpdate, dpoDelete)
	if err != nil {
		message := fmt.Sprintf("Error bulk Updating dpo rules. err: %s", err.Error())
		errSlice = append(errSlice, message)
		log.Error().Msg(message)
	}

	err = UpsertLogs(ctx, dpoUpdate)
	if err != nil {
		message := fmt.Sprintf("Error Upserting logs into db. err: %s", err)
		errSlice = append(errSlice, message)
		log.Error().Msg(message)
	}

	if len(errSlice) != 0 {
		return errors.New(strings.Join(errSlice, "\n"))
	}
	return nil
}

// Utils
func (worker *Worker) InitializeValues(ctx context.Context, conf config.StringMap) error {
	errSlice := make([]string, 0)
	var err error
	var cronExists bool

	worker.LogSeverity, err = conf.GetIntValueWithDefault(config.LogSeverityKey, int(zerolog.InfoLevel))
	if err != nil {
		message := fmt.Sprintf("failed to set Log severity, err: %s", err)
		errSlice = append(errSlice, message)
	}
	zerolog.SetGlobalLevel(zerolog.Level(worker.LogSeverity))

	worker.httpClient = httpclient.New(true)

	historyModule := history.NewHistoryClient()
	worker.bulkService = bulk.NewBulkService(historyModule)
	worker.dpoService = core.NewDPOService(historyModule)
	worker.demandPartnerService = core.NewDemandPartnerService(historyModule)

	worker.Slack, err = messager.NewSlackModule()
	if err != nil {
		message := fmt.Sprintf("failed to initalize Slack module, err: %s", err)
		errSlice = append(errSlice, message)
	}

	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		message := fmt.Sprintf("failed to initalize Postgres DB. err: %s", err)
		errSlice = append(errSlice, message)
	}

	worker.Cron, cronExists = conf.GetStringValue("cron")
	if !cronExists {
		message := fmt.Sprintf("failed to get Cron. err: %s", err)
		errSlice = append(errSlice, message)
	}

	worker.RevenueThreshold, err = conf.GetFloat64ValueWithDefault("revenue_threshold", 5)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		errSlice = append(errSlice, message)
	}

	worker.DpRevenueThreshold, err = conf.GetFloat64ValueWithDefault("dp_revenue_threshold", 0.05)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		errSlice = append(errSlice, message)
	}

	worker.PlacementRevenueThreshold, err = conf.GetFloat64ValueWithDefault("revenue_threshold", 0.015)
	if err != nil {
		message := fmt.Sprintf("failed to get revenue_threshold. err: %s", err)
		errSlice = append(errSlice, message)
	}

	if len(errSlice) != 0 {
		return errors.New(strings.Join(errSlice, "\n"))
	}
	return nil

}

func (worker *Worker) getDemandPartners(demandData []*dto.DemandPartner) (map[string]*DemandSetup, error) {

	demands := make(map[string]*DemandSetup)

	for _, partner := range demandData {
		demands[partner.AutomationName] = &DemandSetup{
			Name:      partner.AutomationName,
			ApiName:   partner.DemandPartnerID,
			Threshold: partner.Threshold,
		}
	}

	return demands, nil
}

func (worker *Worker) getDpFromDB(ctx context.Context, err error) ([]*dto.DemandPartner, error) {
	filters := core.DemandPartnerGetFilter{
		Automation: filter.NewBoolFilter(true),
	}

	options := core.DemandPartnerGetOptions{
		Filter:     filters,
		Pagination: nil,
		Order:      nil,
		Selector:   "",
	}

	dpoDemand, err := worker.demandPartnerService.GetDemandPartners(ctx, &options)

	if err != nil {
		return nil, fmt.Errorf("Cannot get demand partners from database: %s\n", err)
	}

	return dpoDemand, nil
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
