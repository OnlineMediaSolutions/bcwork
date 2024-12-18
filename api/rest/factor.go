package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
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
// @Security ApiKeyAuth
// @Router /factor/get [post]
func (o *OMSNewPlatform) FactorGetAllHandler(c *fiber.Ctx) error {
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

// FactorPostHandler Update and enable Factor setup
// @Description Update Factor setup (publisher is mandatory, domain is mandatory)
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body constant.FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /factor [post]
func (o *OMSNewPlatform) FactorPostHandler(c *fiber.Ctx) error {
	data := &constant.FactorUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factor payload parsing error", err)
	}

	isInsert, err := o.factorService.UpdateFactor(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Factor table", err)
	}

	err = o.factorService.UpdateMetaData(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for factor", err)
	}

	responseMessage := "Factor successfully updated"
	if isInsert {
		responseMessage = "Factor successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}

// FactorPostHandler Soft deletes a factor
// @Description soft delete factor from Factor table and deletes it from metadata_queue table
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body []string true "options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /factor/delete [delete]
func (o *OMSNewPlatform) FactorDeleteHandler(c *fiber.Ctx) error {
	var rulesIds []string

	if err := c.BodyParser(&rulesIds); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of factor ruleIds to delete", err)
	}

	err := o.bulkService.BulkDeleteFactor(c.Context(), rulesIds)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete Factors", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "factors were deleted")
}
