package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// AdsTxtGetHandler Get Pixalate setup
// @Description Get AdsTxt List
// @Tags AdsTxt
// @Accept json
// @Produce json
// @Param options body core.GetAdsTxtOptions true "options"
// @Success 200 {object} core.AdsTxtSlice
// @Router /ads/txt/get [post]
func AdsTxtGetHandler(c *fiber.Ctx) error {

	data := &core.GetAdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "AdsTxt Request body parsing error", err)
	}

	adsTxt, err := core.GetAdsTxt(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve AdsTxt", err)
	}

	return c.JSON(adsTxt)
}
