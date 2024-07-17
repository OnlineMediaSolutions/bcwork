package bulk

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/api/rest"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/rs/zerolog/log"
	"net/http"
)

type FactorUpdateResponse struct {
	Status string `json:"status"`
}

// FactorBulkPostHandler Update and enable Bulk insert Factor setup
// @Description Update Factor setup in bulk (publisher, factor, device and country fields are mandatory)
// @Tags Factor in Bulk
// @Accept json
// @Produce json
// @Param options body []FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /bulk/factor [post]
func FactorBulkPostHandler(c *fiber.Ctx) error {
	var requests []bulk.FactorUpdateRequest
	if err := c.BodyParser(&requests); err != nil {
		log.Error().Err(err).Msg("Error parsing request body for bulk update")
		return c.Status(http.StatusBadRequest).JSON(&rest.Response{Status: "error", Message: "error when parsing request body for bulk update"})
	}

	chunks, err := bulk.MakeChunks(requests)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create chunks for factor updates",
		})
	}

	for i, chunk := range chunks {
		if err := bulk.InsertChunk(c, chunk); err != nil {
			log.Error().Err(err).Msgf("Failed to process bulk update for chunk %d", i)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process factor updates",
			})
		}
	}

	return c.Status(http.StatusOK).JSON(&rest.Response{Status: "ok", Message: "Bulk update successfully processed"})
}
