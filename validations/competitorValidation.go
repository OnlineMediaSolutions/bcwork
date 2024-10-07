package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

var validate = validator.New()

func ValidateCompetitorURL(c *fiber.Ctx) error {
	var requests []constant.CompetitorUpdateRequest

	if err := c.BodyParser(&requests); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
		})
	}

	var validationErrors []fiber.Map

	for _, request := range requests {
		if err := validate.Struct(request); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				var message string
				switch err.Field() {
				case "URL":
					message = fmt.Sprintf("Competitor '%s': URL must be valid and start with either 'http' or 'https'.", request.Name)
				default:
					message = fmt.Sprintf("Competitor '%s': %s is required, validation failed", request.Name, err.Field())
				}

				validationErrors = append(validationErrors, fiber.Map{
					"competitor": request.Name,
					"field":      err.Field(),
					"message":    message,
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": validationErrors,
		})
	}

	return c.Next()
}
