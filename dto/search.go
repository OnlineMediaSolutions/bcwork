package dto

import (
	"sort"
	"strings"

	"github.com/m6yf/bcwork/models"
)

const (
	PublisherSectionType       = "Publisher list"
	DomainSectionType          = "Publisher / domain list"
	DomainDashboardSectionType = "Publisher / domain - Dashboard"
	FactorSectionType          = "Targeting - Bidder"
	JSTargetingSectionType     = "Targeting - JS"
	FloorsSectionType          = "Floors"
	PublisherDemandSectionType = "Publisher / domain - Demand"
	DPOSectionType             = "DPO Rule"
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
}

type SearchRequest struct {
	Query       string `json:"query"`
	SectionType string `json:"section_type"`
}

type SearchResult struct {
	PublisherID   *string `json:"publisher_id"`
	PublisherName *string `json:"publisher_name"`
	Domain        *string `json:"domain"`
}

func PrepareSearchResults(mods models.SearchViewSlice, query, reqSectionType string) map[string][]SearchResult {
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
			PublisherID:   getStringPointer(mod.PublisherID.String),
			PublisherName: getStringPointer(mod.PublisherName.String),
			Domain:        getStringPointer(mod.Domain.String),
		})
	}

	for i := range results {
		results[i] = sortSearchResults(results[i], query)
	}

	return results
}

func sortSearchResults(result []SearchResult, query string) []SearchResult {
	sort.SliceStable(result, func(i, j int) bool {
		var (
			publisherIDI   = strings.ToLower(getStringFromPointer(result[i].PublisherID))
			publisherIDJ   = strings.ToLower(getStringFromPointer(result[j].PublisherID))
			publisherNameI = strings.ToLower(getStringFromPointer(result[i].PublisherName))
			publisherNameJ = strings.ToLower(getStringFromPointer(result[j].PublisherName))
			domainI        = strings.ToLower(getStringFromPointer(result[i].Domain))
			domainJ        = strings.ToLower(getStringFromPointer(result[j].Domain))

			query = strings.ToLower(query)
		)

		isPublisherIDHasPrefixI := strings.HasPrefix(publisherIDI, query)
		isPublisherIDHasPrefixJ := strings.HasPrefix(publisherIDJ, query)
		isPublisherNameHasPrefixI := strings.HasPrefix(publisherNameI, query)
		isPublisherNameHasPrefixJ := strings.HasPrefix(publisherNameJ, query)
		isDomainHasPrefixI := strings.HasPrefix(domainI, query)
		isDomainHasPrefixJ := strings.HasPrefix(domainJ, query)

		isPublisherIDContainsI := strings.Contains(publisherIDI, query)
		isPublisherIDContainsJ := strings.Contains(publisherIDJ, query)
		isPublisherNameContainsI := strings.Contains(publisherNameI, query)
		isPublisherNameContainsJ := strings.Contains(publisherNameJ, query)
		isDomainContainsI := strings.Contains(domainI, query)
		isDomainContainsJ := strings.Contains(domainJ, query)

		maskI := getSearchBitmask(
			isPublisherIDContainsI, isPublisherNameContainsI, isDomainContainsI,
			isPublisherIDHasPrefixI, isPublisherNameHasPrefixI, isDomainHasPrefixI,
		)
		maskJ := getSearchBitmask(
			isPublisherIDContainsJ, isPublisherNameContainsJ, isDomainContainsJ,
			isPublisherIDHasPrefixJ, isPublisherNameHasPrefixJ, isDomainHasPrefixJ,
		)

		return maskI > maskJ
	})
	return result
}

func getSearchBitmask(
	isPublisherIDContains, isPublisherNameContains, isDomainContains bool,
	isPublisherIDHasPrefix, isPublisherNameHasPrefix, isDomainHasPrefix bool,
) int {
	const (
		containsPublisherID   = 1 << iota // 1
		prefixPublisherID                 // 2
		containsPublisherName             // 4
		prefixPublisherName               // 8
		containsDomain                    // 16
		prefixDomain                      // 32
	)

	return bool2int(isPublisherIDContains)*containsPublisherID +
		bool2int(isPublisherIDHasPrefix)*prefixPublisherID +
		bool2int(isPublisherNameContains)*containsPublisherName +
		bool2int(isPublisherNameHasPrefix)*prefixPublisherName +
		bool2int(isDomainContains)*containsDomain +
		bool2int(isDomainHasPrefix)*prefixDomain
}

func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func getStringFromPointer(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func bool2int(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
