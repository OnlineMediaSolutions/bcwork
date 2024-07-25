package core

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type PixalateUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate"`
}

type PixalateUpdateRespose struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func UpdatePixalate(c *fiber.Ctx, data *PixalateUpdateRequest) error {

	updatedPixalate := models.Pixalate{
		PublisherID: data.Publisher,
		PixalateKey: data.Hash,
		Rate:        data.Rate,
		Domain:      data.Domain,
	}

	return updatedPixalate.Upsert(c.Context(), bcdb.DB(), true, []string{models.PixalateColumns.PublisherID, models.PixalateColumns.Domain}, boil.Infer(), boil.Infer())
}

func UpdateMetaDataQueueWithPixalate(c *fiber.Ctx, data *PixalateUpdateRequest) error {

	val, err := json.Marshal(data.Hash)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Pixalate failed to parse hash value")
	}

	mod := models.MetadataQueue{
		Key:           "pixalate:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata_queue with Pixalate")
	}
	return nil
}
