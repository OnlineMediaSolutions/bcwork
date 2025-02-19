package rpm_decrease

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/m6yf/bcwork/workers/email_reports"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	HourlyPercentage          = 0.4
	RPMPubImpsThreshold int64 = 50
)

type RPMReport struct {
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

type RReport struct {
	Data struct {
		Result []RPMReport `json:"result"`
	} `json:"data"`
}

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
	amDomainData := make(map[string][]AggregatedReport)

	lastCompleteHour := time.Now().In(email_reports.Location).Truncate(time.Hour)
	yesterdayHour := lastCompleteHour.Add(-24 * time.Hour)
	sevenDaysAgoHour := lastCompleteHour.Add(-7 * 24 * time.Hour)

	lastCompleteHourUnix := lastCompleteHour.Unix() / 100
	yesterdayHourUnix := yesterdayHour.Unix() / 100
	sevenDaysAgoHourUnix := sevenDaysAgoHour.Unix() / 100

	for key, aggs := range aggregated {
		for _, agg := range aggs {
			if agg.DataStamp == lastCompleteHourUnix || agg.DataStamp == yesterdayHourUnix || agg.DataStamp == sevenDaysAgoHourUnix {
				amDomainData[key] = append(amDomainData[key], agg)
			}
		}
	}

	repo := AlertsEmails{}
	emailReports := compareResults(amDomainData, repo, worker)

	avgDataMap := make(map[string][]AlertsEmails)
	for _, repo := range emailReports {
		avgDataMap[repo.Email] = append(avgDataMap[repo.Email], repo)
	}

	return avgDataMap
}

func compareResults(amDomainData map[string][]AggregatedReport, repo AlertsEmails, worker *Worker) []AlertsEmails {
	var emailReports []AlertsEmails

	for key, reports := range amDomainData {
		if len(reports) < 3 {
			continue
		}

		sort.Slice(reports, func(i, j int) bool {
			return reports[i].DataStamp < reports[j].DataStamp
		})

		if reports[2].RPM < HourlyPercentage*(reports[1].RPM) {
			latestReport := reports[len(reports)-1]
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

func getReport() ([]AggregatedReport, error) {
	compassClient := compass.NewCompass()

	requestData := getRequestData()
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data for today and yesterday: %w", err)
	}

	reportData, err := compassClient.Request("/report-dashboard/report-new-bidder", "POST", data, true)
	if err != nil {
		return nil, fmt.Errorf("error getting report data for today and yesterday: %w", err)
	}

	var report RReport
	err = json.Unmarshal(reportData, &report)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling report data for today and yesterday: %w", err)
	}

	formatter := &helpers.FormatValues{}

	aggregatedReports := prepareReport(report, formatter)

	aggregatedReportsSevenDays, err := get7DaysAgoData(err, compassClient, formatter)
	if err != nil {
		return nil, err
	}

	aggregatedReportsMap := append(aggregatedReports, aggregatedReportsSevenDays...)

	return aggregatedReportsMap, nil
}

func prepareReport(report RReport, formatter *helpers.FormatValues) []AggregatedReport {
	aggregatedReports := make([]AggregatedReport, len(report.Data.Result))
	for i, r := range report.Data.Result {
		if r.PubImps >= RPMPubImpsThreshold {
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

func get7DaysAgoData(err error, compassClient *compass.Compass, formatter *helpers.FormatValues) ([]AggregatedReport, error) {
	requestDataSevenDaysAgo := GetRequestDataSevenDaysAgo()
	dataSevenDaysAgo, err := json.Marshal(requestDataSevenDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data for 7 days ago: %w", err)
	}

	reportDataSevenDaysAgo, err := compassClient.Request("/report-dashboard/report-new-bidder", "POST", dataSevenDaysAgo, true)
	if err != nil {
		return nil, fmt.Errorf("error getting report data for 7 days ago: %w", err)
	}

	var reportSevenDays RReport
	err = json.Unmarshal(reportDataSevenDaysAgo, &reportSevenDays)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling report data for 7 days ago: %w", err)
	}

	aggregatedReportsSevenDays := make([]AggregatedReport, len(reportSevenDays.Data.Result))

	for i, r := range reportSevenDays.Data.Result {
		aggregatedReportsSevenDays[i] = AggregatedReport{
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

	return aggregatedReportsSevenDays, nil
}

func getRequestData() email_reports.RequestData {
	yesterday := time.Now().In(email_reports.Location).Truncate(time.Hour).Add(-24 * time.Hour)
	today := time.Now().In(email_reports.Location).Truncate(time.Hour)

	yesterdayStr := yesterday.Format(time.DateTime)
	todayStr := today.Format(time.DateTime)

	requestData := email_reports.RequestData{
		Data: email_reports.RequestDetails{
			Date: email_reports.Date{
				Range: []string{
					yesterdayStr,
					todayStr,
				},
				Interval: "hour",
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
			},
		},
	}

	return requestData
}

func GetRequestDataSevenDaysAgo() email_reports.RequestData {
	sevenDaysAgoHour := time.Now().In(email_reports.Location).Truncate(time.Hour).Add(-7 * 24 * time.Hour).Format(time.DateTime)
	requestData := email_reports.RequestData{
		Data: email_reports.RequestDetails{
			Date: email_reports.Date{
				Range: []string{
					sevenDaysAgoHour,
					sevenDaysAgoHour,
				},
				Interval: "hour",
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
			},
		},
	}

	return requestData
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
                <td>{{.CPM}}</td>
                <td>{{.Cost}}</td>
                <td class="{{if eq $index 2}}no-changes{{end}}">{{.RPM}}</td>
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
