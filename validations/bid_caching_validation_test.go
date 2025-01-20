package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateBidCaching(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.BidCaching
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "test",
					Domain:     "example.com",
					BidCaching: 1,
					ControlPercentage: func() *float64 {
						var n float64 = 0.5
						return &n
					}(),
				},
			},
			want: []string{},
		},
		{
			name: "valid_withoutControlPercentage",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "publisher",
					Domain:     "domain",
					BidCaching: 5,
				},
			},
			want: []string{},
		},
		{
			name: "valid_withMinControlPercentage",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "publisher",
					Domain:     "domain",
					BidCaching: 5,
					ControlPercentage: func() *float64 {
						var n float64 = 0
						return &n
					}(),
				},
			},
			want: []string{},
		},
		{
			name: "valid_withMaxControlPercentage",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "publisher",
					Domain:     "domain",
					BidCaching: 5,
					ControlPercentage: func() *float64 {
						var n float64 = 1
						return &n
					}(),
				},
			},
			want: []string{},
		},
		{
			name: "invalid_bidCachingLessThanMinimalValue",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "publisher",
					Domain:     "domain",
					BidCaching: 0,
					ControlPercentage: func() *float64 {
						var n float64 = 1
						return &n
					}(),
				},
			},
			want: []string{
				"Bid caching value not allowed, it should be >= 1",
			},
		},
		{
			name: "invalid_withControlPercentage_wrongValue",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "publisher",
					Domain:     "domain",
					BidCaching: 5,
					ControlPercentage: func() *float64 {
						n := 1.01
						return &n
					}(),
				},
			},
			want: []string{
				bidCachingControlPercentageErrorMessage,
			},
		},
		{
			name: "invalid_wrongDevice",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "1234",
					Domain:     "example.com",
					BidCaching: 1,
					Country:    "us",
					Device:     "mm",
				},
			},
			want: []string{
				"Device should be in the allowed list",
			},
		},
		{
			name: "invalid_wrongCountry",
			args: args{
				request: &dto.BidCaching{
					Publisher:  "test",
					Domain:     "example.com",
					BidCaching: 1,
					Country:    "USA",
					Device:     "tablet",
				},
			},
			want: []string{
				"Country code must be 2 characters long and should be in the allowed list",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateBidCache(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
