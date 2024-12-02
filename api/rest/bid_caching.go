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

// BidCachingGetAllHandler Get bid_caching setup
// @Description Get bid_caching setup
// @Tags BidCaching
// @Accept json
// @Produce json
// @Param options body core.GetBidCachingOptions true "options"
// @Success 200 {object} core.BidCaching
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

// BidCachingPostHandler Update and enable BidCaching setup
// @Description Update BidCaching setup
// @Tags BidCaching
// @Accept json
// @Produce json
// @Param options body dto.BidCachingUpdateRequest true "BidCaching update Options"
// @Success 200 {object} BidCachingUpdateResponse
// @Security ApiKeyAuth
// @Router /bid_caching [post]
func (o *OMSNewPlatform) BidCachingPostHandler(c *fiber.Ctx) error {
	data := &dto.BidCachingUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Bid Caching payload parsing error", err)
	}

	isInsert, err := o.bidCachingService.UpdateBidCaching(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Bid Caching table", err)
	}

	err = core.UpdateBidCachingMetaData(*data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for bid caching", err)
	}

	responseMessage := "Bid Caching successfully updated"
	if isInsert {
		responseMessage = "Bid Caching successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}
