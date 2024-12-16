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
		Bcc:     data.Recipients[0],
		Body:    string(data.Content),
		IsHTML:  false,
	}

	err := modules.SendEmail(email)
	if err != nil {
		// Wrap the error with additional context
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Return nil if the email was sent successfully
	return nil
}
