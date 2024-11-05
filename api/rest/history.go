package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// HistoryGetHandler Get history data.
// @Description Get history data.
// @Tags History
// @Param options body core.HistoryOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.History
// @Security ApiKeyAuth
// @Router /history/get [post]
func (o *OMSNewPlatform) HistoryGetHandler(c *fiber.Ctx) error {
	data := &core.HistoryOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting history data", err)
	}

	history, err := o.historyService.GetHistory(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get history data", err)
	}

	return c.JSON(history)
}
