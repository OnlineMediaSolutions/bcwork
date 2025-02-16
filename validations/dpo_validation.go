package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/constant"
)

func ValidateDPO(c *fiber.Ctx) error {
	body := new(dto.DPORuleUpdateRequest)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for dpo. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"country":   "Country code must be 2 characters long and should be in the allowed list",
		"factorDpo": fmt.Sprintf("Factor value not allowed, it should be >= %s and <= %s", fmt.Sprintf("%d", constant.MinDPOFactorValue), fmt.Sprintf("%d", constant.MaxDPOFactorValue)),
		"all":       "'all' is not allowed to the following parameters: OS, Browser and Placement_type",
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
