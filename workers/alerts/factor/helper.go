package factor

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"time"
)

func ConvertReportsToCSV(reports []Report) (string, error) {
	var csvBuffer bytes.Buffer
	writer := csv.NewWriter(&csvBuffer)

	err := writer.Write(Header)
	if err != nil {
		return "Error writing header to  csv factor logs", err
	}

	s, err := generateCSVValues(reports, writer)
	if err != nil {
		return s, err
	}

	writer.Flush()

	if err = writer.Error(); err != nil {
		return "Error in writer factor logs", err
	}

	return csvBuffer.String(), nil
}

func generateCSVValues(reports []Report, writer *csv.Writer) (string, error) {
	for _, report := range reports {
		row := []string{
			report.Time,
			report.EvalTime,
			fmt.Sprintf("%d", report.PubImps),
			fmt.Sprintf("%d", report.SoldImps),
			fmt.Sprintf("%.2f", report.Cost),
			fmt.Sprintf("%.2f", report.Revenue),
			fmt.Sprintf("%.2f", report.GP),
			fmt.Sprintf("%.2f", report.GPP),
			report.Publisher,
			report.Domain,
			report.Country,
			report.Device,
			fmt.Sprintf("%.2f", report.OldFactor),
			fmt.Sprintf("%.2f", report.NewFactor),
			fmt.Sprintf("%d", report.ResponseStatus),
			fmt.Sprintf("%.2f", report.Increase),
		}
		err := writer.Write(row)
		if err != nil {
			return "Error in writing function factor logs", err
		}
	}
	return "", nil
}

func getDataFromDB(ctx context.Context, db *sqlx.DB) (string, error) {

	timeDuration := 30
	records := make(models.PriceFactorLogSlice, 0)
	timeString := getTimeValue(timeDuration)

	sql := ` SELECT * FROM public.price_factor_log
	         WHERE eval_time >= TO_TIMESTAMP('%s', 'YYYY-MM-DD HH24:MI:SS') 
	         AND response_status!= 200  AND response_status!= 0;`

	query := fmt.Sprintf(sql, timeString)

	fmt.Println(`Factor logs Query`, query)
	fmt.Println(`Executing task for factor logs for date`, timeString)

	err := queries.Raw(query).Bind(ctx, db, &records)
	if err != nil {
		return "Error in select query factor logs", err
	}

	data, err := json.Marshal(records)

	if err != nil {
		return "Error in marshalling data factor logs", err
	}

	return string(data), nil
}

func getTimeValue(timeDuration int) string {
	evalTime := time.Now().UTC().Add(-time.Duration(timeDuration) * time.Minute).Truncate(time.Duration(timeDuration) * time.Minute)
	timeString := evalTime.Format("2006-01-02 15:04:05")
	return timeString
}
