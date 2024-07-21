package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
	"time"
)

type FloorUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Floor     float64 `json:"floor"`
	Country   string  `json:"country"`
}

type FloorUpdateResponse struct {
	Status string `json:"status"`
}

// FloorGetHandler Get floor setup
// @Description Get floor setup
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body core.GetFloorOptions true "options"
// @Success 200 {object} core.FloorSlice
// @Router /floor/get [post]
func FloorGetAllHandler(c *fiber.Ctx) error {

	data := &core.GetFloorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Request body parsing error")
	}

	pubs, err := core.GetFloors(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to retrieve floors")
	}
	return c.JSON(pubs)
}

// FloorPostHandler Update and enable Floor setup
// @Description Update Floor setup (publisher, floor, device and country fields are mandatory)
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} FloorUpdateResponse
// @Security ApiKeyAuth
// @Router /floor [post]
func FloorPostHandler(c *fiber.Ctx) error {
	data := &FloorUpdateRequest{}

	err := c.BodyParser(&data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Floor payload parsing error")
	}

	err = updateFloors(c, data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update Floor table")
	}

	err = updateFloorMetaData(c, data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Failed to update metadata table for floor, %s", err))
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Floor and Metadata tables successfully updated")
}

func updateFloorMetaData(c *fiber.Ctx, data *FloorUpdateRequest) error {
	_, err := json.Marshal(data)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to parse hash value for floor")
	}

	metadataKey := utils.MetadataKey{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
	}

	key := utils.CreateMetadataKey(metadataKey, "price:floor")

	floor := strconv.FormatFloat(data.Floor, 'f', 2, 64)
	mod := models.MetadataQueue{
		Key:           key,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(floor),
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to insert metadata update to queue")
	}

	return nil
}

func updateFloors(c *fiber.Ctx, data *FloorUpdateRequest) error {

	modConf := models.Floor{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
		Floor:     data.Floor,
		Country:   data.Country,
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.FloorColumns.Publisher, models.FloorColumns.Domain, models.FloorColumns.Device, models.FloorColumns.Country}, boil.Infer(), boil.Infer())
}
