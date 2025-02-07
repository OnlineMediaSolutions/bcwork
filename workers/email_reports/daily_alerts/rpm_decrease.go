package daily_alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	HOURLY_PERCENTAGE = 0.4
)

type RPMDecreaseReport struct{}

func (r *RPMDecreaseReport) Aggregate(reports []Result) map[string][]AggregatedReport {
	aggregated := make(map[string][]AggregatedReport)

	for _, r := range reports {
		key := fmt.Sprintf("%s|%s|%s", r.AM, r.Domain, r.Publisher)
		if aggList, exists := aggregated[key]; exists {
			aggregated[key] = append(aggList, AggregatedReport{
				AM:           r.AM,
				Domain:       r.Domain,
				Publisher:    r.Publisher,
				PaymentType:  r.PaymentType,
				Date:         r.Date,
				PubImps:      r.PubImps,
				DataStamp:    r.DataStamp,
				RPM:          r.RPM,
				Ratio:        r.Ratio,
				LoopingRatio: r.LoopingRatio,
				CPM:          r.CPM,
				Cost:         r.Cost,
				DPRPM:        r.DPRPM,
				Revenue:      r.Revenue,
				GP:           r.GP,
				GPP:          r.GPP,
			})
		} else {
			aggregated[key] = []AggregatedReport{{
				AM:           r.AM,
				Domain:       r.Domain,
				Publisher:    r.Publisher,
				PaymentType:  r.PaymentType,
				Date:         r.Date,
				PubImps:      r.PubImps,
				DataStamp:    r.DataStamp,
				RPM:          r.RPM,
				Ratio:        r.Ratio,
				LoopingRatio: r.LoopingRatio,
				CPM:          r.CPM,
				Cost:         r.Cost,
				DPRPM:        r.DPRPM,
				Revenue:      r.Revenue,
				GP:           r.GP,
				GPP:          r.GPP,
			}}
		}
	}

	return aggregated
}

func (r *RPMDecreaseReport) ComputeAverage(aggregated map[string][]AggregatedReport, worker *Worker) map[string][]AlertsEmailRepo {
	amDomainData := make(map[string][]AggregatedReport)
	last12HoursReports := make(map[string][]AggregatedReport)

	lastCompleteHour := worker.CurrentTime.Truncate(time.Hour)
	yesterdayHour := lastCompleteHour.Add(-24 * time.Hour)

	lastCompleteHourUnix := lastCompleteHour.Unix() / 100
	yesterdayHourUnix := yesterdayHour.Unix() / 100

	for key, aggs := range aggregated {
		for _, agg := range aggs {
			if agg.DataStamp == lastCompleteHourUnix || agg.DataStamp == yesterdayHourUnix {
				amDomainData[key] = append(amDomainData[key], agg)
				last12HoursReports[key] = append(last12HoursReports[key], agg)

			}
		}
	}

	alerts := make(map[string]AggregatedReport)
	repo := AlertsEmailRepo{}
	var emailReports []AlertsEmailRepo

	for key, reports := range amDomainData {
		if len(reports) < 2 {
			continue
		}

		sort.Slice(reports, func(i, j int) bool {
			return reports[i].DataStamp > reports[j].DataStamp
		})

		if reports[0].RPM < HOURLY_PERCENTAGE*(reports[1].RPM) {
			alerts[key] = reports[0]
			repo = AlertsEmailRepo{
				Email:        "sonai@onlinemediasolutions.com", //worker.UserData[key]
				AM:           key,
				FirstReport:  reports[0],
				SecondReport: last12HoursReports[key],
			}

			emailReports = append(emailReports, repo)

		}

	}

	avgDataMap := make(map[string][]AlertsEmailRepo)

	for _, repo := range emailReports {
		avgDataMap[repo.Email] = append(avgDataMap[repo.Email], repo)
	}
	return avgDataMap
}

func (r *RPMDecreaseReport) PrepareAndSendEmail(reportData map[string][]AlertsEmailRepo, worker *Worker) error {
	if len(reportData) > 0 {
		now := time.Now()
		today := now.Format(constant.PostgresTimestamp)

		for email, alerts := range reportData {
			subject := fmt.Sprintf("RPM decrease alert for %s", today)
			message := fmt.Sprintf("Dear %s,\n\nRPM decrease alert for %s.\n\nPlease review the details below.", email, today)

			err := r.SendCustomHTMLEmail(email, "sonai@onlinemediasolutions.com", subject, message, alerts)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to send email to %s", email)
			}
		}
	}
	return nil
}

func (r *RPMDecreaseReport) SendCustomHTMLEmail(to, bcc, subject string, body string, report []AlertsEmailRepo) error {
	toRecipients := strings.Split(to, ",")
	bccString := strings.Split(bcc, ",")
	emailData := EmailData{
		Body:   body,
		Report: report,
	}

	htmlBody, err := r.GenerateHTMLTableWithTemplate(emailData.Report, emailData.Body)
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

func (r *RPMDecreaseReport) Request(worker *Worker) ([]Result, error) {
	compassClient := compass.NewCompass()

	requestData := r.GetRequestData(worker)

	data, err := json.Marshal(requestData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling request data")
	}

	reportData, err := compassClient.Request("/report-dashboard/report-query/merged", "POST", data, true)

	if err != nil {
		return nil, fmt.Errorf("error getting report data")
	}

	var report Report
	err = json.Unmarshal(reportData, &report)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling report data: %w", err)
	}

	return report.Data.Result, nil
}

func (r *RPMDecreaseReport) GetRequestData(worker *Worker) RequestData {
	requestData := RequestData{
		Data: RequestDetails{
			Date: Date{
				Range: []string{
					worker.Start,
					worker.End,
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
				"looping_ratio",
				"mergedEstRpm",
				"PubImps",
				"estRatio",
				"est_cpm",
				"EstCost",
				"estDpRpm",
				"EstRevenue",
				"mergedEstGp",
				"mergedEstGpp",
			},
		},
	}

	return requestData
}

func (r *RPMDecreaseReport) GenerateHTMLTableWithTemplate(report []AlertsEmailRepo, body string) (string, error) {
	const tpl = `
<html>
    <head>
        <title>RPM decrease alert</title>
        <style>
            table { width: 100%; border-collapse: collapse; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>{{.Body}}</h3>
        
        <!-- First Report Section -->
        <h4>RPM decrease alert</h4>
        <table>
            <tr>
                <th>Date</th>
                <th>Account Manager</th>
                <th>Publisher</th>
                <th>Domain</th>
            </tr>
            {{if .FirstReport}}
            <tr>
                <td>{{.FirstReport.Date}}</td>
                <td>{{.FirstReport.AM}}</td>
                <td>{{.FirstReport.Publisher}}</td>
                <td>{{.FirstReport.Domain}}</td>
            </tr>
            {{end}}
        </table>
        
        <!-- Second Report Section -->
        <h4>RPM by domain</h4>
        <table>
            <tr>
                <th>Date</th>
                <th>Publisher</th>
                <th>Domain</th>
                <th>Payment Type</th>
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
            {{range .SecondReport}}
            <tr>
                <td>{{.Date}}</td>
                <td>{{.Publisher}}</td>
                <td>{{.Domain}}</td>
                <td>{{.PaymentType}}</td>
                <td>{{.PubImps}}</td>
                <td>{{.LoopingRatio}}</td>
                <td>{{.Ratio}}</td>
                <td>{{.CPM}}</td>
                <td>{{.Cost}}</td>
                <td>{{.RPM}}</td>
                <td>{{.DPRPM}}</td>
                <td>{{.Revenue}}</td>
                <td>{{.GP}}</td>
                <td>{{.GPP}}</td>
            </tr>
            {{end}}
        </table>
    </body>
</html>
`

	// Ensure `report` has at least one element before accessing `report[0]`
	if len(report) == 0 {
		return "", fmt.Errorf("no reports available")
	}

	data := struct {
		Body         string
		FirstReport  AggregatedReport
		SecondReport []AggregatedReport
	}{
		Body:         body,
		FirstReport:  report[0].FirstReport,
		SecondReport: report[0].SecondReport, // This is a slice, correctly handled in the template
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
