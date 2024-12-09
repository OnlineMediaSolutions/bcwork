package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type BidCachingUpdateResponse struct {
	Status string `json:"status"`
}

// BidCachingGetAllHandler Get bid caching setup
// @Description Get bid_caching setup
// @Tags BidCaching
// @Accept json
// @Produce json
// @Param options body core.GetBidCachingOptions true "options"
// @Success 200 {object} dto.BidCaching
// @Security ApiKeyAuth
// @Router /bid_caching/get [post]
func (o *OMSNewPlatform) BidCachingGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetBidCachingOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.bidCachingService.GetBidCaching(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve bid caching", err)
	}
	return c.JSON(pubs)
}

// BidCachingUpdateHandler Update BidCaching setup
// @Description Update BidCaching setup
// @Tags BidCaching
// @Accept json
// @Produce json
// @Param options body dto.BidCachingUpdRequest true "BidCaching update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /bid_caching/update [post]
func (o *OMSNewPlatform) BidCachingUpdateHandler(c *fiber.Ctx) error {
	data := &dto.BidCachingUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Bid Caching payload parsing error", err)
	}

	err = o.bidCachingService.UpdateBidCaching(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Bid Caching table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Bid Caching successfully updated")
}

// BidCachingNewHandler Create BidCaching setup
// @Description Create BidCaching setup
// @Tags BidCaching
// @Accept json
// @Produce json
// @Param options body dto.BidCachingUpdateRequest true "BidCaching create Options"
// @Success 200 {object} BidCachingUpdateResponse
// @Security ApiKeyAuth
// @Router /bid_caching/set [post]
func (o *OMSNewPlatform) BidCachingSetHandler(c *fiber.Ctx) error {
	data := &dto.BidCachingUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Bid Caching payload parsing error", err)
	}

	err = o.bidCachingService.CreateBidCaching(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create Bid Caching", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Bid Caching successfully created")
}

// BidCachingDeleteHandler Delete bid caching.
// @Description Delete bid chaching.
// @Tags BidCaching
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param options body []string true "options"
// @Router /bid_caching/delete [delete]
func (o *OMSNewPlatform) BidCachingDeleteHandler(c *fiber.Ctx) error {
	var bidCaching []string

	err := c.BodyParser(&bidCaching)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of bid caching  to delete", err)
	}

	err = o.bidCachingService.DeleteBidCaching(c.Context(), bidCaching)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete from Bid Caching table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Bid Caching successfully deleted")
}
