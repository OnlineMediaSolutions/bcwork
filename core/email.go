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

	fmt.Println(data, "data")

	email := modules.EmailRequest{
		To:      data.To,
		Subject: data.Subject,
		Body:    data.Body,
		IsHTML:  false,
	}

	fmt.Println(email, "email")

	err := modules.SendEmail(email)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
