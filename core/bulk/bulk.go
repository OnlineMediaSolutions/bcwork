package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/m6yf/bcwork/dto"

	"github.com/m6yf/bcwork/modules/history"
	"github.com/rs/zerolog/log"
)

type Bulker interface {
	BulkInsertDPO(ctx context.Context, requests []dto.DPORuleUpdateRequest) error
	BulkInsertFactors(ctx context.Context, requests []FactorUpdateRequest) error
	BulkInsertGlobalFactors(ctx context.Context, requests []GlobalFactorRequest) error
	BulkDeleteFactor(ctx context.Context, ids []string) error
	BulkDeleteFloor(ctx context.Context, ids []string) error
}

type Adjuster interface {
	AdjustFactors(ctx context.Context, data dto.AdjustRequest) error
	AdjustFloors(ctx context.Context, data dto.AdjustRequest) error
}

type BulkService struct {
	historyModule history.HistoryModule
}

func NewBulkService(historyModule history.HistoryModule) *BulkService {
	return &BulkService{
		historyModule: historyModule,
	}
}

type bulkInsertRequest struct {
	tableName       string
	columns         []string
	valueStrings    []string
	args            []interface{}
	conflictColumns []string
	updateColumns   []string
}

func bulkInsert(ctx context.Context, tx *sql.Tx, req *bulkInsertRequest) error {
	query := prepareBulkInsertQuery(req)

	log.Debug().Msgf("executing bulk insert for %s: %s", req.tableName, query)
	if _, err := tx.ExecContext(ctx, query, req.args...); err != nil {
		return fmt.Errorf("failed to insert into %s in bulk: %w", req.tableName, err)
	}

	return nil
}

func prepareBulkInsertQuery(req *bulkInsertRequest) string {
	const excluded = " = EXCLUDED."

	columnNames := strings.Join(req.columns, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", req.tableName, columnNames, strings.Join(req.valueStrings, ","))

	if req.conflictColumns != nil && req.updateColumns != nil {
		for i := range req.updateColumns {
			req.updateColumns[i] = req.updateColumns[i] + excluded + req.updateColumns[i]
		}
		query += fmt.Sprintf(
			" ON CONFLICT (%s) DO UPDATE SET %s",
			strings.Join(req.conflictColumns, ", "), strings.Join(req.updateColumns, ", "),
		)
	}

	return query
}
