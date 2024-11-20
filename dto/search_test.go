package dto

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestPrepareSearchResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mods models.SearchViewSlice
		want map[string][]SearchResult
	}{
		{
			name: "valid",
			mods: models.SearchViewSlice{
				{
					SectionType:   null.StringFrom(PublisherSectionType),
					PublisherID:   null.StringFrom("1"),
					PublisherName: null.StringFrom("publisher_1"),
				},
				{
					SectionType:   null.StringFrom(PublisherSectionType),
					PublisherID:   null.StringFrom("2"),
					PublisherName: null.StringFrom("publisher_2"),
				},
				{
					SectionType:   null.StringFrom(JSTargetingSectionType),
					PublisherID:   null.StringFrom("1"),
					PublisherName: null.StringFrom("publisher_1"),
					Domain:        null.StringFrom("domain_1"),
				},
			},
			want: map[string][]SearchResult{
				PublisherSectionType: {
					{
						PublisherID:   getStringPointer("1"),
						PublisherName: getStringPointer("publisher_1"),
					},
					{
						PublisherID:   getStringPointer("2"),
						PublisherName: getStringPointer("publisher_2"),
					},
				},
				JSTargetingSectionType: {
					{
						PublisherID:   getStringPointer("1"),
						PublisherName: getStringPointer("publisher_1"),
						Domain:        getStringPointer("domain_1"),
					},
				},
				DomainSectionType:          {},
				DomainDashboardSectionType: {},
				FactorSectionType:          {},
				FloorsSectionType:          {},
				PublisherDemandSectionType: {},
				DPOSectionType:             {},
				DemandPartnerSectionType:   {},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := PrepareSearchResults(tt.mods, "")
			assert.Equal(t, tt.want, got)
		})
	}
}
