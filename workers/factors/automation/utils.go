package factors_autmation

import (
	"errors"
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
	_, exists := GppAreas[record.Domain]
	if exists {
		GppOffset = GppAreas[record.Domain] - 0.33
	} else {
		GppOffset = 0
	}

	//Calculate new factor
	if record.Gpp >= (0.56 + GppOffset) {
		updatedFactor = oldFactor * 1.3
	} else if record.Gpp >= (0.51 + GppOffset) {
		updatedFactor = oldFactor * 1.25
	} else if record.Gpp >= (0.46 + GppOffset) {
		updatedFactor = oldFactor * 1.2
	} else if record.Gpp >= (0.39 + GppOffset) {
		updatedFactor = oldFactor * 1.1
	} else if record.Gpp >= (0.28 + GppOffset) {
		updatedFactor = oldFactor // KEEP
	} else if record.Gpp <= (-0.04 + GppOffset) {
		updatedFactor = oldFactor * 0.5
	} else if record.Gpp <= (0.06 + GppOffset) {
		updatedFactor = oldFactor * 0.7
	} else if record.Gpp <= (0.16 + GppOffset) {
		updatedFactor = oldFactor * 0.8
	} else if record.Gpp <= (0.28 + GppOffset) {
		updatedFactor = oldFactor * 0.875
	} else {
		return roundFloat(oldFactor), errors.New(fmt.Sprintf("unable to calculate factor: no matching condition Key: %s", record.Key()))
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

func (w *Worker) CheckDomain(targetDomain string) bool {
	if len(w.Domains) == 0 {
		return true
	}
	for _, item := range w.Domains {
		if targetDomain == item {
			return true
		}
	}
	return false
}

func roundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func (w *Worker) GenerateTimes(minutes int) {
	w.End = time.Now().UTC().Truncate(time.Duration(minutes) * time.Minute)
	w.Start = w.End.Add(-time.Duration(minutes) * time.Minute)

}
