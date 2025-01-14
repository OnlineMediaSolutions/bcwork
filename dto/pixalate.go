package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

type PixalateUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Rate      float64 `json:"rate"`
	Active    bool    `json:"active"`
}

type PixalateUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pixalate struct {
	PixalateKey string     `boil:"pixalate_key" json:"pixalate_key,omitempty" toml:"pixalate_key" yaml:"pixalate_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain      *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        *float64   `boil:"rate" json:"rate,omitempty" toml:"rate" yaml:"rate"`
	Active      *bool      `boil:"active" json:"active,omitempty" toml:"active" yaml:"active"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (pixalate *Pixalate) FromModel(mod *models.Pixalate) error {
	pixalate.PublisherID = mod.PublisherID
	pixalate.CreatedAt = &mod.CreatedAt
	pixalate.UpdatedAt = mod.UpdatedAt.Ptr()
	pixalate.Domain = &mod.Domain
	pixalate.Rate = &mod.Rate
	pixalate.PixalateKey = mod.ID
	pixalate.Active = &mod.Active

	return nil
}

func (pixalate *Pixalate) FromModelToPixalateWIthoutDomains(slice models.PixalateSlice) error {
	for _, mod := range slice {
		if len(mod.Domain) == 0 {
			pixalate.PublisherID = mod.PublisherID
			pixalate.CreatedAt = &mod.CreatedAt
			pixalate.UpdatedAt = mod.UpdatedAt.Ptr()
			pixalate.Domain = &mod.Domain
			pixalate.Rate = &mod.Rate
			pixalate.PixalateKey = mod.ID
			pixalate.Active = &mod.Active
			break
		}
	}

	return nil
}

func (newPixalate *Pixalate) createPixalate(pixalate models.Pixalate) {
	newPixalate.PublisherID = pixalate.PublisherID
	newPixalate.CreatedAt = &pixalate.CreatedAt
	newPixalate.UpdatedAt = pixalate.UpdatedAt.Ptr()
	newPixalate.Domain = &pixalate.Domain
	newPixalate.Rate = &pixalate.Rate
	newPixalate.PixalateKey = pixalate.ID
	newPixalate.Active = &pixalate.Active
}

type PixalateSlice []*Pixalate

func (cs *PixalateSlice) FromModel(slice models.PixalateSlice) error {
	for _, mod := range slice {
		c := Pixalate{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}
