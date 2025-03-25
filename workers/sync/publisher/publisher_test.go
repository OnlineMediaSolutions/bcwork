package publisher

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Test_LoadedPublisher_ToModel(t *testing.T) {
	t.Parallel()

	type want struct {
		publisher *models.Publisher
		domains   models.PublisherDomainSlice
		blacklist boil.Columns
	}

	tests := []struct {
		name            string
		loadedPublisher *LoadedPublisher
		want            want
	}{
		{
			name: "maxBlacklistColumnLength",
			loadedPublisher: &LoadedPublisher{
				Id:         "1",
				Name:       "publisher",
				MediaBuyer: &field{Id: "media_buyer_id"},
				PausedDate: 1000,
				Site:       []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:    "1",
					Name:           "publisher",
					MediaBuyerID:   null.String{Valid: true, String: "media_buyer_id"},
					PauseTimestamp: null.Int64{Valid: true, Int64: 1000},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
						models.PublisherColumns.AccountManagerID,
						models.PublisherColumns.CampaignManagerID,
						models.PublisherColumns.OfficeLocation,
						models.PublisherColumns.ReactivateTimestamp,
						models.PublisherColumns.StartTimestamp,
					},
				},
			},
		},
		{
			name: "minBlacklistColumnLength",
			loadedPublisher: &LoadedPublisher{
				Id:              "1",
				Name:            "publisher",
				MediaBuyer:      &field{Id: "media_buyer_id"},
				StartDate:       500,
				PausedDate:      1000,
				ReactivatedDate: 2000,
				AccountManager:  &field{Id: "account_manager_id"},
				CampaignManager: &field{Id: "campaign_manager_id"},
				OfficeLocation:  "office",
				Site:            []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:         "1",
					Name:                "publisher",
					MediaBuyerID:        null.String{Valid: true, String: "media_buyer_id"},
					StartTimestamp:      null.Int64{Valid: true, Int64: 500},
					PauseTimestamp:      null.Int64{Valid: true, Int64: 1000},
					ReactivateTimestamp: null.Int64{Valid: true, Int64: 2000},
					AccountManagerID:    null.String{Valid: true, String: "account_manager_id"},
					CampaignManagerID:   null.String{Valid: true, String: "campaign_manager_id"},
					OfficeLocation:      null.String{Valid: true, String: "office"},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
					},
				},
			},
		},
		{
			name: "managerIDFromMap",
			loadedPublisher: &LoadedPublisher{
				Id:         "1",
				Name:       "publisher",
				MediaBuyer: &field{Id: "62de259de6e2871c098001e9"},
				PausedDate: 1000,
				Site:       []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:    "1",
					Name:           "publisher",
					MediaBuyerID:   null.String{Valid: true, String: "18"},
					PauseTimestamp: null.Int64{Valid: true, Int64: 1000},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
						models.PublisherColumns.AccountManagerID,
						models.PublisherColumns.CampaignManagerID,
						models.PublisherColumns.OfficeLocation,
						models.PublisherColumns.ReactivateTimestamp,
						models.PublisherColumns.StartTimestamp,
					},
				},
			},
		},
		{
			name: "mirroredPublisherIDFromDomainOptions",
			loadedPublisher: &LoadedPublisher{
				Id:         "1",
				Name:       "publisher",
				MediaBuyer: &field{Id: "62de259de6e2871c098001e9"},
				PausedDate: 1000,
				Site:       []string{"1.com", "2.com", "3.com"},
				DomainOptions: []*domainsOptions{
					{Domain: "1.com", IntegrationType: "Both", MirrorPublisher: ""},
					{Domain: "3.com", IntegrationType: "Compass", MirrorPublisher: ""},
					{Domain: "2.com", IntegrationType: "Both", MirrorPublisher: "4"},
				},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:    "1",
					Name:           "publisher",
					MediaBuyerID:   null.String{Valid: true, String: "18"},
					PauseTimestamp: null.Int64{Valid: true, Int64: 1000},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com", MirrorPublisherID: null.StringFrom("4")},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
						models.PublisherColumns.AccountManagerID,
						models.PublisherColumns.CampaignManagerID,
						models.PublisherColumns.OfficeLocation,
						models.PublisherColumns.ReactivateTimestamp,
						models.PublisherColumns.StartTimestamp,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			publisher, domains, blacklist := tt.loadedPublisher.ToModel(getMockManagersMap())
			assert.Equal(t, tt.want.publisher, publisher)
			assert.Equal(t, tt.want.domains, domains)
			assert.Equal(t, tt.want.blacklist, blacklist)
		})
	}
}

func Test_getManagerID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "idFromMap",
			args: args{
				id: "62de259de6e2871c098001e9",
			},
			want: "18",
		},
		{
			name: "initialId",
			args: args{
				id: "someunknownid",
			},
			want: "someunknownid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getManagerID(tt.args.id, getMockManagersMap())
			assert.Equal(t, tt.want, got)
		})
	}
}

func getMockManagersMap() map[string]string {
	return map[string]string{
		"62de259de6e2871c098001e9": "18",
	}
}
