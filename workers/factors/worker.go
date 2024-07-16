package factors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/quest"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"math"
	"net/http"
	"time"
)

type Worker struct {
	Sleep          time.Duration `json:"sleep"`
	DatabaseEnv    string        `json:"dbenv"`
	Cron           string        `json:"cron"`
	MaxGPP         float64       `json:"max_gpp"`
	MinGPP         float64       `json:"min_gpp"`
	TrendThreshold float64       `json:"trend_threshold"`
	FactorStep     float64       `json:"factor_step"`
	LowFactor      float64       `json:"low_factor"`
	MinFactor      float64       `json:"min_factor"`
	Domains        []string      `json:"domains"`
	FilterExists   bool          `json:"filter_exists"`
}

type FactorChanges struct {
	Time       time.Time `json:"time"`
	EvalTime   time.Time `json:"t1"`
	GP         float64   `json:"gp"`
	GPP        float64   `json:"gpp"`
	Publisher  string    `json:"publisher"`
	Domain     string    `json:"domain"`
	Country    string    `json:"country"`
	Device     string    `json:"device"`
	OldFactor  float64   `json:"old_factor"`
	NewFactor  float64   `json:"new_factor"`
	RespStatus int       `json:"response_status"`
}

//type FactorChanges struct {
//	Time       time.Time `json:"time"`
//	T1         time.Time `json:"t1"`
//	T2         time.Time `json:"t2"`
//	T1GP       float64   `json:"t1_gp"`
//	T2GP       float64   `json:"t2_gp"`
//	T2GPP      float64   `json:"t2_gpp"`
//	Publisher  string    `json:"publisher"`
//	Domain     string    `json:"domain"`
//	Country    string    `json:"country"`
//	Device     string    `json:"device"`
//	OldFactor  float64   `json:"old_factor"`
//	NewFactor  float64   `json:"new_factor"`
//	RespStatus float64   `json:"response_status"`
//}

func (record *FactorChanges) Key() string {
	return fmt.Sprintf("%s - %s - %s - %s", record.Publisher, record.Domain, record.Country, record.Device)
}

type FactorReport struct {
	Time                 time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string    `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain               string    `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country              string    `boil:"country" json:"country" toml:"country" yaml:"country"`
	DeviceType           string    `boil:"device_type" json:"device_type" toml:"device_type" yaml:"device_type"`
	Revenue              float64   `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64   `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	DemandPartnerFee     float64   `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	SoldImpressions      int64     `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions int64     `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	DataFee              float64   `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
	Gp                   float64   `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	Gpp                  float64   `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
}
type Factor struct {
	Publisher string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain    string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Device    string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor    float64 `boil:"factor" json:"factor" toml:"factor" yaml:"factor"`
	Country   string  `boil:"country" json:"country" toml:"country" yaml:"country"`
}

func (rec *Factor) Key() string {
	return fmt.Sprint(rec.Publisher, rec.Domain, rec.Country, rec.Device)
}

func (rec *FactorReport) Key() string {
	return fmt.Sprint(rec.PublisherID, rec.Domain, rec.Country, rec.DeviceType)
}

func (rec *FactorReport) CalculateGP() {
	Gp := rec.Revenue - rec.Cost - rec.DemandPartnerFee - rec.DataFee
	rec.Gpp = roundFloat(Gp / rec.Revenue)
	rec.Gp = roundFloat(Gp)
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var err error

	err = quest.InitDB("quest" + conf.GetStringValueWithDefault("quest", "2"))
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	w.Cron, _ = conf.GetStringValue("cron")

	//w.MaxGPP, err = conf.GetFloat64ValueWithDefault("max_gpp", 0.65)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch MaxGPP config value (will use default 0.65)")
	//}
	//w.MinGPP, err = conf.GetFloat64ValueWithDefault("min_gpp", 0.3)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch MinGPP config value (will use default 0.3)")
	//}
	//w.TrendThreshold, err = conf.GetFloat64ValueWithDefault("trend_threshold", 0.1)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch TrendThreshold config value (will use default 0.1)")
	//}
	//w.FactorStep, err = conf.GetFloat64ValueWithDefault("factor_step", 1.1)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch FactorStep config value (will use default 1.1)")
	//}
	//w.LowFactor, err = conf.GetFloat64ValueWithDefault("low_factor", 0.75)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch FactorStep config value (will use default 1.1)")
	//}
	//w.MinFactor, err = conf.GetFloat64ValueWithDefault("min_factor", 0.01)
	//if err != nil {
	//	log.Warn().Err(err).Msg("failed to fetch FactorStep config value (will use default 1.1)")
	//}

	w.Domains, w.FilterExists = conf.GetStringSlice("domains", ",")
	if !w.FilterExists {
		log.Warn().Msg("Factors calculation is running on full system")
	}

	return nil

}

func (w *Worker) Do(ctx context.Context) error {

	var RecordsMap map[string]*FactorReport
	var factors map[string]*Factor

	var err error

	start, end := generateTimes(30)

	log.Info().Msg("fetch records from QuestDB")
	RecordsMap, err = FetchFromQuest(ctx, start, end)
	if err != nil {
		return err
	}

	log.Info().Msg("fetch records from Factors API")
	factors, err = FetchFactors()
	if err != nil {
		return err
	}

	var newFactors = make(map[string]FactorChanges)
	for _, record := range RecordsMap {
		// Check if the key exists on the first half as well
		if !w.CheckDomain(record.Domain) {
			continue
		}

		// Check if the key exists on the first half as well
		key := record.Key()

		_, exists := factors[key]
		if !exists {
			continue
		}

		oldFactor := factors[key].Factor // get current factor record
		var updatedFactor float64

		updatedFactor, err = w.CalculateFactor(record, oldFactor)
		if err != nil {
			log.Err(err).Msg("failed to calculate factor")
		}
		newFactors[key] = FactorChanges{
			Time:      end,
			EvalTime:  start,
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

	for _, r := range newFactors {

		if r.NewFactor != r.OldFactor {
			err := r.updateFactor()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("Error Updating factor for key: Publisher=%s, Domain=%s, Country=%s, Device=%s. ResponseStatus: %d", r.Publisher, r.Domain, r.Country, r.Device, r.RespStatus))
			}
		}
		logJSON, err := json.Marshal(r)
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error marshalling log for key:%v entry: %v", r.Key(), err))
		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

	}

	return nil
}

func (w *Worker) GetSleep() int {
	if w.Cron != "" {
		return bccron.Next(w.Cron)
	}
	return 0
}

func FetchFromQuest(ctx context.Context, start time.Time, stop time.Time) (map[string]*FactorReport, error) {
	var records []*FactorReport

	startString := start.Format("2006-01-02T15:04:05Z")
	stopString := stop.Format("2006-01-02T15:04:05Z")

	q := fmt.Sprintf(`SELECT to_date('%s','yyyy-MM-ddTHH:mm:ssZ') time,
       publisher publisher_id,
       domain,
       country,
       dtype device_type,
       sum(dbpr)/1000 revenue,
       sum(sbpr)/1000 cost,
       sum(dpfee)/1000 demand_partner_fee,
       count(1) sold_impressions,
       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions,
       sum(CASE WHEN uidsrc='iiq' THEN dbpr/1000 ELSE 0 END) * 0.045 data_fee
FROM impression
WHERE timestamp >= '%s'
  AND timestamp < '%s'
  AND publisher IS NOT NULL
  AND domain IS NOT NULL
  AND country IS NOT NULL
  AND dtype IS NOT NULL
GROUP BY 1, 2, 3, 4, 5`, startString, startString, stopString)
	log.Info().Str("q", q).Msg("processImpressionsCounters")
	err := queries.Raw(q).Bind(ctx, quest.DB(), &records)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query impressions from questdb")
	}

	var RecordsMap = make(map[string]*FactorReport)
	for _, r := range records {
		r.CalculateGP()
		RecordsMap[r.Key()] = r
	}

	return RecordsMap, nil
}

func (record *FactorChanges) updateFactor() error {
	requestBody := map[string]interface{}{
		"publisher": record.Publisher,
		"domain":    record.Domain,
		"country":   record.Country,
		"device":    record.Device,
		"factor":    record.NewFactor,
	}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrapf(err, "Error creating factors request body")
	}

	// Perform the HTTP request
	resp, err := http.Post("http://localhost:8000/factor", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Error updating factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	record.RespStatus = resp.StatusCode

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
	}
	return nil
}

func FetchFactors() (map[string]*Factor, error) {
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"pagination": map[string]interface{}{
			"page":      0,
			"page_size": 3000,
		}}

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating factors request body")
	}

	// Perform the HTTP request
	resp, err := http.Post("http://localhost:8000/factor/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching factors from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error Fetching factors from API. Request failed with status code: %d", resp.StatusCode))
	}

	var factorsResponse []*Factor
	if err := json.NewDecoder(resp.Body).Decode(&factorsResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing factors from API")
	}

	// Convert the response slice to a map
	factorsMap := make(map[string]*Factor)
	for _, item := range factorsResponse {
		factorsMap[item.Key()] = item
	}

	return factorsMap, nil
}

func generateTimes(minutes int) (time.Time, time.Time) {
	end := time.Now().UTC().Truncate(time.Duration(minutes) * time.Minute)
	start := end.Add(-time.Duration(minutes) * time.Minute)

	return start, end
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

func (w *Worker) CalculateFactor(record *FactorReport, oldFactor float64) (float64, error) {
	var updatedFactor float64

	if record.Gpp > 0.5 {
		updatedFactor = oldFactor * 1.3
	} else if record.Gpp > 0.45 {
		updatedFactor = oldFactor * 1.25
	} else if record.Gpp > 0.4 {
		updatedFactor = oldFactor * 1.2
	} else if record.Gpp > 0.33 {
		updatedFactor = oldFactor * 1.1
	} else if record.Gpp > 0.21 {
		updatedFactor = oldFactor // KEEP
	} else if record.Gpp < -0.1 {
		updatedFactor = oldFactor * 0.5
	} else if record.Gpp < 0 {
		updatedFactor = oldFactor * 0.7
	} else if record.Gpp < 0.1 {
		updatedFactor = oldFactor * 0.8
	} else if record.Gpp < 0.21 {
		updatedFactor = oldFactor * 0.875
	} else {
		return roundFloat(oldFactor), errors.New(fmt.Sprintf("unable to calculate factor: no matching condition Key: %s", record.Key()))
	}

	return roundFloat(updatedFactor), nil
}

func roundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

//func (w *Worker) CalculateFactor(T2record *FactorReport, T1record *FactorReport, oldFactor float64) (float64, error) {
//	var updatedFactor float64
//
//	if T2record == nil || T1record == nil {
//		return oldFactor, errors.New("T2record or T1record cannot be nil")
//	}
//
//	if T2record.Gpp > w.MaxGPP { //IF above Max GP (0.65) => Increase factor
//		updatedFactor = oldFactor * w.FactorStep
//	} else if T2record.Gpp < w.MinGPP && T2record.Gpp > 0 { //IF above Min GP (0.30) => Reduce factor
//		updatedFactor = oldFactor / w.FactorStep
//	} else if T2record.Gp < 0 && T1record.Gp > 0 { //IF T2 < 0 AND T1 > 0 => Set to Default low (0.5)
//		updatedFactor = w.LowFactor
//	} else if T2record.Gp > 0 && T1record.Gp < 0 { //IF T2 > 0 AND T1 < 0 => Keep Factor
//		updatedFactor = oldFactor
//	} else if T2record.Gp < 0 && T1record.Gp < 0 { //IF T2 < 0 AND T1 < 0 => Set to Default Min (0.01)
//		updatedFactor = w.MinFactor
//	} else if ((T2record.Gp - T1record.Gp) / T2record.Gp) >= w.TrendThreshold { //IF within thresholds & GP$ Positive trend => Increase factor
//		updatedFactor = oldFactor * w.FactorStep
//	} else if ((T2record.Gp - T1record.Gp) / T2record.Gp) <= (-1 * w.TrendThreshold) { // IF within thresholds & GP$ Negative trend => Decrease factor
//		updatedFactor = oldFactor / w.FactorStep
//	} else if (((T2record.Gp - T1record.Gp) / T2record.Gp) <= w.TrendThreshold) && ((T2record.Gp-T1record.Gp)/T2record.Gp) >= (-1*w.TrendThreshold) {
//		return oldFactor, nil
//	} else {
//		return oldFactor, errors.New(fmt.Sprintf("unable to calculate factor: no matching condition Key: %s", T1record.Key()))
//	}
//
//	return updatedFactor, nil
//}

//func FetchFromQuest(ctx context.Context, start time.Time, middle time.Time, stop time.Time) (map[string]*FactorReport, map[string]*FactorReport, error) {
//	var records []*FactorReport
//
//	startString := start.Format("2006-01-02T15:04:05Z")
//	middleString := middle.Format("2006-01-02T15:04:05Z")
//	stopString := stop.Format("2006-01-02T15:04:05Z")
//
//	q := fmt.Sprintf(`SELECT CASE WHEN timestamp < '%s' THEN to_date('%s','yyyy-MM-ddTHH:mm:ssZ')  ELSE to_date('%s','yyyy-MM-ddTHH:mm:ssZ') END time,
//       publisher publisher_id,
//       domain,
//       country,
//       dtype device_type,
//       sum(dbpr)/1000 revenue,
//       sum(sbpr)/1000 cost,
//       sum(dpfee)/1000 demand_partner_fee,
//       count(1) sold_impressions,
//       sum(CASE WHEN loop=false THEN 1 ELSE 0 END) publisher_impressions,
//       sum(CASE WHEN uidsrc='iiq' THEN dbpr/1000 ELSE 0 END) * 0.045 data_fee
//FROM impression
//WHERE timestamp >= '%s'
//  AND timestamp < '%s'
//  AND publisher IS NOT NULL
//  AND domain IS NOT NULL
//  AND country IS NOT NULL
//  AND dtype IS NOT NULL
//GROUP BY 1, 2, 3, 4, 5`, middleString, startString, middleString, startString, stopString)
//	log.Info().Str("q", q).Msg("processImpressionsCounters")
//	err := queries.Raw(q).Bind(ctx, quest.DB(), &records)
//	if err != nil {
//		return nil, nil, errors.Wrapf(err, "failed to query impressions from questdb")
//	}
//
//	var T1RecordsMap = make(map[string]*FactorReport)
//	var T2RecordsMap = make(map[string]*FactorReport)
//
//	for _, r := range records {
//		r.CalculateGP()
//		if r.Time.Equal(start) {
//			T1RecordsMap[r.Key()] = r
//		} else if r.Time.Equal(middle) {
//			T2RecordsMap[r.Key()] = r
//		}
//	}
//
//	return T1RecordsMap, T2RecordsMap, nil
//}

//func generateTimes(minutes int) (time.Time, time.Time, time.Time) {
//	end := time.Now().UTC().Truncate(time.Duration(minutes) * time.Minute)
//	middle := end.Add(-time.Duration(minutes) * time.Minute)
//	start := end.Add(-2 * time.Duration(minutes) * time.Minute)
//
//	return start, middle, end
//}
