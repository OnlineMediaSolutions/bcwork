package bulk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
)

type GlobalFactorRequest struct {
	Key       string  `json:"key" validate:"globalFactorKey"`
	Publisher string  `json:"publisher_id"`
	Value     float64 `json:"value"`
}

func BulkInsertGlobalFactors(ctx context.Context, requests []GlobalFactorRequest) error {
	chunks, err := makeGlobalFactorsChunks(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for global factor bulk update: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		// TODO: get old global factors
		globalFactors := prepareGlobalFactorsData(chunk)

		if err := bulkInsertGlobalFactors(ctx, tx, globalFactors); err != nil {
			log.Error().Err(err).Msgf("failed to process global factor bulk update for chunk %d", i)
			return fmt.Errorf("failed to process global factor bulk update for chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in global factor bulk update")
		return fmt.Errorf("failed to commit transaction in global factor bulk update: %w", err)
	}

	return nil
}

func makeGlobalFactorsChunks(requests []GlobalFactorRequest) ([][]GlobalFactorRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
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

func prepareGlobalFactorsData(chunk []GlobalFactorRequest) []*models.GlobalFactor {
	var globalFactors []*models.GlobalFactor
	for _, data := range chunk {
		globalFactors = append(globalFactors, &models.GlobalFactor{
			Key:         data.Key,
			PublisherID: data.Publisher,
			Value:       null.Float64From(data.Value),
		})
	}

	return globalFactors
}

func bulkInsertGlobalFactors(ctx context.Context, tx *sql.Tx, globalFactors []*models.GlobalFactor) error {
	req := prepareBulkInsertGlobalFactorsRequest(globalFactors)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertGlobalFactorsRequest(globalFactors []*models.GlobalFactor) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.GlobalFactor,
		columns: []string{
			models.GlobalFactorColumns.Key,
			models.GlobalFactorColumns.PublisherID,
			models.GlobalFactorColumns.Value,
			models.GlobalFactorColumns.CreatedAt,
			models.GlobalFactorColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.GlobalFactorColumns.Key,
			models.GlobalFactorColumns.PublisherID,
		},
		updateColumns: []string{
			models.GlobalFactorColumns.Value,
			models.GlobalFactorColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(globalFactors)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(globalFactors)*multiplier)

	for i, globalFactor := range globalFactors {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5),
		)
		req.args = append(req.args,
			globalFactor.Key,
			globalFactor.PublisherID,
			globalFactor.Value,
			currentTime,
			currentTime,
		)
	}

	return req
}
