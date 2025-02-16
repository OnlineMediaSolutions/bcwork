package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// FactorAdjusterHandler Update Factors based on their Domain
// @Description Update Factor based on Domain
// @Tags Adjust
// @Accept json
// @Produce json
// @Param options body dto.AdjustRequest true "Factor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /adjust/factor [post]
func (o *OMSNewPlatform) FactorAdjusterHandler(c *fiber.Ctx) error {
	data := dto.AdjustRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Factor Adjust payload parsing error", err)
	}

	err = o.adjustService.AdjustFactors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed update factors", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Adjusted Factor values")
}

// FactorAdjusterHandler Update Floor based on their Domain
// @Description Update Floor based on Domain
// @Tags Adjust
// @Accept json
// @Produce json
// @Param options body dto.AdjustRequest true "Floor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /adjust/floor [post]
func (o *OMSNewPlatform) FloorAdjusterHandler(c *fiber.Ctx) error {
	data := dto.AdjustRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Floor Adjust payload parsing error", err)
	}

	err = o.adjustService.AdjustFloors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch Floor", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Adjusted Floor values")
}
