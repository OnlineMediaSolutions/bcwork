package dto

type DPORuleUpdateRequest struct {
	RuleId        string  `json:"rule_id"`
	DemandPartner string  `json:"demand_partner_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain,omitempty"`
	Country       string  `json:"country,omitempty" validate:"country"`
	Browser       string  `json:"browser,omitempty" validate:"all"`
	OS            string  `json:"os,omitempty" validate:"all"`
	DeviceType    string  `json:"device_type,omitempty"`
	PlacementType string  `json:"placement_type,omitempty" validate:"all"`
	Factor        float64 `json:"factor" validate:"required,gte=0,factorDpo"`
}

type DPORuleDeleteRequest struct {
	DemandPartner string
	RuleId        string
}
