package core

import (
	"fmt"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/modules"
	"golang.org/x/net/context"
)

type EmailService struct {
	SendEmailReport func(ctx context.Context, data dto.EmailData) error
}

func NewEmailService(ctx context.Context) *EmailService {
	return &EmailService{
		SendEmailReport: sendEmailReport,
	}
}

func sendEmailReport(ctx context.Context, data dto.EmailData) error {
	email := modules.EmailRequest{
		To:      data.Recipients,
		Subject: data.Subject,
		Body:    string(data.Content),
		IsHTML:  data.IsHtml,
	}

	err := modules.SendEmail(email)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
