package utils

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func CreateHtml(c *fiber.Ctx) error {
	emailReq := EmailRequest{
		To:      "sonai@onlinemediasolutions.com",
		Bcc:     "sonai@onlinemediasolutions.com",
		Subject: "Test Email",
		Body:    "<h1>Hello, World!</h1><p>This is a test email.</p>",
		IsHTML:  false,
	}

	err := SendEmail(emailReq)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Email sent successfully")

	return c.SendString("Email sent successfully")
}
