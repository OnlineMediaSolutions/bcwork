package dto

import "github.com/m6yf/bcwork/models"

const (
	PublisherSectionType       = "Publisher list"
	DomainSectionType          = "Publisher / domain list"
	DomainDashboardSectionType = "Publisher / domain - Dashboard"
	FactorSectionType          = "Targeting - Bidder"
	JSTargetingSectionType     = "Targeting - JS"
	FloorsSectionType          = "Floors"
	PublisherDemandSectionType = "Publisher / domain - Demand"
	DPOSectionType             = "DPO Rule"
	DemandPartnerSectionType   = "Demand - Demand"
)

var sections []string = []string{
	PublisherSectionType,
	DomainSectionType,
	DomainDashboardSectionType,
	FactorSectionType,
	JSTargetingSectionType,
	FloorsSectionType,
	PublisherDemandSectionType,
	DPOSectionType,
	DemandPartnerSectionType,
}

type SearchRequest struct {
	Query       string `json:"query"`
	SectionType string `json:"section_type"`
}

type SearchResult struct {
	PublisherID       *string `json:"publisher_id"`
	PublisherName     *string `json:"publisher_name"`
	Domain            *string `json:"domain"`
	DemandPartnerName *string `json:"demand_partner_name"`
}

func PrepareSearchResults(mods models.SearchViewSlice, reqSectionType string) map[string][]SearchResult {
	results := make(map[string][]SearchResult)

	if reqSectionType == "" {
		for _, section := range sections {
			results[section] = make([]SearchResult, 0, len(mods)/len(sections))
		}
	} else {
		results[reqSectionType] = make([]SearchResult, 0, len(mods))
	}

	for _, mod := range mods {
		results[mod.SectionType.String] = append(results[mod.SectionType.String], SearchResult{
			PublisherID:       getStringPointer(mod.PublisherID.String),
			PublisherName:     getStringPointer(mod.PublisherName.String),
			Domain:            getStringPointer(mod.Domain.String),
			DemandPartnerName: getStringPointer(mod.DemandPartnerName.String),
		})
	}

	return results
}

func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
