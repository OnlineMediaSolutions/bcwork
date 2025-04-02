package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func AdsTxtUpdateValidation(c *fiber.Ctx) error {
	var request *dto.AdsTxtUpdateRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for ads txt update. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateAdsTxtUpdate(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate ads.txt update request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateAdsTxtUpdate(request *dto.AdsTxtUpdateRequest) []string {
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

	if request.DemandStatus == nil && request.DomainStatus == nil {
		validationErrors = append(validationErrors, "domain and demand statuses are nil")
	}

	return validationErrors
}
