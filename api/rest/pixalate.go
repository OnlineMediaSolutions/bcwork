package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// PixalatePostHandler Update and enable Pixalate setup
// @Description Update and enable Pixalate setup (publisher is mandatory, domain is optional)
// @Tags Pixalate
// @Accept json
// @Produce json
// @Param options body core.PixalateUpdateRequest true "Pixalate update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /pixalate [post]
func (o *OMSNewPlatform) PixalatePostHandler(c *fiber.Ctx) error {
	data := &core.PixalateUpdateRequest{}
	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Pixalate payload parsing error", err)
	}

	err = o.pixalateService.UpdateMetaDataQueueWithPixalate(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for Pixalate", err)
	}

	err = o.pixalateService.UpdatePixalateTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Pixalate table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Pixalate and Metadata tables successfully updated")
}

// PixalateGetAllHandler Get Pixalate setup
// @Description Get Pixalate setup
// @Tags Pixalate
// @Accept json
// @Produce json
// @Param options body core.GetPixalateOptions true "options"
// @Success 200 {object} core.PixalateSlice
// @Security ApiKeyAuth
// @Router /pixalate/get [post]
func (o *OMSNewPlatform) PixalateGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetPixalateOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Pixalate Request body parsing error", err)
	}

	pubs, err := o.pixalateService.GetPixalate(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve pixalates", err)
	}

	return c.JSON(pubs)
}

// PixalateDeleteHandler Delete Pixalate.
// @Description Delete Pixalate - soft delete.
// @Tags Pixalate
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param options body []string true "options"
// @Router /pixalate/delete [delete]
func (o *OMSNewPlatform) PixalateDeleteHandler(c *fiber.Ctx) error {
	var pixalateKeys []string
	if err := c.BodyParser(&pixalateKeys); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of pixalate keys to delete", err)
	}

	err := o.pixalateService.SoftDeletePixalateInMetaData(c.Context(), pixalateKeys)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to soft delete metadata table for Pixalate", err)
	}

	err = o.pixalateService.SoftDeletePixalates(c.Context(), pixalateKeys)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error in soft delete pixalate", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Pixalates were deleted")
}
