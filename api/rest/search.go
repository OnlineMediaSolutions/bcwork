package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// SearchHandler Global search.
// @Description Search for publisher_id, publisher_name, domain, demand partner name.
// @Tags Search
// @Param request body dto.SearchRequest true "Request"
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]dto.SearchResult
// @Security ApiKeyAuth
// @Router /search [post]
func (o *OMSNewPlatform) SearchHandler(c *fiber.Ctx) error {
	data := &dto.SearchRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse search request", err)
	}

	searchResults, err := o.searchService.Search(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve search results", err)
	}

	return c.JSON(searchResults)
}
