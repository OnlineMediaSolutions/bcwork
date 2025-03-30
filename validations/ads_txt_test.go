package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_validateAdsTxtUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.AdsTxtUpdateRequest
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid_updateDomainStatus",
			args: args{
				request: &dto.AdsTxtUpdateRequest{
					Domain:       []string{"test.com"},
					DomainStatus: helpers.GetPointerToString(dto.DomainStatusActive),
				},
			},
			want: []string{},
		},
		{
			name: "valid_updateDemandStatus",
			args: args{
				request: &dto.AdsTxtUpdateRequest{
					Domain:          []string{"test.com"},
					DemandPartnerID: helpers.GetPointerToString("demand_partner_id"),
					DemandStatus:    helpers.GetPointerToString(dto.DPStatusApproved),
				},
			},
			want: []string{},
		},
		{
			name: "invalid_noDomains",
			args: args{
				request: &dto.AdsTxtUpdateRequest{
					Domain:       []string{},
					DomainStatus: helpers.GetPointerToString("some_domain_status"),
					DemandStatus: helpers.GetPointerToString("some_demand_status"),
				},
			},
			want: []string{
				"Domain is mandatory, validation failed",
				adsTxtDomainStatusErrorMessage,
				adsTxtDemandStatusErrorMessage,
			},
		},
		{
			name: "invalid_emptyRequest",
			args: args{
				request: &dto.AdsTxtUpdateRequest{},
			},
			want: []string{
				"Domain is mandatory, validation failed",
				"domain and demand statuses are nil",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateAdsTxtUpdate(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
