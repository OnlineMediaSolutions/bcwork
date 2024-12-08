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
}

type BidCaching struct {
	RuleId        string `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain        string `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string `boil:"device" json:"device" toml:"device" yaml:"device"`
	BidCaching    int16  `boil:"bid_caching" json:"bid_caching,omitempty" toml:"bid_caching" yaml:"bid_caching,omitempty"`
	Browser       string `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
	Active        string `boil:"actvie" json:"actvie" toml:"actvie" yaml:"actvie"`
}

type BidCachingUpdRequest struct {
	RuleId     string `json:"rule_id"`
	BidCaching int16  `json:"bid_caching"`
}

type BidCachingRealtimeRecord struct {
	Rule       string `json:"rule"`
	BidCaching int16  `json:"bid_caching"`
	RuleID     string `json:"rule_id"`
}

func (f BidCachingUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f BidCachingUpdateRequest) GetDomain() string        { return f.Domain }
func (f BidCachingUpdateRequest) GetDevice() string        { return f.Device }
func (f BidCachingUpdateRequest) GetCountry() string       { return f.Country }
func (f BidCachingUpdateRequest) GetBrowser() string       { return f.Browser }
func (f BidCachingUpdateRequest) GetOS() string            { return f.OS }
func (f BidCachingUpdateRequest) GetPlacementType() string { return f.PlacementType }
