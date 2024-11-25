package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
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
	data := &core.GetFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.factorService.GetFactors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve factors", err)
	}
	return c.JSON(pubs)
}
