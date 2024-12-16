package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

func (o *OMSNewPlatform) SendEmailReport(c *fiber.Ctx) error {
	data := &dto.EmailData{}

	if err := c.BodyParser(data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email data payload parsing error", err)
	}

	err := o.emailService.SendEmailReport(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send email", err)
	}

	return c.JSON(fiber.Map{"message": "Email report sent successfully"})
}
