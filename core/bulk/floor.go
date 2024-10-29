package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func BulkInsertFloors(ctx context.Context, requests []constant.FloorUpdateRequest) error {
	chunks, err := makeChunksFloor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for floors updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		floors, metaDataQueue, err := prepareFloorsData(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msgf("failed to prepare floor data for chunk %d", i)
			return fmt.Errorf("failed to prepare floor data for chunk %d: %w", i, err)
		}

		if err := bulkInsertFloor(ctx, tx, floors); err != nil {
			log.Error().Err(err).Msgf("failed to process floor bulk update for chunk %d", i)
			return fmt.Errorf("failed to process floor bulk update for chunk %d: %w", i, err)
		}

		if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process floor metadata queue for chunk %d", i)
			return fmt.Errorf("failed to process floor metadata queue for chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in floor bulk update")
		return fmt.Errorf("failed to commit transaction in floor bulk update: %w", err)
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

func prepareFloorsData(ctx context.Context, chunk []constant.FloorUpdateRequest) ([]models.Floor, []models.MetadataQueue, error) {
	metaData, err := prepareFloorsMetadata(ctx, chunk)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot prepare floor metadata: %w", err)
	}

	return prepareFloors(chunk), metaData, nil
}

func prepareFloorsMetadata(ctx context.Context, chunk []constant.FloorUpdateRequest) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		modFloor, err := core.FloorQuery(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("cannot get floors for publisher [%v] and domain [%v]: %w",
				data.Publisher, data.Domain, err)
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
			Publisher: floor.Publisher,
			Domain:    floor.Domain,
			Floor:     floor.Floor,
			RuleID:    floor.GetRuleID(),
		})
	}

	return floors
}

func bulkInsertFloor(ctx context.Context, tx *sql.Tx, floors []models.Floor) error {
	req := prepareBulkInsertFloorsRequest(floors)

	return bulkInsert(ctx, tx, req)
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
