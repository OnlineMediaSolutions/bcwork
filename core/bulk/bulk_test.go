package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
)

func Test_prepareBulkInsertQuery(t *testing.T) {
	t.Parallel()

	type args struct {
		req *bulkInsertRequest
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				req: &bulkInsertRequest{
					tableName: models.TableNames.GlobalFactor,
					columns: []string{
						models.GlobalFactorColumns.Key,
						models.GlobalFactorColumns.PublisherID,
						models.GlobalFactorColumns.Value,
					},
					conflictColumns: []string{
						models.GlobalFactorColumns.Key,
						models.GlobalFactorColumns.PublisherID,
					},
					updateColumns: []string{
						models.GlobalFactorColumns.Value,
					},
					valueStrings: []string{
						"($1, $2, $3)",
						"($4, $5, $6)",
					},
					args: []interface{}{"key_1", "1", 0.1, "key_2", "2", 0.05},
				},
			},
			want: `INSERT INTO global_factor (key, publisher_id, value) VALUES ` +
				`($1, $2, $3),($4, $5, $6)` +
				` ON CONFLICT (key, publisher_id) DO UPDATE SET value = EXCLUDED.value`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertQuery(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}
