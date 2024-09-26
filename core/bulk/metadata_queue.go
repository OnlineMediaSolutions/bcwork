package bulk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m6yf/bcwork/models"
)

func bulkInsertMetaDataQueue(ctx context.Context, tx *sql.Tx, metaDataQueue []models.MetadataQueue) error {
	req := prepareBulkInsertMetaDataQueueRequest(metaDataQueue)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertMetaDataQueueRequest(metaDataQueue []models.MetadataQueue) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.MetadataQueue,
		columns: []string{
			models.MetadataQueueColumns.Key,
			models.MetadataQueueColumns.TransactionID,
			models.MetadataQueueColumns.Value,
			models.MetadataQueueColumns.CommitedInstances,
			models.MetadataQueueColumns.CreatedAt,
			models.MetadataQueueColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(metaDataQueue)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(metaDataQueue)*multiplier)

	for i, metaData := range metaDataQueue {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6),
		)
		req.args = append(req.args,
			metaData.Key,
			metaData.TransactionID,
			metaData.Value,
			metaData.CommitedInstances,
			currentTime,
			currentTime,
		)
	}

	return req
}
