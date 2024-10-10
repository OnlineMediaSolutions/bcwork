package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/constant"
)

type errorResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func ValidateTargeting(c *fiber.Ctx) error {
	var request *constant.Targeting
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Targeting. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateTargeting(request)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Status:  errorStatus,
			Message: "could not validate Targeting request",
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateTargeting(request *constant.Targeting) []string {
	var errorMessages = map[string]string{
		countriesValidationKey:           countryValidationErrorMessage,
		devicesValidationKey:             deviceValidationErrorMessage,
		targetingPriceModelValidationKey: targetingCostModelValidationErrorMessage,
		targetingStatusValidationKey:     targetingStatusValidationErrorMessage,
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

	if request.PriceModel == constant.TargetingPriceModelCPM &&
		(request.Value < constant.TargetingMinValueCostModelCPM || request.Value > constant.TargetingMaxValueCostModelCPM) {
		validationErrors = append(validationErrors,
			fmt.Sprintf("CPM Value should be between %v and %v",
				constant.TargetingMinValueCostModelCPM, constant.TargetingMaxValueCostModelCPM,
			),
		)
	}

	if request.PriceModel == constant.TargetingPriceModelRevShare &&
		(request.Value < constant.TargetingMinValueCostModelRevShare || request.Value > constant.TargetingMaxValueCostModelRevShare) {
		validationErrors = append(validationErrors,
			fmt.Sprintf("Rev Share Value should be between %v and %v",
				constant.TargetingMinValueCostModelRevShare, constant.TargetingMaxValueCostModelRevShare,
			),
		)
	}

	return validationErrors
}
