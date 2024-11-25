package dto

type AdjustRequest struct {
	Domain []string `json:"domain" validate:"required"`
	Value  float64  `json:"value" validate:"required,floor"`
}
