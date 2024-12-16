package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

func (o *OMSNewPlatform) SendEmailReport(c *fiber.Ctx) error {
	data := &dto.EmailData{}

	// Read the request body and unmarshal it into the data struct
	if err := c.BodyParser(data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email data payload parsing error", err)
	}

	// Call the email service to send the email report
	err := o.emailService.SendEmailReport(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send email", err)
	}

	// Return a success response (optional)
	return c.JSON(fiber.Map{"message": "Email report sent successfully"})
}
