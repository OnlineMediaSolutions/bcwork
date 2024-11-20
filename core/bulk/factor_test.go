package bulk

import (
	"testing"

	"github.com/volatiletech/null/v8"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
)

func Test_prepareBulkInsertFactorsRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		factors []*models.Factor
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				factors: []*models.Factor{
					{
						Publisher: "publisher_1",
						Domain:    "1.com",
						Device:    null.StringFrom("mobile"),
						Country:   null.StringFrom("IL"),
						Factor:    0.1,
					},
					{
						Publisher: "publisher_2",
						Domain:    "2.com",
						Device:    null.StringFrom("mobile"),
						Country:   null.StringFrom("US"),
						Factor:    0.05,
					},
					{
						Publisher: "publisher_3",
						Domain:    "3.com",
						Device:    null.StringFrom("mobile"),
						Country:   null.StringFrom("RU"),
						Factor:    0.15,
					},
				},
			},
			want: &bulkInsertRequest{
				tableName: models.TableNames.Factor,
				columns: []string{
					models.FactorColumns.Publisher,
					models.FactorColumns.Domain,
					models.FactorColumns.Device,
					models.FactorColumns.Country,
					models.FactorColumns.Factor,
					models.FactorColumns.RuleID,
					models.FactorColumns.CreatedAt,
					models.FactorColumns.UpdatedAt,
				},
				conflictColumns: []string{
					models.FactorColumns.RuleID,
				},
				updateColumns: []string{
					models.FactorColumns.Factor,
					models.FactorColumns.UpdatedAt,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8)",
					"($9, $10, $11, $12, $13, $14, $15, $16)",
					"($17, $18, $19, $20, $21, $22, $23, $24)",
				},
				args: []interface{}{
					"publisher_1", "1.com", null.String{String: "mobile", Valid: true}, null.String{String: "IL", Valid: true}, 0.1, "", currentTime, currentTime,
					"publisher_2", "2.com", null.String{String: "mobile", Valid: true}, null.String{String: "US", Valid: true}, 0.05, "", currentTime, currentTime,
					"publisher_3", "3.com", null.String{String: "mobile", Valid: true}, null.String{String: "RU", Valid: true}, 0.15, "", currentTime, currentTime,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertFactorsRequest(tt.args.factors)
			assert.Equal(t, tt.want, got)
		})
	}
}
