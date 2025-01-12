package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// PublisherDemandGetHandler Get publisher/demand/domain setup
// @Description Get PublisherDemandResponse List
// @Tags Publisher Demand Domain
// @Accept json
// @Produce json
// @Param options body core.GetPublisherDemandOptions true "options"
// @Success 200 {object} core.PublisherDemandResponseSlice
// @Security ApiKeyAuth
// @Router /publisher/demand/get [post]
func PublisherDemandGetHandler(c *fiber.Ctx) error {

	data := &core.GetPublisherDemandOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "PublisherDemandResponse Request body parsing error", err)
	}

	publisherDemands, err := core.GetPublisherDemandData(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve PublisherDemandResponse", err)
	}

	return c.JSON(publisherDemands)
}

// PublisherDemandUpdate update publisher demand table
// @Description Get PublisherDemandResponse List
// @Tags Publisher Demand Domain
// @Accept json
// @Produce json
// @Param options body dto.PublisherDomainRequest true "options"
// @Success 200 {object} core.PublisherDemandResponseSlice
// @Security ApiKeyAuth
// @Router /publisher/demand/udpate [post]
func PublisherDemandUpdate(c *fiber.Ctx) error {
	data := &dto.PublisherDomainRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse publisher_demand payload", err)
	}

	now := time.Now()
	err := bulk.InsertDataToAdsTxt(c, *data, now)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update publisher_demand table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Publisher_demand bulk update successfully processed")
}
