package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// PublisherDetailsGetHandler Get Publishers with information about domains setup
// @Description Get Publishers with information about domains setup
// @Tags publisher
// @Accept json
// @Produce json
// @Param options body core.GetPublisherDetailsOptions true "options"
// @Success 200 {object} core.PublisherDetailsSlice
// @Security ApiKeyAuth
// @Router /publisher/details/get [post]
func (o *OMSNewPlatform) PublisherDetailsGetHandler(c *fiber.Ctx) error {
	data := &core.GetPublisherDetailsOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body for publisher details parsing error", err)
	}

	pubs, err := o.publisherService.GetPublisherDetails(c.Context(), data)
	activityStatus, err := o.publisherService.GetPubImpsPerPublisherDomain(c.Context(), data)

	print(activityStatus)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve publisher details", err)
	}
	return c.JSON(pubs)
}
