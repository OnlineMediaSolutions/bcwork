package missing_publishers

import (
	"bytes"
	"text/template"
)

func GenerateHTMLFromMissingPublishers(statusMap map[string]MissingPublisherInfo) (string, error) {
	const tpl = `
<html>
    <head>
        <title>Missing PublisherId's in Sellers</title>
        <style>
            table { width: 100%; border-collapse: collapse; }
            th, td { border: 1px solid black; padding: 8px; text-align: left; }
            th { background-color: #f2f2f2; }
            .no-changes { color: red; font-weight: bold; }
        </style>
    </head>
    <body>
        <h3>Missing PublisherId's in Sellers</h3>
        {{if (eq (len .) 0)}}
            <p class="no-changes">There are no missing publishers.</p>
        {{else}}
            <table>
                <tr>
                    <th>Publisher ID</th>
                    <th>Publisher Name</th>
                    <th>Status</th>
                    <th>Seat Owner</th>
                </tr>
                {{ range . }}
                    <tr>
                        <td>{{.PublisherId}}</td>
                        <td>{{.PublisherName}}</td>
                        <td>{{.Status}}</td>
                        <td>{{.SeatOwner}}</td>
                    </tr>
                {{ end }}
            </table>
        {{ end }}
    </body>
</html>
`
	// Convert the statusMap to a slice for template processing
	var missingPublishers []MissingPublisherInfo
	for _, info := range statusMap {
		missingPublishers = append(missingPublishers, info)
	}

	t, err := template.New("missingPublishers").Parse(tpl)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := t.Execute(&tplBuffer, missingPublishers); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}
