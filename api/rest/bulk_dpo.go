package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// DemandPartnerOptimizationBulkPostHandler Update and DPO in bulks
// @Description Update DPO setup in bulk (publisher, factor and demand_partner_id fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []dto.DPORuleUpdateRequest true "DPO update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /bulk/dpo [post]
func (o *OMSNewPlatform) DemandPartnerOptimizationBulkPostHandler(c *fiber.Ctx) error {
	var requests []dto.DPORuleUpdateRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse metadata for DPO bulk", err)
	}

	if err := o.bulkService.BulkInsertDPO(c.Context(), requests); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process dpo_rule bulk updates", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Dpo_rule with Metadata_queue bulk update successfully processed")
}
