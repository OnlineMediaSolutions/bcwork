package factors_autmation

import (
	"fmt"
	"time"
)

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
	RuleId     string    `json:"rule_id"`
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
	BidRequests          float64   `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	DataFee              float64   `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
	TechFee              float64   `boil:"tech_fee" json:"tech_fee" toml:"tech_fee" yaml:"tech_fee"`
	TamFee               float64   `boil:"tam_fee" json:"tam_fee" toml:"tam_fee" yaml:"tam_fee"`
	ConsultantFee        float64   `boil:"consultant_fee" json:"consultant_fee" toml:"consultant_fee" yaml:"consultant_fee"`
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
	RuleId    string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
}

// Automation API Struct
type DomainSetup struct {
	Domain    string  `json:"domain"`
	GppTarget float64 `json:"gpp_target"`
}

// Domain API Struct
type AutomationApi struct {
	PublisherId string    `json:"publisher_id"`
	Domain      string    `json:"domain"`
	Automation  bool      `json:"automation"`
	GppTarget   float64   `json:"gpp_target"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

func (rec *FactorReport) SetupKey() string {
	return fmt.Sprint(rec.PublisherID, rec.Domain)
}

func (rec *AutomationApi) Key() string {
	return fmt.Sprint(rec.PublisherId, rec.Domain)
}

func (rec *FactorReport) CalculateGP(fees map[string]float64, consultantFees map[string]float64) {
	rec.TamFee = RoundFloat(fees["tam_fee"] * rec.Cost)
	rec.TechFee = RoundFloat(fees["tech_fee"] * rec.BidRequests / 1000000)
	rec.ConsultantFee = 0.0
	value, exists := consultantFees[rec.PublisherID]
	if exists {
		rec.ConsultantFee = rec.Cost * value
	}

	rec.Gp = RoundFloat(rec.Revenue - rec.Cost - rec.DemandPartnerFee - rec.DataFee - rec.TamFee - rec.TechFee - rec.ConsultantFee)
	rec.Gpp = 0
	if rec.Revenue != 0 {
		rec.Gpp = RoundFloat(rec.Gp / rec.Revenue)
	}
}

func (record *FactorChanges) UpdateResponseStatus(status int) {
	record.RespStatus = status
}
