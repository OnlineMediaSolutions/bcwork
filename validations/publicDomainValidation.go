package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PublisherDomain struct {
	PublisherID     string   `json:"publisher_id" validate:"required"`
	Domain          string   `json:"domain" validate:"required"`
	GppTarget       *float64 `json:"gpp_target"`
	IntegrationType []string `json:"integration_type"`
	Automation      bool     `json:"automation"`
}

func PublisherDomainValidation(c *fiber.Ctx) error {
	body := new(PublisherDomain)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"integrationType":  "integration type can be one or more than the following list: JS Tags (Compass), JS Tags (NP), Prebid.js, Prebid Server or oRTB EP",
		"activeValidation": "active sholud be true or flase",
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
