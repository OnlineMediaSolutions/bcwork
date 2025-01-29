package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"time"
)

func UpdateMetaDataQueue(ctx context.Context, data *dto.PriceOverrideRequest) error {
	var value types.JSON
	var err error

	priceOverride, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:override:"+data.Domain), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to fetch price override from metadata: %w", err)
	}
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

	err = mod.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert metadata update to queue: %w", err)
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
