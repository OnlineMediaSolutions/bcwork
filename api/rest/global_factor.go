package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// GlobalFactorGetHandler Get Global Factor Fees
// @Description Get confiant setup
// @Tags Global Factor
// @Accept json
// @Produce json
// @Param options body core.GetGlobalFactorOptions true "options"
// @Success 200 {object} core.GlobalFactorSlice
// @Router /global/factor/get [post]
func GlobalFactorGetHandler(c *fiber.Ctx) error {

	data := &core.GetGlobalFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Global Factor Request body parsing error", err)
	}

	pubs, err := core.GetGlobalFactor(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve confiants", err)
	}

	return c.JSON(pubs)
}

// GlobalFactorPostHandler Update Global Factor params
// @Description Update Global Factors Fees
// @Tags Global Factor
// @Accept json
// @Produce json
// @Param options body core.GlobalFactorRequest true "Global Factor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /global/factor [post]
func GlobalFactorPostHandler(c *fiber.Ctx) error {

	data := &core.GlobalFactorRequest{}
	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gobal Factor payload parsing error", err)
	}

	err = core.UpdateGlobalFactor(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Global Factor table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Global Factor table successfully updated")
}
