package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
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
// @Success 200 {object} core.FloorSlice
// @Router /floor/get [post]
func FloorGetAllHandler(c *fiber.Ctx) error {

	data := &core.GetFloorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error")
	}

	pubs, err := core.GetFloors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve floors")
	}
	return c.JSON(pubs)
}

// FloorPostHandler Update and enable Floor setup
// @Description Update Floor setup (publisher, floor, device and country fields are mandatory)
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body core.FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /floor [post]
func FloorPostHandler(c *fiber.Ctx) error {
	data := &core.FloorUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Floor payload parsing error %s", err))
	}

	err = core.UpdateFloors(c, data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update Floor table %s", err))
	}

	err = core.UpdateFloorMetaData(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to update metadata table for floor, %s", err))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Floor and Metadata tables successfully updated")
}
