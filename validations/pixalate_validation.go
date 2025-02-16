package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Pixalate struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate" validate:"rate"`
	Active    *bool   `json:"active,omitempty" validate:"active"`
}

func ValidatePixalate(c *fiber.Ctx) error {
	body := new(Pixalate)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"rate":   "Rate should be between 0 and 100",
		"active": "Active is mandatory (true of false)",
	}

	err = Validator.Struct(body)
	if err != nil {
		errorResponse := map[string]string{
			"status": "error",
		}
		for _, err := range err.(validator.ValidationErrors) {
			if msg, ok := errorMessages[err.Tag()]; ok {
				errorResponse["message"] = msg
			} else {
				errorResponse["message"] = fmt.Sprintf("%s is mandatory, validation failed", err.Field())
			}
			break
		}

		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return c.Next()
}
