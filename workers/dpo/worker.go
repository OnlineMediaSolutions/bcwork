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
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"time"
)

type Worker struct {
	Sleep       time.Duration           `json:"sleep"`
	DatabaseEnv string                  `json:"dbenv"`
	Cron        string                  `json:"cron"`
	Demands     map[string]*DemandSetup `json:"domains"`
	Start       time.Time               `json:"start"`
	End         time.Time               `json:"end"`
	Slack       *messager.SlackModule   `json:"slack_instances"`
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
	// PLAN:
	//FETCH REPORT
	//GROUP BY KEY
	//GROUP BY DP
	//
	//ITERATE THROUGH REPORT
	//DO 4 CHECKS, PUSH TO CHANGES
	//
	//APPLY ALL CHANGES
	//
	//CHECK EXISTING KEYS AND REACTIVATE? NOT NEEDED

	var newRules map[string]*DpoChanges
	var recordsMap map[string]*DpoReport
	var placementMap map[string]*PlacementReport
	var dpMap map[string]*DpReport
	var dpoApi map[string]*DpoApi

	worker.GenerateTimes()

	recordsMap, placementMap, dpMap, dpoApi, err := worker.FetchData(ctx)
	if err != nil {
		message := fmt.Sprintf("failed to fetch data at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
	}

	newRules, err = worker.CalculateRules(recordsMap, placementMap, dpMap, dpoApi)
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
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds. V1.3.4", bccron.Next(worker.Cron)))
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

// Function to calculate the new factors
func (worker *Worker) CalculateRules(recordsMap map[string]*DpoReport, placementMap map[string]*PlacementReport, dpMap map[string]*DpReport, dpoApi map[string]*DpoApi) (map[string]*DpoChanges, error) {
	var DpoUpdates = make(map[string]*DpoChanges)

	for _, record := range recordsMap {
		if worker.CheckDemand(record.DP) {
			oldFactor := 0.0
			key := record.Key()

			revenueFlag := record.Revenue < 5
			demandFlag := record.Revenue < (0.05 * dpMap[record.DP].Revenue)
			placementFlag := record.Revenue < (0.015 * placementMap[record.PlacementKey()].Revenue)
			erpmFlag := record.Erpm < worker.Demands[record.DP].Threshold

			item, exists := dpoApi[key]
			if exists {
				oldFactor = item.Factor
			}

			if revenueFlag && demandFlag && placementFlag && erpmFlag && oldFactor != 90 {
				DpoUpdates[key] = &DpoChanges{
					Time:       record.Time,
					EvalTime:   record.EvalTime,
					DP:         record.DP,
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
					DP:         record.DP,
					Domain:     record.Domain,
					Country:    record.Country,
					Os:         record.Os,
					Revenue:    record.Revenue,
					BidRequest: record.BidRequest,
					Erpm:       record.Erpm,
					OldFactor:  oldFactor,
					NewFactor:  1,
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

// Update the factors via API and push logs
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
	worker.Demands = make(map[string]*DemandSetup)
	worker.Demands["adaptmx"] = &DemandSetup{
		Name:      "adaptmx",
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
