package dto

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestPrepareSearchResults(t *testing.T) {
	t.Parallel()

	type args struct {
		mods        models.SearchViewSlice
		query       string
		sectionType string
	}

	tests := []struct {
		name string
		args args
		want map[string][]SearchResult
	}{
		{
			name: "valid",
			args: args{
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
				query:       "publisher",
				sectionType: "",
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
			},
		},
		{
			name: "sortingResults",
			args: args{
				mods: models.SearchViewSlice{
					{
						// 100100
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20890"),
						PublisherName: null.StringFrom("Test_Dev_Dates"),
						Domain:        null.StringFrom("dev-team.com"),
					},
					{
						// 010100
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20798"),
						PublisherName: null.StringFrom("Test_DEV_Dates"),
						Domain:        null.StringFrom("team-dev.com"),
					},
					{
						// 010000
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20799"),
						PublisherName: null.StringFrom("Test_Dates"),
						Domain:        null.StringFrom("team-dev.com"),
					},
					{
						// 100000
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20798"),
						PublisherName: null.StringFrom("Test_Dates"),
						Domain:        null.StringFrom("dev-team.com"),
					},
					{
						// 011000
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20797"),
						PublisherName: null.StringFrom("Dev_Test_Dates"),
						Domain:        null.StringFrom("team-dev.com"),
					},
					{
						// 101000
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20712"),
						PublisherName: null.StringFrom("DevelopersTeamTestPublisher"),
						Domain:        null.StringFrom("dev-team.com"),
					},
					{
						// 000100
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20891"),
						PublisherName: null.StringFrom("Test_Dev_Dates"),
						Domain:        null.StringFrom("test.com"),
					},
					{
						// 0010
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20952"),
						PublisherName: null.StringFrom("DevelopmentTest2"),
						Domain:        null.StringFrom("test.com"),
					},
				},
				query:       "Dev",
				sectionType: DomainSectionType,
			},
			want: map[string][]SearchResult{
				DomainSectionType: {
					{
						// 101000
						PublisherID:   getStringPointer("20712"),
						PublisherName: getStringPointer("DevelopersTeamTestPublisher"),
						Domain:        getStringPointer("dev-team.com"),
					},
					{
						// 100100
						PublisherID:   getStringPointer("20890"),
						PublisherName: getStringPointer("Test_Dev_Dates"),
						Domain:        getStringPointer("dev-team.com"),
					},
					{
						// 100000
						PublisherID:   getStringPointer("20798"),
						PublisherName: getStringPointer("Test_Dates"),
						Domain:        getStringPointer("dev-team.com"),
					},
					{
						// 011000
						PublisherID:   getStringPointer("20797"),
						PublisherName: getStringPointer("Dev_Test_Dates"),
						Domain:        getStringPointer("team-dev.com"),
					},
					{
						// 010100
						PublisherID:   getStringPointer("20798"),
						PublisherName: getStringPointer("Test_DEV_Dates"),
						Domain:        getStringPointer("team-dev.com"),
					},
					{
						// 010000
						PublisherID:   getStringPointer("20799"),
						PublisherName: getStringPointer("Test_Dates"),
						Domain:        getStringPointer("team-dev.com"),
					},
					{
						// 001000
						PublisherID:   getStringPointer("20952"),
						PublisherName: getStringPointer("DevelopmentTest2"),
						Domain:        getStringPointer("test.com"),
					},
					{
						// 000100
						PublisherID:   getStringPointer("20891"),
						PublisherName: getStringPointer("Test_Dev_Dates"),
						Domain:        getStringPointer("test.com"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := PrepareSearchResults(tt.args.mods, tt.args.query, tt.args.sectionType)
			assert.Equal(t, tt.want, got)
		})
	}
}
