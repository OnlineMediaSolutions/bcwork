package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	RuleId  string `json:"rule_id"`
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	resp := Response{
		Status:  "error",
		Message: message,
	}
	return c.Status(statusCode).JSON(resp)
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string) error {
	resp := Response{
		Status:  "success",
		Message: message,
	}
	return c.Status(statusCode).JSON(resp)
}

func CustomSuccessResponse(c *fiber.Ctx, statusCode int, ruleId string, message string) error {
	resp := Response{
		Status:  "success",
		RuleId:  ruleId,
		Message: message,
	}
	return c.Status(statusCode).JSON(resp)
}
