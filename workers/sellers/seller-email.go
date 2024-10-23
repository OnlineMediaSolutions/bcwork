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
	Name            string
	URL             string
	PublisherDomain []PublisherDomain
	Position        int64
}

type PublisherDomain struct {
	Publisher  string
	Domain     string
	SellerType string
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
                        <th>Publisher - Domain - SellerType</th>
                    </tr>
                    {{range $index, $competitor := .CompetitorsData}}
                        {{ $publisherDomainCount := len $competitor.PublisherDomain }}
                        {{ if eq $publisherDomainCount 1 }}
                            <tr>
                                <td>{{$competitor.Name}}</td>
                                <td>{{$competitor.URL}}</td>
                                <td>{{range $competitor.PublisherDomain}}{{.Publisher}} - {{.Domain}} - {{.SellerType}}{{end}}</td>
                            </tr>
                        {{ else }}
                            {{range $publisherIndex, $publisherDomain := $competitor.PublisherDomain}}
                                {{ if eq $publisherIndex 0 }}
                                    <tr>
                                        <td rowspan="{{$publisherDomainCount}}">{{$competitor.Name}}</td>
                                        <td rowspan="{{$publisherDomainCount}}">{{$competitor.URL}}</td>
                                        <td>{{$publisherDomain.Publisher}} - {{$publisherDomain.Domain}} - {{$publisherDomain.SellerType}}</td>
                                    </tr>
                                {{ else }}
                                    <tr>
                                        <td>{{$publisherDomain.Publisher}} - {{$publisherDomain.Domain}} - {{$publisherDomain.SellerType}}</td>
                                    </tr>
                                {{ end }}
                            {{end}}
                        {{ end }}
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
