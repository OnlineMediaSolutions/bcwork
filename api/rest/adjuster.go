package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// FactorAdjusterHandler Update and enable Floor setup
// @Description Update Factor based on Domain
// @Tags Adjust
// @Accept json
// @Produce json
// @Param options body dto.AdjustRequest true "Floor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /adjust/factor [post]
func (o *OMSNewPlatform) FactorAdjusterHandler(c *fiber.Ctx) error {
	data := dto.AdjustRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factor Adjust payload parsing error", err)
	}

	factors, err := o.adjustService.GetFactors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch Factors", err)
	}

	err = o.adjustService.UpdateFactors(c.Context(), factors, data.Value, o.bulkFactorService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed update factors", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Adjusted Factor values")
}
