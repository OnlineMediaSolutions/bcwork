package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

type CompetitorUpdateResponse struct {
	Status string `json:"status"`
}

// CompetitorGetHandler Get Competitor setup
// @Description Get Competitor setup
// @Tags Competitor
// @Accept json
// @Produce json
// @Param options body core.GetCompetitorOptions true "options"
// @Success 200 {object} core.CompetitorSlice
// @Security ApiKeyAuth
// @Router /competitor/get [post]
func CompetitorGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetCompetitorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := core.GetCompetitors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve  competitors", err)
	}

	return c.JSON(pubs)
}

// CompetitorPostHandler Update and enable Competitor setup
// @Description Update Competitor setup (name is mandatory, url is mandatory)
// @Tags Competitor
// @Accept json
// @Produce json
// @Param options body core.CompetitorUpdateRequest true "Competitor update Options"
// @Success 200 {object} CompetitorUpdateResponse
// @Security ApiKeyAuth
// @Router /competitor [post]
func CompetitorPostHandler(c *fiber.Ctx) error {
	data := &core.CompetitorUpdateRequest{}

	err := c.BodyParser(&data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Competitor payload parsing error", err)
	}

	err = core.UpdateCompetitor(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Competitor table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Competitor  tables successfully updated")
}
