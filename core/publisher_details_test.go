package core

import (
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"testing"
)

func Test_PublisherDetail_FromModel(t *testing.T) {
	t.Parallel()

	type args struct {
		mod *dto.PublisherDetailModel
	}

	tests := []struct {
		name    string
		args    args
		want    dto.PublisherDetail
		wantErr bool
	}{
		{
			name: "valid_withoutFactorsAndActivityStatusEmpty",
			args: args{
				mod: &dto.PublisherDetailModel{
					Publisher: models.Publisher{
						Name:             "publisher",
						PublisherID:      "id",
						AccountManagerID: null.String{Valid: true, String: "am_id"},
					},
					PublisherDomain: models.PublisherDomain{
						Domain:     "domain",
						Automation: true,
						GPPTarget:  null.Float64{Valid: true, Float64: 0.1},
					},
				},
			},
			want: dto.PublisherDetail{
				Name:             "publisher",
				PublisherID:      "id",
				AccountManagerID: "am_id",
				Domain:           "domain",
				Automation:       true,
				GPPTarget:        0.1,
				ActivityStatus:   "Paused",
				BidCaching:       make([]dto.BidCaching, 0),
				RefreshCache:     make([]dto.RefreshCache, 0),
			},
		},
		{
			name: "valid_withoutFactorsAndActivityStatusLow",
			args: args{
				mod: &dto.PublisherDetailModel{
					Publisher: models.Publisher{
						Name:             "publisher",
						PublisherID:      "123",
						AccountManagerID: null.String{Valid: true, String: "am_id"},
					},
					PublisherDomain: models.PublisherDomain{
						Domain:     "domain.com",
						Automation: true,
						GPPTarget:  null.Float64{Valid: true, Float64: 0.1},
					},
				},
			},
			want: dto.PublisherDetail{
				Name:             "publisher",
				PublisherID:      "123",
				AccountManagerID: "am_id",
				Domain:           "domain.com",
				Automation:       true,
				GPPTarget:        0.1,
				ActivityStatus:   "Active",
				BidCaching:       make([]dto.BidCaching, 0),
				RefreshCache:     make([]dto.RefreshCache, 0),
			},
		},
		{
			name: "valid_withUserFullname",
			args: args{
				mod: &dto.PublisherDetailModel{
					Publisher: models.Publisher{
						Name:             "publisher",
						PublisherID:      "456",
						AccountManagerID: null.String{Valid: true, String: "1"},
					},
					PublisherDomain: models.PublisherDomain{
						Domain:     "domain.com",
						Automation: true,
						GPPTarget:  null.Float64{Valid: true, Float64: 0.1},
					},
					User: dto.UserModelCompact{
						ID:        1,
						FirstName: null.StringFrom("first"),
						LastName:  null.StringFrom("last"),
					},
				},
			},
			want: dto.PublisherDetail{
				Name:                   "publisher",
				PublisherID:            "456",
				AccountManagerID:       "1",
				AccountManagerFullName: "first last",
				Domain:                 "domain.com",
				Automation:             true,
				GPPTarget:              0.1,
				ActivityStatus:         "Low",
				BidCaching:             make([]dto.BidCaching, 0),
				RefreshCache:           make([]dto.RefreshCache, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pd := dto.PublisherDetail{}

			activeStatus := make(map[string]map[string]dto.ActivityStatus)
			activeStatus["domain.com"] = make(map[string]dto.ActivityStatus)
			activeStatus["domain.com"]["123"] = dto.ActivityStatus(2)
			activeStatus["domain.com"]["456"] = dto.ActivityStatus(1)

			err := pd.FromModel(tt.args.mod, activeStatus, models.Confiant{}, models.Pixalate{}, []models.BidCaching{}, []models.RefreshCache{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			// object mutates
			assert.Equal(t, tt.want, pd)
		})
	}
}

func Test_withConfiantAndBidcacheData(t *testing.T) {
	t.Parallel()
	publisherId := "456"
	domain := "domain.com"
	confiantKey := "1234"
	rate := 10.0
	bidCaching := make([]dto.BidCaching, 0)
	bidCaching = append(bidCaching, dto.BidCaching{
		RuleID:     "999",
		Publisher:  publisherId,
		Domain:     domain,
		Country:    "IL",
		BidCaching: 30,
		Active:     true,
	})

	type args struct {
		mod *dto.PublisherDetailModel
	}

	tests := []struct {
		name    string
		args    args
		want    dto.PublisherDetail
		wantErr bool
	}{
		{
			name: "valid_withUserFullnameAndBidcacheData",
			args: args{
				mod: &dto.PublisherDetailModel{
					Publisher: models.Publisher{
						Name:             "publisher",
						PublisherID:      "456",
						AccountManagerID: null.String{Valid: true, String: "1"},
					},
					PublisherDomain: models.PublisherDomain{
						Domain:     "domain.com",
						Automation: true,
						GPPTarget:  null.Float64{Valid: true, Float64: 0.1},
					},
					User: dto.UserModelCompact{
						ID:        1,
						FirstName: null.StringFrom("Idan"),
						LastName:  null.StringFrom("Finkiel"),
					},
				},
			},

			want: dto.PublisherDetail{
				Name:                   "publisher",
				PublisherID:            publisherId,
				AccountManagerID:       "1",
				AccountManagerFullName: "Idan Finkiel",
				Domain:                 domain,
				Automation:             true,
				GPPTarget:              0.1,
				ActivityStatus:         "Low",
				Confiant: dto.Confiant{
					ConfiantKey: &confiantKey,
					PublisherID: publisherId,
					Domain:      &domain,
					Rate:        &rate,
				},
				Pixalate: dto.Pixalate{
					PublisherID: publisherId,
					Domain:      &domain,
					Rate:        &rate,
				},
				BidCaching:   bidCaching,
				RefreshCache: make([]dto.RefreshCache, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pd := dto.PublisherDetail{}

			activeStatus := make(map[string]map[string]dto.ActivityStatus)
			activeStatus["domain.com"] = make(map[string]dto.ActivityStatus)
			activeStatus["domain.com"]["123"] = dto.ActivityStatus(2)
			activeStatus["domain.com"]["456"] = dto.ActivityStatus(1)

			err := pd.FromModel(tt.args.mod, activeStatus,
				models.Confiant{
					ConfiantKey: confiantKey,
					PublisherID: publisherId,
					Domain:      domain,
					Rate:        rate,
				}, models.Pixalate{
					ID:          "1234",
					PublisherID: publisherId,
					Domain:      domain,
					Rate:        rate,
				},
				[]models.BidCaching{models.BidCaching{
					Publisher: publisherId,
					RuleID:    "999",
					Domain: null.String{
						String: domain,
						Valid:  true,
					},
					BidCaching: 30,
					Country: null.String{
						String: "IL",
						Valid:  true,
					},
					Active: true,
				}}, []models.RefreshCache{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			// object mutates
			assert.Equal(t, tt.want.Confiant.ConfiantKey, pd.Confiant.ConfiantKey)
			assert.Equal(t, tt.want.Confiant.PublisherID, pd.Confiant.PublisherID)
			assert.Equal(t, tt.want.Confiant.Domain, pd.Confiant.Domain)
			assert.Equal(t, tt.want.GPPTarget, pd.GPPTarget)
			assert.Equal(t, tt.want.Pixalate.Rate, pd.Pixalate.Rate)
			assert.Equal(t, tt.want.BidCaching[0].BidCaching, pd.BidCaching[0].BidCaching)
			assert.Equal(t, tt.want.BidCaching[0].RuleID, pd.BidCaching[0].RuleID)
			assert.Equal(t, tt.want.BidCaching[0].Country, pd.BidCaching[0].Country)
		})
	}
}
