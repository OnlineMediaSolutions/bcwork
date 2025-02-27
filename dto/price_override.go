package dto

import "time"

type PriceOverrideRequest struct {
	Domain string `json:"domain"`
	Ips    []Ips  `json:"ips" validate:"duplicateIps,overridePriceKey"`
}

type Ips struct {
	IP    string    `json:"ip"`
	Date  time.Time `json:"date"`
	Price float64   `json:"price"`
}
