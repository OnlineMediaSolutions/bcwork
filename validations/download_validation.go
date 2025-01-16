package validations

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func ValidateDownload(c *fiber.Ctx) error {
	var body dto.DownloadRequest
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Download. Please ensure it's a valid JSON.",
		})
	}

	if len(body.Data) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("no data to create %v file.", body.FileFormat),
		})
	}

	return c.Next()
}
