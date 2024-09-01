package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Publisher struct {
	Name              string   `json:"name"  validate:"required"`
	AccountManagerID  string   `json:"account_manager_id"`
	MediaBuyerID      string   `json:"media_buyer_id"`
	CampaignManagerID string   `json:"campaign_manager_id"`
	OfficeLocation    string   `json:"office_location"`
	Status            string   `json:"status"`
	IntegrationType   []string `json:"integration_type"  validate:"integrationType"`
}

func PublisherValidation(c *fiber.Ctx) error {
	body := new(Publisher)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"integrationType": "integration type can be one or more than the following list: JS Tags (Compass), JS Tags (NP), Prebid.js, Prebid Server or oRTB EP",
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
