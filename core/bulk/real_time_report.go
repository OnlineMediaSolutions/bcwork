package bulk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m6yf/bcwork/models"
)

func BulkInsertRealTimeReport(ctx context.Context, tx *sql.Tx, report []models.RealTimeReport) error {
	req := prepareBulkInsertRealTimeReport(report)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertRealTimeReport(report []models.RealTimeReport) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.RealTimeReport,
		columns: []string{
			models.RealTimeReportColumns.Time,
			models.RealTimeReportColumns.Publisher,
			models.RealTimeReportColumns.PublisherID,
			models.RealTimeReportColumns.Domain,
			models.RealTimeReportColumns.BidRequests,
			models.RealTimeReportColumns.BidResponses,
			models.RealTimeReportColumns.Device,
			models.RealTimeReportColumns.Country,
			models.RealTimeReportColumns.Revenue,
			models.RealTimeReportColumns.Cost,
			models.RealTimeReportColumns.SoldImpressions,
			models.RealTimeReportColumns.PublisherImpressions,
			models.RealTimeReportColumns.PubFillRate,
			models.RealTimeReportColumns.CPM,
			models.RealTimeReportColumns.RPM,
			models.RealTimeReportColumns.DPRPM,
			models.RealTimeReportColumns.GPP,
			models.RealTimeReportColumns.GP,
			models.RealTimeReportColumns.ConsultantFee,
			models.RealTimeReportColumns.TamFee,
			models.RealTimeReportColumns.TechFee,
			models.RealTimeReportColumns.DemandPartnerFee,
			models.RealTimeReportColumns.DataFee,
		},

		valueStrings: make([]string, 0, len(report)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(report)*multiplier)

	for i, report := range report {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8,
				offset+9, offset+10, offset+11, offset+12, offset+13, offset+14, offset+15, offset+16,
				offset+17, offset+18, offset+19, offset+20, offset+21, offset+22, offset+23),
		)
		req.args = append(req.args,
			report.Time,
			report.Publisher,
			report.PublisherID,
			report.Domain,
			report.BidRequests,
			report.BidResponses,
			report.Device,
			report.Country,
			report.Revenue,
			report.Cost,
			report.SoldImpressions,
			report.PublisherImpressions,
			report.PubFillRate,
			report.CPM,
			report.RPM,
			report.DPRPM,
			report.GP,
			report.GPP,
			report.ConsultantFee,
			report.TamFee,
			report.TechFee,
			report.DemandPartnerFee,
			report.DataFee,
		)
	}

	return req
}
