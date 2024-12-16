package dto

import "encoding/json"

type EmailData struct {
	Recipients []string        `json:"recipients"`
	Bcc        string          `json:"bcc"`
	Content    json.RawMessage `json:"content"`
	Subject    string          `json:"subject"`
}
