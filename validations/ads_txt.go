package validations

import "github.com/gofiber/fiber/v2"

// TODO: add validation
func AdsTxtValidation(c *fiber.Ctx) error {
	return c.Next()
}
