package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
)

// MetadataUpdateRespose
type MetadataUpdateRespose struct {
	// in: body
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// MetadataPostHandler Update new bidder metadata.
// @Description Update new bidder metadata.
// @Tags MetaData
// @Accept json
// @Produce json
// @Param options body core.MetadataUpdateRequest true "Metadata update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /metadata/update [post]
func MetadataPostHandler(c *fiber.Ctx) error {
	data := &core.MetadataUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse metadata update payload", err)
	}

	value, err := createValue(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create json for adstxt", err)
	}

	now := time.Now()
	err = core.InsertDataToMetaData(c, *data, value, now)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Metadata_queue table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "metadata queue was updated successfully")
}

func createValue(c *fiber.Ctx, data *core.MetadataUpdateRequest) ([]byte, error) {
	if data.Key == "" {
		log.Error().Str("body", string(c.Body())).Msg("empty key on metadata update request")
		return nil, c.SendStatus(http.StatusBadRequest)
	}

	value, err := json.Marshal(data.Data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to marshal metadata update payload")
		return nil, c.SendStatus(http.StatusBadRequest)
	}

	return value, nil
}
