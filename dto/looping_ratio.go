package dto

type LoopingRatioUpdateRequest struct {
	Publisher     string `json:"publisher"`
	Domain        string `json:"domain"`
	Device        string `json:"device"`
	LoopingRatio  int16  `json:"looping_ratio"`
	Country       string `json:"country"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	PlacementType string `json:"placement_type"`
}

func (lr LoopingRatioUpdateRequest) GetPublisher() string     { return lr.Publisher }
func (lr LoopingRatioUpdateRequest) GetDomain() string        { return lr.Domain }
func (lr LoopingRatioUpdateRequest) GetDevice() string        { return lr.Device }
func (lr LoopingRatioUpdateRequest) GetCountry() string       { return lr.Country }
func (lr LoopingRatioUpdateRequest) GetBrowser() string       { return lr.Browser }
func (lr LoopingRatioUpdateRequest) GetOS() string            { return lr.OS }
func (lr LoopingRatioUpdateRequest) GetPlacementType() string { return lr.PlacementType }
