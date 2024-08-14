package factors_monitor

import (
	"context"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/workers/factors/automation"
	"github.com/rs/zerolog/log"
	"time"
)

type Worker struct {
	Sleep            time.Duration                             `json:"sleep"`
	DatabaseEnv      string                                    `json:"dbenv"`
	Cron             string                                    `json:"cron"`
	Domains          map[string]*factors_autmation.DomainSetup `json:"domains"`
	StopLoss         float64                                   `json:"stop_loss"`
	Quest            []string                                  `json:"quest_instances"`
	DefaultFactor    float64                                   `json:"default_factor"`
	Start            time.Time                                 `json:"start"`
	End              time.Time                                 `json:"end"`
	AutomationWorker factors_autmation.Worker                  `json:"automation_worker"`
}

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
	var recordsMap map[string]*factors_autmation.FactorReport
	var factors map[string]*factors_autmation.Factor
	var err error

	worker.GenerateTimes()
	worker.InitializeAutomationWorker()

	recordsMap, factors, err = worker.AutomationWorker.FetchData(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch data")
	}

	recordsMap["tetetete"] = &factors_autmation.FactorReport{
		Time:                 time.Now(), // Set this to the current time or any other time as required
		PublisherID:          "te",       // Replace with actual PublisherID
		Domain:               "te",       // Replace with actual Domain
		Country:              "te",       // Replace with actual Country
		DeviceType:           "te",       // Replace with actual DeviceType
		Revenue:              100.0,      // Replace with actual Revenue
		Cost:                 50.0,       // Replace with actual Cost
		DemandPartnerFee:     5.0,        // Replace with actual DemandPartnerFee
		SoldImpressions:      1000,       // Replace with actual SoldImpressions
		PublisherImpressions: 1200,       // Replace with actual PublisherImpressions
		DataFee:              2.0,        // Replace with actual DataFee
		Gp:                   -12.0,      // Replace with actual Gp
		Gpp:                  -15.0,      // This is the key requirement
	}

	factors["tetetete"] = &factors_autmation.Factor{
		Publisher: "te",
		Domain:    "te",
		Country:   "te",
		Device:    "te",
		Factor:    6,
	}

	newFactors := worker.CalculateFactors(recordsMap, factors)

	err = factors_autmation.UpdateAndLogChanges(ctx, newFactors)
	if err != nil {
		return errors.Wrapf(err, "failed to update factors and log changes")
	}
	return nil
}

func (worker *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds", bccron.Next(worker.Cron)))
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}

func (worker *Worker) GenerateTimes() {
	worker.Start = time.Now().UTC().Truncate(time.Duration(30) * time.Minute)
	worker.End = time.Now().UTC()
}

func (worker *Worker) InitializeAutomationWorker() {
	worker.AutomationWorker = factors_autmation.Worker{
		DatabaseEnv:   worker.DatabaseEnv,
		Cron:          worker.Cron,
		Domains:       worker.Domains,
		StopLoss:      worker.StopLoss,
		Quest:         worker.Quest,
		DefaultFactor: worker.DefaultFactor,
		Start:         worker.Start,
		End:           worker.End,
	}
}

func (worker *Worker) CalculateFactors(recordsMap map[string]*factors_autmation.FactorReport, factors map[string]*factors_autmation.Factor) map[string]*factors_autmation.FactorChanges {
	newFactors := make(map[string]*factors_autmation.FactorChanges)
	for _, record := range recordsMap {
		// Check if the key exists on factors
		key := record.Key()
		_, exists := factors[key]
		if !exists {
			continue
		}

		//Check if the record is below the stop loss
		if record.Gp > worker.StopLoss {
			continue
		}

		newFactors[key] = &factors_autmation.FactorChanges{
			Time:      worker.End,
			EvalTime:  worker.Start,
			Pubimps:   record.PublisherImpressions,
			Soldimps:  record.SoldImpressions,
			Cost:      factors_autmation.RoundFloat(record.Cost + record.DataFee + record.DemandPartnerFee),
			Revenue:   factors_autmation.RoundFloat(record.Revenue),
			GP:        record.Gp,
			GPP:       record.Gpp,
			Publisher: factors[key].Publisher,
			Domain:    factors[key].Domain,
			Country:   factors[key].Country,
			Device:    factors[key].Device,
			OldFactor: factors[key].Factor,
			NewFactor: worker.DefaultFactor,
			Source:    "monitor-system",
		}
	}
	return newFactors
}
