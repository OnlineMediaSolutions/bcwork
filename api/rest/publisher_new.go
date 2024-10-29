package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/rotisserie/eris"
)

// PublisherNewResponse result of the request
type PublisherNewResponse struct {
	PublisherID string `json:"publisher_id"`
	Status      string `json:"status"`
}

// PublisherNewHandler Create a publisher
// @Description Create a publisher
// @Tags publisher
// @Produce json
// @Param options body core.PublisherCreateValues true "create publisher values"
// @Success 200 {object} PublisherNewResponse
// @Security ApiKeyAuth
// @Router /publisher/new [post]
func PublisherNewHandler(ctx *fiber.Ctx) error {

	data := &core.PublisherCreateValues{}
	if err := ctx.BodyParser(&data); err != nil {
		return eris.Wrap(err, "error when parsing request body")
	}

	publisherID, err := core.CreatePublisher(ctx.Context(), *data)
	if err != nil {
		return eris.Wrap(err, "failed to create publisher")
	}

	return ctx.JSON(PublisherNewResponse{Status: "success", PublisherID: publisherID})
}
