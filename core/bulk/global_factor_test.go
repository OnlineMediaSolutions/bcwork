package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_prepareBulkInsertGlobalFactorsRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		globalFactors []*models.GlobalFactor
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				globalFactors: []*models.GlobalFactor{
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
					"key_1", "1", null.Float64{Valid: true, Float64: 0.1}, constant.PostgresCurrentTime, constant.PostgresCurrentTime,
					"key_2", "2", null.Float64{Valid: true, Float64: 0.05}, constant.PostgresCurrentTime, constant.PostgresCurrentTime,
					"key_3", "3", null.Float64{Valid: true, Float64: 0.15}, constant.PostgresCurrentTime, constant.PostgresCurrentTime,
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
