package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	const logSizeLimit = 250000

	start := time.Now()

	requestID := bcguid.NewFromf(time.Now())
	c.Locals(constant.RequestIDContextKey, requestID)

	logger := log.Logger.With().
		Str(constant.RequestIDContextKey, requestID).
		Str("method", string(c.Request().Header.Method())).
		Str("url", c.Request().URI().String()).
		Caller().
		Logger()
	c.Locals(constant.LoggerContextKey, &logger)

	err := c.Next()
	if err != nil {
		return err
	}

	reqSize := len(c.Request().Body())
	respSize := len(c.Response().Body())
	duration := time.Since(start)

	if reqSize+respSize <= logSizeLimit {
		logger.Info().
			Str("request", string(c.Request().Body())).
			Str("response", string(c.Response().Body())).
			Interface("duration", duration).
			Str("duration_readable", duration.String()).
			Msg("logging middleware")
	}

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
