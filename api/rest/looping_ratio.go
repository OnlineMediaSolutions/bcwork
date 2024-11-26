package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type LoopingRatioUpdateResponse struct {
	Status string `json:"status"`
}

// LoopingRatioGetAllHandler Get looping ratio setup
// @Description Get looping ratio setup
// @Tags LoopingRatio
// @Accept json
// @Produce json
// @Param options body core.GetLoopingRatioOptions true "options"
// @Success 200 {object} core.LoopingRatio
// @Security ApiKeyAuth
// @Router /looping_ratio/get [post]
func (o *OMSNewPlatform) LoopingRatioGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetLoopingRatioOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.loopingRatioService.GetLoopingRatio(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve looping ratio", err)
	}
	return c.JSON(pubs)
}

// LoopingRatioPostHandler Update and enable Looping Ratio setup
// @Description Update Looping Ratio setup
// @Tags LoopingRatio
// @Accept json
// @Produce json
// @Param options body dto.LoopingRatioUpdateRequest true "Looping Ratio update Options"
// @Success 200 {object} LoopingRatioUpdateResponse
// @Security ApiKeyAuth
// @Router /looping_ratio [post]
func (o *OMSNewPlatform) LoopingRatioPostHandler(c *fiber.Ctx) error {
	data := &dto.LoopingRatioUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Looping Ratio payload parsing error", err)
	}

	isInsert, err := o.loopingRatioService.UpdateLoopingRatio(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Loopiong ratio table", err)
	}

	err = o.loopingRatioService.UpdateMetaData(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for looping ratio", err)
	}

	responseMessage := "Looping ratio successfully updated"
	if isInsert {
		responseMessage = "Looping ratio successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}
