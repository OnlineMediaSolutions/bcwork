package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func ValidateBlocks(c *fiber.Ctx) error {
	var request *dto.BlockUpdateRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Targeting. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateBlocks(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate blocks request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateBlocks(request *dto.BlockUpdateRequest) []string {
	var errorMessages = map[string]string{}

	validationErrors := make([]string, 0)

	err := Validator.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if msg, ok := errorMessages[err.Tag()]; ok {
				validationErrors = append(validationErrors, msg)
			} else {
				validationErrors = append(validationErrors,
					fmt.Sprintf("%s is mandatory, validation failed", err.Field()))
			}
		}
	}

	return validationErrors
}
