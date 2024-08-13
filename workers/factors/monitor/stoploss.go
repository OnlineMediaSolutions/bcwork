package factors_monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

type Worker struct {
	Sleep         time.Duration          `json:"sleep"`
	DatabaseEnv   string                 `json:"dbenv"`
	Cron          string                 `json:"cron"`
	Domains       map[string]DomainSetup `json:"domains"`
	StopLoss      float64                `json:"stop_loss"`
	Quest         []string               `json:"quest_instances"`
	DefaultFactor float64                `json:"default_factor"`
}

// Changes applied on factors struct
type FactorChanges struct {
	Time       time.Time `json:"time"`
	EvalTime   time.Time `json:"eval_time"`
	Pubimps    int       `json:"pubimps"`
	Soldimps   int       `json:"soldimps"`
	Revenue    float64   `json:"revenue"`
	Cost       float64   `json:"cost"`
	GP         float64   `json:"gp"`
	GPP        float64   `json:"gpp"`
	Publisher  string    `json:"publisher"`
	Domain     string    `json:"domain"`
	Country    string    `json:"country"`
	Device     string    `json:"device"`
	OldFactor  float64   `json:"old_factor"`
	NewFactor  float64   `json:"new_factor"`
	RespStatus int       `json:"response_status"`
	Source     string    `json:"source"`
}

// Report from Quest struct
type FactorReport struct {
	Time                 time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string    `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain               string    `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country              string    `boil:"country" json:"country" toml:"country" yaml:"country"`
	DeviceType           string    `boil:"device_type" json:"device_type" toml:"device_type" yaml:"device_type"`
	Revenue              float64   `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64   `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	DemandPartnerFee     float64   `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	SoldImpressions      int       `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions int       `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	DataFee              float64   `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
	Gp                   float64   `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	Gpp                  float64   `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
}

// Factors API struct
type Factor struct {
	Publisher string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain    string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Device    string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor    float64 `boil:"factor" json:"factor" toml:"factor" yaml:"factor"`
	Country   string  `boil:"country" json:"country" toml:"country" yaml:"country"`
}

// Automation API Struct
type AutomationApi struct {
	PublisherId string    `json:"publisher_id"`
	Domain      string    `json:"domain"`
	Automation  bool      `json:"automation"`
	GppTarget   float64   `json:"gpp_target"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Automation API Struct
type DomainSetup struct {
	Domain    string  `json:"domain"`
	GppTarget float64 `json:"gpp_target"`
}

// Key functions for each struct
func (record *FactorChanges) Key() string {
	return fmt.Sprintf("%s - %s - %s - %s", record.Publisher, record.Domain, record.Country, record.Device)
}

func (rec *Factor) Key() string {
	return fmt.Sprint(rec.Publisher, rec.Domain, rec.Country, rec.Device)
}

func (rec *FactorReport) Key() string {
	return fmt.Sprint(rec.PublisherID, rec.Domain, rec.Country, rec.DeviceType)
}

func (rec *AutomationApi) Key() string {
	return fmt.Sprint(rec.PublisherId, rec.Domain)
}

func (rec *FactorReport) CalculateGP() {
	Gp := rec.Revenue - rec.Cost - rec.DemandPartnerFee - rec.DataFee
	rec.Gpp = roundFloat(Gp / rec.Revenue)
	rec.Gp = roundFloat(Gp)
}

// Worker functions
func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var err error
	var questExist bool

	w.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		w.Quest = []string{"amsquest2", "nycquest2"}
	}

	w.StopLoss, err = conf.GetFloat64ValueWithDefault("stoploss", -10)
	if err != nil {
		return errors.Wrapf(err, "failed to get stoploss value")
	}

	w.DefaultFactor, err = conf.GetFloat64ValueWithDefault("default_factor", 0.75)
	if err != nil {
		return errors.Wrapf(err, "failed to get stoploss value")
	}

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	w.Cron, _ = conf.GetStringValue("cron")

	return nil

}

func (w *Worker) Do(ctx context.Context) error {

	var RecordsMap map[string]*FactorReport
	var factors map[string]*Factor
	var err error

	w.Domains, err = FetchAutomationDomains()
	if err != nil {
		return errors.Wrapf(err, "failed to fetch automation domains")
	}

	start, end := generateTimes()

	log.Info().Msg("fetch records from QuestDB")
	RecordsMap, err = w.FetchFromQuest(ctx, start, end)
	if err != nil {
		return err
	}

	log.Info().Msg("fetch records from Factors API")
	factors, err = w.FetchFactors()
	if err != nil {
		return err
	}

	var newFactors = make(map[string]FactorChanges)
	for _, record := range RecordsMap {
		// Check if the key exists on factors
		key := record.Key()
		_, exists := factors[key]
		if !exists {
			continue
		}

		//Check if the record is below the stop loss
		if record.Gp > w.StopLoss {
			continue
		}

		newFactors[key] = FactorChanges{
			Time:      end,
			EvalTime:  start,
			Pubimps:   record.PublisherImpressions,
			Soldimps:  record.SoldImpressions,
			Cost:      roundFloat(record.Cost + record.DataFee + record.DemandPartnerFee),
			Revenue:   roundFloat(record.Revenue),
			GP:        record.Gp,
			GPP:       record.Gpp,
			Publisher: factors[key].Publisher,
			Domain:    factors[key].Domain,
			Country:   factors[key].Country,
			Device:    factors[key].Device,
			OldFactor: factors[key].Factor,
			NewFactor: w.DefaultFactor,
			Source:    "safety stop loss",
		}
	}

	for _, r := range newFactors {
		if r.NewFactor != r.OldFactor {
			err := r.updateFactor()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("Error Updating factor for key: Publisher=%s, Domain=%s, Country=%s, Device=%s. ResponseStatus: %d", r.Publisher, r.Domain, r.Country, r.Device, r.RespStatus))
			}
		}

		logJSON, err := json.Marshal(r) //Create log json to log it
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error marshalling log for key:%v entry: %v", r.Key(), err))
		}
		log.Info().Msg(fmt.Sprintf("%s", logJSON))

		mod, err := r.ToModel()
		if err != nil {
			log.Error().Err(err).Msg("failed to convert to model")
		}

		err = mod.Upsert(ctx, bcdb.DB(), true, Columns, boil.Infer(), boil.Infer())
		if err != nil {
			log.Error().Err(err).Msg("failed to push log to postgres")
		}
	}

	return nil
}

func (w *Worker) GetSleep() int {
	if w.Cron != "" {
		return bccron.Next(w.Cron)
	}
	return 0
}

// Fetch performance data from quest
func (w *Worker) FetchFromQuest(ctx context.Context, start time.Time, stop time.Time) (map[string]*FactorReport, error) {
	var records []*FactorReport

	startString := start.Format("2006-01-02T15:04:05Z")
	stopString := stop.Format("2006-01-02T15:04:05Z")
	domains := w.AutomationDomains()

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
  AND domain in('%s')
GROUP BY 1, 2, 3, 4, 5`, startString, startString, stopString, strings.Join(domains, "', '"))
	log.Info().Str("q", q).Msg("Fetch impressions data from quest")

	var RecordsMap = make(map[string]*FactorReport)
	for _, instance := range w.Quest {
		err := quest.InitDB(instance)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to initialize Quest instance: %s", instance))
		}

		err = queries.Raw(q).Bind(ctx, quest.DB(), &records)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query impressions from questdb")
		}

		// Check if the key exists on factors
		for _, record := range records {

			key := record.Key()
			item, exists := RecordsMap[key]
			if exists {
				mergedItem := &FactorReport{
					Time:                 record.Time,
					PublisherID:          record.PublisherID,
					Domain:               record.Domain,
					Country:              record.Country,
					DeviceType:           record.DeviceType,
					Revenue:              record.Revenue + item.Revenue,
					Cost:                 record.Cost + item.Cost,
					DemandPartnerFee:     record.DemandPartnerFee + item.DemandPartnerFee,
					SoldImpressions:      record.SoldImpressions + item.SoldImpressions,
					PublisherImpressions: record.PublisherImpressions + item.PublisherImpressions,
					DataFee:              record.DataFee + item.DataFee,
				}
				mergedItem.CalculateGP()
				RecordsMap[key] = mergedItem
			} else {
				record.CalculateGP()
				RecordsMap[key] = record
			}

		}

		records = nil
		err = quest.CloseDB()
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("Failed to close Quest instance: %s", instance))
		}

	}

	return RecordsMap, nil
}

// Fetch the factors from the Factors API
func FetchAutomationDomains() (map[string]DomainSetup, error) {

	// Create the request body using a map
	requestBody := map[string]interface{}{
		"filter": map[string][]string{
			"automation": {"true"},
		},
	}

	fmt.Println(requestBody)

	// Marshal the request body to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating automation setup request body")
	}

	// Perform the HTTP request
	resp, err := http.Post("http://localhost:8000/publisher/domain/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "Error Fetching automation setup from API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error Fetching automation setup from API. Request failed with status code: %d", resp.StatusCode))
	}

	var AutomationResponse []*AutomationApi
	if err := json.NewDecoder(resp.Body).Decode(&AutomationResponse); err != nil {
		return nil, errors.Wrapf(err, "Error parsing automation setup from API")
	}

	// Append the domains to the list of active domains & get Targets
	domainsMap := make(map[string]DomainSetup)
	for _, item := range AutomationResponse {
		gppTarget := item.GppTarget
		if gppTarget != 0 {
			gppTarget /= 100
		}

		domainsMap[item.Key()] = DomainSetup{
			Domain:    item.Domain,
			GppTarget: gppTarget,
		}

	}

	return domainsMap, nil
}

// Fetch the factors from the Factors API
func (w *Worker) FetchFactors() (map[string]*Factor, error) {
	// Create the request body using a map
	requestBody := map[string]interface{}{
		"filter": map[string][]string{"domain": w.AutomationDomains()},
		"pagination": map[string]interface{}{
			"page":      0,
			"page_size": 10000,
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

// Update a factor via the API
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
		Increase:       roundFloat((record.NewFactor / record.OldFactor) - 1),
		Source:         record.Source,
	}

	return model, nil
}

func roundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func generateTimes() (time.Time, time.Time) {
	start := time.Now().UTC().Truncate(time.Duration(30) * time.Minute)
	end := time.Now().UTC()

	return start, end
}

// Columns variable to update on the factors log table
var Columns = []string{
	models.PriceFactorLogColumns.Time,
	models.PriceFactorLogColumns.Publisher,
	models.PriceFactorLogColumns.Domain,
	models.PriceFactorLogColumns.Country,
	models.PriceFactorLogColumns.Device,
}

func (w *Worker) AutomationDomains() []string {
	var domains []string
	for _, item := range w.Domains {
		domains = append(domains, item.Domain)
	}
	return domains
}
