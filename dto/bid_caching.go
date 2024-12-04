package dto

type BidCachingUpdateRequest struct {
	RuleId        string `json:"rule_id"`
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	BidCaching    int16  `json:"bid_caching"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
	Active        bool   `json:"active"`
}

func (f BidCachingUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f BidCachingUpdateRequest) GetDomain() string        { return f.Domain }
func (f BidCachingUpdateRequest) GetDevice() string        { return f.Device }
func (f BidCachingUpdateRequest) GetCountry() string       { return f.Country }
func (f BidCachingUpdateRequest) GetBrowser() string       { return f.Browser }
func (f BidCachingUpdateRequest) GetOS() string            { return f.OS }
func (f BidCachingUpdateRequest) GetPlacementType() string { return f.PlacementType }
