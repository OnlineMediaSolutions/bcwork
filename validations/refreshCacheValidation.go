package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

type RefreshCache struct {
	Publisher     string `json:"publisher" validate:"required"`
	Device        string `json:"device" validate:"device"`
	Country       string `json:"country" validate:"country"`
	PlacementType string `json:"placement_type" validate:"placement_type"`
	OS            string `json:"os" validate:"os"`
	Browser       string `json:"browser" validate:"browser"`
	RefreshCache  int16  `json:"refresh_cache" validate:"refresh_cache"`
	Domain        string `json:"domain"`
}

func ValidateRefreshCache(c *fiber.Ctx) error {
	body := new(RefreshCache)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for refresh cache. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"country":       "Country code must be 2 characters long and should be in the allowed list",
		"device":        "Device should be in the allowed list",
		"refresh_cache": fmt.Sprintf("Refresh cache value not allowed, it should be <= %s and >= %s", fmt.Sprintf("%d", constant.MaxRefreshCacheValue), fmt.Sprintf("%d", constant.MinRefreshCacheValue)),
	}

	err = Validator.Struct(body)
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
	return c.Next()
}