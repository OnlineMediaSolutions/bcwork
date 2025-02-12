package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
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
// @Param options body dto.PublisherCreateValues true "create publisher values"
// @Success 200 {object} PublisherNewResponse
// @Security ApiKeyAuth
// @Router /publisher/new [post]
func (o *OMSNewPlatform) PublisherNewHandler(c *fiber.Ctx) error {
	data := &dto.PublisherCreateValues{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "error when parsing request body", err)
	}

	publisherID, err := o.publisherService.CreatePublisher(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create publisher", err)
	}

	return c.JSON(PublisherNewResponse{Status: "success", PublisherID: publisherID})
}
