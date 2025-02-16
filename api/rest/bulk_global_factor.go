package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
)

type GlobalFactorUpdateResponse struct {
	Status string `json:"status"`
}

// GlobalFactorBulkPostHandler Update and enable Bulk insert Global Factor Fees
// @Description Update Global Factor Fees in bulk
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []bulk.GlobalFactorRequest true "Global Factor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /bulk/global/factor [post]
func (o *OMSNewPlatform) GlobalFactorBulkPostHandler(c *fiber.Ctx) error {
	var requests []bulk.GlobalFactorRequest
	if err := c.BodyParser(&requests); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "error parsing request body for global factor bulk update", err)
	}

	if err := o.bulkService.BulkInsertGlobalFactors(c.Context(), requests); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "failed to process bulk global factor updates", err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "global factor bulk update successfully processed")
}
