package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateAdsTxt(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.AdsTxt
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &dto.AdsTxt{
					DomainStatus: dto.DomainStatusActive,
					DemandStatus: dto.DPStatusApproved,
				},
			},
			want: []string{},
		},
		{
			name: "invalid",
			args: args{
				request: &dto.AdsTxt{
					DomainStatus: "some_domain_status",
					DemandStatus: "some_demand_status",
				},
			},
			want: []string{
				adsTxtDomainStatusErrorMessage,
				adsTxtDemandStatusErrorMessage,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateAdsTxt(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
