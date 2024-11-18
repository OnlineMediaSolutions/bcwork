package dto

type AdjustRequest struct {
	Domain []string `json:"domain"`
	Value  float64  `json:"Value"`
}
