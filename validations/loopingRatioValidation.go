package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

type LoopingRatio struct {
	Publisher     string `json:"publisher" validate:"required"`
	Device        string `json:"device" validate:"device"`
	Country       string `json:"country" validate:"country"`
	PlacementType string `json:"placement_type" validate:"placement_type"`
	OS            string `json:"os" validate:"os"`
	Browser       string `json:"browser" validate:"browser"`
	LoopingRatio  int16  `json:"looping_ratio" validate:"looping_ratio"`
	Domain        string `json:"domain"`
}

func ValidateLoopingRatio(c *fiber.Ctx) error {
	body := new(LoopingRatio)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for looping ratio. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"country":       "Country code must be 2 characters long and should be in the allowed list",
		"device":        "Device should be in the allowed list",
		"looping_ratio": fmt.Sprintf("Looping ratio value not allowed, it should be <= %s", fmt.Sprintf("%d", constant.MaxLoopingRatioValue)),
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
