package dto

import (
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/types"
	"time"
)

type DpApiRequest struct {
	DateStamp     string `json:"date_stamp"`
	DemandPartner string `json:"demand_partner"`
	Domain        string `json:"domain"`
	SoldImps      string `json:"sold_imps"`
	Revenue       string `json:"revenue"`
	UpdatedAt     string `json:"updated_at"`
}

type DpiApi struct {
	DateStamp     string        `json:"date_stamp"`
	Date          string        `json:"date"`
	DemandPartner string        `json:"demand_partner"`
	LastFullDay   string        `json:"last_full_day"`
	LastUpdated   string        `json:"last_updated"`
	Counter       int64         `json:"counter"`
	Domain        string        `json:"domain"`
	SoldImps      int64         `json:"sold_imps"`
	Revenue       types.Decimal `json:"revenue"`
	UpdatedAt     string        `json:"updated_at"`
}

type DpApiSlice []*DpiApi

func (dpApi *DpiApi) FromModel(mod *models.DPAPIReport) error {
	dpApi.DemandPartner = mod.DemandPartner
	dpApi.Domain = mod.Domain
	dpApi.SoldImps = mod.SoldImps
	dpApi.Revenue = mod.Revenue
	dpApi.DateStamp = mod.DateStamp.Format(time.DateTime)

	return nil
}

func (dpApi *DpApiSlice) FromModel(slice models.DPAPIReportSlice) error {
	for _, mod := range slice {
		c := DpiApi{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*dpApi = append(*dpApi, &c)
	}

	return nil
}

func (dpApi *DpiApi) ToModel() *models.DPAPIReport {
	mod := models.DPAPIReport{
		DemandPartner: dpApi.DemandPartner,
		Domain:        dpApi.Domain,
		SoldImps:      dpApi.SoldImps,
		Revenue:       dpApi.Revenue,
	}

	return &mod
}
