package validations

import (
	"fmt"
	"testing"

	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
)

func Test_validateTargeting(t *testing.T) {
	t.Parallel()

	type args struct {
		request *constant.Targeting
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         5,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{},
		},
		{
			name: "whenNoAllowedCostModel_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    "unknown_price_model",
					Value:         5,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{
				targetingCostModelValidationErrorMessage,
			},
		},
		{
			name: "whenNoAllowedValueForRevShareCostModel_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelRevShare,
					Value:         5,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{
				fmt.Sprintf("Rev Share Value should be between %v and %v",
					constant.TargetingMinValueCostModelRevShare, constant.TargetingMaxValueCostModelRevShare,
				),
			},
		},
		{
			name: "whenNoAllowedValueForCPMCostModel_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         55,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{
				fmt.Sprintf("CPM Value should be between %v and %v",
					constant.TargetingMinValueCostModelCPM, constant.TargetingMaxValueCostModelCPM,
				),
			},
		},
		{
			name: "whenNoAllowedCountry_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"xx"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         5,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{
				countryValidationErrorMessage,
			},
		},
		{
			name: "whenNoAllowedDevice_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"new_device"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         5,
					Status:        constant.TargetingStatusActive,
				},
			},
			want: []string{
				deviceValidationErrorMessage,
			},
		},
		{
			name: "whenNoAllowedStatus_ThenError",
			args: args{
				request: &constant.Targeting{
					PublisherID:   "publisher",
					Domain:        "1.com",
					UnitSize:      "1X1",
					PlacementType: "placement_type",
					Country:       []string{"il"},
					DeviceType:    []string{"mobile"},
					Browser:       []string{"firefox"},
					OS:            []string{"linux"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         5,
					Status:        "unknown",
				},
			},
			want: []string{
				targetingStatusValidationErrorMessage,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateTargeting(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}