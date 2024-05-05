package core

import (
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"time"
)

type Confiant struct {
	ConfiantKey string     `boil:"confiant_key" json:"confiant_key" toml:"confiant_key" yaml:"confiant_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain      string     `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        float64    `boil:"rate" json:"rate" toml:"rate" yaml:"rate"`
	CreatedAt   time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

type ConfiantSlice []*Confiant

func (confiant *Confiant) FromModel(mod *models.Confiant) error {

	confiant.PublisherID = mod.PublisherID
	confiant.CreatedAt = mod.CreatedAt
	confiant.Domain = mod.Domain.String
	confiant.Rate = mod.Rate
	confiant.ConfiantKey = mod.ConfiantKey

	return nil
}

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
