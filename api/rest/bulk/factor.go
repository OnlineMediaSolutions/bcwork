package bulk

import (
	"github.com/m6yf/bcwork/utils/constant"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
)

type FactorUpdateResponse struct {
	Status string `json:"status"`
}

// FactorBulkPostHandler Update and enable Bulk insert Factor setup
// @Description Update Factor setup in bulk (publisher, factor, device and country fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /bulk/factor [post]
func FactorBulkPostHandler(c *fiber.Ctx) error {
	var requests []constant.FactorUpdateRequest
	if err := c.BodyParser(&requests); err != nil {
		log.Error().Err(err).Msg("error parsing request body for factor bulk update")
		return utils.ErrorResponse(c, http.StatusBadRequest, "error parsing request body for factor bulk update", err)
	}

	if err := bulk.BulkInsertFactors(c.Context(), requests); err != nil {
		log.Error().Err(err).Msg("failed to process bulk factor updates")
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to process bulk factor updates", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "factor bulk update successfully processed")
}
