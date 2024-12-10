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

// RefreshCacheGetAllHandler Get refresh cache setup
// @Description Get refresh cache setup
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Param options body core.GetRefreshCacheOptions true "options"
// @Success 200 {object} dto.RefreshCache
// @Security ApiKeyAuth
// @Router /refresh_cache/get [post]
func (o *OMSNewPlatform) RefreshCacheGetAllHandler(c *fiber.Ctx) error {
	data := &core.GetRefreshCacheOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error", err)
	}

	rc, err := o.refreshCacheService.GetRefreshCache(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve refresh cache", err)
	}
	return c.JSON(rc)
}

// RefreshCacheSetHandler Create Refresh Cache setup
// @Description Create Refresh Cache setup
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Param options body dto.RefreshCacheUpdateRequest true "Refresh Cache create Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /refresh_cache/set [post]
func (o *OMSNewPlatform) RefreshCacheSetHandler(c *fiber.Ctx) error {
	data := &dto.RefreshCacheUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh cache payload parsing error", err)
	}

	err = o.refreshCacheService.CreateRefreshCache(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Refresh cache table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Refresh cache successfully created")
}

// RefreshCacheUpdateHandler Update Refresh Cache setup
// @Description Update Refresh Cache setup
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Param options body dto.RefreshCacheUpdRequest true "Refresh Cache update Options"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /refresh_cache/update [post]
func (o *OMSNewPlatform) RefreshCacheUpdateHandler(c *fiber.Ctx) error {
	data := &dto.RefreshCacheUpdRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Refresh cache payload parsing error", err)
	}

	err = o.refreshCacheService.UpdateRefreshCache(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Refresh Cache table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Refresh Cache successfully updated")
}

// RefreshCacheDeleteHandler Delete refresh cache.
// @Description Delete refresh  cache.
// @Tags RefreshCache
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param options body []string true "options"
// @Router /refresh_cache/delete [delete]
func (o *OMSNewPlatform) RefreshCacheDeleteHandler(c *fiber.Ctx) error {
	var refreshCache []string

	err := c.BodyParser(&refreshCache)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to parse array of refresh cache  to delete", err)
	}

	err = o.refreshCacheService.DeleteRefreshCache(c.Context(), refreshCache)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete from  Refresh cache table", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Refresh cache successfully deleted")
}
