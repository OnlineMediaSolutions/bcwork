package bulk

import (
	"fmt"
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
)

func Test_prepareBulkInsertFactorsRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		factors []models.Factor
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				factors: []models.Factor{
					{
						Publisher: "publisher_1",
						Domain:    "1.com",
						Device:    "mobile",
						Country:   "IL",
						Factor:    0.1,
					},
					{
						Publisher: "publisher_2",
						Domain:    "2.com",
						Device:    "mobile",
						Country:   "US",
						Factor:    0.05,
					},
					{
						Publisher: "publisher_3",
						Domain:    "3.com",
						Device:    "mobile",
						Country:   "RU",
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
					models.FactorColumns.CreatedAt,
					models.FactorColumns.UpdatedAt,
				},
				conflictColumns: []string{
					models.FactorColumns.Publisher,
					models.FactorColumns.Domain,
					models.FactorColumns.Device,
					models.FactorColumns.Country,
				},
				updateColumns: []string{
					models.FactorColumns.Factor,
					models.FactorColumns.UpdatedAt,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5, $6, $7)",
					"($8, $9, $10, $11, $12, $13, $14)",
					"($15, $16, $17, $18, $19, $20, $21)",
				},
				args: []interface{}{
					"publisher_1", "1.com", "mobile", "IL", 0.1, currentTime, currentTime,
					"publisher_2", "2.com", "mobile", "US", 0.05, currentTime, currentTime,
					"publisher_3", "3.com", "mobile", "RU", 0.15, currentTime, currentTime,
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

func Test_prepareFactorsData(t *testing.T) {
	t.Parallel()

	type args struct {
		chunk []FactorUpdateRequest
	}

	type want struct {
		factors  []models.Factor
		metadata []models.MetadataQueue
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "valid",
			args: args{
				chunk: []FactorUpdateRequest{
					{
						Publisher: "publisher_1",
						Domain:    "1.com",
						Device:    "mobile",
						Factor:    0.1,
						Country:   "il",
					},
					{
						Publisher: "publisher_2",
						Domain:    "2.com",
						Device:    "web",
						Factor:    0.05,
						Country:   "us",
					},
				},
			},
			want: want{
				factors: []models.Factor{
					{
						Publisher: "publisher_1",
						Domain:    "1.com",
						Device:    "mobile",
						Factor:    0.1,
						Country:   "il",
					},
					{
						Publisher: "publisher_2",
						Domain:    "2.com",
						Device:    "web",
						Factor:    0.05,
						Country:   "us",
					},
				},
				metadata: []models.MetadataQueue{
					{
						Key:           "mobile:price:factor:v2:publisher_1:1.com:il",
						TransactionID: "uuid_1",
						Value:         []byte("0.10"),
					},
					{
						Key:           "price:factor:v2:publisher_2:2.com:us",
						TransactionID: "uuid_2",
						Value:         []byte("0.05"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			factors, metadata := prepareFactorsData(tt.args.chunk)
			// skipping transaction key because of depending on current time
			for i := range metadata {
				metadata[i].TransactionID = fmt.Sprintf("uuid_%v", i+1)
			}
			assert.Equal(t, tt.want.factors, factors)
			assert.Equal(t, tt.want.metadata, metadata)
		})
	}
}
