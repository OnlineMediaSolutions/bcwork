package dto

const (
	AllSectionType             = "All"
	PublisherSectionType       = "Publisher list"
	DomainSectionType          = "Publisher / domain list"
	DomainDashboardSectionType = "Publisher / domain - Dashboard"
	FactorSectionType          = "Publisher / domain - Targeting - Bidder"
	JSTargetingSectionType     = "Publisher / domain - Targeting - JS"
	FloorsSectionType          = "Publisher /domain - Floors"
	PublisherDemandSectionType = "Publisher /domain - Demand"
	DPOSectionType             = "Demand / Publisher / Domain - DPO"
	DemandPartnerSectionType   = "Demand - Demand"
	// AdsTxtSectionType          = "Demand - Ads.txt dashboard"
)

type SearchResult struct {
	SectionType string `json:"section_type"`
	// TODO: expand
}
