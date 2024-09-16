package constant

type FloorUpdateRequest struct {
	RuleId        string  `json:"rule_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Floor         float64 `json:"floor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
}

func (f FloorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FloorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FloorUpdateRequest) GetDevice() string        { return f.Device }
func (f FloorUpdateRequest) GetCountry() string       { return f.Country }
func (f FloorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FloorUpdateRequest) GetOS() string            { return f.OS }
func (f FloorUpdateRequest) GetPlacementType() string { return f.PlacementType }
