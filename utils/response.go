package utils

import (
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

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	resp := BaseResponse{
		Status:  "error",
		Message: message,
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
