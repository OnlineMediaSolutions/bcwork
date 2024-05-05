package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"time"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type MetadataUpdateRequest struct {
	Key     string      `json:"key"`
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

// MetadataUpdateRespose
type MetadataUpdateRespose struct {
	// in: body
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

// MetadataPostHandler Update new bidder metadata.
// @Description Update new bidder metadata.
// @Tags metadata
// @Accept json
// @Produce json
// @Param options body MetadataUpdateRequest true "Metadata update Options"
// @Success 200 {object} MetadataUpdateRespose
// @Security ApiKeyAuth
// @Router /metadata/update [post]
func MetadataPostHandler(c *fiber.Ctx) error {

	data := &MetadataUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to parse metadata update payload")

		return c.SendStatus(http.StatusBadRequest)
	}

	if data.Key == "" {
		log.Error().Str("body", string(c.Body())).Msg("empty key on metadata update request")
		return c.SendStatus(http.StatusBadRequest)
	}

	//log.Info().Interface("update", data).Msg("metadata update parsed")

	value, err := json.Marshal(data.Data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to marshal metadata update payload")
		return c.SendStatus(http.StatusBadRequest)
	}

	now := time.Now()
	mod := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(data.Key, now),
		Key:           data.Key,
		Value:         value,
		CreatedAt:     now,
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(MetadataUpdateRespose{
		Status:        "ok",
		TransactionID: mod.TransactionID,
	})
}
