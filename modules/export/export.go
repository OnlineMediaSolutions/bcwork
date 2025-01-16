package export

import (
	"context"
	"encoding/json"
	"time"

	"github.com/m6yf/bcwork/dto"
	"github.com/xuri/excelize/v2"
)

const (
	// column styles (formats)
	// dates
	HourColumnStyle    = "hour"    // 2025-01-14 08:00:00
	DayColumnStyle     = "day"     // 2025-01-14
	WeekColumnStyle    = "week"    // 02, 2025
	MonthColumnStyle   = "month"   // January, 2025
	QuarterColumnStyle = "quarter" // Q1, 2025
	YearColumnStyle    = "year"    // 2025
	// numbers
	IntColumnStyle        = "int"
	FloatColumnStyle      = "float"
	PercentageColumnStyle = "percentage"
	CurrencyColumnStyle   = "currency"
	Currency3ColumnStyle  = "currency3"
	// booleans - based on true/false we show ‘string’ which was provided in "boolean_replacement" field
	BooleanColumnStyle = "boolean"
)

type Exporter interface {
	ExportCSV(ctx context.Context, srcs []json.RawMessage) ([]byte, error)
	ExportXLSX(ctx context.Context, req *dto.DownloadRequest) ([]byte, error)
}

type ExportModule struct {
	stylesMap         map[string]*excelize.Style
	datetimeLayoutMap map[string]string
}

func NewExportModule() *ExportModule {
	return &ExportModule{
		stylesMap: map[string]*excelize.Style{
			IntColumnStyle:        {CustomNumFmt: func() *string { s := "#,##0"; return &s }()},
			FloatColumnStyle:      {CustomNumFmt: func() *string { s := "#,##0.00"; return &s }()},
			PercentageColumnStyle: {CustomNumFmt: func() *string { s := "#,##0.00%"; return &s }()},
			CurrencyColumnStyle:   {CustomNumFmt: func() *string { s := "$#,##0.00"; return &s }()},
			Currency3ColumnStyle:  {CustomNumFmt: func() *string { s := "$#,##0.000"; return &s }()},
			HourColumnStyle:       {CustomNumFmt: func() *string { s := "yyyy-mm-dd hh:mm:ss"; return &s }()},
			DayColumnStyle:        {CustomNumFmt: func() *string { s := "yyyy-mm-dd"; return &s }()},
			// WeekColumnStyle:       {CustomNumFmt: func() *string { s := ""; return &s }()},
			MonthColumnStyle: {CustomNumFmt: func() *string { s := "mmmm, yyyy"; return &s }()},
			// QuarterColumnStyle:    {CustomNumFmt: func() *string { s := ""; return &s }()},
			YearColumnStyle: {CustomNumFmt: func() *string { s := "yyyy"; return &s }()},
		},
		datetimeLayoutMap: map[string]string{
			HourColumnStyle:  "2006-01-02 15:00:00",
			DayColumnStyle:   time.DateOnly,
			MonthColumnStyle: "January, 2006",
			YearColumnStyle:  "2006",
		},
	}
}
