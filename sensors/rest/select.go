package rest

import (
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/sensors/core"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func SelectHandler(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		c.SendString("key is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	log.Info().Msg("loading hourly file ")
	hourly, err := core.LoadHourlySensors(time.Now().UTC())
	if err != nil {
		return errors.Wrapf(err, "failed to load hourly file")
	}

	log.Info().Str("hour", hourly.Hour).Int("len", len(hourly.Sensors)).Msg("hourly file loaded")
	table := core.SensorTableFromFile(hourly)

	res := table.Select(key, 10)

	c.Response().Header.Set("Content-Type", "application/json")
	return c.JSON(res)
}

func SumCountHandler(c *fiber.Ctx) error {
	var err error
	key := c.Query("key")
	if key == "" {
		c.SendString("key is mandatory")
		return c.SendStatus(http.StatusBadRequest)
	}

	t := time.Now().UTC()
	hour := c.Query("hour")
	if hour != "" {
		t, err = time.Parse("2006010215", hour)
		if err != nil {
			return errors.Wrapf(err, "failed to parse hour(hour:%s)", hour)
		}
	}

	log.Info().Msg("loading hourly file ")
	hourly, err := core.LoadHourlySensors(t)
	if err != nil {
		return errors.Wrapf(err, "failed to load hourly file")
	}

	log.Info().Str("hour", hourly.Hour).Int("len", len(hourly.Sensors)).Msg("hourly file loaded")
	table := core.SensorTableFromFile(hourly)

	res := table.SumCount(key)

	c.Response().Header.Set("Content-Type", "application/json")
	return c.JSON(res)
}
