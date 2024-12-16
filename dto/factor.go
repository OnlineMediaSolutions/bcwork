package dto

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type FactorUpdateRequest struct {
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Factor        float64 `json:"factor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
}

type Factor struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string  `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Factor        float64 `boil:"factor" json:"factor,omitempty" toml:"factor" yaml:"factor,omitempty"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        bool    `boil:"active" json:"active" toml:"active" yaml:"active"`
}

type FactorSlice []*Factor

func (factor *Factor) FromModel(mod *models.Factor) error {
	factor.RuleId = mod.RuleID
	factor.Publisher = mod.Publisher
	factor.Domain = mod.Domain
	factor.Factor = mod.Factor
	factor.RuleId = mod.RuleID
	factor.Active = mod.Active

	if mod.Os.Valid {
		factor.OS = mod.Os.String
	}

	if mod.Country.Valid {
		factor.Country = mod.Country.String
	}

	if mod.Device.Valid {
		factor.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		factor.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		factor.Browser = mod.Browser.String
	}

	return nil
}

func (cs *FactorSlice) FromModel(slice models.FactorSlice) error {
	for _, mod := range slice {
		c := Factor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (factor *Factor) GetFormula() string {
	p := factor.Publisher
	if p == "" {
		p = "*"
	}

	d := factor.Domain
	if d == "" {
		d = "*"
	}

	c := factor.Country
	if c == "" {
		c = "*"
	}

	os := factor.OS
	if os == "" {
		os = "*"
	}

	dt := factor.Device
	if dt == "" {
		dt = "*"
	}

	pt := factor.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := factor.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (factor *Factor) GetRuleID() string {
	if len(factor.RuleId) > 0 {
		return factor.RuleId
	} else {
		return bcguid.NewFrom(factor.GetFormula())
	}
}

func (factor *Factor) ToModel() *models.Factor {

	mod := models.Factor{
		RuleID:    factor.GetRuleID(),
		Factor:    factor.Factor,
		Publisher: factor.Publisher,
		Domain:    factor.Domain,
		Active:    factor.Active,
	}

	if factor.Country != "" {
		mod.Country = null.StringFrom(factor.Country)
	} else {
		mod.Country = null.String{}
	}

	if factor.OS != "" {
		mod.Os = null.StringFrom(factor.OS)
	} else {
		mod.Os = null.String{}
	}

	if factor.Device != "" {
		mod.Device = null.StringFrom(factor.Device)
	} else {
		mod.Device = null.String{}
	}

	if factor.PlacementType != "" {
		mod.PlacementType = null.StringFrom(factor.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if factor.Browser != "" {
		mod.Browser = null.StringFrom(factor.Browser)
	}

	return &mod

}

func (f FactorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FactorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FactorUpdateRequest) GetDevice() string        { return f.Device }
func (f FactorUpdateRequest) GetCountry() string       { return f.Country }
func (f FactorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FactorUpdateRequest) GetOS() string            { return f.OS }
func (f FactorUpdateRequest) GetPlacementType() string { return f.PlacementType }
