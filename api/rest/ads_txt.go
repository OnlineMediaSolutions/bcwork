package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// AdsTxtMainHandler Get ads.txt main table.
// @Description Get ads.txt main table.
// @Tags AdsTxt
// @Param options body core.AdsTxtOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AdsTxt
// @Security ApiKeyAuth
// @Router /ads_txt/main [post]
func (o *OMSNewPlatform) AdsTxtMainHandler(c *fiber.Ctx) error {
	data := &core.AdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting ads.txt main table data", err)
	}

	adsTxt, err := o.adsTxtService.GetMainAdsTxtTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve ads.txt main table data", err)
	}

	return c.JSON(adsTxt)
}

// AdsTxtGroupByDPHandler Get ads.txt group by DP table.
// @Description Get ads.txt group by DP table.
// @Tags AdsTxt
// @Param options body core.AdsTxtOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AdsTxtGroupedByDPData
// @Security ApiKeyAuth
// @Router /ads_txt/group_by_dp [post]
func (o *OMSNewPlatform) AdsTxtGroupByDPHandler(c *fiber.Ctx) error {
	data := &core.AdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting ads.txt group by dp table data", err)
	}

	adsTxt, err := o.adsTxtService.GetGroupByDPAdsTxtTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve ads.txt group by dp table data", err)
	}

	return c.JSON(adsTxt)
}

// AdsTxtAMHandler Get ads.txt AM table.
// @Description Get ads.txt AM table.
// @Tags AdsTxt
// @Param options body core.AdsTxtOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AdsTxt
// @Security ApiKeyAuth
// @Router /ads_txt/am [post]
func (o *OMSNewPlatform) AdsTxtAMHandler(c *fiber.Ctx) error {
	data := &core.AdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting ads.txt AM table data", err)
	}

	adsTxt, err := o.adsTxtService.GetAMAdsTxtTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve ads.txt AM table data data", err)
	}

	return c.JSON(adsTxt)
}

// AdsTxtCMHandler Get ads.txt CM table.
// @Description Get ads.txt CM table.
// @Tags AdsTxt
// @Param options body core.AdsTxtOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AdsTxt
// @Security ApiKeyAuth
// @Router /ads_txt/cm [post]
func (o *OMSNewPlatform) AdsTxtCMHandler(c *fiber.Ctx) error {
	data := &core.AdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting ads.txt CM table data", err)
	}

	adsTxt, err := o.adsTxtService.GetCMAdsTxtTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve ads.txt CM table data data", err)
	}

	return c.JSON(adsTxt)
}

// AdsTxtMBHandler Get ads.txt MB table.
// @Description Get ads.txt MB table.
// @Tags AdsTxt
// @Param options body core.AdsTxtOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AdsTxt
// @Security ApiKeyAuth
// @Router /ads_txt/mb [post]
func (o *OMSNewPlatform) AdsTxtMBHandler(c *fiber.Ctx) error {
	data := &core.AdsTxtOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting ads.txt MB table data", err)
	}

	adsTxt, err := o.adsTxtService.GetMBAdsTxtTable(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve ads.txt MB table data data", err)
	}

	return c.JSON(adsTxt)
}

// AdsTxtUpdateHandler Update ads.txt.
// @Description Update ads.txt.
// @Tags AdsTxt
// @Accept json
// @Produce json
// @Param adstxt body dto.AdsTxtUpdateRequest true "AdsTxt"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /ads_txt/update [post]
func (o *OMSNewPlatform) AdsTxtUpdateHandler(c *fiber.Ctx) error {
	data := &dto.AdsTxtUpdateRequest{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for updating ads.txt", err)
	}

	err := o.adsTxtService.UpdateAdsTxt(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update ads.txt", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "ads.txt successfully updated")
}
