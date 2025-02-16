package dto

import (
	"fmt"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

const (
	BidCachingControlPercentageMin = 0
	BidCachingControlPercentageMax = 1
)

type BidCachingUpdateRequest struct {
	RuleId            string   `json:"rule_id" validate:"required"`
	Publisher         string   `json:"publisher"`
	Domain            string   `json:"domain"`
	Device            string   `json:"device"`
	BidCaching        int16    `json:"bid_caching" validate:"bid_caching"`
	Country           string   `json:"country"`
	Browser           string   `json:"browser"`
	OS                string   `json:"os"`
	PlacementType     string   `json:"placement_type"`
	ControlPercentage *float64 `json:"control_percentage,omitempty" validate:"bccp"`
}

type BidCaching struct {
	RuleID            string     `json:"rule_id"`
	Publisher         string     `json:"publisher" validate:"required"`
	Domain            string     `json:"domain"`
	DemandPartnerID   string     `json:"demand_partner_id,omitempty"`
	Country           string     `json:"country" validate:"country"`
	Device            string     `json:"device" validate:"device"`
	BidCaching        int16      `json:"bid_caching" validate:"bid_caching"`
	Browser           string     `json:"browser" validate:"browser"`
	OS                string     `json:"os" validate:"os"`
	PlacementType     string     `json:"placement_type" validate:"placement_type"`
	Active            bool       `json:"active"`
	ControlPercentage *float64   `json:"control_percentage,omitempty" validate:"bccp"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

type BidCachingSlice []*BidCaching

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
	bc.RuleID = mod.RuleID
	bc.Publisher = mod.Publisher
	bc.Domain = mod.Domain.String
	bc.DemandPartnerID = mod.DemandPartnerID
	bc.BidCaching = mod.BidCaching
	bc.Active = mod.Active
	bc.CreatedAt = mod.CreatedAt
	bc.ControlPercentage = mod.ControlPercentage.Ptr()

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

	if mod.UpdatedAt.Valid {
		bc.UpdatedAt = mod.UpdatedAt.Ptr()
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
	if len(bc.RuleID) > 0 {
		return bc.RuleID
	} else {
		return bcguid.NewFrom(bc.GetFormula())
	}
}

func (bc *BidCaching) ToModel() *models.BidCaching {
	mod := models.BidCaching{
		RuleID:            bc.GetRuleID(),
		BidCaching:        bc.BidCaching,
		ControlPercentage: null.Float64FromPtr(bc.ControlPercentage),
		Publisher:         bc.Publisher,
		Active:            true,
	}

	if bc.Domain != "" {
		mod.Domain = null.StringFrom(bc.Domain)
	} else {
		mod.Domain = null.StringFrom("")
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
