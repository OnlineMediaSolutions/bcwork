package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FloorUpdateRequest struct {
	Publisher string  `json:"publisher" validate:"required"`
	Device    string  `json:"device" validate:"device"`
	Country   string  `json:"country" validate:"country"`
	Floor     float64 `json:"floor" validate:"required,floor"`
	Domain    string  `json:"domain"`
}

func ValidateBulkFloor(c *fiber.Ctx) error {
	var requests []FloorUpdateRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON for floor.",
		})
	}

	var errorMessages = map[string]string{
		"country": "Country code must be 2 characters long and should be in the allowed list",
		"device":  "Device should be in the allowed list",
		"floor":   "Floor should not be negative value",
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
