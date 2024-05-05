package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils/omserr"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Debug   interface{} `json:"debug,omitempty"`
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	errctx := omserr.ExtractErrContext(err)
	log.Error().Err(err).Interface("errctx", errctx).Interface("error_details", eris.ToJSON(err, true)).Msg("Rest Error")
	code := errctx.GetNumericValue("http", fiber.StatusInternalServerError)
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	msg := errctx.GetValue("msg")
	if msg == "" {
		if code >= 500 {
			msg = "Something went wrong"
		} else if code == 400 {
			msg = "Bad input"
		} else if code == 401 || code == 403 {
			msg = "Unauthorized"
		}
	}

	return ctx.Status(code).JSON(Response{
		Status: "error",
		//		Debug:   eris.ToJSON(err, true),
		Message: msg,
	})
}
