package utils

import (
	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	To      string `json:"to"`
	Bcc     string `json:"bcc"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  bool   `json:"is_html"`
}

func SendEmail(emailReq EmailRequest) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "sources@onlinemediasolutions.com")
	mailer.SetHeader("To", emailReq.To)
	mailer.SetHeader("Bcc", emailReq.Bcc)
	mailer.SetHeader("Subject", emailReq.Subject)

	if emailReq.IsHTML {
		mailer.SetBody("text/html", emailReq.Body)
	} else {
		mailer.SetBody("text/plain", emailReq.Body)
	}

	dialer := gomail.NewDialer("smtp.gmail.com", 465, "smtp@onlinemediasolutions.com", "sqmrvlxfljsjkyhh")

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
