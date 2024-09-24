package sellers

import (
	"bytes"
	"github.com/m6yf/bcwork/utils"
	"strings"
	"text/template"
)

type CompetitorData struct {
	Name       string
	URL        string
	Publishers []string
	Domains    []string
}

type EmailData struct {
	Body        string
	Competitors []CompetitorData
}

func GenerateHTMLTableWithTemplate(competitorsData []CompetitorData, body string) (string, error) {
	const tpl = `
	<html>
		<head>
			<title>Sellers JSON Updates</title>
			<style>
				table { width: 100%; border-collapse: collapse; }
				th, td { border: 1px solid black; padding: 8px; text-align: left; }
				th { background-color: #f2f2f2; }
			</style>
		</head>
		<body>
		   <h3>{{.Body}}</h3>
			<table>
				<tr>
					<th>Competitor Name</th>
					<th>Competitor URL</th>
					<th>Added Publishers</th>
					<th>Added Domains</th>
				</tr>
				{{range .CompetitorsData}}
					<tr>
						<td>{{.Name}}</td>
						<td>{{.URL}}</td>
						<td>{{join .Publishers ", "}}</td>
						<td>{{join .Domains ", "}}</td>
					</tr>
				{{else}}
					<tr>
						<td colspan="3">No data available</td>
					</tr>
				{{end}}
			</table>
		</body>
	</html>
	`
	data := struct {
		Body            string
		CompetitorsData []CompetitorData
	}{
		Body:            body,
		CompetitorsData: competitorsData,
	}

	t, err := template.New("emailTemplate").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(tpl)

	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, data); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func SendCustomHTMLEmail(to, bcc, subject string, body string, competitorsData []CompetitorData) error {
	emailData := EmailData{
		Body:        body,
		Competitors: competitorsData,
	}

	htmlBody, err := GenerateHTMLTableWithTemplate(emailData.Competitors, emailData.Body)
	if err != nil {
		return err
	}

	emailReq := utils.EmailRequest{
		To:      to,
		Bcc:     bcc,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return utils.SendEmail(emailReq)
}
