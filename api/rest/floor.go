package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type FloorUpdateResponse struct {
	Status string `json:"status"`
}

// FloorGetHandler Get floor setup
// @Description Get floor setup
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body core.GetFloorOptions true "options"
// @Success 200 {object} dto.FloorSlice
// @Security ApiKeyAuth
// @Router /floor/get [post]
func (o *OMSNewPlatform) FloorGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetFloorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.floorService.GetFloors(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve floors", err)
	}

	return c.JSON(pubs)
}

// FloorPostHandler Update and enable Floor setup
// @Description Update Floor setup (publisher, floor, device and country fields are mandatory)
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body dto.FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /floor [post]
func (o *OMSNewPlatform) FloorPostHandler(c *fiber.Ctx) error {
	data := dto.FloorUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Floor payload parsing error", err)
	}

	isInsert, err := o.floorService.UpdateFloors(c.Context(), data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Floor table", err)
	}

	err = o.floorService.UpdateFloorMetaData(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for floor", err)
	}

	responseMessage := "Floor successfully updated"
	if isInsert {
		responseMessage = "Floor successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}

// FloorPostHandler Soft deletes floor
// @Description soft delete floor from Floor table and deletes it from metadata_queue table
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body []string true "options"
// @Success 200 {object} FloorUpdateResponse
// @Security ApiKeyAuth
// @Router /floor/delete [delete]
func (o *OMSNewPlatform) FloorDeleteHandler(c *fiber.Ctx) error {
	var rulesIds []string

	if err := c.BodyParser(&rulesIds); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of floor ruleIds to delete", err)
	}

	err := o.bulkService.BulkDeleteFloor(c.Context(), rulesIds)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete Floors", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "floors were deleted")
}
