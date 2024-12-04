package dto

type RefreshCacheUpdateRequest struct {
	RuleId        string `json:"rule_id"`
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	RefreshCache  int16  `json:"refresh_cache"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
}

type RefreshCacheUpdRequest struct {
	RuleId       string `json:"rule_id"`
	RefreshCache int16  `json:"refresh_cache"`
}

func (rc RefreshCacheUpdateRequest) GetPublisher() string     { return rc.Publisher }
func (rc RefreshCacheUpdateRequest) GetDomain() string        { return rc.Domain }
func (rc RefreshCacheUpdateRequest) GetDevice() string        { return rc.Device }
func (rc RefreshCacheUpdateRequest) GetCountry() string       { return rc.Country }
func (rc RefreshCacheUpdateRequest) GetBrowser() string       { return rc.Browser }
func (rc RefreshCacheUpdateRequest) GetOS() string            { return rc.OS }
func (rc RefreshCacheUpdateRequest) GetPlacementType() string { return rc.PlacementType }
