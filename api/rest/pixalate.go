package rest

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// ConfiantPostHandler Update and enable Pixalate setup
// @Description Update and enable Pixalate setup (publisher is mandatory, domain is optional)
// @Tags Pixalate
// @Accept json
// @Produce json
// @Param options body core.PixalateUpdateRequest true "Pixalate update Options"
// @Success 200 {object} utils.Response
// @Security ApiKeyAuth
// @Router /pixalate [post]
func PixalatePostHandler(c *fiber.Ctx) error {

	data := &core.PixalateUpdateRequest{}
	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Pixalate payload parsing error")
	}

	err = core.UpdateMetaDataQueueWithPixalate(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to update metadata table for Pixalate, %s", err))
	}

	err = core.UpdatePixalate(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Pixalate table")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Pixalate and Metadata tables successfully updated")
}
