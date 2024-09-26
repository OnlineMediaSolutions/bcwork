package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
)

var validate = validator.New()

func ValidateCompetitorURL(c *fiber.Ctx) error {
	var request core.CompetitorUpdateRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON.",
		})
	}

	err = validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			fmt.Printf("err", err)
			switch err.Field() {
			case "URL":
				message = "URL must be valid and start with either 'http' or 'https'."
			default:
				message = fmt.Sprintf("%s is required, validation failed", err.Field())
			}

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": message,
			})
		}
	}

	return c.Next()
}
