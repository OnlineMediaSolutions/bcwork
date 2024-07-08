package rest

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
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
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: "error when parsing request body for /floor/get"})
	}

	pubs, err := core.GetFloors(c.Context(), data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: "failed to retrieve floors"})
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
	var body FloorUpdateRequest

	if c.Body() == nil || len(c.Body()) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request", "message": "Invalid JSON payload"})
	}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request", "message": "Invalid JSON payload"})
	}

	err = updateFloors(c, &body)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to update Floor table with the following")
		return c.SendStatus(http.StatusInternalServerError)
	}

	errMessage := updateFloorMetaData(c, &body)
	if len(errMessage) != 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: errMessage})
	}

	return c.Status(http.StatusOK).JSON(Response{Status: "ok", Message: "Floor and metadata tables successfully updated"})
}

func updateFloorMetaData(c *fiber.Ctx, data *FloorUpdateRequest) string {
	utils.ConvertingAllValues(data)
	val, err := json.Marshal(data)

	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse hash value for floor")
		return "Failed to parse hash value"
	}

	mod := models.MetadataQueue{
		Key:           "price:floor:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	if data.Device == "mobile" {
		mod.Key = "mobile:" + mod.Key
	} else if data.Device == "desktop" {
		mod.Key = "desktop:" + mod.Key
	} else if data.Device == "tablet" {
		mod.Key = "tablet:" + mod.Key
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to insert metadata update to queue for floor")
		return "failed to insert metadata update to queue for floor"
	}
	return ""
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
