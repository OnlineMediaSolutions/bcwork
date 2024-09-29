package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func BulkInsertDPO(ctx context.Context, requests []core.DPOUpdateRequest) error {
	chunks, err := makeChunksDPO(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for dpos updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		dpos, metaDataQueue, err := prepareDPOData(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare dpo data for chunk %d", i)
			return fmt.Errorf("failed to prepare dpo data for chunk %d: %w", i, err)
		}

		if err := bulkInsertDPO(ctx, tx, dpos); err != nil {
			log.Error().Err(err).Msgf("failed to process dpos bulk update for chunk %d", i)
			return fmt.Errorf("failed to process dpos bulk update for chunk %d: %w", i, err)
		}

		if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process dpos metadata queue for chunk %d", i)
			return fmt.Errorf("failed to process dpos metadata queue for chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in dpos bulk update")
		return fmt.Errorf("failed to commit transaction in dpos bulk update: %w", err)
	}

	return nil
}

func makeChunksDPO(requests []core.DPOUpdateRequest) ([][]core.DPOUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]core.DPOUpdateRequest

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

func prepareDPOData(ctx context.Context, chunk []core.DPOUpdateRequest) ([]models.DpoRule, []models.MetadataQueue, error) {
	metaData, err := prepareDPOMetadata(ctx, chunk)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot prepare dpo metadata: %w", err)
	}

	return prepareDPO(chunk), metaData, nil
}

func prepareDPO(chunk []core.DPOUpdateRequest) []models.DpoRule {
	var dpos []models.DpoRule

	for _, data := range chunk {
		DPOOptimizationRule := core.DemandPartnerOptimizationRule{
			DemandPartner: data.DemandPartner,
			Publisher:     data.Publisher,
			Domain:        data.Domain,
			Country:       data.Country,
			OS:            data.OS,
			DeviceType:    data.DeviceType,
			PlacementType: data.PlacementType,
			Browser:       data.Browser,
			Factor:        data.Factor,
			RuleID:        data.RuleId,
		}
		dpos = append(dpos, *DPOOptimizationRule.ToModel())
	}

	return dpos
}

func prepareDPOMetadata(ctx context.Context, chunk []core.DPOUpdateRequest) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		modDpos, err := models.DpoRules(models.DpoRuleWhere.DemandPartnerID.EQ(data.DemandPartner)).All(ctx, bcdb.DB())
		if err != nil {
			return nil, fmt.Errorf("cannot get dpo rules for demand partner id [%v]: %w", data.DemandPartner, err)
		}

		dpos := make(core.DemandPartnerOptimizationRuleSlice, 0, len(modDpos))
		dpos.FromModel(modDpos)

		dposRT := core.DpoRT{
			DemandPartnerID: data.DemandPartner,
			IsInclude:       false,
		}

		for _, dpo := range dpos {
			dposRT.Rules = append(dposRT.Rules, dpo.ToRtRule())
		}
		dposRT.Rules.Sort()

		b, err := json.Marshal(dposRT)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON for dpo metadata value: %w", err)
		}

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           utils.DPOMetaDataKeyPrefix + ":" + data.DemandPartner,
			Value:         b,
		})
	}
	return metaDataQueue, nil
}

func bulkInsertDPO(ctx context.Context, tx *sql.Tx, dpos []models.DpoRule) error {
	req := prepareBulkInsertDPORequest(dpos)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertDPORequest(dpos []models.DpoRule) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.DpoRule,
		columns: []string{
			models.DpoRuleColumns.RuleID,
			models.DpoRuleColumns.DemandPartnerID,
			models.DpoRuleColumns.Publisher,
			models.DpoRuleColumns.Domain,
			models.DpoRuleColumns.Country,
			models.DpoRuleColumns.Browser,
			models.DpoRuleColumns.Os,
			models.DpoRuleColumns.DeviceType,
			models.DpoRuleColumns.PlacementType,
			models.DpoRuleColumns.Factor,
			models.DpoRuleColumns.CreatedAt,
			models.DpoRuleColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.DpoRuleColumns.RuleID,
		},
		updateColumns: []string{
			models.DpoRuleColumns.Country,
			models.DpoRuleColumns.Factor,
			models.DpoRuleColumns.DeviceType,
			models.DpoRuleColumns.Domain,
			models.DpoRuleColumns.PlacementType,
			models.DpoRuleColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(dpos)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(dpos)*multiplier)

	for i, dpo := range dpos {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6,
				offset+7, offset+8, offset+9, offset+10, offset+11, offset+12),
		)
		req.args = append(req.args,
			dpo.RuleID,
			dpo.DemandPartnerID,
			dpo.Publisher,
			dpo.Domain,
			dpo.Country,
			dpo.Browser,
			dpo.Os,
			dpo.DeviceType,
			dpo.PlacementType,
			dpo.Factor,
			currentTime,
			currentTime,
		)
	}

	return req
}
