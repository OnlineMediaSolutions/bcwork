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
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: "failed to retrieve floors"})
	}
	return c.JSON(pubs)
}

// FloorPostHandler Update and enable Floor setup
// @Description Update Floor setup (publisher is mandatory)
// @Tags Floor
// @Accept json
// @Produce json
// @Param options body FloorUpdateRequest true "Floor update Options"
// @Success 200 {object} FloorUpdateResponse
// @Security ApiKeyAuth
// @Router /floor [post]
func FloorPostHandler(c *fiber.Ctx) error {
	data := &FloorUpdateRequest{}
	done := false
	err := c.BodyParser(&data)

	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Error when parsing floor payload")
		return c.SendStatus(http.StatusBadRequest)
	}

	err, done = validateInputs(c, data)
	if done {
		return err
	}

	err = updateFloor(c, data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to update Floor table with the following")
		return c.SendStatus(http.StatusInternalServerError)
	}

	errMessage := updateMetaData(c, data)
	if len(errMessage) != 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: errMessage})
	}

	return c.Status(http.StatusOK).JSON(Response{Status: "Ok", Message: "Floor and metadata tables successfully updated"})
}

func validateInputs(c *fiber.Ctx, data *FloorUpdateRequest) (error, bool) {

	if data.Country != "all" && (len(data.Country) > 2) {
		c.SendString(fmt.Sprintf("Country must be a 2-letter country code or 'all'"))
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	if data.Publisher == "" {
		c.SendString(fmt.Sprintf("Publisher is mandatory"))
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	if data.Domain == "" {
		c.SendString(fmt.Sprintf("Domain is mandatory"))
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	// TODO need to add floor validation

	if data.Device != "" && !allowedDevices(data.Device) {
		c.SendString(fmt.Sprintf("'%s' not allowed as device  name", data.Device))
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	return nil, false
}

func updateMetaData(c *fiber.Ctx, data *FloorUpdateRequest) string {
	val, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse floor data")
		return "Failed to parse data"
	}

	mod := models.MetadataQueue{
		Key:           "floor:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to insert floor metadata update to queue")
		return "failed to insert floor metadata update to queue"
	}
	return ""
}

func updateFloor(c *fiber.Ctx, data *FloorUpdateRequest) error {

	modConf := models.Floor{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
		Floor:     data.Floor,
		Country:   data.Country,
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.FloorColumns.Publisher, models.FloorColumns.Domain}, boil.Infer(), boil.Infer())
}

func allowedDevices(device string) bool {
	_, isExists := utils.AllowedDevices[device]
	return isExists
}
