package dto

import "encoding/json"

type EmailData struct {
	Recipients []string        `json:"recipients"`
	Content    json.RawMessage `json:"content"`
	Subject    string          `json:"subject"`
	IsHtml     bool            `json:"isHtml"`
}
