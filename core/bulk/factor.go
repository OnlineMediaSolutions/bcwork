package bulk

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
	"time"
)

type FactorUpdateRequest struct {
	utils.MetadataKey
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Factor    float64 `json:"factor"`
	Country   string  `json:"country"`
}

func MakeChunks(requests []FactorUpdateRequest) ([][]FactorUpdateRequest, error) {
	chunkSize := viper.GetInt("chunkSize")
	var chunks [][]FactorUpdateRequest

	for i := 0; i < len(requests); i += chunkSize {
		end := i + chunkSize
		if end > len(requests) {
			end = len(requests)
		}
		chunk := requests[i:end]
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func InsertChunk(c *fiber.Ctx, chunk []FactorUpdateRequest) error {
	factors, metaDataQueue := prepareData(chunk)

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := bulkInsertFactors(c, tx, factors); err != nil {
		return err
	}

	if err := bulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction in factor bulk update")
		return fmt.Errorf("failed to commit transaction in factor bulk: %w", err)
	}

	return nil
}

func prepareData(chunk []FactorUpdateRequest) ([]models.Factor, []models.MetadataQueue) {
	var factors []models.Factor
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		factors = append(factors, models.Factor{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.Device,
			Factor:    data.Factor,
			Country:   data.Country,
		})

		key := utils.CreateMetadataKey(data.MetadataKey, "price:factor")
		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			Key:           key,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(data.Factor, 'f', 2, 64)),
		})
	}

	return factors, metaDataQueue
}

func bulkInsertFactors(c *fiber.Ctx, tx *sql.Tx, factors []models.Factor) error {
	columns := []string{"publisher", "domain", "device", "factor", "country", "created_at", "updated_at"}
	conflictColumns := []string{"publisher", "domain", "device", "country"}
	updateColumns := []string{"factor = EXCLUDED.factor", "created_at = EXCLUDED.created_at", "updated_at = EXCLUDED.updated_at"}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, factor := range factors {
		values = append(values, factor.Publisher, factor.Domain, factor.Device, factor.Factor, factor.Country, currTime, currTime)
	}

	return InsertInBulk(c, tx, "factor", columns, values, conflictColumns, updateColumns)
}

func bulkInsertMetaDataQueue(c *fiber.Ctx, tx *sql.Tx, metaDataQueue []models.MetadataQueue) error {
	columns := []string{"key", "transaction_id", "value", "commited_instances", "created_at", "updated_at"}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, metaData := range metaDataQueue {
		values = append(values, metaData.Key, metaData.TransactionID, metaData.Value, 0, currTime, currTime)
	}

	return InsertInBulk(c, tx, "metadata_queue", columns, values, nil, nil)
}
