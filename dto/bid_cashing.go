package dto

type BidCashingUpdateRequest struct {
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	BidCashing    int64  `json:"bid_cashing"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
}

func (f BidCashingUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f BidCashingUpdateRequest) GetDomain() string        { return f.Domain }
func (f BidCashingUpdateRequest) GetDevice() string        { return f.Device }
func (f BidCashingUpdateRequest) GetCountry() string       { return f.Country }
func (f BidCashingUpdateRequest) GetBrowser() string       { return f.Browser }
func (f BidCashingUpdateRequest) GetOS() string            { return f.OS }
func (f BidCashingUpdateRequest) GetPlacementType() string { return f.PlacementType }
