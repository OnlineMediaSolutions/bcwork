package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

type ConfiantUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate"`
}

type ConfiantUpdateRespose struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Confiant struct {
	ConfiantKey *string    `boil:"confiant_key" json:"confiant_key,omitempty" toml:"confiant_key" yaml:"confiant_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain      *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        *float64   `boil:"rate" json:"rate,omitempty" toml:"rate" yaml:"rate"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (confiant *Confiant) FromModel(mod *models.Confiant) error {
	confiant.PublisherID = mod.PublisherID
	confiant.CreatedAt = &mod.CreatedAt
	confiant.UpdatedAt = mod.UpdatedAt.Ptr()
	confiant.Domain = &mod.Domain
	confiant.Rate = &mod.Rate
	confiant.ConfiantKey = &mod.ConfiantKey

	return nil
}

func (confiant *Confiant) FromModelToCOnfiantWIthoutDomains(slice models.ConfiantSlice) error {
	for _, mod := range slice {
		if len(mod.Domain) == 0 {
			confiant.PublisherID = mod.PublisherID
			confiant.CreatedAt = &mod.CreatedAt
			confiant.UpdatedAt = mod.UpdatedAt.Ptr()
			confiant.Domain = &mod.Domain
			confiant.Rate = &mod.Rate
			confiant.ConfiantKey = &mod.ConfiantKey
			break
		}
	}

	return nil
}

type ConfiantSlice []*Confiant

func (cs *ConfiantSlice) FromModel(slice models.ConfiantSlice) error {
	for _, mod := range slice {
		c := Confiant{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}
