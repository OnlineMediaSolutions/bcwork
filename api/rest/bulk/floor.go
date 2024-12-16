package bulk

import (
	"github.com/m6yf/bcwork/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
)

type FloorUpdateResponse struct {
	Status string `json:"status"`
}

// FloorBulkPostHandler Update and enable Bulk insert Floor setup
// @Description Update Floor setup in bulk (publisher, floor, device and country fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []dto.FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} FloorUpdateResponse
// @Security ApiKeyAuth
// @Router /bulk/floor [post]
func FloorBulkPostHandler(c *fiber.Ctx) error {
	var requests []dto.FloorUpdateRequest
	if err := c.BodyParser(&requests); err != nil {
		log.Error().Err(err).Msg("error parsing request body for floor bulk update")
		return utils.ErrorResponse(c, http.StatusBadRequest, "error parsing request body for floor bulk update", err)
	}

	if err := bulk.BulkInsertFloors(c.Context(), requests); err != nil {
		log.Error().Err(err).Msg("failed to process bulk floor updates")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to process bulk floor updates", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "floor bulk update successfully processed")
}
