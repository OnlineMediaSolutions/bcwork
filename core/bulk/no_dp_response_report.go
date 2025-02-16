package bulk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
)

func BulkInsertNoDPResponseReport(ctx context.Context, tx *sql.Tx, report []*dto.NoDPResponseReport) error {
	req := prepareBulkInsertNoDPResponseReport(report)
	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertNoDPResponseReport(report []*dto.NoDPResponseReport) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.NoDPResponseReport,
		columns: []string{
			models.NoDPResponseReportColumns.Time,
			models.NoDPResponseReportColumns.DemandPartnerID,
			models.NoDPResponseReportColumns.PublisherID,
			models.NoDPResponseReportColumns.Domain,
			models.NoDPResponseReportColumns.BidRequests,
		},
		valueStrings: make([]string, 0, len(report)),
		conflictColumns: []string{
			models.NoDPResponseReportColumns.Time,
			models.NoDPResponseReportColumns.DemandPartnerID,
			models.NoDPResponseReportColumns.PublisherID,
			models.NoDPResponseReportColumns.Domain,
		},
		updateColumns: []string{
			models.NoDPResponseReportColumns.BidRequests,
		},
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(report)*multiplier)

	for i, report := range report {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5),
		)
		req.args = append(req.args,
			report.Time,
			report.DPID,
			report.PubID,
			report.Domain,
			report.BidRequests,
		)
	}

	return req
}
