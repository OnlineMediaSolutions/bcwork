package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// DemandPartnerGetHandler Get demand partners.
// @Description Get demand partners.
// @Tags DemandPartner
// @Param options body core.DemandPartnerGetOptions true "options"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dp/get [post]
func (o *OMSNewPlatform) DemandPartnerGetHandler(c *fiber.Ctx) error {
	data := &core.DemandPartnerGetOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting demand partners data", err)
	}

	demandPartners, err := o.demandPartnerService.GetDemandPartners(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve demand partners data", err)
	}

	return c.JSON(demandPartners)
}

// DemandPartnerSetHandler Create new demand partner.
// @Description Create new demand partner.
// @Tags DemandPartner
// @Accept json
// @Produce json
// @Param user body dto.DemandPartner true "DemandPartner"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /dp/set [post]
func (o *OMSNewPlatform) DemandPartnerSetHandler(c *fiber.Ctx) error {
	data := &dto.DemandPartner{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for creating demand partner", err)
	}

	err := o.demandPartnerService.CreateDemandPartner(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create demand partner", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "demand partner successfully created")
}

// DemandPartnerUpdateHandler Update demand partner.
// @Description Update demand partner.
// @Tags DemandPartner
// @Accept json
// @Produce json
// @Param user body dto.DemandPartner true "DemandPartner"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /dp/update [post]
func (o *OMSNewPlatform) DemandPartnerUpdateHandler(c *fiber.Ctx) error {
	data := &dto.DemandPartner{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for updating demand partner", err)
	}

	err := o.demandPartnerService.UpdateDemandPartner(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update demand partner", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "demand partner successfully updated")
}

// DemandPartnerGetSeatOwnersHandler Get seat owners for demand partners.
// @Description Get seat owners for demand partners.
// @Tags DemandPartner
// @Param options body core.SeatOwnerGetOptions true "options"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /dp/get [post]
func (o *OMSNewPlatform) DemandPartnerGetSeatOwnersHandler(c *fiber.Ctx) error {
	data := &core.SeatOwnerGetOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting seat owners data", err)
	}

	seatOwners, err := o.demandPartnerService.GetSeatOwners(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve seat owners data", err)
	}

	return c.JSON(seatOwners)
}
