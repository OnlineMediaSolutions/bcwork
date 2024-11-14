package rest

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	const digitalOceanPingUrl = "http://cloud.digitalocean.com/"
	logSizeLimit := viper.GetInt(config.LogSizeLimitKey)

	start := time.Now()

	requestID := bcguid.NewFromf(time.Now())
	url := c.Request().URI().String()
	c.Locals(constant.RequestIDContextKey, requestID)
	logger := log.Logger.With().
		Str(constant.RequestIDContextKey, requestID).
		Str("method", string(c.Request().Header.Method())).
		Str("url", url).
		Caller().
		Logger()
	c.Locals(constant.LoggerContextKey, &logger)
	err := c.Next()
	if err != nil {
		return err
	}

	// inner checks from digitalocean
	if url == digitalOceanPingUrl {
		return nil
	}

	// don't log responses of getting values from config manager
	if strings.HasSuffix(url, "/config/get") {
		return nil
	}

	reqSize := len(c.Request().Body())
	respSize := len(c.Response().Body())
	duration := time.Since(start)

	resultLogger := logger.With().Logger()
	if reqSize+respSize <= logSizeLimit {
		resultLogger = logger.With().
			Str("request", string(c.Request().Body())).
			Str("response", string(c.Response().Body())).
			Logger()
	}

	resultLogger.Info().
		Int("request_size", reqSize).
		Int("response_size", respSize).
		Interface("duration", duration).
		Str("duration_readable", duration.String()).
		Msg("logging middleware")

	return nil
}

// func ProfileMiddleware(c *fiber.Ctx) error {
// 	f, err := os.Create("cpu.prof")
// 	if err != nil {
// 		return err
// 	}

// 	// runtime.SetCPUProfileRate(1000)

// 	err = pprof.StartCPUProfile(f)
// 	if err != nil {
// 		return err
// 	}
// 	defer pprof.StopCPUProfile()

// 	err = c.Next()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
