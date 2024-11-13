package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// GlobalSearchHandler Global search.
// @Description Search for publisher_id, publisher_name, domain, demand partner.
// @Tags Search
// @Param options body core.UserOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.SearchResults
// @Security ApiKeyAuth
// @Router /user/get [post]
func (o *OMSNewPlatform) GlobalSearchHandler(c *fiber.Ctx) error {
	data := &core.SearchRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for global search", err)
	}

	searchResults, err := o.searchService.Search(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve global search results", err)
	}

	return c.JSON(searchResults)
}
