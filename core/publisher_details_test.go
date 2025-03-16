package core

import (
	"testing"

	"github.com/m6yf/bcwork/dto"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
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
