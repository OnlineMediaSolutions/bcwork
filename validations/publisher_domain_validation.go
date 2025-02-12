package validations

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func PublisherDomainValidation(c *fiber.Ctx) error {
	var request *dto.PublisherDomainUpdateRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for update publisher domain request. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validatePublisherDomain(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate update publisher domain request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validatePublisherDomain(request *dto.PublisherDomainUpdateRequest) []string {
	var errorMessages = map[string]string{
		intergrationTypeValidationKey: intergrationTypeErrorMessage + ": " + strings.Join(integrationTypes, ","),
		mediaTypeValidationKey:        mediaTypeErrorMessage + ": " + strings.Join(mediaTypes, ","),
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
