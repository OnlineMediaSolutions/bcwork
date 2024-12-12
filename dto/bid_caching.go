package dto

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type BidCachingUpdateRequest struct {
	RuleId        string `json:"rule_id"`
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	BidCaching    int16  `json:"bid_caching"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
}

type BidCaching struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	BidCaching    int16  `boil:"bid_caching" json:"bid_caching,omitempty" toml:"bid_caching" yaml:"bid_caching,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        string `boil:"active" json:"active" toml:"active" yaml:"active"`
}

type BidCachingSlice []*BidCaching

type BidCachingUpdRequest struct {
	RuleId     string `json:"rule_id"`
	BidCaching int16  `json:"bid_caching"`
}

type BidCachingRealtimeRecord struct {
	Rule       string `json:"rule"`
	BidCaching int16  `json:"bid_caching"`
	RuleID     string `json:"rule_id"`
}

func (f BidCachingUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f BidCachingUpdateRequest) GetDomain() string        { return f.Domain }
func (f BidCachingUpdateRequest) GetDevice() string        { return f.Device }
func (f BidCachingUpdateRequest) GetCountry() string       { return f.Country }
func (f BidCachingUpdateRequest) GetBrowser() string       { return f.Browser }
func (f BidCachingUpdateRequest) GetOS() string            { return f.OS }
func (f BidCachingUpdateRequest) GetPlacementType() string { return f.PlacementType }

func (cs *BidCachingSlice) FromModel(slice models.BidCachingSlice) error {
	for _, mod := range slice {
		c := BidCaching{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (bc *BidCaching) FromModel(mod *models.BidCaching) error {
	bc.RuleId = mod.RuleID
	bc.Publisher = mod.Publisher
	bc.Domain = mod.Domain
	bc.BidCaching = mod.BidCaching
	bc.Active = fmt.Sprintf("%t", mod.Active)

	if mod.Os.Valid {
		bc.OS = mod.Os.String
	}

	if mod.Country.Valid {
		bc.Country = mod.Country.String
	}

	if mod.Device.Valid {
		bc.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		bc.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		bc.Browser = mod.Browser.String
	}

	return nil
}

func (bc *BidCaching) GetFormula() string {
	p := bc.Publisher
	if p == "" {
		p = "*"
	}

	d := bc.Domain
	if d == "" {
		d = "*"
	}

	c := bc.Country
	if c == "" {
		c = "*"
	}

	os := bc.OS
	if os == "" {
		os = "*"
	}

	dt := bc.Device
	if dt == "" {
		dt = "*"
	}

	pt := bc.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := bc.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (bc *BidCaching) GetRuleID() string {
	if len(bc.RuleId) > 0 {
		return bc.RuleId
	} else {
		return bcguid.NewFrom(bc.GetFormula())
	}
}

func (bc *BidCaching) ToModel() *models.BidCaching {

	mod := models.BidCaching{
		RuleID:     bc.GetRuleID(),
		BidCaching: bc.BidCaching,
		Publisher:  bc.Publisher,
		Domain:     bc.Domain,
	}

	if bc.Country != "" {
		mod.Country = null.StringFrom(bc.Country)
	} else {
		mod.Country = null.String{}
	}

	if bc.OS != "" {
		mod.Os = null.StringFrom(bc.OS)
	} else {
		mod.Os = null.String{}
	}

	if bc.Device != "" {
		mod.Device = null.StringFrom(bc.Device)
	} else {
		mod.Device = null.String{}
	}

	if bc.PlacementType != "" {
		mod.PlacementType = null.StringFrom(bc.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if bc.Browser != "" {
		mod.Browser = null.StringFrom(bc.Browser)
	}

	return &mod

}
