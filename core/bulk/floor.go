package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var softDeleteFloorsQuery = `UPDATE floor
SET active = false
WHERE rule_id in (%s)`

func BulkInsertFloors(ctx context.Context, requests []dto.FloorUpdateRequest) error {
	chunks, err := makeChunksFloor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for floors updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	pubDomains := make(map[string]struct{})
	err = handleBulkFloor(ctx, chunks, pubDomains, tx)
	if err != nil {
		return err
	}

	err = handleMetaDataFloorRules(ctx, pubDomains, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in floor bulk update: %w", err)
	}

	return nil
}

func (f *BulkService) BulkDeleteFloor(ctx context.Context, ids []string) error {
	mods, err := models.Floors(models.FloorWhere.RuleID.IN(ids)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed getting floors for soft deleting: %w", err)
	}

	oldMods := make([]any, 0, len(mods))
	newMods := make([]any, 0, len(mods))

	pubDomains := make(map[string]struct{})
	for _, mod := range mods {
		oldMods = append(oldMods, mod)
		newMods = append(newMods, nil)
		pubDomains[mod.Publisher+":"+mod.Domain] = struct{}{}
	}

	deleteQuery := utils.CreateDeleteQuery(ids, softDeleteFloorsQuery)

	_, err = queries.Raw(deleteQuery).Exec(bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed soft deleting floor rules: %w", err)
	}

	err = updateFloorInMetaData(ctx, pubDomains)
	if err != nil {
		return err
	}

	f.historyModule.SaveAction(ctx, oldMods, newMods, nil)

	return nil
}

func updateFloorInMetaData(ctx context.Context, pubDomains map[string]struct{}) error {
	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for delete floor in metadata_queue: %w", err)
	}
	defer tx.Rollback()

	err = handleMetaDataFloorRules(ctx, pubDomains, tx)
	if err != nil {
		return fmt.Errorf("failed to update RT metadata for delete floors: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in delete floors in metadata_queue: %w", err)
	}

	return nil
}

func handleBulkFloor(ctx context.Context, chunks [][]dto.FloorUpdateRequest, pubDomains map[string]struct{}, tx *sql.Tx) error {
	for i, chunk := range chunks {
		floors := createFloorData(chunk, pubDomains)

		if err := bulkInsertFloor(ctx, tx, floors); err != nil {
			return fmt.Errorf("failed to process floor bulk update for chunk %d: %w", i, err)
		}
	}

	return nil
}

func handleMetaDataFloorRules(ctx context.Context, pubDomains map[string]struct{}, tx *sql.Tx) error {
	metaDataQueue, err := prepareMetaDataWithFloors(ctx, pubDomains, tx)
	if err != nil {
		return err
	}

	if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
		return fmt.Errorf("failed to process floors metadata queue for chunk: %w", err)
	}

	return nil
}

func createFloorData(chunk []dto.FloorUpdateRequest, pubDomain map[string]struct{}) []models.Floor {
	var floors []models.Floor

	for _, data := range chunk {
		floor := &dto.Floor{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Country:   data.Country,
			Device:    data.Device,
			Floor:     data.Floor,
		}

		if len(data.RuleId) > 0 {
			floor.RuleId = data.RuleId
		} else {
			floor.RuleId = floor.GetRuleID()
		}

		floors = append(floors, *floor.ToModel())
		pubDomain[data.Publisher+":"+data.Domain] = struct{}{}
	}

	return floors
}
func makeChunksFloor(requests []dto.FloorUpdateRequest) ([][]dto.FloorUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]dto.FloorUpdateRequest

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

func prepareMetaDataWithFloors(ctx context.Context, pubDomains map[string]struct{}, tx *sql.Tx) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for pubDomain := range pubDomains {
		pubDomainSplit := strings.Split(pubDomain, ":")

		modFloor, err := models.Floors(models.FloorWhere.Publisher.EQ(pubDomainSplit[0]), models.FloorWhere.Domain.EQ(pubDomainSplit[1]), models.FloorWhere.Active.EQ(true)).All(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("cannot get floor rules for publisher + demand partner id [%v]: %w", pubDomainSplit, err)
		}

		key := utils.FloorMetaDataKeyPrefix + ":" + pubDomain

		var finalRules []core.FloorRealtimeRecord
		if len(modFloor) > 0 {
			finalRules = core.CreateFloorMetadata(modFloor, finalRules)
		} else {
			finalRules = []core.FloorRealtimeRecord{}
		}

		finalOutput := struct {
			Rules []core.FloorRealtimeRecord `json:"rules"`
		}{Rules: finalRules}

		value, err := json.Marshal(finalOutput)
		if err != nil {
			return nil, eris.Wrap(err, "failed to marshal bulk floorRT to JSON")
		}

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           key,
			Value:         value,
		})
	}

	return metaDataQueue, nil
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
			constant.PostgresCurrentTime,
			constant.PostgresCurrentTime,
		)
	}

	return req
}
