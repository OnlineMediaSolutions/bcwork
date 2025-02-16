package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
)

// FactorBulkPostHandler Update and enable Bulk insert Factor setup
// @Description Update Factor setup in bulk (publisher, factor, device and country fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []bulk.FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /bulk/factor [post]
func (o *OMSNewPlatform) FactorBulkPostHandler(c *fiber.Ctx) error {
	var requests []bulk.FactorUpdateRequest
	if err := c.BodyParser(&requests); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "error parsing request body for factor bulk update", err)
	}

	if err := o.bulkService.BulkInsertFactors(c.Context(), requests); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to process bulk factor updates", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "factor bulk update successfully processed")
}
