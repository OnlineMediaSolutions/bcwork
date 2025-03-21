package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// DpAPIGetHandler Get demand partner api setup
// @Description Get dp api setup
// @Tags Dp API
// @Accept json
// @Produce json
// @Param options body core.GetDpApiOptions true "options"
// @Success 200 {object} dto.DpAPISlice
// @Security ApiKeyAuth
// @Router /dp_api/get [post]
func (o *OMSNewPlatform) DpAPIGetHandler(c *fiber.Ctx) error {
	data := &core.GetDPApiOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.dpApiService.GetDpApiReport(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve dp api report", err)
	}

	return c.JSON(pubs)
}
