package dpo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"time"
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
}

// Worker functions
func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
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
	var data DpoData
	var newRules map[string]*DpoChanges
	var err error

	worker.GenerateTimes()

	data = worker.FetchData(ctx)
	if data.Error != nil {
		message := fmt.Sprintf("failed to fetch data at %s: %s", worker.End.Format(constant.PostgresTimestampLayout), data.Error.Error())
		worker.Alert(message)
		return errors.Wrap(data.Error, message)
	}

	newRules, err = worker.CalculateRules(data)
	if err != nil {
		return err
	}

	err = UpdateAndLogChanges(ctx, newRules)
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds. V1.0", bccron.Next(worker.Cron)))
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

func (worker *Worker) CalculateRules(data DpoData) (map[string]*DpoChanges, error) {
	var DpoUpdates = make(map[string]*DpoChanges)

	for _, record := range data.DpoReport {
		if worker.CheckDemand(record.DP) && record.Country != "other" {
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
				DpoUpdates[key] = &DpoChanges{
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
			} else if exists && item.Factor != 0 && !erpmFlag {
				DpoUpdates[key] = &DpoChanges{
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
					NewFactor:  0.1,
				}
			}
		}
	}

	return DpoUpdates, nil
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
func UpdateAndLogChanges(ctx context.Context, newRules map[string]*DpoChanges) error {
	stringErrors := make([]string, 0)
	for _, record := range newRules {
		err := record.UpdateFactor()
		if err != nil {
			message := fmt.Sprintf("Error Updating factor for key: Publisher=%s, Domain=%s, Country=%s, Demand=%s. ResponseStatus: %d. err: %s", record.Publisher, record.Domain, record.Country, record.DP, record.RespStatus, err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)
		}

		logJSON, err := json.Marshal(record) //Create log json to log it
		if err != nil {
			message := fmt.Sprintf("Error marshalling log for key:%v entry: %v", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)

		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

		mod, err := record.ToModel()
		if err != nil {
			message := fmt.Sprintf("failed to convert to model for key:%v. error: %v", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Msg(message)
		}

		err = mod.Upsert(ctx, bcdb.DB(), true, Columns, boil.Infer(), boil.Infer())
		if err != nil {
			message := fmt.Sprintf("failed to push log to postgres for key %s. Err: %s", record.Key(), err)
			stringErrors = append(stringErrors, message)
			log.Error().Err(err).Msg(message)
		}
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
