package export

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/xuri/excelize/v2"
)

func (e *ExportModule) ExportXLSX(ctx context.Context, req *dto.DownloadRequest) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	defaultSheetName := time.Now().Format(time.DateOnly)
	idx, err := f.NewSheet(defaultSheetName)
	if err != nil {
		return nil, fmt.Errorf("cannot create new sheet: %w", err)
	}

	f.SetActiveSheet(idx)

	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("cannot delete default sheet: %w", err)
	}

	stylesIDMap, err := e.prepareStyles(f)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare styles: %w", err)
	}

	// prepare header and columns styles (formats)
	for idx, columm := range req.Columns {
		styleID := stylesIDMap[columm.Style]

		columnNumber := idx + 1
		columnName, err := excelize.ColumnNumberToName(columnNumber)
		if err != nil {
			return nil, fmt.Errorf("cannot translate column number [%v] to column name: %w", columnNumber, err)
		}

		err = f.SetColStyle(defaultSheetName, columnName, styleID)
		if err != nil {
			return nil, fmt.Errorf("cannot set column [%v] style: %w", columnName, err)
		}

		cellName, err := excelize.CoordinatesToCellName(columnNumber, 1)
		if err != nil {
			return nil, fmt.Errorf("cannot translate coordinates to cell name: %w", err)
		}

		displayName := columm.DisplayName
		if displayName == "" {
			displayName = columm.Name
		}

		err = f.SetCellValue(defaultSheetName, cellName, displayName)
		if err != nil {
			return nil, fmt.Errorf("cannot set cell [%v] value to %v: %w", cellName, displayName, err)
		}
	}

	// writing subsequent rows
	for i := 0; i < len(req.Data); i++ {
		rowNumber := i + 2

		var temp map[string]interface{}
		err := json.Unmarshal(req.Data[i], &temp)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal %v to map[string]interface{}: %w", req.Data[i], err)
		}

		for j, column := range req.Columns {
			columnNumber := j + 1

			cellName, err := excelize.CoordinatesToCellName(columnNumber, rowNumber)
			if err != nil {
				return nil, fmt.Errorf("cannot translate coordinates to cell name: %w", err)
			}

			cellValue := e.getCellValue(temp, column)
			err = f.SetCellValue(defaultSheetName, cellName, cellValue)
			if err != nil {
				return nil, fmt.Errorf("cannot set cell [%v] value to %v: %w", cellName, cellValue, err)
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *ExportModule) prepareStyles(f *excelize.File) (map[string]int, error) {
	stylesIDMap := make(map[string]int)

	for styleName, style := range e.stylesMap {
		styleID, err := f.NewStyle(style)
		if err != nil {
			return nil, err
		}

		stylesIDMap[styleName] = styleID
	}

	return stylesIDMap, nil
}

func (e *ExportModule) getCellValue(data map[string]interface{}, column *dto.Column) interface{} {
	value := data[column.Name]

	switch column.Style {
	case BooleanColumnStyle:
		boolValue, ok := value.(bool)
		if ok {
			return column.GetBooleanReplacementValue(boolValue)
		}
	case HourColumnStyle, DayColumnStyle, WeekColumnStyle, MonthColumnStyle, QuarterColumnStyle, YearColumnStyle:
		return e.getFormattedDatetimeString(value, column.Style)
	case IntColumnStyle, FloatColumnStyle, PercentageColumnStyle, PercentageFloatColumnStyle, CurrencyColumnStyle, Currency3ColumnStyle:
		switch number := value.(type) {
		case int:
			return float64(number) * column.GetMultiply()
		case float64:
			return number * column.GetMultiply()
		}
	}

	return value
}

func (e *ExportModule) getFormattedDatetimeString(value interface{}, style string) string {
	dateStr, ok := value.(string)
	if ok {
		t, err := time.Parse(time.RFC3339Nano, dateStr)
		if err != nil {
			return dateStr
		}

		switch style {
		case WeekColumnStyle:
			year, week := t.ISOWeek()
			if week < 10 {
				return fmt.Sprintf("0%v %v", week, year)
			}
			return fmt.Sprintf("%v %v", week, year)
		case QuarterColumnStyle:
			return fmt.Sprintf("%v %v", helpers.GetQuarter(t.Month()), t.Year())
		default:
			return t.Format(e.datetimeLayoutMap[style])
		}
	}

	return dateStr
}
