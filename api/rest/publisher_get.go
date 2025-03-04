package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
)

// PublisherCountHandler Count publishers
// @Summary Count publishers
// @Tags publisher
// @Produce json
// @Param options body core.GetPublisherOptions true "options"
// @Success 200 {object} dto.PublisherSlice
// @Security ApiKeyAuth
// @Router /publisher/get [post]
func (o *OMSNewPlatform) PublisherGetHandler(c *fiber.Ctx) error {
	data := &core.GetPublisherOptions{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(500).JSON(Response{Status: "error", Message: "error when parsing request body"})
	}

	pubs, err := o.publisherService.GetPublisher(c.Context(), data)
	if err != nil {
		return c.Status(400).JSON(Response{Status: "error", Message: "failed to retrieve publishers"})
	}

	return c.JSON(pubs)
}
