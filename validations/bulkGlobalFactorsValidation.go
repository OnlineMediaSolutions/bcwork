package validations

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
)

type errorBulkResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

const (
	globalFactorConsultantFeeType = "consultant_fee"

	errorStatus                  = "error"
	validationError              = "couldn't validate some of the requests"
	keyValidationError           = "key most be one of the following: 'tech_fee', 'consultant_fee' or 'tam_fee'"
	publisherValidationError     = "only 'consultant_fee' can have publisher"
	valueValidationError         = "value must be positive"
	consultantFeeValidationError = "'consultant_fee' must have publisher"
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

	validationErrors := validateBulkGlobalFactor(requests)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorBulkResponse{
			Status:  errorStatus,
			Message: validationError,
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateBulkGlobalFactor(requests []*core.GlobalFactorRequest) map[string][]string {
	var errorMessages = map[string]string{
		"globalFactorKey": keyValidationError,
		"gte":             valueValidationError,
	}

	validationErrors := make(map[string][]string)
	for idx, request := range requests {
		key := fmt.Sprintf("request %v", idx+1)

		if request.Key != globalFactorConsultantFeeType && request.Publisher != "" {
			validationErrors[key] = append(validationErrors[key], publisherValidationError)
		}

		if request.Key == globalFactorConsultantFeeType && request.Publisher == "" {
			validationErrors[key] = append(validationErrors[key], consultantFeeValidationError)
		}

		err := Validator.Struct(request)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				log.Println(err.Tag())
				if msg, ok := errorMessages[err.Tag()]; ok {
					validationErrors[key] = append(validationErrors[key], msg)
				} else {
					validationErrors[key] = append(validationErrors[key],
						fmt.Sprintf("%s is mandatory, validation failed", err.Field()))
				}
			}
		}
	}

	return validationErrors
}
