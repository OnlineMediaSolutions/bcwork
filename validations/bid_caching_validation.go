package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/constant"
)

func ValidateBidCaching(c *fiber.Ctx) error {
	body := new(dto.BidCaching)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for bid caching. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateBidCache(body)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate bid cache",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func ValidateUpdateBidCaching(c *fiber.Ctx) error {
	body := new(dto.BidCachingUpdateRequest)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for bid caching update. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateBidCachingUpdateRequest(body)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate bid cache update",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateBidCache(request *dto.BidCaching) []string {
	var errorMessages = map[string]string{
		"country":                      "Country code must be 2 characters long and should be in the allowed list",
		"device":                       "Device should be in the allowed list",
		"bid_caching":                  fmt.Sprintf("Bid caching value not allowed, it should be >= %s", fmt.Sprintf("%d", constant.MinBidCachingValue)),
		bidCachingControlPercentageKey: bidCachingControlPercentageErrorMessage,
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

func validateBidCachingUpdateRequest(request *dto.BidCachingUpdateRequest) []string {
	var errorMessages = map[string]string{
		"bid_caching":                  fmt.Sprintf("Bid caching value not allowed, it should be >= %s", fmt.Sprintf("%d", constant.MinBidCachingValue)),
		bidCachingControlPercentageKey: bidCachingControlPercentageErrorMessage,
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
