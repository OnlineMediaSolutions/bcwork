package rest

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type TargetingTagsResponse struct {
	utils.BaseResponse
	Tags []dto.Tags `json:"tags"`
}

// TargetingGetHandler Get targeting data.
// @Description Get targeting data.
// @Tags Targeting
// @Param options body core.TargetingOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.Targeting
// @Security ApiKeyAuth
// @Router /targeting/get [post]
func (o *OMSNewPlatform) TargetingGetHandler(c *fiber.Ctx) error {
	data := &core.TargetingOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting targeting data", err)
	}

	targeting, err := o.targetingService.GetTargetings(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve targeting data", err)
	}

	return c.JSON(targeting)
}

// TargetingSetHandler Create new targeting.
// @Description Create new targeting.
// @Tags Targeting
// @Accept json
// @Produce json
// @Param options body dto.Targeting true "Targeting"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /targeting/set [post]
func (o *OMSNewPlatform) TargetingSetHandler(c *fiber.Ctx) error {
	data := &dto.Targeting{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for creating targeting", err)
	}

	mod, err := o.targetingService.CreateTargeting(c.Context(), data)
	if err != nil {
		if errors.Is(err, dto.ErrFoundDuplicate) {
			return utils.ErrorFoundDuplicateResponse(c, "found duplicate while creating targeting", err, mod)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create targeting", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "targeting successfully added")
}

// TargetingUpdateHandler Update targeting in terms of cost model, value and status.
// @Description Update targeting in terms of cost model, value and status.
// @Tags Targeting
// @Accept json
// @Produce json
// @Param options body dto.Targeting true "Targeting"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /targeting/update [post]
func (o *OMSNewPlatform) TargetingUpdateHandler(c *fiber.Ctx) error {
	data := &dto.Targeting{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for updated targeting", err)
	}

	mod, err := o.targetingService.UpdateTargeting(c.Context(), data)
	if err != nil {
		if errors.Is(err, dto.ErrFoundDuplicate) {
			return utils.ErrorFoundDuplicateResponse(c, "found duplicate while updating targeting", err, mod)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update targeting", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "targeting successfully updated")
}

// TargetingExportTagsHandler Export one or multiple tags.
// @Description Export one or multiple tags.
// @Tags Targeting
// @Accept json
// @Produce json
// @Param request body core.ExportTagsRequest true "Export tags request"
// @Success 200 {object} TargetingTagsResponse
// @Security ApiKeyAuth
// @Router /targeting/tags [post]
func (o *OMSNewPlatform) TargetingExportTagsHandler(c *fiber.Ctx) error {
	data := &core.ExportTagsRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for export tags", err)
	}

	tags, err := o.targetingService.ExportTags(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to export tags", err)
	}

	if len(tags) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "failed to export tags", fmt.Errorf("no tags found for ids %v", data.IDs))
	}

	resp := TargetingTagsResponse{
		BaseResponse: utils.BaseResponse{
			Status:  utils.ResponseStatusSuccess,
			Message: "tags successfully exported",
		},
		Tags: tags,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
