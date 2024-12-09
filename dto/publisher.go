package dto

type Publisher struct {
	PublisherId string `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Name        string `boil:"name" json:"name" toml:"name" yaml:"name"`
}
