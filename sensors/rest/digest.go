package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/sensors/core"
	"github.com/rs/zerolog/log"
)

func DigestSensorsHandler(c *fiber.Ctx) error {
	var body = make([]byte, len(c.Body()))
	copy(body, c.Body())
	core.SensorQ <- body
	log.Info().Int("body", len(body)).Msg("digest sensors")
	return nil
}
