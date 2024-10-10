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
				models.TargetingColumns.UpdatedAt,
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
					KV:         map[string]string{"key_1": "value_old"},
				},
				currentData: &models.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Value:      4,
					Status:     constant.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
					KV:         null.JSONFrom([]byte(`{"key_1": "value_new"}`)),
				},
			},
			want: []string{
				models.TargetingColumns.UpdatedAt,
				models.TargetingColumns.KV,
				models.TargetingColumns.Value,
				models.TargetingColumns.DailyCap,
			},
		},
		{
			name: "updateToNullKV",
			args: args{
				newData: &constant.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Status:     constant.TargetingStatusActive,
				},
				currentData: &models.Targeting{
					PriceModel: constant.TargetingPriceModelCPM,
					Status:     constant.TargetingStatusActive,
					KV:         null.JSONFrom([]byte(`{"key_1": "value_new"}`)),
				},
			},
			want: []string{
				models.TargetingColumns.UpdatedAt,
				models.TargetingColumns.KV,
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

func Test_isDuplicate(t *testing.T) {
	t.Parallel()

	type args struct {
		mod  *models.Targeting
		data *constant.Targeting
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "noDuplicate_modNotFound",
			args: args{
				mod:  nil,
				data: &constant.Targeting{ID: 0},
			},
			want: false,
		},
		{
			name: "duplicate_modFound_creatingNewTargeting",
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &constant.Targeting{ID: 0},
			},
			want: true,
		},
		{
			name: "noDuplicate_modFound_equalIDs", // when updating same entity
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &constant.Targeting{ID: 2},
			},
			want: false,
		},
		{
			name: "duplicate_modFound_differentIDs", // conflict when updating entity
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &constant.Targeting{ID: 5},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isDuplicate(tt.args.mod, tt.args.data)
			assert.Equal(t, tt.want, got)
		})
	}
}
