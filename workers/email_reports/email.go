package email_reports

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"reflect"
	"strings"
)

func SendCustomHTMLEmail(to, bcc, subject string, htmlBody string, report []RealTimeReport) error {
	toRecipients := strings.Split(to, ",")

	csvData, err := createCSVData(report)
	if err != nil {
		return err
	}

	emailReq := modules.EmailRequest{
		To:       toRecipients,
		Bcc:      bcc,
		Subject:  subject,
		Body:     htmlBody,
		IsHTML:   false,
		Attach:   csvData,
		Filename: "real_time_report.csv",
	}

	return modules.SendEmail(emailReq)
}

func createCSVData(report []RealTimeReport) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	reportType := reflect.TypeOf(RealTimeReport{})

	var columnNames []string
	for i := 0; i < reportType.NumField(); i++ {
		field := reportType.Field(i)
		columnNames = append(columnNames, field.Name)
	}

	if err := writer.Write(columnNames); err != nil {
		return nil, err
	}

	for _, record := range report {
		formattedTime := formatDate(record.Time)
		row := []string{
			formattedTime,
			record.PublisherID,
			record.Domain,
			fmt.Sprintf("%.2f", record.BidRequests),
			record.Device,
			record.Country,
			fmt.Sprintf("%.2f", record.Revenue),
			fmt.Sprintf("%.2f", record.Cost),
			fmt.Sprintf("%.2f", record.SoldImpressions),
			fmt.Sprintf("%.2f", record.PublisherImpressions),
			fmt.Sprintf("%.2f", record.PubFillRate),
			fmt.Sprintf("%.2f", record.CPM),
			fmt.Sprintf("%.2f", record.RPM),
			fmt.Sprintf("%.2f", record.DpRPM),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return &buf, nil
}
