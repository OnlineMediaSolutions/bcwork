package factors_monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	factors_automation "github.com/m6yf/bcwork/workers/factors/automation"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	Sleep            time.Duration                              `json:"sleep"`
	DatabaseEnv      string                                     `json:"dbenv"`
	Cron             string                                     `json:"cron"`
	Domains          map[string]*factors_automation.DomainSetup `json:"domains"`
	StopLoss         float64                                    `json:"stop_loss"`
	Quest            []string                                   `json:"quest_instances"`
	DefaultFactor    float64                                    `json:"default_factor"`
	Start            time.Time                                  `json:"start"`
	End              time.Time                                  `json:"end"`
	AutomationWorker factors_automation.Worker                  `json:"automation_worker"`
	Slack            *messager.SlackModule                      `json:"slack_instances"`
	httpClient       httpclient.Doer
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

	worker.Slack, err = messager.NewSlackModule()
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("failed to initalize Slack module, err: %s", err))
	}

	worker.httpClient = httpclient.New(true)

	return nil

}

func (worker *Worker) Do(ctx context.Context) error {
	var recordsMap map[string]*factors_automation.FactorReport
	var factors map[string]*factors_automation.Factor
	var err error

	worker.GenerateTimes()
	worker.InitializeAutomationWorker()

	recordsMap, factors, err = worker.AutomationWorker.FetchData(ctx)
	if err != nil {
		worker.AutomationWorker.Alert(fmt.Sprintf("FACTOR MONITORING: failed to fetch data at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error()))
		return errors.Wrapf(err, "failed to fetch data")
	}

	newFactors := worker.CalculateFactors(recordsMap, factors)
	if len(newFactors) == 0 {
		log.Info().Msg("No stop loss cases were found")
		return nil
	}

	alert, err := GenerateStopLossAlerts(newFactors)
	if err != nil {
		worker.AutomationWorker.Alert(fmt.Sprintf("Could not generate stoploss alerts, but there are stoploss cases. Err: %s", err.Error()))
	} else {
		worker.AutomationWorker.Alert(alert)
	}

	err = worker.AutomationWorker.UpdateAndLogChanges(ctx, newFactors)
	if err != nil {
		worker.AutomationWorker.Alert(fmt.Sprintf("FACTOR MONITORING: error updating and log changes at %s: %s", worker.End.Format("2006-01-02T15:04:05Z"), err.Error()))
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
	worker.AutomationWorker = factors_automation.Worker{
		DatabaseEnv:   worker.DatabaseEnv,
		Cron:          worker.Cron,
		Domains:       worker.Domains,
		StopLoss:      worker.StopLoss,
		Quest:         worker.Quest,
		DefaultFactor: worker.DefaultFactor,
		Start:         worker.Start,
		End:           worker.End,
		Slack:         worker.Slack,
		HttpClient:    worker.httpClient,
	}
}

func (worker *Worker) CalculateFactors(recordsMap map[string]*factors_automation.FactorReport, factors map[string]*factors_automation.Factor) map[string]*factors_automation.FactorChanges {
	newFactors := make(map[string]*factors_automation.FactorChanges)
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

		newFactors[key] = &factors_automation.FactorChanges{
			Time:      worker.End,
			EvalTime:  worker.Start,
			Pubimps:   record.PublisherImpressions,
			Soldimps:  record.SoldImpressions,
			Cost:      factors_automation.RoundFloat(record.Cost + record.DataFee + record.DemandPartnerFee),
			Revenue:   factors_automation.RoundFloat(record.Revenue),
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

func GenerateStopLossAlerts(changesMap map[string]*factors_automation.FactorChanges) (string, error) {
	changesArr := make([]string, 0)
	changesArr = append(changesArr, "FACTOR MONITOR - STOP LOSS WAS HIT IN THE FOLLOWING CASES:")
	for _, item := range changesMap {
		logJSON, err := json.Marshal(item) //Create json to log it
		if err != nil {
			return "", errors.Wrapf(err, "error generating alerts array")
		}
		changesArr = append(changesArr, fmt.Sprintf("%s", logJSON))
	}
	return strings.Join(changesArr, "\n"), nil
}
