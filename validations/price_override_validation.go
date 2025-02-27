package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
)

func ValidatePriceOverride(c *fiber.Ctx) error {
	body := new(dto.PriceOverrideRequest)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for bid caching. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validatePriceOverride(body)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate price override",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validatePriceOverride(request *dto.PriceOverrideRequest) []string {
	var errorMessages = map[string]string{
		ipsKey:           duplicateIpsErrorMessage,
		overridePriceKey: overridePriceErrorMessage,
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
