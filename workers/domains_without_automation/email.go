package domains_without_automation

import (
	"bytes"
	"encoding/json"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules"
	"github.com/rs/zerolog/log"
	"strings"
	"text/template"
)

const subject = "Domains without automation - Last 7 days"
const realTimeReport = "real_time_report"

func CreateEmailBody(accountManagerData []Result) (string, error) {
	const tpl = `
<!DOCTYPE html>
<html>
    <head>
        <title>Domains without automation - Last 7 days</title>
        <style>
            table { width: 100%; border-collapse: collapse; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>Domains without automation</h3>
                <table>
                    <tr>
                        <th>Publisher</th>
                        <th>Domain</th>
                        <th>Account Manager</th>
                        <th>Pub Imps</th>
                        <th>Looping Ratio</th>
                        <th>Cost</th>
                        <th>CPM</th>
                        <th>Revenue</th>
                        <th>RPM</th>
                        <th>DP RPM</th>
                        <th>GP</th>
                        <th>GP%</th>
                    </tr>
					{{range .}}
					<tr>
						<td>{{.Publisher}}</td>
						<td>{{.Domain}}</td>
						<td>{{.AccountManager}}</td>
						<td>{{.PubImps}}</td>
						<td>{{printf "%.2f" .LoopingRatio}}</td>
						<td>{{printf "%.2f" .Cost}}</td>
						<td>{{printf "%.2f" .CPM}}</td>
						<td>{{printf "%.2f" .Revenue}}</td>
						<td>{{printf "%.2f" .RPM}}</td>
						<td>{{printf "%.2f" .DpRPM}}</td>
						<td>{{printf "%.2f" .GP}}</td>
						<td>{{printf "%.2f" .GPP}}</td>
					</tr>
					{{end}}
                </table>
    </body>
</html>
`
	// Parse the HTML template
	t, err := template.New("table").Parse(tpl)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, accountManagerData); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func SendEmails(emails map[string]string, domainsPerAccountManager map[string][]Result, managerList []Result) {
	//Send email to Manager (Maayan)
	managersEmail := emails[managerEmail]
	emailCredsMap, _ := config.FetchConfigValues([]string{realTimeReport})

	var emailProperties EmailProperties
	if err := json.Unmarshal([]byte(emailCredsMap[realTimeReport]), &emailProperties); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal email credentials")
	}
	bccEmails := strings.Split(emailProperties.BCC, ",")
	bccEmails = append(bccEmails, managersEmail)

	//Send email per account manager
	for _, accountManager := range domainsPerAccountManager {
		body, _ := CreateEmailBody(accountManager)
		sendEmail(body, emails[accountManager[0].AccountManager], bccEmails)
		log.Info().Msg("Email sent to AM: " + emails[accountManager[0].AccountManager])
	}
}

func sendEmail(body string, email string, bccEmails []string) {

	emailReq := modules.EmailRequest{
		To:      strings.Split(email, ","),
		Bcc:     bccEmails,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}

	modules.SendEmail(emailReq)
	log.Debug().Msg("Email to " + email + " was sent")
}
