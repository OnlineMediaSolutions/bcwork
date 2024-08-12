package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// PublisherDomainGetHandler Get Publisher Domain setup
// @Description Get Publisher Domain setup
// @Tags Publisher Domain
// @Accept json
// @Produce json
// @Param options body core.GetPublisherDomainOptions true "options"
// @Success 200 {object} core.PublisherDomainSlice
// @Router /publisher/domain/get [post]
func PublisherDomainGetHandler(c *fiber.Ctx) error {
	data := &core.GetPublisherDomainOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body for publisher domain parsing error")
	}

	pubs, err := core.GetPublisherDomain(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve factors")
	}
	return c.JSON(pubs)
}

// PublisherDomainPostHandler Update and enable Publisher Domain setup
// @Description Update and enable Publisher Domain setup (publisher is mandatory, domain is optional)
// @Tags Publisher Domain
// @Accept json
// @Produce json
// @Param options body core.PublisherDomainUpdateRequest true "Publishe Domain update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /publisher/domain [post]
func PublisherDomainPostHandler(c *fiber.Ctx) error {

	data := &core.PublisherDomainUpdateRequest{}
	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Publisher Domain payload parsing error")
	}

	err = core.UpdatePublisherDomain(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Publisher Domain table")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Publisher Domain table successfully updated")
}
