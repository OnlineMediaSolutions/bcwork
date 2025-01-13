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
				ActivityStatus:   "",
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
				ActivityStatus:   "Low",
			},
		},
		{
			name: "valid_withUserFullname",
			args: args{
				mod: &dto.PublisherDetailModel{
					Publisher: models.Publisher{
						Name:             "publisher",
						PublisherID:      "123",
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
				PublisherID:            "123",
				AccountManagerID:       "1",
				AccountManagerFullName: "first last",
				Domain:                 "domain.com",
				Automation:             true,
				GPPTarget:              0.1,
				ActivityStatus:         "Low",
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

			err := pd.FromModel(tt.args.mod, activeStatus)
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
