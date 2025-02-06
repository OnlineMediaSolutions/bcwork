package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func AdsTxtValidation(c *fiber.Ctx) error {
	var request *dto.AdsTxt
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for User. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateAdsTxt(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate ads.txt request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateAdsTxt(request *dto.AdsTxt) []string {
	var errorMessages = map[string]string{
		adsTxtDemandStatusValidationKey: adsTxtDemandStatusErrorMessage,
		adsTxtDomainStatusValidationKey: adsTxtDomainStatusErrorMessage,
	}

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
