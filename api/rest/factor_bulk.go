package rest

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FactorBulkPostHandler Update and enable Bulk insert Factor setup
// @Description Update Factor setup in bulk (publisher, factor, device and country fields are mandatory)
// @Tags Factor in Bulk
// @Accept json
// @Produce json
// @Param options body []FactorUpdateRequest true "Factor update Options"
// @Success 200 {object} FactorUpdateResponse
// @Security ApiKeyAuth
// @Router /factor/bulk [post]
func FactorBulkPostHandler(c *fiber.Ctx) error {
	var requests []FactorUpdateRequest
	if err := c.BodyParser(&requests); err != nil {
		log.Error().Err(err).Msg("Error parsing request body for bulk update")
		return c.Status(http.StatusBadRequest).JSON(Response{Status: "error", Message: "error when parsing request body for bulk update"})
	}

	err, done := makeChunks(c, requests)
	if done {
		return err
	}

	return c.Status(http.StatusOK).JSON(Response{Status: "ok", Message: "Bulk update successfully processed"})
}

func makeChunks(c *fiber.Ctx, requests []FactorUpdateRequest) (error, bool) {
	chunkSize := viper.GetInt("chunkSize")
	for i := 0; i < len(requests); i += chunkSize {
		end := i + chunkSize
		if end > len(requests) {
			end = len(requests)
		}
		chunk := requests[i:end]

		if err := processChunk(c, chunk); err != nil {
			log.Error().Err(err).Msgf("Failed to process bulk update for chunk %d-%d", i, end)
			continue // Skip the failed chunk and continue with the next
		}
	}
	return nil, false
}

func processChunk(c *fiber.Ctx, chunk []FactorUpdateRequest) error {
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

		key := createMetadataKey(data)
		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			Key:           key,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(data.Factor, 'f', 2, 64)),
		})
	}

	return factors, metaDataQueue
}

func createMetadataKey(data FactorUpdateRequest) string {
	key := "price:factor:" + data.Publisher
	if data.Domain != "" {
		key = key + ":" + data.Domain
	}
	if data.Device == "mobile" {
		key = "mobile:" + key
	}
	return key
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

	return bulkInsert(c, tx, "factor", columns, values, conflictColumns, updateColumns)
}

func bulkInsertMetaDataQueue(c *fiber.Ctx, tx *sql.Tx, metaDataQueue []models.MetadataQueue) error {
	columns := []string{"key", "transaction_id", "value", "commited_instances", "created_at", "updated_at"}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, metaData := range metaDataQueue {
		values = append(values, metaData.Key, metaData.TransactionID, metaData.Value, 0, currTime, currTime)
	}

	return bulkInsert(c, tx, "metadata_queue", columns, values, nil, nil)
}

func bulkInsert(c *fiber.Ctx, tx *sql.Tx, tableName string, columns []string, values []interface{}, conflictColumns, updateColumns []string) error {
	columnCount := len(columns)
	valueStrings := make([]string, 0, len(values)/columnCount)
	valueArgs := make([]interface{}, 0, len(values))

	for i := 0; i < len(values)/columnCount; i++ {
		placeholders := make([]string, columnCount)
		for j := 0; j < columnCount; j++ {
			placeholders[j] = fmt.Sprintf("$%d", i*columnCount+j+1)
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		valueArgs = append(valueArgs, values[i*columnCount:(i+1)*columnCount]...)
	}

	columnNames := strings.Join(columns, ", ")
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, columnNames, strings.Join(valueStrings, ","))

	if conflictColumns != nil && updateColumns != nil {
		stmt += fmt.Sprintf(" ON CONFLICT (%s) DO UPDATE SET %s", strings.Join(conflictColumns, ", "), strings.Join(updateColumns, ", "))
	}

	log.Info().Msgf("Executing bulk insert for %s: %s", tableName, stmt)
	if _, err := tx.ExecContext(c.Context(), stmt, valueArgs...); err != nil {
		log.Error().Err(err).Msgf("Failed to execute bulk insert for %s: %s", tableName, stmt)
		return fmt.Errorf("failed to insert into %s in bulk: %w", tableName, err)
	}

	return nil
}
