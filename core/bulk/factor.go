package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type FactorUpdateRequest struct {
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Factor    float64 `json:"factor"`
	Country   string  `json:"country"`
}

func BulkInsertFactors(ctx context.Context, requests []FactorUpdateRequest) error {
	chunks, err := makeChunksFactor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for factor updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		factors, metaDataQueue := prepareFactorsData(chunk)

		if err := bulkInsertFactors(ctx, tx, factors); err != nil {
			log.Error().Err(err).Msgf("failed to process factor bulk update for chunk %d", i)
			return fmt.Errorf("failed to process factor bulk update for chunk %d: %w", i, err)
		}

		if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process factor metadata queue for chunk %d", i)
			return fmt.Errorf("failed to process factor metadata queue for chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in factor bulk update")
		return fmt.Errorf("failed to commit transaction in factor bulk update: %w", err)
	}

	return nil
}

func makeChunksFactor(requests []FactorUpdateRequest) ([][]FactorUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
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

func prepareFactorsData(chunk []FactorUpdateRequest) ([]models.Factor, []models.MetadataQueue) {
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

		key := utils.CreateMetadataOldKey(metadataKey, utils.FactorMetaDataKeyPrefix)

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			Key:           key,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(data.Factor, 'f', 2, 64)),
		})
	}

	return factors, metaDataQueue
}

func bulkInsertFactors(ctx context.Context, tx *sql.Tx, factors []models.Factor) error {
	req := prepareBulkInsertFactorsRequest(factors)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertFactorsRequest(factors []models.Factor) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.Factor,
		columns: []string{
			models.FactorColumns.Publisher,
			models.FactorColumns.Domain,
			models.FactorColumns.Device,
			models.FactorColumns.Country,
			models.FactorColumns.Factor,
			models.FactorColumns.CreatedAt,
			models.FactorColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.FactorColumns.Publisher,
			models.FactorColumns.Domain,
			models.FactorColumns.Device,
			models.FactorColumns.Country,
		},
		updateColumns: []string{
			models.FactorColumns.Factor,
			models.FactorColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(factors)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(factors)*multiplier)

	for i, factor := range factors {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7),
		)
		req.args = append(req.args,
			factor.Publisher,
			factor.Domain,
			factor.Device,
			factor.Country,
			factor.Factor,
			currentTime,
			currentTime,
		)
	}

	return req
}
