package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

type PublisherDomainRequest struct {
	DemandParnerId string `json:"demand_partner_id" validate:"required"`
	Data           []Data `json:"data"`
}

type Data struct {
	PubId        string `json:"pubId" `
	Domain       string `json:"domain" `
	AdsTxtStatus bool   `json:"ads_txt_status"`
}

const (
	ErrorStatus                 = "error"
	ValidationError             = "couldn't validate some of the requests"
	AdsTxtStatusValidationError = "adsTxtStatus is missing"
	DomainValidationError       = "Domain is missing"
	PublisherIdValidationError  = "PublisherId is missing"
)

func ValidateBulkPublisherDemands(c *fiber.Ctx) error {
	var requests PublisherDomainRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body for Publisher Demand. Please ensure it's a valid JSON.",
		})
	}

	validationErrors := validateBulkPublisherDemand(requests.Data)

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorBulkResponse{
			Status:  ErrorStatus,
			Message: ValidationError,
			Errors:  validationErrors,
		})
	}

	return c.Next()
}

func validateBulkPublisherDemand(requests []Data) map[string][]string {
	var errorMessages = map[string]string{}

	validationErrors := make(map[string][]string)
	for idx, request := range requests {
		key := fmt.Sprintf("request %v", idx+1)

		if request.PubId == "" {
			validationErrors[key] = append(validationErrors[key], PublisherIdValidationError)
		}
		if request.Domain == "" {
			validationErrors[key] = append(validationErrors[key], DomainValidationError)
		}
		if strconv.FormatBool(request.AdsTxtStatus) == "" {
			validationErrors[key] = append(validationErrors[key], AdsTxtStatusValidationError)
		}

		err := Validator.Struct(request)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				log.Println(err.Tag())
				if msg, ok := errorMessages[err.Tag()]; ok {
					validationErrors[key] = append(validationErrors[key], msg)
				} else {
					validationErrors[key] = append(validationErrors[key],
						fmt.Sprintf("%s is mandatory, validation failed", err.Field()))
				}
			}
		}
	}

	return validationErrors
}
