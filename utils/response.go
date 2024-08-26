package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DpoResponse struct {
	BaseResponse
	RuleId string `json:"rule_id"`
}

type ErrorMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorResponse(c *fiber.Ctx, statusCode int, customMessage string, errorMessage error) error {
	resp := ErrorMessage{
		Status:  "error",
		Message: customMessage,
		Error:   fmt.Sprintf("%s", errorMessage),
	}
	return c.Status(statusCode).JSON(resp)
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string) error {
	resp := BaseResponse{
		Status:  "success",
		Message: message,
	}
	return c.Status(statusCode).JSON(resp)
}

func DpoSuccessResponse(c *fiber.Ctx, statusCode int, ruleId string, message string) error {
	resp := DpoResponse{
		BaseResponse: BaseResponse{
			Status:  "success",
			Message: message,
		},
		RuleId: ruleId,
	}
	return c.Status(statusCode).JSON(resp)
}
