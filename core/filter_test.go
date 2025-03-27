package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterService_GetFilterFields(t *testing.T) {
	t.Parallel()

	type args struct {
		filterName string
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filterName: "ads_txt_group_by_dp",
			},
			want: []string{
				"publisher_id",
				"account_manager_id",
				"campaign_manager_id",
				"domain",
				"media_type",
				"demand_status",
				"domain_status",
				"demand_manager_id",
				"demand_partner_name",
				"is_ready_to_go_live",
			},
		},
		{
			name: "unknownFilter",
			args: args{
				filterName: "somefilter",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewFilterService()

			got, err := s.GetFilterFields(context.Background(), tt.args.filterName)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
