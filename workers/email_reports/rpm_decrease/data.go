package rpm_decrease

import (
	"bytes"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/workers/email_reports"
	"github.com/rs/zerolog/log"
	"strings"
	"text/template"
	"time"
)

func compareResults(amDomainData map[string][]email_reports.AggregatedReport, percentage float64, userData map[string]string) map[string][]AlertsEmails {
	yesterday, startOfLastWeek, today := email_reports.GetTimestampsForDateRange()

	var emailReports []AlertsEmails
	for key, reports := range amDomainData {
		totalYesterday := 0.0
		totalLastWeek := 0.0
		countYesterday := 0
		countLastWeek := 0
		totalToday := 0.0
		countToday := 0

		for _, report := range reports {
			if report.DataStamp == yesterday {
				totalYesterday += report.RPM
				countYesterday++
			} else if report.DataStamp == startOfLastWeek {
				totalLastWeek += report.RPM
				countLastWeek++
			} else if report.DataStamp == today {
				totalToday += report.RPM
				countToday++
			}
		}

		if countToday == 0 || (countYesterday == 0 && countLastWeek == 0) {
			continue
		}

		if (totalToday < percentage*totalYesterday) && (totalToday < percentage*totalLastWeek) {
			latestReport := reports[len(reports)-1]
			emailKey := strings.Split(key, "|")
			emailReports = append(emailReports, AlertsEmails{
				Email:        userData[emailKey[0]],
				AM:           key,
				FirstReport:  latestReport,
				SecondReport: reports,
			})
		}
	}

	avgDataMap := make(map[string][]AlertsEmails)
	for _, report := range emailReports {
		avgDataMap[report.Email] = append(avgDataMap[report.Email], report)
	}

	return avgDataMap
}

func prepareAndSendEmail(reportData map[string][]AlertsEmails, worker *Worker) error {
	if len(reportData) > 0 {
		today := time.Now().Format(time.DateOnly)

		for email, alerts := range reportData {
			subject := fmt.Sprintf("RPM decrease alert for %s", today)
			message := fmt.Sprintf("Dear %s,\n\nRPM decrease alert for %s.\n\nPlease review the details below.", email, today)

			err := sendCustomHTMLEmail(email, worker.BCC, subject, message, alerts)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to send email to %s", email)
			}
		}
	}

	return nil
}

func sendCustomHTMLEmail(to, bcc, subject string, body string, report []AlertsEmails) error {
	toRecipients := strings.Split(to, ",")
	bccString := strings.Split(bcc, ",")
	emailData := EmailData{
		Body:   body,
		Report: report,
	}

	htmlBody, err := generateHTMLTableWithTemplate(emailData.Report, emailData.Body)
	if err != nil {
		return err
	}

	emailReq := modules.EmailRequest{
		To:      toRecipients,
		Bcc:     bccString,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return modules.SendEmail(emailReq)
}

func generateHTMLTableWithTemplate(report []AlertsEmails, body string) (string, error) {
	const tpl = `
<html>
   <head>
       <title>RPM Report</title>
       <style>
           table { width: 100%; border-collapse: collapse; }
           th, td { border: 1px solid black; padding: 8px; text-align: left; }
           th { background-color: #f2f2f2; }
           .no-changes { color: red; font-weight: bold; }
       </style>
   </head>
   <body>
       <h3>{{.Body}}</h3>
       {{range .Reports}}
        <!-- First Report Section -->
       <h4>RPM decrease alert for {{.FirstReport.Domain}}</h4>
       <table>
           <tr>
              <th>Date</th>
              <th>Account Manager</th>
              <th>Publisher</th>
              <th>Domain</th>
           </tr>
           <tr>
               <td>{{.FirstReport.Date}}</td>
               <td>{{.FirstReport.AM}}</td>
               <td>{{.FirstReport.Publisher}}</td>
               <td>{{.FirstReport.Domain}}</td>
           </tr>
       </table>
       <!-- Second Report Section -->
       <h4>RPM for Yesterday, 7 days ago and Today</h4>
       <table>
           <tr>
               <th>Date</th>
               <th>Publisher</th>
               <th>Domain</th>
               <th>Payment Type</th>
               <th>PublisherBidRequests</th>
               <th>Imps</th>
               <th>LR</th>
               <th>Ratio</th>
               <th>CPM</th>
               <th>Cost</th>
               <th>RPM</th>
               <th>DP RPM</th>
               <th>Revenue</th>
               <th>GP$</th>
               <th>GP%</th>
           </tr>
           {{range $index, $report := .SecondReport}}
           <tr>
               <td>{{.Date}}</td>
               <td>{{.Publisher}}</td>
               <td>{{.Domain}}</td>
               <td>{{.PaymentType}}</td>
               <td>{{.PublisherBidRequests}}</td>
               <td>{{.PubImps}}</td>
               <td>{{.LoopingRatio}}</td>
               <td>{{.Ratio}}</td>
               <td>${{.CPM}}</td>
               <td>${{.Cost}}</td>
               <td class="{{if eq $index 2}}no-changes{{end}}">${{.RPM}}</td>
               <td>${{.DpRPM}}</td>
               <td>${{.Revenue}}</td>
               <td>${{.GP}}</td>
               <td>{{.GPP}}%</td>
           </tr>
           {{end}}
       </table>
       <br>
<hr>
       {{end}}
   </body>
</html>
`

	var reportsList []struct {
		FirstReport  email_reports.AggregatedReport   `json:"FirstReport"`
		SecondReport []email_reports.AggregatedReport `json:"SecondReport"`
	}

	for _, reportGroup := range report {
		reportsList = append(reportsList, struct {
			FirstReport  email_reports.AggregatedReport   `json:"FirstReport"`
			SecondReport []email_reports.AggregatedReport `json:"SecondReport"`
		}{

			FirstReport:  reportGroup.FirstReport,
			SecondReport: reportGroup.SecondReport,
		})
	}

	data := struct {
		Body    string
		Reports []struct {
			FirstReport  email_reports.AggregatedReport   `json:"FirstReport"`
			SecondReport []email_reports.AggregatedReport `json:"SecondReport"`
		}
	}{
		Body:    body,
		Reports: reportsList,
	}

	t, err := template.New("emailTemplate").Parse(tpl)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, data); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}
