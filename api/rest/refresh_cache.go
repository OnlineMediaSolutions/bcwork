package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

type RefreshCacheUpdateResponse struct {
	Status string `json:"status"`
}

// RefreshCacheGetAllHandler Get looping ratio setup
// @Description Get refresh cache setup
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Param options body core.GetRefreshCacheOptions true "options"
// @Success 200 {object} core.RefreshCache
// @Security ApiKeyAuth
// @Router /refresh_cache/get [post]
func (o *OMSNewPlatform) RefreshCacheGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetRefreshCacheOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	pubs, err := o.refreshCacheService.GetRefreshCache(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve refresh cache", err)
	}
	return c.JSON(pubs)
}

// RefreshCachePostHandler Update and enable Refresh Cache setup
// @Description Update Refresh Cache setup
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Param options body dto.RefreshCacheUpdateRequest true "Refresh Cache update Options"
// @Success 200 {object} RefreshCacheUpdateResponse
// @Security ApiKeyAuth
// @Router /refresh_cache [post]
func (o *OMSNewPlatform) RefreshCachePostHandler(c *fiber.Ctx) error {
	data := &dto.RefreshCacheUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh cache payload parsing error", err)
	}

	isInsert, err := o.refreshCacheService.UpdateRefreshCache(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Refresh cache table", err)
	}

	err = o.refreshCacheService.UpdateMetaData(c.Context(), *data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata table for refresh cache", err)
	}

	responseMessage := "Refresh cache successfully updated"
	if isInsert {
		responseMessage = "Refresh cache successfully created"
	}

	return utils.SuccessResponse(c, fiber.StatusOK, responseMessage)
}
