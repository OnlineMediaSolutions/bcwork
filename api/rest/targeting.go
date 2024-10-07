package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
)

// TargetingGetHandler Get targeting data.
// @Description Get targeting data.
// @Tags Targeting
// @Param options body core.TargetingOptions true "Options"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /targeting/get [post]
func TargetingGetHandler(c *fiber.Ctx) error {
	data := &core.TargetingOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting targeting data", err)
	}

	targeting, err := core.GetTargetings(c.Context(), data)
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
// @Param options body constant.Targeting true "Targeting"
// @Security ApiKeyAuth
// @Router /targeting/set [post]
func TargetingSetHandler(c *fiber.Ctx) error {
	data := &constant.Targeting{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for creating targeting", err)
	}

	err := core.CreateTargeting(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create targeting", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "targeting successfully added")
}

// TargetingUpdateHandler Update targeting in terms of cost model, value and status.
// @Description Update targeting in terms of cost model, value and status.
// @Tags Targeting
// @Accept json
// @Produce json
// @Param id query int true "Targeting ID"
// @Param options body constant.Targeting true "Targeting"
// @Security ApiKeyAuth
// @Router /targeting/update [post]
func TargetingUpdateHandler(c *fiber.Ctx) error {
	data := &constant.Targeting{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for updated targeting", err)
	}

	id := c.QueryInt("id", 0)
	data.ID = id

	err := core.UpdateTargeting(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update targeting", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "targeting successfully updated")
}
