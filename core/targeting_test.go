package core

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_getColumnsToUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		newData     *dto.Targeting
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
				newData: &dto.Targeting{
					Country:       []string{"il", "us"},
					DeviceType:    []string{"mobile"},
					OS:            []string{"linux"},
					Browser:       []string{"firefox"},
					PlacementType: "rectangle",
					KV:            map[string]string{"key_1": "value_1"},
					PriceModel:    dto.TargetingPriceModelCPM,
					Value:         5,
					Status:        dto.TargetingStatusPaused,
					DailyCap:      func() *int { i := 5000; return &i }(),
				},
				currentData: &models.Targeting{
					PriceModel: dto.TargetingPriceModelRevShare,
					Value:      0.5,
					Status:     dto.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
				},
			},
			want: []string{
				models.TargetingColumns.RuleID,
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
				newData: &dto.Targeting{
					PriceModel: dto.TargetingPriceModelCPM,
					Value:      5,
					Status:     dto.TargetingStatusActive,
					DailyCap:   func() *int { i := 5000; return &i }(),
					KV:         map[string]string{"key_1": "value_old"},
				},
				currentData: &models.Targeting{
					PriceModel: dto.TargetingPriceModelCPM,
					Value:      4,
					Status:     dto.TargetingStatusActive,
					DailyCap:   null.IntFrom(3000),
					KV:         null.JSONFrom([]byte(`{"key_1": "value_new"}`)),
				},
			},
			want: []string{
				models.TargetingColumns.RuleID,
				models.TargetingColumns.UpdatedAt,
				models.TargetingColumns.KV,
				models.TargetingColumns.Value,
				models.TargetingColumns.DailyCap,
			},
		},
		{
			name: "updateToNullKV",
			args: args{
				newData: &dto.Targeting{
					PriceModel: dto.TargetingPriceModelCPM,
					Status:     dto.TargetingStatusActive,
				},
				currentData: &models.Targeting{
					PriceModel: dto.TargetingPriceModelCPM,
					Status:     dto.TargetingStatusActive,
					KV:         null.JSONFrom([]byte(`{"key_1": "value_new"}`)),
				},
			},
			want: []string{
				models.TargetingColumns.RuleID,
				models.TargetingColumns.UpdatedAt,
				models.TargetingColumns.KV,
			},
		},
		{
			name: "updateDailyCapToNull",
			args: args{
				newData: &dto.Targeting{
					DailyCap: nil,
				},
				currentData: &models.Targeting{
					DailyCap: null.IntFrom(3000),
				},
			},
			want: []string{
				models.TargetingColumns.RuleID,
				models.TargetingColumns.UpdatedAt,
				models.TargetingColumns.DailyCap,
			},
		},
		{
			name: "nothingToUpdate",
			args: args{
				newData: &dto.Targeting{
					Country:       []string{"il", "us"},
					DeviceType:    []string{"mobile"},
					OS:            []string{"linux"},
					Browser:       []string{"firefox"},
					PlacementType: "rectangle",
					KV:            map[string]string{"key_1": "value_1"},
					PriceModel:    dto.TargetingPriceModelCPM,
					Value:         5,
					Status:        dto.TargetingStatusActive,
					DailyCap:      func() *int { i := 5000; return &i }(),
				},
				currentData: &models.Targeting{
					Country:       []string{"il", "us"},
					DeviceType:    []string{"mobile"},
					Os:            []string{"linux"},
					Browser:       []string{"firefox"},
					PlacementType: null.StringFrom("rectangle"),
					KV:            null.JSONFrom([]byte(`{"key_1": "value_1"}`)),
					PriceModel:    dto.TargetingPriceModelCPM,
					Value:         5,
					Status:        dto.TargetingStatusActive,
					DailyCap:      null.IntFrom(5000),
				},
			},
			want: []string{
				models.TargetingColumns.RuleID,
				models.TargetingColumns.UpdatedAt,
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
		data *dto.Targeting
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
				data: &dto.Targeting{ID: 0},
			},
			want: false,
		},
		{
			name: "duplicate_modFound_creatingNewTargeting",
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &dto.Targeting{ID: 0},
			},
			want: true,
		},
		{
			name: "noDuplicate_modFound_equalIDs", // when updating same entity
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &dto.Targeting{ID: 2},
			},
			want: false,
		},
		{
			name: "duplicate_modFound_differentIDs", // conflict when updating entity
			args: args{
				mod:  &models.Targeting{ID: 2},
				data: &dto.Targeting{ID: 5},
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

func Test_createTargetingMetaData(t *testing.T) {
	t.Parallel()

	type args struct {
		mods      models.TargetingSlice
		publisher string
		domain    string
	}

	tests := []struct {
		name    string
		args    args
		want    *models.MetadataQueue
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				mods: models.TargetingSlice{
					{
						RuleID:        "rule_1",
						PublisherID:   "publisher_1",
						Domain:        "1.com",
						UnitSize:      "300x200",
						PlacementType: null.StringFrom("placement"),
						Country:       []string{"il", "ru", "uk", "us"},
						DeviceType:    []string{"desktop", "mobile"},
						Browser:       []string{"chrome", "edge", "firefox"},
						Os:            []string{"linux", "macos", "windows"},
						KV:            null.JSONFrom([]byte(`{"key1": "value1", "key2": "value2", "key3": "value3"}`)),
						Status:        dto.TargetingStatusActive,
						PriceModel:    dto.TargetingPriceModelCPM,
						Value:         0.01,
						DailyCap:      null.IntFrom(5),
					},
					{
						RuleID:        "rule_2",
						PublisherID:   "publisher_1",
						Domain:        "1.com",
						UnitSize:      "400x100",
						PlacementType: null.StringFrom("rectangle"),
						Country:       []string{"cn", "fr"},
						DeviceType:    []string{"desktop"},
						Browser:       []string{"edge", "opera"},
						Os:            []string{"linux", "macos", "windows"},
						Status:        dto.TargetingStatusPaused,
						PriceModel:    dto.TargetingPriceModelRevShare,
					},
				},
				publisher: "publisher_1",
				domain:    "1.com",
			},
			want: &models.MetadataQueue{
				Key:   utils.JSTagMetaDataKeyPrefix + ":publisher_1:1.com",
				Value: []byte(`[{"rule_id":"rule_1","rule":"p=publisher_1__d=1.com__s=300x200__c=(il|ru|uk|us)__os=(linux|macos|windows)__dt=(desktop|mobile)__pt=placement__b=(chrome|edge|firefox)__key1=value1__key2=value2__key3=value3","price_model":"cpm","value":0.01,"daily_cap":5},{"rule_id":"rule_2","rule":"p=publisher_1__d=1.com__s=400x100__c=(cn|fr)__os=(linux|macos|windows)__dt=(desktop)__pt=rectangle__b=(edge|opera)","price_model":"revshare","value":0,"daily_cap":null}]`),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := createTargetingMetaData(tt.args.mods, tt.args.publisher, tt.args.domain)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			got.TransactionID = "" // because it depends on current time
			assert.Equal(t, tt.want, got)
		})
	}
}
