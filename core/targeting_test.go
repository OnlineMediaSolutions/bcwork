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
		name string
		args args
		want []string
	}{
		{
			name: "updateAllFields",
			args: args{
				newData: &constant.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Value:      5,
					Status:     constant.TargetingStatusPaused,
					DailyCap:   5000,
				},
				currentData: &models.Targeting{
					PriceModel: constant.TargetingPriceModelRevShare,
					Value:      0.5,
					Status:     constant.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
				},
			},
			want: []string{
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

			got := getColumnsToUpdate(tt.args.newData, tt.args.currentData)
			assert.Equal(t, tt.want, got)
		})
	}
}
