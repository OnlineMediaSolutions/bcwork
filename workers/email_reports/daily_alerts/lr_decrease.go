package daily_alerts

import (
	"bytes"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"math"
	"strings"
	"text/template"
	"time"
)

const (
	PERCENTAGE     = 50
	COUNT_HOURS    = 3
	HOURLY_PUBIMPS = 500
)

type LoopingRationDecreaseReport struct{}

func (l *LoopingRationDecreaseReport) Aggregate(reports []Result) map[string][]AggregatedReport {
	aggregated := make(map[string][]AggregatedReport)

	for _, r := range reports {
		key := fmt.Sprintf("%s|%s|%s", r.AM, r.Domain, r.Publisher)

		if aggList, exists := aggregated[key]; exists {
			aggregated[key] = append(aggList, AggregatedReport{
				AM:           r.AM,
				Domain:       r.Domain,
				Publisher:    r.Publisher,
				Date:         r.Date,
				PubImps:      r.PubImps,
				DataStamp:    r.DataStamp,
				LoopingRatio: r.LoopingRatio,
			})
		} else {
			aggregated[key] = []AggregatedReport{{
				AM:           r.AM,
				Domain:       r.Domain,
				Publisher:    r.Publisher,
				Date:         r.Date,
				DataStamp:    r.DataStamp,
				PubImps:      r.PubImps,
				LoopingRatio: r.LoopingRatio,
			}}
		}
	}

	return aggregated
}

func (l *LoopingRationDecreaseReport) ComputeAverage(aggregated map[string][]AggregatedReport, worker *Worker) map[string][]AlertsEmailRepo {
	amDomainData := make(map[string][]AggregatedReport)
	last12HoursReports := make(map[string][]AggregatedReport)
	for key, aggs := range aggregated {
		for _, agg := range aggs {
			if agg.DataStamp >= worker.ThreeHoursAgo && agg.PubImps > HOURLY_PUBIMPS {
				amDomainData[key] = append(amDomainData[key], agg)
			}

			last12HoursReports[key] = append(last12HoursReports[key], agg)
		}
	}

	alerts := make(map[string]AggregatedReport)
	repo := AlertsEmailRepo{}
	var emailReports []AlertsEmailRepo

	for key, reports := range amDomainData {
		var total float64
		count := 0
		for _, report := range reports {
			total += report.LoopingRatio
			count++
		}

		if count >= COUNT_HOURS {
			avgLast3H := total / float64(count)

			latestReport := reports[len(reports)-1]
			if latestReport.LoopingRatio < avgLast3H {
				percentageChange := math.Abs((latestReport.LoopingRatio - avgLast3H) / avgLast3H * 100)

				if percentageChange > PERCENTAGE {
					alerts[key] = latestReport
					repo = AlertsEmailRepo{
						Email:        worker.UserData[key],
						AM:           key,
						FirstReport:  latestReport,
						SecondReport: last12HoursReports[key],
					}

					emailReports = append(emailReports, repo)

				}
			}
		}
	}

	avgDataMap := make(map[string][]AlertsEmailRepo)

	for _, repo := range emailReports {
		avgDataMap[repo.Email] = append(avgDataMap[repo.Email], repo)
	}

	return avgDataMap
}

func (l *LoopingRationDecreaseReport) PrepareAndSendEmail(reportData map[string][]AlertsEmailRepo, worker *Worker) error {
	if len(reportData) > 0 {
		now := time.Now()
		today := now.Format(constant.PostgresTimestamp)

		for email, alerts := range reportData {
			subject := fmt.Sprintf("Looping ratio decrease alert for %s", today)
			message := fmt.Sprintf("Dear %s,\n\nLooping ratio decrease alert for %s.\n\nPlease review the details below.", email, today)

			err := l.SendCustomHTMLEmail(email, "sonai@onlinemediasolutions.com", subject, message, alerts)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to send email to %s", email)
			}
		}
	}
	return nil
}

func (l *LoopingRationDecreaseReport) SendCustomHTMLEmail(to, bcc, subject string, body string, report []AlertsEmailRepo) error {
	toRecipients := strings.Split(to, ",")
	bccString := strings.Split(bcc, ",")
	emailData := EmailData{
		Body:   body,
		Report: report,
	}

	htmlBody, err := l.GenerateHTMLTableWithTemplate(emailData.Report, emailData.Body)
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

func (l *LoopingRationDecreaseReport) GenerateHTMLTableWithTemplate(report []AlertsEmailRepo, body string) (string, error) {
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
        
        <!-- First Report Section -->
        <h4>Looping ratio decrease alert</h4>
        <table>
            <tr>
               <th>Date</th>
               <th>Account Manager</th>
                <th>Domain</th>
                <th>Publisher</th>
                <th>Looping Ratio</th>
            </tr>
            {{if .FirstReport}}
            <tr>
                <td>{{.FirstReport.Date}}</td>
                <td>{{.FirstReport.AM}}</td>
                <td>{{.FirstReport.Domain}}</td>
                <td>{{.FirstReport.Publisher}}</td>
                <td>{{.FirstReport.LoopingRatio}}</td>
            </tr>
            {{end}}
        </table>
        
        <!-- Second Report Section -->
        <h4>Last 12 hours Looping ratio by domain</h4>
        <table>
            <tr>
                <th>Date</th>
                <th>Looping Ratio</th>
            </tr>
            {{range .SecondReport}}
            <tr>
                <td>{{.Date}}</td>
                <td>{{.LoopingRatio}}</td>
            </tr>
            {{end}}
        </table>
    </body>
</html>
`

	data := struct {
		Body         string
		FirstReport  AggregatedReport
		SecondReport []AggregatedReport
	}{
		Body:         body,
		FirstReport:  report[0].FirstReport,
		SecondReport: report[0].SecondReport,
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
