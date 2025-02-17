package export

import (
	"context"
	"encoding/json"
	"time"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
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
	IntColumnStyle             = "int"
	FloatColumnStyle           = "float"
	PercentageColumnStyle      = "percentage"
	PercentageFloatColumnStyle = "percentage_float"
	CurrencyColumnStyle        = "currency"
	Currency3ColumnStyle       = "currency3"
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
			IntColumnStyle:             {CustomNumFmt: helpers.GetPointerToString("#,##0")},
			FloatColumnStyle:           {CustomNumFmt: helpers.GetPointerToString("#,##0.00")},
			PercentageColumnStyle:      {CustomNumFmt: helpers.GetPointerToString("#,##0%")},
			PercentageFloatColumnStyle: {CustomNumFmt: helpers.GetPointerToString("#,##0.00%")},
			CurrencyColumnStyle:        {CustomNumFmt: helpers.GetPointerToString("$#,##0.00")},
			Currency3ColumnStyle:       {CustomNumFmt: helpers.GetPointerToString("$#,##0.000")},
			HourColumnStyle:            {CustomNumFmt: helpers.GetPointerToString("yyyy-mm-dd hh:mm:ss")},
			DayColumnStyle:             {CustomNumFmt: helpers.GetPointerToString("yyyy-mm-dd")},
			MonthColumnStyle:           {CustomNumFmt: helpers.GetPointerToString("mmmm, yyyy")},
			YearColumnStyle:            {CustomNumFmt: helpers.GetPointerToString("yyyy")},
		},
		datetimeLayoutMap: map[string]string{
			HourColumnStyle:  "2006-01-02 15:00:00",
			DayColumnStyle:   time.DateOnly,
			MonthColumnStyle: "January, 2006",
			YearColumnStyle:  "2006",
		},
	}
}
