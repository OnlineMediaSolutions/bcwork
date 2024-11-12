package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/volatiletech/null/v8"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type PublisherDomain struct {
	Publisher string
	Domain    string
}

func BulkInsertFactors(ctx context.Context, requests []constant.FactorUpdateRequest) error {
	chunks, err := makeChunksFactor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for factor updates: %w", err)
	}
	for i, chunk := range chunks {
		tx, err := bcdb.DB().BeginTx(ctx, nil)

		factors, err := prepareFactorsData(chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare factors data for chunk %d", i)
			return fmt.Errorf("failed to prepare factors data for chunk %d: %w", i, err)
		}

		if err := bulkInsertFactors(ctx, tx, factors); err != nil {
			tx.Rollback()

			log.Error().Err(err).Msgf("failed to process factor bulk update for chunk %d", i)
			return fmt.Errorf("failed to process factor bulk update for chunk %d: %w", i, err)
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("failed to commit factor transaction")
			return fmt.Errorf("failed to commit factor transaction: %w", err)
		}

		publisherDomainMap := make(map[string]PublisherDomain)
		var publishers []PublisherDomain
		for _, factor := range factors {
			key := fmt.Sprintf("%s-%s", factor.Publisher, factor.Domain)
			publisherDomainMap[key] = PublisherDomain{
				Publisher: factor.Publisher,
				Domain:    factor.Domain,
			}
			publishers = append(publishers, publisherDomainMap[key])
		}

		metaDataQueue, err := prepareFactorMetadata(ctx, publishers)

		if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process factor metadata queue for chunk %d", i)
			return fmt.Errorf("failed to process factor metadata queue for chunk %d: %w", i, err)
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("failed to commit factor transaction")
			return fmt.Errorf("failed to commit factor transaction: %w", err)
		}
	}

	return nil
}

func makeChunksFactor(requests []constant.FactorUpdateRequest) ([][]constant.FactorUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]constant.FactorUpdateRequest

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

func prepareFactorsData(chunk []constant.FactorUpdateRequest) ([]models.Factor, error) {
	var factors []models.Factor
	for _, data := range chunk {
		factor := core.Factor{
			Publisher:     data.Publisher,
			Domain:        data.Domain,
			Country:       data.Country,
			Device:        data.Device,
			Factor:        data.Factor,
			Browser:       data.Browser,
			OS:            data.OS,
			PlacementType: data.PlacementType,
		}

		factors = append(factors, models.Factor{
			Publisher:     factor.Publisher,
			Domain:        factor.Domain,
			Country:       null.NewString(factor.Country, factor.Country != ""),
			Device:        null.NewString(factor.Device, factor.Device != ""),
			Factor:        factor.Factor,
			Browser:       null.NewString(factor.Browser, factor.Browser != ""),
			Os:            null.NewString(factor.OS, factor.OS != ""),
			PlacementType: null.NewString(factor.PlacementType, factor.PlacementType != ""),
			RuleID:        factor.GetRuleID(),
		})
	}

	return factors, nil
}

func prepareFactorMetadata(ctx context.Context, publishers PublisherDomain) ([]models.MetadataQueue, error) {
	metaData, err := prepareFactorsMetadata(ctx, publishers)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare factor metadata: %w", err)
	}
	return metaData, nil

}

func prepareFactorsMetadata(ctx context.Context, chunk []constant.FactorUpdateRequest) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		modFactor, err := core.FactorQuery(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("error querying factor: %w", err)
		}
		var finalRules []core.FactorRealtimeRecord

		finalRules = core.CreateFactorMetadata(modFactor, finalRules)

		finalOutput := struct {
			Rules []core.FactorRealtimeRecord `json:"rules"`
		}{Rules: finalRules}

		value, err := json.Marshal(finalOutput)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON for factor metadata value: %w", err)
		}

		key := utils.GetMetadataObject(data)
		metadataKey := utils.CreateMetadataKey(key, utils.FactorMetaDataKeyPrefix)
		metadataValue := utils.CreateMetadataObject(data, metadataKey, value)
		metaDataQueue = append(metaDataQueue, metadataValue)
	}

	return metaDataQueue, nil
}

func bulkInsertFactors(ctx context.Context, tx *sql.Tx, factors []models.Factor) error {
	req := prepareBulkInsertFactorsRequest(factors)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertFactorsRequest(factors []models.Factor) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.Factor,
		columns: []string{
			models.FactorColumns.RuleID,
			models.FactorColumns.Publisher,
			models.FactorColumns.Domain,
			models.FactorColumns.Country,
			models.FactorColumns.Browser,
			models.FactorColumns.Os,
			models.FactorColumns.Device,
			models.FactorColumns.PlacementType,
			models.FactorColumns.Factor,
			models.FactorColumns.CreatedAt,
			models.FactorColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.FactorColumns.RuleID,
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
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5,
				offset+6, offset+7, offset+8, offset+9, offset+10, offset+11),
		)
		req.args = append(req.args,
			factor.RuleID,
			factor.Publisher,
			factor.Domain,
			factor.Country,
			factor.Browser,
			factor.Os,
			factor.Device,
			factor.PlacementType,
			factor.Factor,
			currentTime,
			currentTime,
		)
	}

	return req
}
