package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AdjustRequest struct {
	Domain []string `json:"domain" validate:"required"`
	Value  float64  `json:"Value" validate:"required,floor"`
}

func ValidateAdjusterURL(c *fiber.Ctx) error {
	body := new(Factor)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for adjust. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"floor": fmt.Sprintf("Value is mandatory and can't be negative"),
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
