package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/m6yf/bcwork/config"
	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	To       []string `json:"to"`
	Bcc      []string `json:"bcc"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	IsHTML   bool     `json:"is_html"`
	Attach   *bytes.Buffer
	Filename string
}

type EmailCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	From     string `json:"from"`
	Port     int    `json:"port"`
}

func GetEmailCredsByKey(configKey string) (*EmailCreds, error) {
	emailCredsMap, err := config.FetchConfigValues([]string{configKey})
	if err != nil {
		return nil, fmt.Errorf("error fetching config values: %w", err)
	}

	credsRaw, found := emailCredsMap[configKey]
	if !found {
		return nil, fmt.Errorf("config key not found for email")
	}

	var emailCreds EmailCreds
	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		return nil, fmt.Errorf("error unmarshalling email credentials: %w", err)
	}

	return &emailCreds, nil
}

func SendEmail(emailReq EmailRequest) error {
	emailCreds, err := GetEmailCredsByKey("email")
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailCreds.From)
	mailer.SetHeader("To", emailReq.To...)
	mailer.SetHeader("Bcc", emailReq.Bcc...)
	mailer.SetHeader("Subject", emailReq.Subject)

	if emailReq.IsHTML {
		mailer.SetBody("text/html", emailReq.Body)
	} else {
		mailer.SetBody("text/plain", emailReq.Body)
	}

	if emailReq.Attach != nil {
		mailer.Attach(emailReq.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := emailReq.Attach.WriteTo(w)

			return err
		}))
	}

	dialer := gomail.NewDialer(emailCreds.Host, emailCreds.Port, emailCreds.Username, emailCreds.Password)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
