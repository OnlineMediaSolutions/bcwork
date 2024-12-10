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
		want SearchResponse
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
			want: SearchResponse{
				FactorSectionType: {},
				DPOSectionType:    {},
				JSTargetingSectionType: {
					{
						PublisherID:   "1",
						PublisherName: "publisher_1",
						Domain:        "domain_1",
					},
				},
				DomainSectionType: {},
				FloorsSectionType: {},
				PublisherSectionType: {
					{
						PublisherID:   "1",
						PublisherName: "publisher_1",
					},
					{
						PublisherID:   "2",
						PublisherName: "publisher_2",
					},
				},
				DomainDashboardSectionType: {},
				PublisherDemandSectionType: {},
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
						// 001000
						SectionType:   null.StringFrom(DomainSectionType),
						PublisherID:   null.StringFrom("20952"),
						PublisherName: null.StringFrom("DevelopmentTest2"),
						Domain:        null.StringFrom("test.com"),
					},
				},
				query:       "Dev",
				sectionType: DomainSectionType,
			},
			want: SearchResponse{
				DomainSectionType: {
					{
						// 101000
						PublisherID:   "20712",
						PublisherName: "DevelopersTeamTestPublisher",
						Domain:        "dev-team.com",
					},
					{
						// 100100
						PublisherID:   "20890",
						PublisherName: "Test_Dev_Dates",
						Domain:        "dev-team.com",
					},
					{
						// 100000
						PublisherID:   "20798",
						PublisherName: "Test_Dates",
						Domain:        "dev-team.com",
					},
					{
						// 011000
						PublisherID:   "20797",
						PublisherName: "Dev_Test_Dates",
						Domain:        "team-dev.com",
					},
					{
						// 010100
						PublisherID:   "20798",
						PublisherName: "Test_DEV_Dates",
						Domain:        "team-dev.com",
					},
					{
						// 010000
						PublisherID:   "20799",
						PublisherName: "Test_Dates",
						Domain:        "team-dev.com",
					},
					{
						// 001000
						PublisherID:   "20952",
						PublisherName: "DevelopmentTest2",
						Domain:        "test.com",
					},
					{
						// 000100
						PublisherID:   "20891",
						PublisherName: "Test_Dev_Dates",
						Domain:        "test.com",
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
