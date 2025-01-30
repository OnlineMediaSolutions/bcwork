package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// PriceOverrideHandler set bid factor per user IP - valid for only 8 hours
// @Description set bid factor per user IP - valid for only 8 hours
// @Tags Price Override
// @Accept json
// @Produce json
// @Param options body dto.PriceOverrideRequest true "options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /price/override [post]
func PriceOverrideHandler(c *fiber.Ctx) error {

	data := &dto.PriceOverrideRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse price override payload", err)
	}

	err := core.UpdateMetaDataQueue(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update metadata_queue with new price by IPs", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "override prices successfully updated")
}
