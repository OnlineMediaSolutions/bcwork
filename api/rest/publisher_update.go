package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/rotisserie/eris"
)

// PublisherUpdateRequest Updates publisher fields (except roles which can be updated only by admin)
type PublisherUpdateRequest struct {
	PublisherID string                     `json:"publisher_id"`
	Options     core.UpdatePublisherValues `json:"updates"`
}

// PublisherUpdateResponse result of the request
type PublisherUpdateResponse struct {
	Status string `json:"status"`
}

// PublisherUpdateHandler Update a sser .
// @Description Updates publisher fields
// @Summary Update publisher.
// @Tags publisher
// @Produce json
// @Param options body PublisherUpdateRequest true "Publisher Update Options"
// @Success 200 {object} PublisherUpdateResponse
// @Security ApiKeyAuth
// @Router /publisher/update [post]
func PublisherUpdateHandler(c *fiber.Ctx) error {

	data := &PublisherUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		return eris.Wrap(err, "error when parsing request body")
	}

	err := core.UpdatePublisher(c.Context(), data.PublisherID, data.Options)
	if err != nil {
		return eris.Wrapf(err, "failed to update publisher fields")
	}

	return c.JSON(PublisherUpdateResponse{Status: "updated"})
}
