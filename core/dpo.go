package core

import (
	"github.com/m6yf/bcwork/models"
	"time"
)

type Dpo struct {
	DemandPartnerID string                             `json:"demand_partner_id"`
	IsInclude       bool                               `json:"is_include"`
	CreatedAt       time.Time                          `json:"created_at"`
	UpdatedAt       *time.Time                         `json:"updated_at"`
	Rules           DemandPartnerOptimizationRuleSlice `json:"rules"`
}

type DpoSlice []*Dpo

func (dpo *Dpo) FromModel(mod *models.Dpo) {

	dpo.DemandPartnerID = mod.DemandPartnerID
	dpo.IsInclude = mod.IsInclude
	dpo.CreatedAt = mod.CreatedAt
	dpo.UpdatedAt = mod.UpdatedAt.Ptr()

	if mod.R != nil {
		if mod.R.DemandPartnerDpoRules != nil {
			dpo.Rules = make(DemandPartnerOptimizationRuleSlice, 0, 0)
			dpo.Rules.FromModel(mod.R.DemandPartnerDpoRules)
		}
	}
}

func (dpos *DpoSlice) FromModel(slice models.DpoSlice) {

	for _, mod := range slice {
		dpo := Dpo{}
		dpo.FromModel(mod)
		*dpos = append(*dpos, &dpo)
	}

}
