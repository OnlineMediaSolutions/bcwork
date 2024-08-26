package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

type FactorUpdateResponse struct {
	Status string `json:"status"`
}

// FactorGetHandler Get factor setup
// @Description Get factor setup
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body core.GetFactorOptions true "options"
// @Success 200 {object} core.FactorSlice
// @Router /factor/get [post]
func FactorGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := core.GetFactors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve factors", err)
	}
	return c.JSON(pubs)
}

// FactorPostHandler Update and enable Factor setup
// @Description Update Factor setup (publisher is mandatory, domain is mandatory)
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body core.FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /factor [post]
func FactorPostHandler(c *fiber.Ctx) error {
	data := &core.FactorUpdateRequest{}

	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factor payload parsing error", err)
	}

	err = core.UpdateFactor(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Factor table", err)
	}

	err = core.UpdateMetaData(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for factor", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Factor and Metadata tables successfully updated")
}
