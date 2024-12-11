package dto

import "time"

type GlobalFactor struct {
	Key         string     `boil:"key" json:"key" toml:"key" yaml:"key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Value       float64    `boil:"value" json:"value" toml:"value" yaml:"value"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}
