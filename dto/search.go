package dto

import (
	"bytes"
	"encoding/json"
	"log"
	"sort"
	"strings"

	"github.com/m6yf/bcwork/models"
)

const (
	FactorSectionType          = "Bidder Targetings"
	DPOSectionType             = "DPO Rules"
	JSTargetingSectionType     = "JS Targetings"
	DomainSectionType          = "Domains list"
	FloorsSectionType          = "Floors list"
	PublisherSectionType       = "Publishers list"
	DomainDashboardSectionType = "Domain - Dashboard"
	PublisherDemandSectionType = "Domain - Demand"
)

var sections []string = []string{
	FactorSectionType,
	DPOSectionType,
	JSTargetingSectionType,
	DomainSectionType,
	FloorsSectionType,
	PublisherSectionType,
	DomainDashboardSectionType,
	PublisherDemandSectionType,
}

type SearchRequest struct {
	Query       string `json:"query"`
	SectionType string `json:"section_type"`
}

type SearchResult struct {
	PublisherID   string `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	Domain        string `json:"domain"`
}

type SearchResponse map[string][]SearchResult

func (s SearchResponse) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	_, err := buf.WriteString(`{`)
	if err != nil {
		return nil, err
	}

	for _, section := range sections {
		searchResult, ok := s[section]
		if !ok {
			continue
		}

		_, err := buf.WriteString(`"` + section + `":`)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(searchResult)
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(data)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteString(`,`)
		if err != nil {
			return nil, err
		}
	}

	buf.Truncate(len(buf.Bytes()) - 1)

	_, err = buf.WriteString(`}`)
	if err != nil {
		return nil, err
	}

	log.Println(buf.String())

	return buf.Bytes(), nil
}

func PrepareSearchResults(mods models.SearchViewSlice, query, reqSectionType string) SearchResponse {
	results := make(SearchResponse)

	if reqSectionType == "" {
		for _, section := range sections {
			results[section] = make([]SearchResult, 0, len(mods)/len(sections))
		}
	} else {
		results[reqSectionType] = make([]SearchResult, 0, len(mods))
	}

	for _, mod := range mods {
		results[mod.SectionType.String] = append(results[mod.SectionType.String], SearchResult{
			PublisherID:   mod.PublisherID.String,
			PublisherName: mod.PublisherName.String,
			Domain:        mod.Domain.String,
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
			publisherIDI   = strings.ToLower(result[i].PublisherID)
			publisherIDJ   = strings.ToLower(result[j].PublisherID)
			publisherNameI = strings.ToLower(result[i].PublisherName)
			publisherNameJ = strings.ToLower(result[j].PublisherName)
			domainI        = strings.ToLower(result[i].Domain)
			domainJ        = strings.ToLower(result[j].Domain)

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

func bool2int(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}

	return i
}
