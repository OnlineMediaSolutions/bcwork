package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_prepareBulkInsertGlobalFactorsRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		globalFactors []models.GlobalFactor
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				globalFactors: []models.GlobalFactor{
					{
						Key:         "key_1",
						PublisherID: "1",
						Value:       null.Float64{Valid: true, Float64: 0.1},
					},
					{
						Key:         "key_2",
						PublisherID: "2",
						Value:       null.Float64{Valid: true, Float64: 0.05},
					},
					{
						Key:         "key_3",
						PublisherID: "3",
						Value:       null.Float64{Valid: true, Float64: 0.15},
					},
				},
			},
			want: &bulkInsertRequest{
				tableName: models.TableNames.GlobalFactor,
				columns: []string{
					models.GlobalFactorColumns.Key,
					models.GlobalFactorColumns.PublisherID,
					models.GlobalFactorColumns.Value,
					models.GlobalFactorColumns.CreatedAt,
					models.GlobalFactorColumns.UpdatedAt,
				},
				conflictColumns: []string{
					models.GlobalFactorColumns.Key,
					models.GlobalFactorColumns.PublisherID,
				},
				updateColumns: []string{
					models.GlobalFactorColumns.Value,
					models.GlobalFactorColumns.UpdatedAt,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5)",
					"($6, $7, $8, $9, $10)",
					"($11, $12, $13, $14, $15)",
				},
				args: []interface{}{
					"key_1", "1", null.Float64{Valid: true, Float64: 0.1}, currentTime, currentTime,
					"key_2", "2", null.Float64{Valid: true, Float64: 0.05}, currentTime, currentTime,
					"key_3", "3", null.Float64{Valid: true, Float64: 0.15}, currentTime, currentTime,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertGlobalFactorsRequest(tt.args.globalFactors)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_prepareGlobalFactorsData(t *testing.T) {
	t.Parallel()

	type args struct {
		chunk []GlobalFactorRequest
	}

	tests := []struct {
		name string
		args args
		want []models.GlobalFactor
	}{
		{
			name: "valid",
			args: args{
				chunk: []GlobalFactorRequest{
					{
						Key:       "key_1",
						Publisher: "1",
						Value:     0.1,
					},
					{
						Key:       "key_2",
						Publisher: "2",
						Value:     0.05,
					},
					{
						Key:       "key_3",
						Publisher: "3",
						Value:     0.15,
					},
				},
			},
			want: []models.GlobalFactor{
				{
					Key:         "key_1",
					PublisherID: "1",
					Value:       null.Float64From(0.1),
				},
				{
					Key:         "key_2",
					PublisherID: "2",
					Value:       null.Float64From(0.05),
				},
				{
					Key:         "key_3",
					PublisherID: "3",
					Value:       null.Float64From(0.15),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareGlobalFactorsData(tt.args.chunk)
			assert.Equal(t, tt.want, got)
		})
	}
}
