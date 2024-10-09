package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DpoResponse struct {
	BaseResponse
	RuleId string `json:"rule_id"`
}

type TagsResponse struct {
	BaseResponse
	Tags string `json:"tags"`
}

type ErrorMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorResponse(c *fiber.Ctx, statusCode int, customMessage string, errorMessage error) error {
	resp := ErrorMessage{
		Status:  ResponseStatusError,
		Message: customMessage,
		Error:   fmt.Sprintf("%s", errorMessage),
	}
	return c.Status(statusCode).JSON(resp)
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string) error {
	resp := BaseResponse{
		Status:  ResponseStatusSuccess,
		Message: message,
	}
	return c.Status(statusCode).JSON(resp)
}

func DpoSuccessResponse(c *fiber.Ctx, statusCode int, ruleId string, message string) error {
	resp := DpoResponse{
		BaseResponse: BaseResponse{
			Status:  ResponseStatusSuccess,
			Message: message,
		},
		RuleId: ruleId,
	}
	return c.Status(statusCode).JSON(resp)
}
