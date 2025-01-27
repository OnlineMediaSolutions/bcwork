package core

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"net/http"
	"time"
)

func UpdateMetaDataQueue(c *fiber.Ctx, data *dto.PriceOverrideRequest) error {

	value, err := buildPriceOvverideValue(data)
	mod := models.MetadataQueue{
		Key:           "price:override:" + data.Domain,
		TransactionID: bcguid.NewFromf(data.Domain, time.Now()),
		Value:         value,
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return c.SendStatus(http.StatusInternalServerError)
	}

	return nil
}

func buildPriceOvverideValue(data *dto.PriceOverrideRequest) (types.JSON, error) {
	currentTime := time.Now()
	ips := []ipPriceDate{}

	for _, userData := range data.Ips {
		ipPriceDate := ipPriceDate{
			IP:    userData.IP,
			Date:  currentTime,
			Price: userData.Price,
		}
		ips = append(ips, ipPriceDate)
	}

	val, err := json.Marshal(ips)
	if err != nil {
		return nil, fmt.Errorf("price override failed to parse hash value: %w", err)
	}

	return val, err
}

type ipPriceDate struct {
	IP    string    `json:"ip"`
	Date  time.Time `json:"date"`
	Price float64   `json:"price"`
}
