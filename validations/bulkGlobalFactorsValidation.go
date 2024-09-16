package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
)

type errorBulkResponse struct {
	Status  string              `json:"status"`
	Message map[string][]string `json:"message"`
}

const (
	errorStatus        = "error"
	keyValidationError = "key most be one of the following: 'tech_fee', 'consultant_fee' or 'tam_fee'"
)

func ValidateBulkGlobalFactor(c *fiber.Ctx) error {
	var requests []*core.GlobalFactorRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Global Factor. Please ensure it's a valid JSON.",
		})
	}

	errorResponse := validateBulkGlobalFactor(requests)

	if errorResponse.Status == errorStatus {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return c.Next()
}

func validateBulkGlobalFactor(requests []*core.GlobalFactorRequest) errorBulkResponse {
	var errorMessages = map[string]string{
		"globalFactorKey": keyValidationError,
	}

	errorResponse := errorBulkResponse{Message: make(map[string][]string)}
	for idx, request := range requests {
		err := Validator.Struct(request)
		if err != nil {
			errorResponse.Status = errorStatus
			key := fmt.Sprintf("request %v", idx+1)
			for _, err := range err.(validator.ValidationErrors) {
				if msg, ok := errorMessages[err.Tag()]; ok {
					errorResponse.Message[key] = append(errorResponse.Message[key], msg)
				} else {
					errorResponse.Message[key] = append(errorResponse.Message[key],
						fmt.Sprintf("%s is mandatory, validation failed", err.Field()))
				}
			}
		}
	}

	return errorResponse
}
