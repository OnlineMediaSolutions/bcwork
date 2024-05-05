package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/rotisserie/eris"
)

// PublisherCountResponse list of proertiese returned
type PublisherCountResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// PublisherCountHandler Count publishers
// @Summary Count publishers
// @Tags publisher
// @Produce json
// @Param options body core.GetPublisherOptions true "options"
// @Success 200 {object} PublisherCountResponse
// @Router /publisher/count [post]
func PublisherCountHandler(c *fiber.Ctx) error {

	data := &core.GetPublisherOptions{}
	if err := c.BodyParser(&data); err != nil {
		return eris.Wrap(err, "error when parsing request body")
	}

	pubs, err := core.PublisherCount(c.Context(), &data.Filter)
	if err != nil {
		return eris.Wrapf(err, "failed to retrieve  bids")
	}

	return c.JSON(PublisherCountResponse{Status: "success", Count: pubs})
}
