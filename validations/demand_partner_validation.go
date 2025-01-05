package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func ValidateDemandPartner(c *fiber.Ctx) error {
	var request *dto.DemandPartner
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Demand Partner. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateDemandPartner(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate Demand Partner request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateDemandPartner(request *dto.DemandPartner) []string {
	var errorMessages = map[string]string{
		approvalProcessKey: approvalProcessErrorMessage,
		dpBlocksKey:        dpBlocksErrorMessage,
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

	for _, child := range request.Children {
		err := Validator.Struct(child)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				if msg, ok := errorMessages[err.Tag()]; ok {
					validationErrors = append(validationErrors, msg)
				} else {
					validationErrors = append(validationErrors,
						fmt.Sprintf("Children: %s is mandatory, validation failed", err.Field()))
				}
			}
		}
	}

	for _, connection := range request.Connections {
		err := Validator.Struct(connection)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				if msg, ok := errorMessages[err.Tag()]; ok {
					validationErrors = append(validationErrors, msg)
				} else {
					validationErrors = append(validationErrors,
						fmt.Sprintf("Connections: %s is mandatory, validation failed", err.Field()))
				}
			}
		}
	}

	return validationErrors
}
