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

type FactorUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Factor    float64 `json:"factor"`
	Country   string  `json:"country"`
}

type FactorUpdateResponse struct {
	Status string `json:"status"`
}

// FactorGetHandler Get factor setup
// @Description Get factor setup
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body core.GetFactorOptions true "options"
// @Success 200 {object} core.FactorSlice
// @Router /factor/get [post]
func FactorGetAllHandler(c *fiber.Ctx) error {

	data := &core.GetFactorOptions{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(Response{Status: "error", Message: "error when parsing request body for /price/factor/get"})
	}

	pubs, err := core.GetFactors(c.Context(), data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: "failed to retrieve factors"})
	}
	return c.JSON(pubs)
}

// FactorPostHandler Update and enable Factor setup
// @Description Update Factor setup (publisher is mandatory, domain is mandatory)
// @Tags Factor
// @Accept json
// @Produce json
// @Param options body FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /factor [post]
func FactorPostHandler(c *fiber.Ctx) error {
	data := &FactorUpdateRequest{}
	done := false
	err := c.BodyParser(&data)

	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Error when parsing factor payload")
		return c.SendStatus(http.StatusBadRequest)
	}

	err, done = validateInputs(c, data)
	if done {
		return err
	}

	err = updateFactor(c, data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to update Factor table with the following")
		return c.SendStatus(http.StatusInternalServerError)
	}

	errMessage := updateMetaData(c, data)
	if len(errMessage) != 0 {
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: errMessage})
	}

	return c.Status(http.StatusOK).JSON(Response{Status: "Ok", Message: "Factor and metadata tables successfully updated"})
}

func validateInputs(c *fiber.Ctx, data *FactorUpdateRequest) (error, bool) {

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

	if data.Factor < 0.01 || data.Factor > 10 {
		c.SendString("'Factor' must be between 0.01 and 10")
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	if data.Device != "" && !allowedDevices(data.Device) {
		c.SendString(fmt.Sprintf("'%s' not allowed as device  name", data.Device))
		c.Status(http.StatusBadRequest)
		return nil, true
	}

	return nil, false
}

func updateMetaData(c *fiber.Ctx, data *FactorUpdateRequest) string {
	val, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse data")
		return "Failed to parse data"
	}

	mod := models.MetadataQueue{
		Key:           "factor:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to insert metadata update to queue")
		return "failed to insert metadata update to queue"
	}
	return ""
}

func updateFactor(c *fiber.Ctx, data *FactorUpdateRequest) error {

	modConf := models.Factor{
		Publisher: data.Publisher,
		Domain:    data.Domain,
		Device:    data.Device,
		Factor:    data.Factor,
		Country:   data.Country,
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.FactorColumns.Publisher, models.FactorColumns.Domain}, boil.Infer(), boil.Infer())
}

func allowedDevices(device string) bool {
	_, isExists := utils.AllowedDevices[device]
	return isExists
}
