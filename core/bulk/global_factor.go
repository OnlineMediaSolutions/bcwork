package bulk

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type GlobalFactorRequest struct {
	Key       string  `json:"key" validate:"globalFactorKey"`
	Publisher string  `json:"publisher_id" `
	Value     float64 `json:"value"`
}

func MakeGlobalFactorsChunks(requests []GlobalFactorRequest) ([][]GlobalFactorRequest, error) {
	chunkSize := viper.GetInt("api.chunkSize")
	var chunks [][]GlobalFactorRequest

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

func ProcessGlobalFactorsChunks(c *fiber.Ctx, chunks [][]GlobalFactorRequest) error {
	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to begin transaction",
		})
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		globalFactors := prepareGlobalFactorsData(chunk)

		if err := bulkInsertGlobalFactors(c, tx, globalFactors); err != nil {
			log.Error().Err(err).Msgf("failed to process global factor bulk update for chunk %d", i)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in global factor bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to commit transaction",
		})
	}

	return nil
}

func prepareGlobalFactorsData(chunk []GlobalFactorRequest) []models.GlobalFactor {
	var globalFactors []models.GlobalFactor

	for _, data := range chunk {
		globalFactors = append(globalFactors, models.GlobalFactor{
			Key:         data.Key,
			PublisherID: data.Publisher,
			Value:       null.Float64From(data.Value),
		})
	}

	return globalFactors
}

func bulkInsertGlobalFactors(c *fiber.Ctx, tx *sql.Tx, globalFactors []models.GlobalFactor) error {
	const excluded = " = EXCLUDED."

	columns := []string{
		models.GlobalFactorColumns.Key,
		models.GlobalFactorColumns.PublisherID,
		models.GlobalFactorColumns.Value,
		models.GlobalFactorColumns.CreatedAt,
	}
	conflictColumns := []string{models.GlobalFactorColumns.Key, models.GlobalFactorColumns.PublisherID}
	updateColumns := []string{
		models.GlobalFactorColumns.Value + excluded + models.GlobalFactorColumns.Value,
		models.GlobalFactorColumns.UpdatedAt + excluded + models.GlobalFactorColumns.CreatedAt,
	}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, globalFactor := range globalFactors {
		values = append(values, globalFactor.Key, globalFactor.PublisherID, globalFactor.Value, currTime)
	}

	return InsertInBulk(c, tx, models.TableNames.GlobalFactor, columns, values, conflictColumns, updateColumns)
}
