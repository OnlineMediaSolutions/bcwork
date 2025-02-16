package validations

import (
	"testing"

	"github.com/m6yf/bcwork/core"
	"github.com/stretchr/testify/assert"
)

func Test_validateBulkGlobalFactor(t *testing.T) {
	t.Parallel()

	type args struct {
		requests []*core.GlobalFactorRequest
	}

	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "valid",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:   "tech_fee",
						Value: 1,
					},
					{
						Key:       "consultant_fee",
						Publisher: "publisher_1",
						Value:     5,
					},
					{
						Key:       "consultant_fee",
						Publisher: "publisher_2",
						Value:     10,
					},
					{
						Key:   "tam_fee",
						Value: 3,
					},
				},
			},
			want: map[string][]string{},
		},
		{
			name: "whenFeeNameIsIncorrect_ThenReturnAnError",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:       "consultant_fee",
						Publisher: "publisher_1",
						Value:     5,
					},
					{
						Key:   "unknown_fee",
						Value: 10,
					},
				},
			},
			want: map[string][]string{
				"request 2": {keyValidationError},
			},
		},
		{
			name: "whenNotConsultantFeeHavePublisher_ThenReturnAnError",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:       "tech_fee",
						Publisher: "publisher_1",
						Value:     5,
					},
					{
						Key:       "tam_fee",
						Publisher: "publisher_2",
						Value:     10,
					},
				},
			},
			want: map[string][]string{
				"request 1": {publisherValidationError},
				"request 2": {publisherValidationError},
			},
		},
		{
			name: "whenValueLessThan0_ThenReturnAnError",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:   "tech_fee",
						Value: -1,
					},
					{
						Key:   "tam_fee",
						Value: 0,
					},
				},
			},
			want: map[string][]string{
				"request 1": {valueValidationError},
			},
		},
		{
			name: "whenConsultantFeeDontHavePublisher_ThenReturnAnError",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:   "consultant_fee",
						Value: 5,
					},
				},
			},
			want: map[string][]string{
				"request 1": {consultantFeeValidationError},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateBulkGlobalFactor(tt.args.requests)
			assert.Equal(t, tt.want, got)
		})
	}
}
