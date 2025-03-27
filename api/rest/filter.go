package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils"
)

// FiltersGetHandler Get filters fields names.
// @Description Get filters fields names.
// @Tags Filter
// @Param string query string true "Filter Name"
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Security ApiKeyAuth
// @Router /filter [get]
func (o *OMSNewPlatform) FiltersGetHandler(c *fiber.Ctx) error {
	filterName := c.Query("filter_name", "")
	if filterName == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to get filter name from query params", fmt.Errorf("failed to get filter name from query params"))
	}

	filterFields, err := o.filterService.GetFilterFields(c.Context(), filterName)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get filter fields", err)
	}

	return c.JSON(filterFields)
}
