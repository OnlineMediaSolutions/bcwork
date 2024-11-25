package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type BidCashingUpdateResponse struct {
	Status string `json:"status"`
}

// BidCashingGetAllHandler Get bid_cashing setup
// @Description Get bid_cashing setup
// @Tags BidCashing
// @Accept json
// @Produce json
// @Param options body core.GetBidCashingOptions true "options"
// @Success 200 {object} core.BidCashing
// @Security ApiKeyAuth
// @Router /bid_cashing/get [post]
func (o *OMSNewPlatform) BidCashingGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetBidCashingOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.bidCashingService.GetBidCashing(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve bid cashing", err)
	}
	return c.JSON(pubs)
}

// BidCashingPostHandler Update and enable BidCashing setup
// @Description Update BidCashing setup
// @Tags BidCashing
// @Accept json
// @Produce json
// @Param options body dto.BidCashingUpdateRequest true "BidCashing update Options"
// @Success 200 {object} BidCashingUpdateResponse
// @Security ApiKeyAuth
// @Router /bid_cashing [post]
func (o *OMSNewPlatform) BidCashingPostHandler(c *fiber.Ctx) error {
	data := &dto.BidCashingUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Bid Cashing payload parsing error", err)
	}

	isInsert, err := o.bidCashingService.UpdateBidCashing(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Bid Cashing table", err)
	}

	err = o.bidCashingService.UpdateMetaData(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for bid cashing", err)
	}

	responseMessage := "Bid Cashing successfully updated"
	if isInsert {
		responseMessage = "Bid Cashing successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}
