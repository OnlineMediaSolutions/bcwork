package dto

import (
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
)

type FloorUpdateRequest struct {
	RuleId        string  `json:"rule_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Floor         float64 `json:"floor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
	Active        bool    `json:"active"`
}
type FloorSlice []*Floor

type Floor struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	PublisherName string  `boil:"publisher_name" json:"publisher_name" toml:"publisher_name" yaml:"publisher_name"`
	Domain        string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Floor         float64 `boil:"floor" json:"floor" toml:"floor" yaml:"floor"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        bool    `boil:"active" json:"active" toml:"active" yaml:"active"`
}

func (f FloorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FloorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FloorUpdateRequest) GetDevice() string        { return f.Device }
func (f FloorUpdateRequest) GetCountry() string       { return f.Country }
func (f FloorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FloorUpdateRequest) GetOS() string            { return f.OS }
func (f FloorUpdateRequest) GetPlacementType() string { return f.PlacementType }

func (floor *Floor) FromModel(mod *models.Floor) error {
	floor.RuleId = mod.RuleID
	floor.Publisher = mod.Publisher
	floor.Domain = mod.Domain
	floor.Floor = mod.Floor
	floor.RuleId = mod.RuleID
	floor.Active = mod.Active

	if mod.R != nil && mod.R.FloorPublisher != nil {
		floor.PublisherName = mod.R.FloorPublisher.Name
	}

	if mod.Os.Valid {
		floor.OS = mod.Os.String
	}

	if mod.Country.Valid {
		floor.Country = mod.Country.String
	}

	if mod.Device.Valid {
		floor.Device = mod.Device.String
	}

	if mod.PlacementType.Valid {
		floor.PlacementType = mod.PlacementType.String
	}

	if mod.Browser.Valid {
		floor.Browser = mod.Browser.String
	}

	return nil
}

func (floor *Floor) GetRuleID() string {
	if len(floor.RuleId) > 0 {
		return floor.RuleId
	} else {
		return bcguid.NewFrom(floor.GetFormula())
	}
}

func (floor *Floor) GetFormula() string {
	p := floor.Publisher
	if p == "" {
		p = "*"
	}

	d := floor.Domain
	if d == "" {
		d = "*"
	}

	c := floor.Country
	if c == "" {
		c = "*"
	}

	os := floor.OS
	if os == "" {
		os = "*"
	}

	dt := floor.Device
	if dt == "" {
		dt = "*"
	}

	pt := floor.PlacementType
	if pt == "" {
		pt = "*"
	}

	b := floor.Browser
	if b == "" {
		b = "*"
	}

	return fmt.Sprintf("p=%s__d=%s__c=%s__os=%s__dt=%s__pt=%s__b=%s", p, d, c, os, dt, pt, b)

}

func (cs *FloorSlice) FromModel(slice models.FloorSlice) error {
	for _, mod := range slice {
		c := Floor{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (floor *Floor) ToModel() *models.Floor {
	mod := models.Floor{
		Floor:     floor.Floor,
		Publisher: floor.Publisher,
		Domain:    floor.Domain,
		RuleID:    floor.RuleId,
	}

	if floor.Country != "" {
		mod.Country = null.StringFrom(floor.Country)
	} else {
		mod.Country = null.String{}
	}

	if floor.OS != "" {
		mod.Os = null.StringFrom(floor.OS)
	} else {
		mod.Os = null.String{}
	}

	if floor.Device != "" {
		mod.Device = null.StringFrom(floor.Device)
	} else {
		mod.Device = null.String{}
	}

	if floor.PlacementType != "" {
		mod.PlacementType = null.StringFrom(floor.PlacementType)
	} else {
		mod.PlacementType = null.String{}
	}

	if floor.Browser != "" {
		mod.Browser = null.StringFrom(floor.Browser)
	} else {
		mod.Browser = null.String{}
	}

	return &mod
}
