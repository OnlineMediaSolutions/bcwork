package bulk

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
)

func DemandPartnerOptimizationBulkPostHandler(c *fiber.Ctx) error {
	var data []core.DemandPartnerOptimizationUpdateRequest

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse metadata update payload")
	}

	var ruleIDs []string

	for _, request := range data {
		dpoRule := core.DemandPartnerOptimizationRule{
			DemandPartner: request.DemandPartner,
			Publisher:     request.Publisher,
			Domain:        request.Domain,
			Country:       request.Country,
			OS:            request.OS,
			DeviceType:    request.DeviceType,
			PlacementType: request.PlacementType,
			Browser:       request.Browser,
			Factor:        request.Factor,
		}

		ruleID, err := dpoRule.Save(c.Context())
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save rule")
		}

		ruleIDs = append(ruleIDs, ruleID)

		go func(demandPartner string) {
			err := core.SendToRT(context.Background(), demandPartner)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update RT metadata for dpo")
			}
		}(request.DemandPartner)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fmt.Sprintf("rule_ids: %v", ruleIDs))
}
