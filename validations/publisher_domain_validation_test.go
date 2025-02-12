package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validatePublisherDomain(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.PublisherDomainUpdateRequest
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &dto.PublisherDomainUpdateRequest{
					PublisherID:     "publisher",
					Domain:          "domain.com",
					IntegrationType: []string{dto.ORTBIntergrationType},
				},
			},
			want: []string{},
		},
		{
			name: "invalid",
			args: args{
				request: &dto.PublisherDomainUpdateRequest{},
			},
			want: []string{
				"PublisherID is mandatory, validation failed",
				"Domain is mandatory, validation failed",
				// "integration type must be in allowed list: oRTB,Prebid Server,Amazon APS",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validatePublisherDomain(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
