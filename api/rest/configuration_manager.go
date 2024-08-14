package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// ConfigurationGetHandler Get factor setup
// @Description Get all configuration from DB
// @Tags Configuration
// @Accept json
// @Produce json
// @Param options body core.ConfigurationPayload true "options"
// @Success 200 {object} utils.BaseResponse
// @Router /config/get [post]
func ConfigurationGetHandler(c *fiber.Ctx) error {
	data := &core.ConfigurationPayload{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error parsing configuration request")
	}

	pubs, err := core.GetConfigurations(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve configurations")
	}
	return c.JSON(pubs)
}

// ConfigurationPostHandler Update and enable Factor setup
// @Description Update Factor setup (publisher is mandatory, domain is mandatory)
// @Tags Configuration
// @Accept json
// @Produce json
// @Param options body core.ConfigurationRequest true "Factor update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /config [post]
func ConfigurationPostHandler(c *fiber.Ctx) error {
	data := &core.ConfigurationRequest{}

	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Configuration payload parsing error")
	}

	err = core.UpdateConfiguration(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Configuration table")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Configuration table successfully updated")
}
