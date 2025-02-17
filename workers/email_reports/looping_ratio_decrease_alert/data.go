package looping_ratio_decrease_alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/m6yf/bcwork/workers/email_reports"
	"github.com/rs/zerolog/log"
	"strings"
	"text/template"
	"time"
)

const (
	PERCENTAGE             = 0.4
	PubImpsThreshold int64 = 3000
)

type LRReport struct {
	Date                 string  `json:"date"`
	DataStamp            int64   `json:"DateStamp"`
	Publisher            string  `json:"Publisher"`
	Domain               string  `json:"Domain"`
	PaymentType          string  `json:"PaymentType"`
	AM                   string  `json:"AM"`
	PubImps              int64   `json:"PubImps"`
	LoopingRatio         float64 `json:"nbLR"`
	Ratio                float64 `json:"nbRatio"`
	CPM                  float64 `json:"nbCpm"`
	Cost                 float64 `json:"Cost"`
	RPM                  float64 `json:"nbRpm"`
	DpRPM                float64 `json:"nbDpRpm"`
	Revenue              float64 `json:"Revenue"`
	GP                   float64 `json:"nbGp"`
	GPP                  float64 `json:"nbGpp"`
	PublisherBidRequests int64   `json:"PublisherBidRequests"`
}

type LReport struct {
	Data struct {
		Result []LRReport `json:"result"`
	} `json:"data"`
}

type LoopingRationDecreaseReport struct{}

func aggregate(reports []AggregatedReport) map[string][]AggregatedReport {
	aggregated := make(map[string][]AggregatedReport)

	for _, r := range reports {
		key := fmt.Sprintf("%s|%s|%s|%s", r.AM, r.Domain, r.Publisher, r.PaymentType)
		aggregated[key] = append(aggregated[key], AggregatedReport{
			AM:                   r.AM,
			Domain:               r.Domain,
			Publisher:            r.Publisher,
			PaymentType:          r.PaymentType,
			Date:                 r.Date,
			PubImps:              r.PubImps,
			DataStamp:            r.DataStamp,
			RPM:                  r.RPM,
			Ratio:                r.Ratio,
			LoopingRatio:         r.LoopingRatio,
			CPM:                  r.CPM,
			Cost:                 r.Cost,
			DpRPM:                r.DpRPM,
			Revenue:              r.Revenue,
			GP:                   r.GP,
			GPP:                  r.GPP,
			PublisherBidRequests: r.PublisherBidRequests,
		})
	}

	return aggregated
}

func computeAverage(aggregated map[string][]AggregatedReport, worker *Worker) map[string][]AlertsEmails {
	now := time.Now().In(email_reports.Location)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := startOfDay.AddDate(0, 0, -1).Unix() / 100
	startOfLastWeek := startOfDay.AddDate(0, 0, -7).Unix() / 100
	today := startOfDay.Unix() / 100

	amDomainData := make(map[string][]AggregatedReport)

	for key, aggs := range aggregated {
		for _, agg := range aggs {
			if agg.DataStamp == yesterday || agg.DataStamp == startOfLastWeek || agg.DataStamp == today {
				amDomainData[key] = append(amDomainData[key], agg)
			}
		}
	}

	alerts := make(map[string]AggregatedReport)
	repo := AlertsEmails{}

	reports := compareResults(amDomainData, yesterday, startOfLastWeek, today, alerts, repo, worker)

	avgDataMap := make(map[string][]AlertsEmails)
	for _, repo := range reports {
		avgDataMap[repo.Email] = append(avgDataMap[repo.Email], repo)
	}

	return avgDataMap
}

func compareResults(amDomainData map[string][]AggregatedReport, yesterday int64, startOfLastWeek int64, today int64, alerts map[string]AggregatedReport, repo AlertsEmails, worker *Worker) []AlertsEmails {
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
				totalYesterday += report.LoopingRatio
				countYesterday++
			} else if report.DataStamp == startOfLastWeek {
				totalLastWeek += report.LoopingRatio
				countLastWeek++
			} else if report.DataStamp == today {
				totalToday += report.LoopingRatio
				countToday++
			}
		}

		if countToday == 0 || (countYesterday == 0 && countLastWeek == 0) {
			continue
		}

		if totalToday < PERCENTAGE*(totalYesterday+totalLastWeek) {
			latestReport := reports[len(reports)-1]
			alerts[key] = latestReport
			emailKey := strings.Split(key, "|")
			repo = AlertsEmails{
				Email:        worker.UserData[emailKey[0]],
				AM:           key,
				FirstReport:  latestReport,
				SecondReport: reports,
			}

			emailReports = append(emailReports, repo)
		}
	}

	return emailReports
}

func getReport() ([]AggregatedReport, error) {
	compassClient := compass.NewCompass()

	requestData := getRequestData()

	data, err := json.Marshal(requestData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling request data, %w", err)
	}

	reportData, err := compassClient.Request("/report-dashboard/report-new-bidder", "POST", data, true)

	if err != nil {
		return nil, fmt.Errorf("error getting report data,%w", err)
	}

	var report LReport
	err = json.Unmarshal(reportData, &report)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling report data: %w", err)
	}

	aggregatedReports := prepareLRReport(report)

	return aggregatedReports, nil
}

func prepareLRReport(report LReport) []AggregatedReport {
	aggregatedReports := make([]AggregatedReport, len(report.Data.Result))
	formatter := &helpers.FormatValues{}
	for i, r := range report.Data.Result {
		if r.PubImps >= PubImpsThreshold {
			aggregatedReports[i] = AggregatedReport{
				Date:                 r.Date,
				DataStamp:            r.DataStamp,
				Publisher:            r.Publisher,
				Domain:               r.Domain,
				PaymentType:          r.PaymentType,
				AM:                   r.AM,
				PubImps:              formatter.PubImps(int(r.PubImps)),
				PublisherBidRequests: formatter.BidRequests(float64(r.PublisherBidRequests)),
				LoopingRatio:         helpers.RoundFloat(r.LoopingRatio),
				Ratio:                helpers.RoundFloat(r.Ratio),
				CPM:                  helpers.RoundFloat(r.CPM),
				Cost:                 helpers.RoundFloat(r.Cost),
				RPM:                  helpers.RoundFloat(r.RPM),
				DpRPM:                helpers.RoundFloat(r.DpRPM),
				Revenue:              helpers.RoundFloat(r.Revenue),
				GP:                   helpers.RoundFloat(r.GP),
				GPP:                  helpers.RoundFloat(r.GPP),
			}
		}
	}

	return aggregatedReports
}

func getRequestData() RequestData {
	currentTime := time.Now().In(email_reports.Location)
	startOfLast7Days := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -7)

	endOfLast7Days := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).AddDate(0, 0, 0)

	startOfLast7DaysStr := startOfLast7Days.Format(time.DateTime)
	endOfLast7DaysStr := endOfLast7Days.Format(time.DateTime)

	requestData := RequestData{
		Data: RequestDetails{
			Date: Date{
				Range: []string{
					startOfLast7DaysStr,
					endOfLast7DaysStr,
				},
				Interval: "day",
			},
			Dimensions: []string{
				"Publisher",
				"Domain",
				"AM",
				"PaymentType",
			},
			Metrics: []string{
				"PublisherBidRequests",
				"nbLR",
				"nbRatio",
				"nbRpm",
				"nbDpRpm",
				"Revenue",
				"nbGp",
				"nbGpp",
				"PubImps",
				"nbCpm",
				"Cost",
			},
		},
	}

	return requestData
}

func prepareAndSendEmail(reportData map[string][]AlertsEmails, worker *Worker) error {
	if len(reportData) > 0 {
		currentTime := time.Now().In(email_reports.Location)
		today := currentTime.Format(time.DateOnly)
		for email, alerts := range reportData {
			subject := fmt.Sprintf("Looping ratio decrease alert for %s", today)
			message := fmt.Sprintf("Dear %s,\n\nLooping ratio decrease alert for %s.\n\nPlease review the details below.", email, today)
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
        <title>Looping Ratio Report</title>
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
        <h4>Looping ratio decrease alert for {{.FirstReport.Domain}}</h4>
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
        <h4>Looping ratio for Yesterday, 7 days ago and Today</h4>
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
                <td class="{{if eq $index 2}}no-changes{{end}}">{{.LoopingRatio}}</td>
                <td>{{.Ratio}}</td>
                <td>{{.CPM}}</td>
                <td>{{.Cost}}</td>
                <td>{{.RPM}}</td>
                <td>{{.DpRPM}}</td>
                <td>{{.Revenue}}</td>
                <td>{{.GP}}</td>
                <td>{{.GPP}}</td>
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
		FirstReport  AggregatedReport   `json:"FirstReport"`
		SecondReport []AggregatedReport `json:"SecondReport"`
	}

	for _, reportGroup := range report {
		reportsList = append(reportsList, struct {
			FirstReport  AggregatedReport   `json:"FirstReport"`
			SecondReport []AggregatedReport `json:"SecondReport"`
		}{

			FirstReport:  reportGroup.FirstReport,
			SecondReport: reportGroup.SecondReport,
		})
	}

	data := struct {
		Body    string
		Reports []struct {
			FirstReport  AggregatedReport   `json:"FirstReport"`
			SecondReport []AggregatedReport `json:"SecondReport"`
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
