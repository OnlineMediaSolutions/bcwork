package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils/constant"
)

func ValidateDPOBulk(c *fiber.Ctx) error {
	var dpoList []core.Dpo
	err := c.BodyParser(&dpoList)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for dpo. Please ensure it's a valid JSON array.",
		})
	}

	var errorMessages = map[string]string{
		"country":   "Country code must be 2 characters long and should be in the allowed list",
		"factorDpo": fmt.Sprintf("Factor value not allowed, it should be >= %d and <= %d", constant.MinDPOFactorValue, constant.MaxDPOFactorValue),
	}

	var validationErrors []fiber.Map

	for i, body := range dpoList {
		err = Validator.Struct(body)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				errorMessage := errorMessages[err.Tag()]
				if errorMessage == "" {
					errorMessage = fmt.Sprintf("%s is mandatory, validation failed", err.Field())
				}
				validationErrors = append(validationErrors, fiber.Map{
					"status":  "error",
					"index":   i,
					"field":   err.Field(),
					"message": errorMessage,
				})
				break
			}
		}
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	return c.Next()
}
