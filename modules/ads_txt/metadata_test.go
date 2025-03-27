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
		resp *dto.AdsTxtGroupByDPResponse
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
				resp: &dto.AdsTxtGroupByDPResponse{
					Data: []*dto.AdsTxtGroupedByDPData{
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10000",
								Domain:          "test.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: true,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10000",
								Domain:          "test.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: true,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10001",
								Domain:          "test1.net",
								DemandPartnerID: "openx",
								IsReadyToGoLive: true,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10001",
								Domain:          "test1.net",
								DemandPartnerID: "openx",
								IsReadyToGoLive: true,
							},
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
				resp: &dto.AdsTxtGroupByDPResponse{
					Data: []*dto.AdsTxtGroupedByDPData{
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "9994",
								Domain:          "reverso.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: false,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "9994",
								Domain:          "reverso.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: false,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10000",
								Domain:          "test.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: false,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "9994",
								Domain:          "reverso.net",
								DemandPartnerID: "yieldmo",
								IsReadyToGoLive: false,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10001",
								Domain:          "test1.net",
								DemandPartnerID: "openx",
								IsReadyToGoLive: false,
							},
						},
						{
							Parent: &dto.AdsTxt{
								PublisherID:     "10001",
								Domain:          "test1.net",
								DemandPartnerID: "openx",
								IsReadyToGoLive: false,
							},
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

			got, err := createAdsTxtMetaData(context.Background(), tt.args.resp)
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
