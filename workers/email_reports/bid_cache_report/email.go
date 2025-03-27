package bid_cache_report

import (
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (worker *Worker) SendEmail(data []*BidCacheData, emailCreds EmailCreds) error {
	toRecipients := strings.Split(emailCreds.TO, ",")
	bccStr := strings.Split(emailCreds.BCC, ",")

	html := generateHTMLTable(data)
	subject := fmt.Sprintf("Bid Caching Report for %s", time.Now().Format("2006-01-02"))

	emailReq := modules.EmailRequest{
		To:      toRecipients,
		Bcc:     bccStr,
		Subject: subject,
		Body:    html,
		IsHTML:  true,
	}
	log.Info().Msg("Sending bid cache email report")

	return modules.SendEmail(emailReq)
}

func generateHTMLTable(data []*BidCacheData) string {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	sb.WriteString("<table border=\"1\" style=\"border-collapse: collapse; width: 100%;\">")
	sb.WriteString("<tr style=\"background-color: blue; color: white;\">")
	sb.WriteString("<th style=\"width: 70%;\">Time</th><th>PublisherID</th><th>Domain</th><th>Target</th>" +
		"<th>Revenue</th><th>Cost</th><th>DP Fee</th><th>Sold Impressions</th><th>Publisher Impressions</th>" +
		"<th>Data Fee</th><th>GP</th><th style=\"width: 70%;\">GP Per PubImp</th></tr>")

	for _, d := range data {
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td>"+
			"<td>%s</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%d</td><td>%d</td>"+
			"<td>%.2f</td><td>%.2f</td><td>%.6f</td></tr>",
			d.Time, d.PublisherID, d.Domain, d.Target, d.Revenue, d.Cost, d.DemandPartnerFee, d.SoldImpressions,
			d.PublisherImpressions, d.DataFee, d.GP, d.GPperPubImp))
	}

	sb.WriteString("</table>")
	sb.WriteString("</body></html>")

	return sb.String()
}
