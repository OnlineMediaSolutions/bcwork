package core

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_PublisherDetail_FromModel(t *testing.T) {
	t.Parallel()

	type args struct {
		mod    *models.Publisher
		domain *models.PublisherDomain
	}

	tests := []struct {
		name    string
		args    args
		want    PublisherDetail
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				mod: &models.Publisher{
					Name:             "publisher",
					PublisherID:      "id",
					AccountManagerID: null.String{Valid: true, String: "am_id"},
				},
				domain: &models.PublisherDomain{
					Domain:     "domain",
					Automation: true,
					GPPTarget:  null.Float64{Valid: true, Float64: 0.1},
				},
			},
			want: PublisherDetail{
				Name:             "publisher",
				PublisherID:      "id",
				AccountManagerID: "am_id",
				Domain:           "domain",
				Automation:       true,
				GPPTarget:        0.1,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pd := PublisherDetail{}

			err := pd.FromModel(tt.args.mod, tt.args.domain)
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
