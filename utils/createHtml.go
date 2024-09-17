package utils

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func CreateHtml(c *fiber.Ctx) error {
	var emailReq EmailRequest

	if err := c.BodyParser(&emailReq); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	err := SendEmail(emailReq)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send email")
	}

	log.Println("Email sent successfully")
	return c.SendString("Email sent successfully")
}
