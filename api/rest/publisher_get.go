package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/rotisserie/eris"
)

// PublisherBidCountResponse list of proertiese returned
type PublisherGetResponse struct {
	Status     string              `json:"status"`
	Publishers core.PublisherSlice `json:"publishers"`
}

// PublisherCountHandler Count publishers
// @Summary Count publishers
// @Tags publisher
// @Produce json
// @Param options body core.GetPublisherOptions true "options"
// @Success 200 {object} PublisherGetResponse
// @Router /publisher/get [post]
func PublisherGetHandler(c *fiber.Ctx) error {

	data := &core.GetPublisherOptions{}
	if err := c.BodyParser(&data); err != nil {
		return eris.Wrap(err, "error when parsing request body")
	}

	pubs, err := core.GetPublisher(c.Context(), data)
	if err != nil {
		return eris.Wrapf(err, "failed to retrieve  bids")
	}

	return c.JSON(PublisherGetResponse{Status: "success", Publishers: pubs})
}
