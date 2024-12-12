package dto

type EmailData struct {
	To      []string `json:"to"`
	From    string   `json:"from"`
	Body    string   `json:"body"`
	Subject string   `json:"subject"`
}
