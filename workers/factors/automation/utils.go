package factors_autmation

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"math"
	"time"
)

// Factor strategy function
func (w *Worker) FactorStrategy(record *FactorReport, oldFactor float64) (float64, error) {
	var updatedFactor float64
	var GppOffset float64

	//STOP LOSS - Higher priority rule
	if record.Gp <= w.StopLoss {
		log.Warn().Msg(fmt.Sprintf("%s factor set to 0.75 because GP hit stop loss. GP: %f Stoploss: %f", record.Key(), record.Gp, w.StopLoss))
		return 0.75, nil //if we are losing more than 10$ in 30 minutes reduce to 0.75
	}

	//Check if the GPP Area is different for this domain
	_, exists := w.Domains[record.SetupKey()]
	if exists && w.Domains[record.SetupKey()].GppTarget != 0 {
		GppOffset = w.Domains[record.SetupKey()].GppTarget - w.GppTarget
	}

	//Calculate new factor
	if record.Gpp >= (w.GppTarget + 0.23 + GppOffset) {
		updatedFactor = oldFactor * 1.3
	} else if record.Gpp >= (w.GppTarget + 0.18 + GppOffset) {
		updatedFactor = oldFactor * 1.25
	} else if record.Gpp >= (w.GppTarget + 0.13 + GppOffset) {
		updatedFactor = oldFactor * 1.2
	} else if record.Gpp >= (w.GppTarget + 0.06 + GppOffset) {
		updatedFactor = oldFactor * 1.1
	} else if record.Gpp >= (w.GppTarget - 0.06 + +GppOffset) {
		updatedFactor = oldFactor // KEEP
	} else if record.Gpp >= (w.GppTarget - 0.17 + +GppOffset) {
		updatedFactor = oldFactor * 0.875
	} else if record.Gpp >= (w.GppTarget - 0.27 + +GppOffset) {
		updatedFactor = oldFactor * 0.8
	} else if record.Gpp >= (w.GppTarget - 0.37 + +GppOffset) {
		updatedFactor = oldFactor * 0.7
	} else {
		updatedFactor = oldFactor * 0.5
	}

	//Factor Ceiling
	if updatedFactor > w.MaxFactor {
		updatedFactor = w.MaxFactor
	}
	return roundFloat(updatedFactor), nil
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
		Increase:       roundFloat((record.NewFactor / record.OldFactor) - 1)}

	return model, nil
}

func (w *Worker) CheckDomain(record *FactorReport) bool {
	_, exists := w.Domains[record.SetupKey()]
	if exists {
		return true
	} else {
		return false
	}
}

func roundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func (w *Worker) GenerateTimes(minutes int) {
	w.End = time.Now().UTC().Truncate(time.Duration(minutes) * time.Minute)
	w.Start = w.End.Add(-time.Duration(minutes) * time.Minute)

}

func (w *Worker) AutomationDomains() []string {
	var domains []string
	for _, item := range w.Domains {
		domains = append(domains, item.Domain)
	}
	return domains
}
