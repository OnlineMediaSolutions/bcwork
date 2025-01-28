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
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"net/http"
	"time"
)

func UpdateMetaDataQueue(c *fiber.Ctx, data *dto.PriceOverrideRequest) error {
	var value types.JSON
	var err error

	priceOverride, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:override:"+data.Domain), qm.OrderBy("updated_at desc")).One(c.Context(), bcdb.DB())
	if priceOverride == nil || priceOverride.UpdatedAt.Time.Before(time.Now().Add(-8*time.Hour)) {
		value, err = buildPriceOvverideValue(data)
	} else {
		value, err = addNewIpToValue(priceOverride.Value, data)
	}

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

func addNewIpToValue(value types.JSON, data *dto.PriceOverrideRequest) (types.JSON, error) {

	var metaDataValue []dto.Ips
	err := json.Unmarshal([]byte(value), &metaDataValue)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal metadata value for price override: %w", err)
	}

	currentTime := time.Now()
	for _, userData := range data.Ips {
		ipPriceDate := dto.Ips{
			IP:    userData.IP,
			Date:  currentTime,
			Price: userData.Price,
		}
		metaDataValue = append(metaDataValue, ipPriceDate)
	}

	val, err := json.Marshal(metaDataValue)
	if err != nil {
		return nil, fmt.Errorf("price override failed to parse hash value: %w", err)
	}

	return val, nil
}

func buildPriceOvverideValue(data *dto.PriceOverrideRequest) (types.JSON, error) {

	ips := []dto.Ips{}
	currentTime := time.Now()
	for _, userData := range data.Ips {
		ipPriceDate := dto.Ips{
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
