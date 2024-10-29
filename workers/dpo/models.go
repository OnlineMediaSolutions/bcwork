package dpo

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"time"
)

// Changes applied on factors struct
type DpoReport struct {
	Time       time.Time `json:"time"`
	EvalTime   time.Time `json:"eval_time"`
	Domain     string    `json:"domain"`
	Publisher  string    `json:"publisher"`
	Os         string    `json:"os"`
	Country    string    `json:"country"`
	DP         string    `json:"dp"`
	DpApiName  string    `json:"dp_api_name"`
	BidRequest int       `json:"bid_request"`
	Revenue    float64   `json:"revenue"`
	Erpm       float64   `json:"erpm"`
}

type PlacementReport struct {
	Domain    string  `json:"domain"`
	Publisher string  `json:"publisher"`
	Os        string  `json:"os"`
	Country   string  `json:"country"`
	Revenue   float64 `json:"revenue"`
}

type DpReport struct {
	DP         string  `json:"dp"`
	Revenue    float64 `json:"revenue"`
	BidRequest int     `json:"bid_request"`
}

type DpoChanges struct {
	Time       time.Time `json:"time"`
	EvalTime   time.Time `json:"eval_time"`
	Domain     string    `json:"domain"`
	Publisher  string    `json:"publisher"`
	Os         string    `json:"os"`
	Country    string    `json:"country"`
	DP         string    `json:"dp"`
	DpApiName  string    `json:"dp_api_name"`
	BidRequest int       `json:"bid_request"`
	Revenue    float64   `json:"revenue"`
	Erpm       float64   `json:"erpm"`
	OldFactor  float64   `json:"old_factor"`
	NewFactor  float64   `json:"new_factor"`
	RespStatus int       `json:"response_status"`
}

type DemandSetup struct {
	Name      string  `json:"name"`
	ApiName   string  `json:"api_name"`
	Threshold float64 `json:"demand_partner_id"`
}

type DpoApi struct {
	DP        string  `json:"demand_partner_id"`
	Domain    string  `json:"domain"`
	Publisher string  `json:"publisher"`
	Os        string  `json:"os"`
	Country   string  `json:"country"`
	Factor    float64 `json:"factor"`
}

type DpoData struct {
	DpoReport       map[string]*DpoReport       `json:"dpo_report"`
	PlacementReport map[string]*PlacementReport `json:"placement_report"`
	DpReport        map[string]*DpReport        `json:"dp_report"`
	DpoApi          map[string]*DpoApi          `json:"dpo_api"`
	Error           error                       `json:"error"`
}

func (record *DpoReport) PlacementKey() string {
	return fmt.Sprint(record.Domain, record.Publisher, record.Os, record.Country)
}

func (record *DpoReport) Key() string {
	return fmt.Sprint(record.DP, record.Domain, record.Publisher, record.Os, record.Country)
}

func (record *DpoReport) ApiKey() string {
	return fmt.Sprint(record.DpApiName, record.Domain, record.Publisher, record.Os, record.Country)
}

func (record *DpoChanges) Key() string {
	return fmt.Sprint(record.DP, record.Domain, record.Publisher, record.Os, record.Country)
}

func (record *DpoApi) Key() string {
	return fmt.Sprint(record.DP, record.Domain, record.Publisher, record.Os, record.Country)
}

func (record *DpoChanges) ToModel() (models.DpoAutomationLog, error) {
	model := models.DpoAutomationLog{
		Time:       record.Time,
		EvalTime:   record.EvalTime,
		DP:         record.DP,
		Publisher:  record.Publisher,
		Domain:     record.Domain,
		Country:    record.Country,
		Os:         record.Os,
		Revenue:    record.Revenue,
		BidRequest: record.BidRequest,
		Erpm:       record.Erpm,
		OldFactor:  record.OldFactor,
		NewFactor:  record.NewFactor,
		RespStatus: record.RespStatus,
	}

	return model, nil
}

func (worker *Worker) getDemandNames() []string {
	var names []string
	for _, Demand := range worker.Demands {
		names = append(names, Demand.Name)
	}
	return names
}

func (worker *Worker) CheckDemand(demand string) bool {
	_, exists := worker.Demands[demand]
	return exists
}
