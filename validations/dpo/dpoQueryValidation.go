package dpo

import (
	"fmt"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/validations"
)

type RequestQuery struct {
	Rid    string  `query:"rid" validate:"required"`
	Factor float64 `query:"factor" validate:"required,factorDpo"`
}

func ValidateQueryParams(c *fiber.Ctx) error {
	query := new(RequestQuery)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid query parameters. Please ensure they are correctly formatted.",
		})
	}

	var errorMessages = map[string]string{
		"Rid":    "'rid' (rule id) is mandatory",
		"Factor": fmt.Sprintf("'Factor' must be a number between %d and %d", constant.MinDPOFactorValue, constant.MaxDPOFactorValue),
	}

	if err := validations.Validator.Struct(query); err != nil {
		errorResponse := map[string]string{
			"status": "error",
		}
		for _, err := range err.(validator.ValidationErrors) {
			if msg, ok := errorMessages[err.Field()]; ok {
				errorResponse["message"] = msg
			} else {
				errorResponse["message"] = fmt.Sprintf("%s is invalid", err.Field())
			}
			break
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	return c.Next()
}
