package factors_autmation

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type Worker struct {
	Sleep         time.Duration           `json:"sleep"`
	DatabaseEnv   string                  `json:"dbenv"`
	Cron          string                  `json:"cron"`
	Domains       map[string]*DomainSetup `json:"domains"`
	StopLoss      float64                 `json:"stop_loss"`
	GppTarget     float64                 `json:"gpp_target"`
	MaxFactor     float64                 `json:"max_factor"`
	Quest         []string                `json:"quest_instances"`
	Start         time.Time               `json:"start"`
	End           time.Time               `json:"end"`
	DefaultFactor float64                 `json:"default_factor"`
}

// Worker functions
func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var err error
	var questExist bool

	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	worker.StopLoss, err = conf.GetFloat64ValueWithDefault("stoploss", -10)
	if err != nil {
		return errors.Wrapf(err, "failed to get stoploss value")
	}

	worker.GppTarget, err = conf.GetFloat64ValueWithDefault("gpp_target", 0.33)
	if err != nil {
		return errors.Wrapf(err, "failed to get GppTarget value")
	}

	worker.MaxFactor, err = conf.GetFloat64ValueWithDefault("max_factor", 10)
	if err != nil {
		return errors.Wrapf(err, "failed to get MaxFactor value")
	}

	worker.DefaultFactor, err = conf.GetFloat64ValueWithDefault("default_factor", 0.75)
	if err != nil {
		return errors.Wrapf(err, "failed to get stoploss value")
	}

	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	worker.Cron, _ = conf.GetStringValue("cron")

	return nil

}

func (worker *Worker) Do(ctx context.Context) error {
	var recordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var newFactors map[string]*FactorChanges
	var err error

	worker.GenerateTimes(30)

	recordsMap, factors, err = worker.FetchData(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to fetch data")
	}

	newFactors, err = worker.CalculateFactors(recordsMap, factors)
	if err != nil {
		return errors.Wrap(err, "failed to calculate factors")
	}

	err = UpdateAndLogChanges(ctx, newFactors)

	return nil
}

func (worker *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds", bccron.Next(worker.Cron)))
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
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

		oldFactor := factors[key].Factor // get current factor record
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
		}
	}

	return newFactors, nil
}

// Update the factors via API and push logs
func UpdateAndLogChanges(ctx context.Context, newFactors map[string]*FactorChanges) error {
	for _, rec := range newFactors {
		if rec.NewFactor != rec.OldFactor {
			err := rec.UpdateFactor()
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error Updating factor for key: Publisher=%s, Domain=%s, Country=%s, Device=%s. ResponseStatus: %d. err: %s", rec.Publisher, rec.Domain, rec.Country, rec.Device, rec.RespStatus, err))
			}
		}

		logJSON, err := json.Marshal(rec) //Create log json to log it
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error marshalling log for key:%v entry: %v", rec.Key(), err))

		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

		mod, err := rec.ToModel()
		if err != nil {
			log.Error().Err(err).Msg("failed to convert to model")
		}

		err = mod.Upsert(ctx, bcdb.DB(), true, Columns, boil.Infer(), boil.Infer())
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to push log to postgres. Err: %s", err))
		}
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
