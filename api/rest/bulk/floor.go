package bulk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
	"net/http"
)

type FloorUpdateResponse struct {
	Status string `json:"status"`
}

// FloorBulkPostHandler Update and enable Bulk insert Floor setup
// @Description Update Floor setup in bulk (publisher, floor, device and country fields are mandatory)
// @Tags Bulk
// @Accept json
// @Produce json
// @Param options body []constant.FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} FloorUpdateResponse
// @Security ApiKeyAuth
// @Router /bulk/floor [post]
func FloorBulkPostHandler(c *fiber.Ctx) error {
	var requests []constant.FloorUpdateRequest

	if err := c.BodyParser(&requests); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Error parsing request body for floor bulk update", err)
	}

	chunks, err := bulk.MakeChunksFloor(requests)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create chunks for bulk Floor", err)
	}

	if err := bulk.ProcessChunksFloor(c, chunks); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process bulk Floor updates", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Floor bulk update successfully processed")
}
