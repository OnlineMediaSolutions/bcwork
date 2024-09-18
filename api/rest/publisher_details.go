package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// PublisherDetailsGetHandler Get Publisher with information about domains and factors setup
// @Description Get Publishers with information about domains and factors setup
// @Tags Publisher Domain Factor
// @Accept json
// @Produce json
// @Param options body core.GetPublisherDetailsOptions true "options"
// @Success 200 {object} core.PublisherDetailsSlice
// @Router /publisher/details/get [post]
func PublisherDetailsGetHandler(c *fiber.Ctx) error {
	data := &core.GetPublisherDetailsOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body for publisher domain parsing error", err)
	}

	pubs, err := core.GetPublisherDetails(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve publisher domain", err)
	}
	return c.JSON(pubs)
}
