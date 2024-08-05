package bulk

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
)

func DemandPartnerOptimizationBulkInsertHandler(c *fiber.Ctx) error {
	var requests []core.DemandPartnerOptimizationUpdateRequest
	err := c.BodyParser(&requests)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse metadata update payload")
	}

	var ruleIDs []string
	var errors []error

	for _, data := range requests {
		dpoRule := core.DemandPartnerOptimizationRule{
			DemandPartner: data.DemandPartner,
			Publisher:     data.Publisher,
			Domain:        data.Domain,
			Country:       data.Country,
			OS:            data.OS,
			DeviceType:    data.DeviceType,
			PlacementType: data.PlacementType,
			Browser:       data.Browser,
			Factor:        data.Factor,
		}

		ruleID, err := dpoRule.Save(c.Context())
		if err != nil {
			errors = append(errors, err)
			continue
		}
		ruleIDs = append(ruleIDs, ruleID)

		go func(demandPartner string) {
			err := core.SendToRT(context.Background(), demandPartner)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
			}
		}(data.DemandPartner)
	}

	// Check if there were any errors during the insert
	if len(errors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Some records failed to insert")
	}

	// Return success with the list of rule IDs
	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("Inserted rule_ids: %v", ruleIDs))
}
