package dto

import "time"

type PriceOverrideRequest struct {
	Domain string `json:"domain"`
	Ips    []Ips  `json:"ips"`
}

type Ips struct {
	IP    string    `json:"ip"`
	Date  time.Time `json:"-" swagger:"-"`
	Price float64   `json:"price"`
}
