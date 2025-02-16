package rest

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Description Check Health of Service
// @Tags Health
// @Accept json
// @Produce html
// @Router /ping [get]
func PingPong(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(Response{Status: "OK", Message: "Service is UP!!!"})
}
