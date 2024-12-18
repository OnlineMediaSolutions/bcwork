package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

type BidCaching struct {
	Publisher     string `json:"publisher" validate:"required"`
	Device        string `json:"device" validate:"device"`
	Country       string `json:"country" validate:"country"`
	PlacementType string `json:"placement_type" validate:"placement_type"`
	OS            string `json:"os" validate:"os"`
	Browser       string `json:"browser" validate:"browser"`
	BidCaching    int16  `json:"bid_caching" validate:"bid_caching"`
	Domain        string `json:"domain"`
}

type BidCachingUpdate struct {
	BidCaching int16  `json:"bid_caching" validate:"bid_caching"`
	RuleId     string `json:"rule_id" validate:"required"`
}

func ValidateBidCaching(c *fiber.Ctx) error {
	body := new(BidCaching)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for bid caching. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"country":     "Country code must be 2 characters long and should be in the allowed list",
		"device":      "Device should be in the allowed list",
		"bid_caching": fmt.Sprintf("Bid caching  value not allowed, it should be >= %s", fmt.Sprintf("%d", constant.MinBidCachingValue)),
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

func ValidateUpdateBidCaching(c *fiber.Ctx) error {
	body := new(BidCachingUpdate)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for bid caching. Please ensure it's a valid JSON.",
		})
	}

	var errorMessages = map[string]string{
		"bid_caching": fmt.Sprintf("Bid caching  value not allowed, it should be >= %s", fmt.Sprintf("%d", constant.MinBidCachingValue)),
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
