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
	chunkSize := viper.GetInt("api.chunkSize")
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

func ProcessChunks(c *fiber.Ctx, chunks [][]FactorUpdateRequest) error {

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to begin transaction",
		})
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		factors, metaDataQueue := prepareData(chunk)

		if err := bulkInsertFactors(c, tx, factors); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for chunk %d", i)
			return err
		}

		if err := BulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in factor bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to commit transaction",
		})
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

		metadataKey := utils.MetadataKey{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.Device,
			Country:   data.Country,
		}
		key := utils.CreateMetadataOldKey(metadataKey, "price:factor")
		fmt.Println(key)

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
