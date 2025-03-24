package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// DpAPIGetDemandPartners Get dpi-report demand partners
// @Description Get dpi-report demand partners
// @Tags Dp API
// @Accept json
// @Produce json
// @Success 200 {object} dto.DpApiSlice
// @Security ApiKeyAuth
// @Router /dp-api-report/demandPartners [post]
func (o *OMSNewPlatform) DpAPIGetDemandPartners(c *fiber.Ctx) error {
	dps, err := o.dpApiService.GetDemandPartners(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve dp api report demand partners", err)
	}

	return c.JSON(dps)
}

// GetDpApiReport Get dpi-report  demand partners report by date range
// @Description Get dpi-report  demand partners report by date range
// @Tags Dp API
// @Accept json
// @Produce json
// @Param options body core.GetDPApiOptions true "options"
// @Success 200 {object} dto.DpApiSlice
// @Security ApiKeyAuth
// @Router /dp-api-report [post]
func (o *OMSNewPlatform) GetDpApiReport(c *fiber.Ctx) error {
	data := &core.GetDPApiOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error for dp-api report", err)
	}

	dps, err := o.dpApiService.GetDpApiReport(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve dp api report by date range", err)
	}

	return c.JSON(dps)
}
