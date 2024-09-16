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
		want errorBulkResponse
	}{
		{
			name: "valid",
			args: args{
				requests: []*core.GlobalFactorRequest{
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
				},
			},
			want: errorBulkResponse{
				Errors: make(map[string][]string),
			},
		},
		{
			name: "validationError",
			args: args{
				requests: []*core.GlobalFactorRequest{
					{
						Key:       "consultant_fee",
						Publisher: "publisher_1",
						Value:     5,
					},
					{
						Key:       "unknown_fee",
						Publisher: "publisher_2",
						Value:     10,
					},
				},
			},
			want: errorBulkResponse{
				Status:  errorStatus,
				Message: validationError,
				Errors: map[string][]string{
					"request 2": {keyValidationError},
				},
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
