package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// ConfiantGetHandler Get confiant setup
// @Description Get confiant setup
// @Tags Confiant
// @Accept json
// @Produce json
// @Param options body core.GetConfiantOptions true "options"
// @Success 200 {object} dto.ConfiantSlice
// @Security ApiKeyAuth
// @Router /confiant/get [post]
func (o *OMSNewPlatform) ConfiantGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetConfiantOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.confiantService.GetConfiants(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve confiants", err)
	}

	return c.JSON(pubs)
}

// ConfiantPostHandler Update and enable Confiant setup
// @Description Update and enable Confiant setup (publisher is mandatory, domain is optional)
// @Tags Confiant
// @Accept json
// @Produce json
// @Param options body dto.ConfiantUpdateRequest true "Confiant update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /confiant [post]
func (o *OMSNewPlatform) ConfiantPostHandler(c *fiber.Ctx) error {
	data := &dto.ConfiantUpdateRequest{}
	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Confiant payload parsing error", err)
	}

	err = o.confiantService.UpdateMetaDataQueue(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for confiant", err)
	}

	err = o.confiantService.UpdateConfiant(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Confiant table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Confiant and Metadata tables successfully updated")
}
