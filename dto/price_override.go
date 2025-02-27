package dto

type PriceOverrideRequest struct {
	Domain string `json:"domain"`
	Ips    []Ips  `json:"ips"`
}

type Ips struct {
	IP    string  `json:"ip"`
	Price float64 `json:"price"`
}
