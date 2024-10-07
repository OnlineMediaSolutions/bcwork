package core

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_getColumnsToUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		newData     *constant.Targeting
		currentData *models.Targeting
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "updateAllFields",
			args: args{
				newData: &constant.Targeting{
					Country:       []string{"il", "us"},
					DeviceType:    []string{"mobile"},
					OS:            []string{"linux"},
					Browser:       []string{"firefox"},
					PlacementType: "rectangle",
					KV:            map[string]string{"key_1": "value_1"},
					PriceModel:    constant.TargetingPriceModelCPM,
					Value:         5,
					Status:        constant.TargetingStatusPaused,
					DailyCap:      5000,
				},
				currentData: &models.Targeting{
					PriceModel: constant.TargetingPriceModelRevShare,
					Value:      0.5,
					Status:     constant.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
				},
			},
			want: []string{
				models.TargetingColumns.Country,
				models.TargetingColumns.DeviceType,
				models.TargetingColumns.Os,
				models.TargetingColumns.Browser,
				models.TargetingColumns.PlacementType,
				models.TargetingColumns.KV,
				models.TargetingColumns.PriceModel,
				models.TargetingColumns.Value,
				models.TargetingColumns.Status,
				models.TargetingColumns.DailyCap,
			},
		},
		{
			name: "updatePartialFields",
			args: args{
				newData: &constant.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Value:      5,
					Status:     constant.TargetingStatusActive,
					DailyCap:   5000,
				},
				currentData: &models.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Value:      4,
					Status:     constant.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
				},
			},
			want: []string{
				models.TargetingColumns.Value,
				models.TargetingColumns.DailyCap,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := getColumnsToUpdate(tt.args.newData, tt.args.currentData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
