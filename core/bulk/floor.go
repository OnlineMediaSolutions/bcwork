package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/utils"
	"github.com/volatiletech/null/v8"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func BulkInsertFloors(ctx context.Context, requests []constant.FloorUpdateRequest) error {
	chunks, err := makeChunksFloor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for floors updates: %w", err)
	}
	for i, chunk := range chunks {
		tx, err := bcdb.DB().BeginTx(ctx, nil)

		floors, err := prepareFloorsData(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare floor data for chunk %d", i)
			return fmt.Errorf("failed to prepare floor data for chunk %d: %w", i, err)
		}

		if err := bulkInsertFloor(ctx, tx, floors); err != nil {
			tx.Rollback()
			log.Error().Err(err).Msgf("failed to process floor bulk update for chunk %d", i)
			return fmt.Errorf("failed to process floor bulk update for chunk %d: %w", i, err)
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("failed to commit floor transaction")
			return fmt.Errorf("failed to commit floor transaction: %w", err)
		}

		tx, err = bcdb.DB().BeginTx(ctx, nil)
		if err != nil {
			log.Error().Err(err).Msg("failed to start transaction for metadata")
			return fmt.Errorf("failed to start transaction for metadata: %w", err)
		}

		metaDataQueue, err := prepareMetadata(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare metadata for chunk %d", i)
			return fmt.Errorf("failed to prepare metadata for chunk %d: %w", i, err)
		}

		if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
			tx.Rollback()
			log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
			return fmt.Errorf("failed to process metadata queue for chunk %d: %w", i, err)
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("failed to commit metadata transaction")
			return fmt.Errorf("failed to commit metadata transaction: %w", err)
		}
	}

	return nil
}

func makeChunksFloor(requests []constant.FloorUpdateRequest) ([][]constant.FloorUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]constant.FloorUpdateRequest

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

func prepareFloorsData(ctx context.Context, chunk []constant.FloorUpdateRequest) ([]models.Floor, error) {
	return prepareFloors(chunk), nil
}

func prepareMetadata(ctx context.Context, chunk []constant.FloorUpdateRequest) ([]models.MetadataQueue, error) {
	metaData, err := prepareFloorsMetadata(ctx, chunk)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare floor metadata: %w", err)
	}
	return metaData, nil

}

func prepareFloorsMetadata(ctx context.Context, chunk []constant.FloorUpdateRequest) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		modFloor, err := core.FloorQuery(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("error querying floor: %w", err)
		}
		var finalRules []core.FloorRealtimeRecord

		finalRules = core.CreateFloorMetadata(modFloor, finalRules)

		finalOutput := struct {
			Rules []core.FloorRealtimeRecord `json:"rules"`
		}{Rules: finalRules}

		value, err := json.Marshal(finalOutput)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON for floor metadata value: %w", err)
		}

		key := utils.GetMetadataObject(data)
		metadataKey := utils.CreateMetadataKey(key, utils.FloorMetaDataKeyPrefix)
		metadataValue := utils.CreateMetadataObject(data, metadataKey, value)
		metaDataQueue = append(metaDataQueue, metadataValue)
	}

	return metaDataQueue, nil
}

func prepareFloors(chunk []constant.FloorUpdateRequest) []models.Floor {
	var floors []models.Floor
	for _, data := range chunk {
		floor := core.Floor{
			Publisher:     data.Publisher,
			Domain:        data.Domain,
			Country:       data.Country,
			Device:        data.Device,
			Floor:         data.Floor,
			Browser:       data.Browser,
			OS:            data.OS,
			PlacementType: data.PlacementType,
		}

		floors = append(floors, models.Floor{
			Publisher:     floor.Publisher,
			Domain:        floor.Domain,
			Country:       null.NewString(floor.Country, floor.Country != ""),
			Device:        null.NewString(floor.Device, floor.Device != ""),
			Floor:         floor.Floor,
			Browser:       null.NewString(floor.Browser, floor.Browser != ""),
			Os:            null.NewString(floor.OS, floor.OS != ""),
			PlacementType: null.NewString(floor.PlacementType, floor.PlacementType != ""),
			RuleID:        floor.GetRuleID(),
		})
	}

	return floors
}

func bulkInsertFloor(ctx context.Context, tx *sql.Tx, floors []models.Floor) error {
	req := prepareBulkInsertFloorsRequest(floors)

	err := bulkInsert(ctx, tx, req)
	if err != nil {
		return err
	}
	return nil
}

func prepareBulkInsertFloorsRequest(floors []models.Floor) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.Floor,
		columns: []string{
			models.FloorColumns.RuleID,
			models.FloorColumns.Publisher,
			models.FloorColumns.Domain,
			models.FloorColumns.Country,
			models.FloorColumns.Browser,
			models.FloorColumns.Os,
			models.FloorColumns.Device,
			models.FloorColumns.PlacementType,
			models.FloorColumns.Floor,
			models.FloorColumns.CreatedAt,
			models.FloorColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.FloorColumns.RuleID,
		},
		updateColumns: []string{
			models.FloorColumns.Floor,
			models.FloorColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(floors)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(floors)*multiplier)

	for i, floor := range floors {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5,
				offset+6, offset+7, offset+8, offset+9, offset+10, offset+11),
		)
		req.args = append(req.args,
			floor.RuleID,
			floor.Publisher,
			floor.Domain,
			floor.Country,
			floor.Browser,
			floor.Os,
			floor.Device,
			floor.PlacementType,
			floor.Floor,
			currentTime,
			currentTime,
		)
	}

	return req
}
