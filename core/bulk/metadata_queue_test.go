package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func Test_prepareBulkInsertMetaDataQueueRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		metaDataQueue []models.MetadataQueue
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				metaDataQueue: []models.MetadataQueue{
					{
						Key:               "key_1",
						TransactionID:     "1",
						Value:             []byte("1"),
						CommitedInstances: 1,
					},
					{
						Key:               "key_2",
						TransactionID:     "2",
						Value:             []byte("2"),
						CommitedInstances: 2,
					},
					{
						Key:               "key_3",
						TransactionID:     "3",
						Value:             []byte("3"),
						CommitedInstances: 3,
					},
				},
			},
			want: &bulkInsertRequest{
				tableName: models.TableNames.MetadataQueue,
				columns: []string{
					models.MetadataQueueColumns.Key,
					models.MetadataQueueColumns.TransactionID,
					models.MetadataQueueColumns.Value,
					models.MetadataQueueColumns.CommitedInstances,
					models.MetadataQueueColumns.CreatedAt,
					models.MetadataQueueColumns.UpdatedAt,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5, $6)",
					"($7, $8, $9, $10, $11, $12)",
					"($13, $14, $15, $16, $17, $18)",
				},
				args: []interface{}{
					"key_1", "1", func() types.JSON { return []byte("1") }(), int64(1), currentTime, currentTime,
					"key_2", "2", func() types.JSON { return []byte("2") }(), int64(2), currentTime, currentTime,
					"key_3", "3", func() types.JSON { return []byte("3") }(), int64(3), currentTime, currentTime,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertMetaDataQueueRequest(tt.args.metaDataQueue)
			assert.Equal(t, tt.want, got)
		})
	}
}
