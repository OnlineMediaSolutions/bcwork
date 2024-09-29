package validations

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func ValidateDownload(c *fiber.Ctx) error {
	var body []json.RawMessage
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Download. Please ensure it's a valid JSON.",
		})
	}

	if len(body) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "no data to create csv file.",
		})
	}

	return c.Next()
}
