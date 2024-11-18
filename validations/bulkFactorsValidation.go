package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

type FactorUpdateRequest struct {
	Publisher string  `json:"publisher" validate:"required"`
	Device    string  `json:"device" validate:"device"`
	Country   string  `json:"country" validate:"country"`
	Factor    float64 `json:"factor" validate:"required,factor"`
	Domain    string  `json:"domain"`
}

func ValidateBulkFactors(c *fiber.Ctx) error {
	var requests []FactorUpdateRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"country": fmt.Sprintf("Country code must be %d characters long and should be in the allowed list", constant.MaxCountryCodeLength),
		"device":  "Device should be in the allowed list",
		"factor":  fmt.Sprintf("Factor value not allowed, it should be >= %s and <= %s", fmt.Sprintf("%.2f", constant.MinFactorValue), fmt.Sprintf("%.2f", constant.MaxFactorValue)),
	}

	for _, request := range requests {
		err = Validator.Struct(request)
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
	}

	return c.Next()
}
