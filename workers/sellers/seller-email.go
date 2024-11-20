package sellers

import (
	"bytes"
	"github.com/m6yf/bcwork/modules"
	"strings"
	"text/template"
)

type EmailData struct {
	Body        string
	Competitors []CompetitorData
}

type CompetitorData struct {
	Name                   string
	URL                    string
	AddedPublisherDomain   []PublisherDomain
	DeletedPublisherDomain []PublisherDomain
	Position               string
}

type PublisherDomain struct {
	Publisher  string
	Domain     string
	SellerType string
	SellerId   string
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
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>{{.Body}}</h3>
        {{if (eq (len .CompetitorsData) 0)}}
            <p class="no-changes">There are no seller changes.</p>
        {{else}}
            {{ $hasUpdates := false }} <!-- Initialize the flag to track if any updates exist -->
            
            <!-- First, check if there are any updates -->
            {{ range .CompetitorsData }}
                {{ $addedCount := len .AddedPublisherDomain }}
                {{ $deletedCount := len .DeletedPublisherDomain }}

                {{ if or (gt $addedCount 0) (gt $deletedCount 0) }}
                    {{ $hasUpdates = true }} <!-- Set the flag to true if there are updates -->
                {{ end }}
            {{ end }}

            <!-- Now, display the table only if there are updates -->
            {{ if $hasUpdates }}
                <table>
                    <tr>
                        <th>Competitor Name</th>
                        <th>Competitor URL</th>
                        <th>Added Publisher - Domain - SellerType</th>
                        <th>Deleted Publisher - Domain - SellerType</th>
                    </tr>
                    {{ range .CompetitorsData }}
                        {{ $addedCount := len .AddedPublisherDomain }}
                        {{ $deletedCount := len .DeletedPublisherDomain }}

                        {{ if or (gt $addedCount 0) (gt $deletedCount 0) }} <!-- Only show competitors with updates -->
                            <tr>
                                <td>{{.Name}}</td>
                                <td>{{.URL}}</td>
                                <td>{{range .AddedPublisherDomain}}{{.Publisher}} - {{.Domain}} - {{.SellerType}}<br>{{end}}</td>
                                <td>{{range .DeletedPublisherDomain}}{{.Publisher}} - {{.Domain}} - {{.SellerType}}<br>{{end}}</td>
                            </tr>
                        {{ end }}
                    {{ end }}
                </table>
            {{ end }}
        {{ end }}
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

func SendCustomHTMLEmail(to, bcc, subject string, body string, competitorsData []CompetitorData) error {
	toRecipients := strings.Split(to, ",")
	emailData := EmailData{
		Body:        body,
		Competitors: competitorsData,
	}

	htmlBody, err := GenerateHTMLTableWithTemplate(emailData.Competitors, emailData.Body)
	if err != nil {
		return err
	}

	emailReq := modules.EmailRequest{
		To:      toRecipients,
		Bcc:     bcc,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return modules.SendEmail(emailReq)
}
