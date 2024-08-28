package factors_autmation

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"math"
	"time"
)

// Factor strategy function
func (worker *Worker) FactorStrategy(record *FactorReport, oldFactor float64) (float64, error) {
	var updatedFactor float64
	var GppOffset float64

	//STOP LOSS - Higher priority rule
	if record.Gp <= worker.StopLoss {
		message := fmt.Sprintf("%s factor set to %f because GP hit stop loss. GP: %f Stoploss: %f", record.Key(), worker.DefaultFactor, record.Gp, worker.StopLoss)
		worker.Alert(message)
		log.Warn().Msg(message)
		return worker.DefaultFactor, nil //if we are losing more than 10$ in 30 minutes reduce to default factor (0.75)
	}

	if worker.CheckInactiveKey(record) {
		return worker.DefaultFactor, nil
	}

	//Check if the GPP Target is different for this domain
	_, exists := worker.Domains[record.SetupKey()]
	if exists && worker.Domains[record.SetupKey()].GppTarget != 0 {
		GppOffset = worker.Domains[record.SetupKey()].GppTarget - worker.GppTarget
	}

	//Calculate new factor
	if record.Gpp >= (worker.GppTarget + 0.23 + GppOffset) {
		updatedFactor = oldFactor * 1.3
	} else if record.Gpp >= (worker.GppTarget + 0.18 + GppOffset) {
		updatedFactor = oldFactor * 1.25
	} else if record.Gpp >= (worker.GppTarget + 0.13 + GppOffset) {
		updatedFactor = oldFactor * 1.2
	} else if record.Gpp >= (worker.GppTarget + 0.06 + GppOffset) {
		updatedFactor = oldFactor * 1.1
	} else if record.Gpp >= (worker.GppTarget - 0.06 + +GppOffset) {
		updatedFactor = oldFactor // KEEP
	} else if record.Gpp >= (worker.GppTarget - 0.17 + +GppOffset) {
		updatedFactor = oldFactor * 0.875
	} else if record.Gpp >= (worker.GppTarget - 0.27 + +GppOffset) {
		updatedFactor = oldFactor * 0.8
	} else if record.Gpp >= (worker.GppTarget - 0.37 + +GppOffset) {
		updatedFactor = oldFactor * 0.7
	} else {
		updatedFactor = oldFactor * 0.5
	}

	//Factor Ceiling
	if updatedFactor > worker.MaxFactor {
		updatedFactor = worker.MaxFactor
	}
	return RoundFloat(updatedFactor), nil
}

func (record *FactorChanges) ToModel() (models.PriceFactorLog, error) {
	model := models.PriceFactorLog{
		Time:           record.Time,
		EvalTime:       record.EvalTime,
		Pubimps:        record.Pubimps,
		Soldimps:       record.Soldimps,
		Cost:           record.Cost,
		Revenue:        record.Revenue,
		GP:             record.GP,
		GPP:            record.GPP,
		Publisher:      record.Publisher,
		Domain:         record.Domain,
		Country:        record.Country,
		Device:         record.Device,
		OldFactor:      record.OldFactor,
		NewFactor:      record.NewFactor,
		ResponseStatus: record.RespStatus,
		Increase:       RoundFloat((record.NewFactor / record.OldFactor) - 1)}

	return model, nil
}

func (worker *Worker) CheckDomain(record *FactorReport) bool {
	_, exists := worker.Domains[record.SetupKey()]
	if exists {
		return true
	} else {
		return false
	}
}

func RoundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func (worker *Worker) GenerateTimes(minutes int) {
	worker.End = time.Now().UTC().Truncate(time.Duration(minutes) * time.Minute)
	worker.Start = worker.End.Add(-time.Duration(minutes) * time.Minute)

}

func (worker *Worker) AutomationDomains() []string {
	var domains []string
	for _, item := range worker.Domains {
		domains = append(domains, item.Domain)
	}
	return domains
}

func (worker *Worker) Alert(message string) {
	err := worker.Slack.SendMessage(message)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error sending slack alert: %s", err))
	}
}

func (worker *Worker) CheckInactiveKey(record *FactorReport) bool {
	for _, key := range worker.InactiveKeys {
		if record.Key() == key {
			return true
		}
	}
	return false
}
