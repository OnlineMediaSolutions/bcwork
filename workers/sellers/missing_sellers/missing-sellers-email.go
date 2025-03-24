package missing_sellers

import (
	"bytes"
	"text/template"
	"time"
)

func GenerateHTMLFromMissingPublishers(statusMap map[string]MissingPublisherInfo) (string, error) {
	currentDate := time.Now().Format(time.DateOnly)

	const tpl = `
<html>
    <head>
        <title>Missing Publishers in seller.json - {{.Date}}</title>
        <style>
            table { width: 100%; border-collapse: collapse; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>Missing Publishers in seller.json</h3>
        {{if (eq (len .) 0)}}
            <p class="no-changes">There are no missing publishers.</p>
        {{else}}
            {{range $partner, $urls := .}}
                <h4>Sellers URL: {{$partner}}</h4>
                {{range $url, $publishers := $urls}}
                    <table>
                        <tr>
                            <th>Publisher Name</th>
                            <th>Publisher ID</th>
                            <th>Status</th>
                        </tr>
                        {{range $publishers}}
                            <tr>
                                <td>{{.PublisherName}}</td>
                                <td>{{.PublisherId}}</td>
                                <td>{{.Status}}</td>
                            </tr>
                        {{end}}
                    </table>
                {{end}}
            {{end}}
        {{end}}
    </body>
</html>
`

	// **Step 1: Group by Partner > URL**
	groupedData := make(map[string]map[string][]MissingPublisherInfo)

	for _, info := range statusMap {
		if _, exists := groupedData[info.SeatOwner]; !exists {
			groupedData[info.SeatOwner] = make(map[string][]MissingPublisherInfo)
		}
		groupedData[info.SeatOwner][info.SeatURL] = append(groupedData[info.SeatOwner][info.SeatURL], info)
	}

	// **Step 2: Prepare data for the template**
	data := struct {
		Date string
		Data map[string]map[string][]MissingPublisherInfo
	}{
		Date: currentDate,
		Data: groupedData,
	}

	// **Step 3: Parse & Execute Template**
	t, err := template.New("missingPublishers").Parse(tpl)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, data); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}
