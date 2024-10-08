package constant

type FactorUpdateRequest struct {
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Factor        float64 `json:"factor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
}

func (f FactorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FactorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FactorUpdateRequest) GetDevice() string        { return f.Device }
func (f FactorUpdateRequest) GetCountry() string       { return f.Country }
func (f FactorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FactorUpdateRequest) GetOS() string            { return f.OS }
func (f FactorUpdateRequest) GetPlacementType() string { return f.PlacementType }
