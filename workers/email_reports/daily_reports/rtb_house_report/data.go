package rtb_house_report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/helpers"
	"strconv"
	"strings"
	"time"
)

func PrepareRequestData() RequestData {
	requestData := RequestData{
		Data: RequestDetails{
			Date: Date{
				Range: []string{
					time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05"),
					time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
				},
				Interval: "day",
			},
			Dimensions: []string{
				"DemandPartner",
			},
			Filters: Filters{
				DemandPartner: []string{"rtb house"},
			},
			Metrics: []string{
				"SoldImps",
				"Revenue",
			},
		},
	}

	return requestData
}

func (worker *Worker) prepareEmail(report map[string]interface{}, err error, emailCreds EmailCreds) {
	body, subject, reportName := GenerateReportDetails(worker)
	var reports []Report
	for _, r := range report {
		if report, ok := r.(Report); ok {
			reports = append(reports, report)
		}
	}

	//if data, ok := report["data"].(map[string]interface{}); ok {
	//	if result, ok := data["result"].([]interface{}); ok {
	//		for _, r := range result {
	//			if reportMap, ok := r.(map[string]interface{}); ok {
	//				revenue, _ := strconv.ParseFloat(reportMap["Revenue"].(string), 64)
	//				soldImps, _ := strconv.Atoi(reportMap["SoldImps"].(string))
	//
	//				dateStamp, _ := strconv.ParseFloat(reportMap["DateStamp"].(string), 64)
	//				time := strconv.FormatFloat(dateStamp, 'f', -1, 64) // Convert float64 to string
	//
	//				report := Report{
	//					DemandPartner: reportMap["DemandPartner"].(string),
	//					Revenue:       revenue,
	//					SoldImps:      soldImps,
	//					Time:          time, // Assign the converted DateStamp to Time
	//				}
	//				reports = append(reports, report)
	//			}
	//		}
	//	}
	//}

	err = SendCustomHTMLEmail(
		emailCreds.TO,
		emailCreds.BCC,
		subject,
		body,
		reports,
		reportName)
}

func GenerateReportDetails(worker *Worker) (string, string, string) {
	body := fmt.Sprintf("Rtb House  %s - %s\n",
		helpers.FormatDate(worker.Start.Format(time.RFC3339)),
		helpers.FormatDate(worker.End.Format(time.RFC3339)))
	subject := fmt.Sprintf("Rtb House reports %s", helpers.FormatDate(worker.End.Format(time.RFC3339)))
	reportName := fmt.Sprintf("Rtb House report_%s.csv", helpers.FormatDate(worker.End.Format(time.RFC3339)))

	return body, subject, reportName
}

func SendCustomHTMLEmail(to, bcc, subject string, htmlBody string, report []Report, reportName string) error {
	toRecipients := strings.Split(to, ",")
	bccStr := strings.Split(bcc, ",")

	csvData, err := createCSVData(report)
	if err != nil {
		return fmt.Errorf("failed to create csv data, %w", err)
	}

	emailReq := modules.EmailRequest{
		To:       toRecipients,
		Bcc:      bccStr,
		Subject:  subject,
		Body:     htmlBody,
		IsHTML:   false,
		Attach:   csvData,
		Filename: reportName,
	}

	return modules.SendEmail(emailReq)
}

func createCSVData(report []Report) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(ColumnNames); err != nil {
		return nil, fmt.Errorf("failed to create csv data, %w", err)
	}
	//formatter := &helpers.FormatValues{}

	for _, record := range report {
		formattedTime := helpers.FormatDate(record.Time)
		row := []string{
			formattedTime,
			strconv.Itoa(record.SoldImps),
			strconv.FormatFloat(record.Revenue, 'f', 6, 64),
			record.DemandPartner,
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
