package dto

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type RefreshCacheUpdateRequest struct {
	RuleId        string `json:"rule_id"`
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	RefreshCache  int16  `json:"refresh_cache"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
}

type RefreshCache struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	RefreshCache  int16  `boil:"refresh_cache" json:"refresh_cache,omitempty" toml:"refresh_cache" yaml:"refresh_cache,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        string `boil:"actvie" json:"actvie" toml:"actvie" yaml:"actvie"`
}

type RefreshCacheUpdRequest struct {
	RuleId       string `json:"rule_id"`
	RefreshCache int16  `json:"refresh_cache"`
}

type RefreshCacheSlice []*RefreshCache

func (rc RefreshCacheUpdateRequest) GetPublisher() string     { return rc.Publisher }
func (rc RefreshCacheUpdateRequest) GetDomain() string        { return rc.Domain }
func (rc RefreshCacheUpdateRequest) GetDevice() string        { return rc.Device }
func (rc RefreshCacheUpdateRequest) GetCountry() string       { return rc.Country }
func (rc RefreshCacheUpdateRequest) GetBrowser() string       { return rc.Browser }
func (rc RefreshCacheUpdateRequest) GetOS() string            { return rc.OS }
func (rc RefreshCacheUpdateRequest) GetPlacementType() string { return rc.PlacementType }

func (rc *RefreshCache) FromModel(mod *models.RefreshCache) error {
	rc.RuleId = mod.RuleID
	rc.Publisher = mod.Publisher
	rc.RefreshCache = mod.RefreshCache
	rc.RuleId = mod.RuleID

	rc.Active = fmt.Sprintf("%t", mod.Active)

	if mod.Os.Valid {
		rc.OS = mod.Os.String
	}

	if mod.Domain.Valid {
		rc.Domain = mod.Domain.String
	}

	if mod.Country.Valid {
		rc.Country = mod.Country.String
	}

	if mod.Device.Valid {
		rc.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		rc.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		rc.Browser = mod.Browser.String
	}

	return nil
}
func (cs *RefreshCacheSlice) FromModel(slice models.RefreshCacheSlice) error {
	for _, mod := range slice {
		c := RefreshCache{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}
func (lr *RefreshCache) GetFormula() string {
	p := lr.Publisher
	if p == "" {
		p = "*"
	}

	d := lr.Domain
	if d == "" {
		d = "*"
	}

	c := lr.Country
	if c == "" {
		c = "*"
	}

	os := lr.OS
	if os == "" {
		os = "*"
	}

	dt := lr.Device
	if dt == "" {
		dt = "*"
	}

	pt := lr.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := lr.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (rc *RefreshCache) GetRuleID() string {
	if len(rc.RuleId) > 0 {
		return rc.RuleId
	} else {
		return bcguid.NewFrom(rc.GetFormula())
	}
}

func (rc *RefreshCache) ToModel() *models.RefreshCache {

	mod := models.RefreshCache{
		RuleID:       rc.GetRuleID(),
		RefreshCache: rc.RefreshCache,
		Publisher:    rc.Publisher,
	}

	if rc.Domain != "" {
		mod.Domain = null.StringFrom(rc.Domain)
	} else {
		mod.Domain = null.String{}
	}

	if rc.Country != "" {
		mod.Country = null.StringFrom(rc.Country)
	} else {
		mod.Country = null.String{}
	}

	if rc.OS != "" {
		mod.Os = null.StringFrom(rc.OS)
	} else {
		mod.Os = null.String{}
	}

	if rc.Device != "" {
		mod.Device = null.StringFrom(rc.Device)
	} else {
		mod.Device = null.String{}
	}

	if rc.PlacementType != "" {
		mod.PlacementType = null.StringFrom(rc.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if rc.Browser != "" {
		mod.Browser = null.StringFrom(rc.Browser)
	}

	return &mod

}
