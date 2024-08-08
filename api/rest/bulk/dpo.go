package bulk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
)

// DemandPartnerOptimizationBulkPostHandler Update and DPO in bulks
// @Description Update DPO setup in bulk (publisher, factor and demand_partner_id fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []core.DPOUpdateRequest true "DPO update Options"
// @Success 200 {object} DPOUpdateResponse
// @Security ApiKeyAuth
// @Router /bulk/dpo [post]
func DemandPartnerOptimizationBulkPostHandler(ctx *fiber.Ctx) error {
	var data []core.DPOUpdateRequest

	err := ctx.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse metadata for DPO bulk")
	}

	chunksDPO, err := bulk.MakeChunksDPO(data)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create chunks for DPO")
	}

	if err := bulk.ProcessChunksDPO(ctx, chunksDPO); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to process dpo_rule bulk updates")
	}

	return utils.SuccessResponse(ctx, fiber.StatusOK, "DPO bulk update successfully processed")
}

type DPOUpdateResponse struct {
	Status string `json:"status"`
}
