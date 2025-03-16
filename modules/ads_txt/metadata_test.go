package adstxt

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/stretchr/testify/assert"
)

func Test_createAdsTxtMetaData(t *testing.T) {
	t.Parallel()

	type args struct {
		data map[string]*dto.AdsTxtGroupedByDPData
	}

	tests := []struct {
		name    string
		args    args
		want    []*models.MetadataQueue
		wantErr bool
	}{
		{
			name: "valid_allReadyToWork",
			args: args{
				data: map[string]*dto.AdsTxtGroupedByDPData{
					"10000:test.net:Yieldmo:inapp": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10000",
							Domain:          "test.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: true,
						},
					},
					"10000:test.net:Yieldmo:video": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10000",
							Domain:          "test.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: true,
						},
					},
					"10001:test1.net:OpenX:inapp": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10001",
							Domain:          "test1.net",
							DemandPartnerID: "openx",
							IsReadyToGoLive: true,
						},
					},
					"10001:test1.net:OpenX:video,banner": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10001",
							Domain:          "test1.net",
							DemandPartnerID: "openx",
							IsReadyToGoLive: true,
						},
					},
				},
			},
			want: []*models.MetadataQueue{
				{
					Key:   fmt.Sprintf(utils.AdsTxtMetaDataKeyTemplate, "openx"),
					Value: []byte(`[{"pubid":"10001","domain":"test1.net"}]`),
				},
				{
					Key:   fmt.Sprintf(utils.AdsTxtMetaDataKeyTemplate, "yieldmo"),
					Value: []byte(`[{"pubid":"10000","domain":"test.net"}]`),
				},
			},
		},
		{
			name: "valid_noDemandPartnerReadyToWork",
			args: args{
				data: map[string]*dto.AdsTxtGroupedByDPData{
					"9994:reverso.net:Yieldmo:video": {
						Parent: &dto.AdsTxt{
							PublisherID:     "9994",
							Domain:          "reverso.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: false,
						},
					},
					"9994:reverso.net:Yieldmo:inapp": {
						Parent: &dto.AdsTxt{
							PublisherID:     "9994",
							Domain:          "reverso.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: false,
						},
					},
					"10000:test.net:Yieldmo:video": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10000",
							Domain:          "test.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: false,
						},
					},
					"9994:reverso.net:Yieldmo:banner": {
						Parent: &dto.AdsTxt{
							PublisherID:     "9994",
							Domain:          "reverso.net",
							DemandPartnerID: "yieldmo",
							IsReadyToGoLive: false,
						},
					},
					"10001:test1.net:OpenX:inapp": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10001",
							Domain:          "test1.net",
							DemandPartnerID: "openx",
							IsReadyToGoLive: false,
						},
					},
					"10001:test1.net:OpenX:video,banner": {
						Parent: &dto.AdsTxt{
							PublisherID:     "10001",
							Domain:          "test1.net",
							DemandPartnerID: "openx",
							IsReadyToGoLive: false,
						},
					},
				},
			},
			want: []*models.MetadataQueue{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := createAdsTxtMetaData(context.Background(), tt.args.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			// sorting to make test stable
			sort.SliceStable(got, func(i, j int) bool { return got[i].Key < got[j].Key })
			for i := range got {
				got[i].TransactionID = "" // because it depends on current time
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
