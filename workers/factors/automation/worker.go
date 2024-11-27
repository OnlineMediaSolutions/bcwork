package factors_autmation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/modules/history"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	Sleep                   time.Duration           `json:"sleep"`
	DatabaseEnv             string                  `json:"dbenv"`
	Cron                    string                  `json:"cron"`
	Domains                 map[string]*DomainSetup `json:"domains"`
	StopLoss                float64                 `json:"stop_loss"`
	GppTarget               float64                 `json:"gpp_target"`
	MaxFactor               float64                 `json:"max_factor"`
	MinFactor               float64                 `json:"min_factor"`
	InactiveDaysThreshold   int                     `json:"inactive_days"`
	InactiveFactorThreshold float64                 `json:"inactive_factor"`
	InactiveKeys            []string                `json:"inactive_keys"`
	Quest                   []string                `json:"quest_instances"`
	Start                   time.Time               `json:"start"`
	End                     time.Time               `json:"end"`
	Fees                    map[string]float64      `json:"fees"`
	ConsultantFees          map[string]float64      `json:"consultant_fees"`
	DefaultFactor           float64                 `json:"default_factor"`
	Slack                   *messager.SlackModule   `json:"slack_instances"`
	HttpClient              httpclient.Doer         `json:"http_client"`
	BulkService             *bulk.BulkService
	skipInitRun             bool
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

	var recordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var newFactors map[string]*FactorChanges
	var err error

	worker.GenerateTimes(30)

	recordsMap, factors, err = worker.FetchData(ctx)
	if err != nil {
		message := fmt.Sprintf("failed to fetch data at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
	}

	newFactors, err = worker.CalculateFactors(recordsMap, factors)
	if err != nil {
		message := fmt.Sprintf("failed to calculate factors at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
	}

	err = worker.UpdateAndLogChanges(ctx, newFactors)
	if err != nil {
		message := fmt.Sprintf("error updating and log changes at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error())
		worker.Alert(message)
		return errors.Wrap(err, message)
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

func (worker *Worker) InitializeValues(conf config.StringMap) error {
	stringErrors := make([]string, 0)
	var err error
	var questExist bool
	var cronExists bool

	worker.HttpClient = httpclient.New(true)

	worker.Slack, err = messager.NewSlackModule()
	if err != nil {
		message := fmt.Sprintf("failed to initalize Slack module, err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	worker.StopLoss, err = conf.GetFloat64ValueWithDefault("stoploss", -10)
	if err != nil {
		message := fmt.Sprintf("failed to get stoploss value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.GppTarget, err = worker.FetchGppTargetDefault()
	if err != nil {
		message := fmt.Sprintf("failed to fetch gpp target value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.InactiveDaysThreshold, err = conf.GetIntValueWithDefault("inactive_days", 3)
	if err != nil {
		message := fmt.Sprintf("failed to get inactive days value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.InactiveFactorThreshold, err = conf.GetFloat64ValueWithDefault("inactive_factor", 0.2)
	if err != nil {
		message := fmt.Sprintf("failed to get inactive days value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.MaxFactor, err = conf.GetFloat64ValueWithDefault("max_factor", 10)
	if err != nil {
		message := fmt.Sprintf("failed to get MaxFactor value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.MinFactor, err = conf.GetFloat64ValueWithDefault("min_factor", 0.5)
	if err != nil {
		message := fmt.Sprintf("failed to get MinFactor value. err: %s", err)
		stringErrors = append(stringErrors, message)
	}

	worker.DefaultFactor, err = conf.GetFloat64ValueWithDefault("default_factor", 0.75)
	if err != nil {
		message := fmt.Sprintf("failed to get stoploss value. err: %s", err)
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

	historyModule := history.NewHistoryClient()
	worker.BulkService = bulk.NewBulkService(historyModule)

	if len(stringErrors) != 0 {
		return errors.New(strings.Join(stringErrors, "\n"))
	}
	return nil

}

// Function to calculate the new factors
func (worker *Worker) CalculateFactors(RecordsMap map[string]*FactorReport, factors map[string]*Factor) (map[string]*FactorChanges, error) {
	var err error
	var newFactors = make(map[string]*FactorChanges)

	for _, record := range RecordsMap {
		// Check if the key exists on the first half as well
		if !worker.CheckDomain(record) {
			continue
		}

		// Check if the key exists on factors
		key := record.Key()
		_, exists := factors[key]
		if !exists {
			continue
		}

		oldFactor := factors[key].Factor // get current factor value
		ruleId := factors[key].RuleId

		var updatedFactor float64
		updatedFactor, err = worker.FactorStrategy(record, oldFactor)
		if err != nil {
			log.Err(err).Msg("failed to calculate factor")
			logJSON, err := json.Marshal(record)
			if err != nil {
				log.Err(err).Msg("failed to parse record to json.")
				return nil, err
			}
			log.Info().Msg(fmt.Sprintf("%s", logJSON))
		}

		newFactors[key] = &FactorChanges{
			Time:      worker.End,
			EvalTime:  worker.Start,
			Pubimps:   record.PublisherImpressions,
			Soldimps:  record.SoldImpressions,
			Cost:      RoundFloat(record.Cost + record.DataFee + record.DemandPartnerFee),
			Revenue:   RoundFloat(record.Revenue),
			GP:        record.Gp,
			GPP:       record.Gpp,
			Publisher: factors[key].Publisher,
			Domain:    factors[key].Domain,
			Country:   factors[key].Country,
			Device:    factors[key].Device,
			OldFactor: factors[key].Factor,
			NewFactor: updatedFactor,
			RuleId:    ruleId,
		}
	}

	return newFactors, nil
}

func (worker *Worker) UpdateAndLogChanges(ctx context.Context, newFactors map[string]*FactorChanges) error {
	stringErrors := make([]string, 0)

	err, newFactors := worker.UpdateFactors(ctx, newFactors)
	if err != nil {
		message := fmt.Sprintf("Error bulk Updating factors. err: %s", err.Error())
		stringErrors = append(stringErrors, message)
		log.Error().Msg(message)
	}

	err = UpsertLogs(ctx, newFactors)
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

// Columns variable to check conflict on the price_factor_log table
var Columns = []string{
	models.PriceFactorLogColumns.Time,
	models.PriceFactorLogColumns.Publisher,
	models.PriceFactorLogColumns.Domain,
	models.PriceFactorLogColumns.Country,
	models.PriceFactorLogColumns.Device,
}
