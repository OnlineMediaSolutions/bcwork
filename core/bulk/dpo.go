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
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
)

func (b *BulkService) BulkInsertDPO(ctx context.Context, requests []dto.DPORuleUpdateRequest) error {
	chunks, err := makeChunksDPO(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for dpos updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	demandPartners := make(map[string]struct{})
	oldMods, newMods, err := handleDpoRuleTable(ctx, tx, chunks, demandPartners, len(requests))
	if err != nil {
		return err
	}

	err = handleMetaDataRules(ctx, demandPartners, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in dpos bulk update: %w", err)
	}

	b.historyModule.SaveAction(ctx, oldMods, newMods, &history.HistoryOptions{Subject: history.DPOSubject, IsMultipleValuesExpected: true})

	return nil
}

func handleDpoRuleTable(
	ctx context.Context,
	tx *sql.Tx,
	chunks [][]dto.DPORuleUpdateRequest,
	demandPartners map[string]struct{},
	amountOfRequests int,
) ([]any, []any, error) {
	oldMods := make([]any, 0, amountOfRequests) // dpos before changes
	newMods := make([]any, 0, amountOfRequests) // dpos after changes

	for i, chunk := range chunks {
		dpos := prepareDPO(chunk, demandPartners)

		oldDpos, err := getOldDPO(ctx, dpos)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get old dpos for chunk %d: %w", i, err)
		}

		err = bulkInsertDPO(ctx, tx, dpos)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to process dpos bulk update for chunk %d: %w", i, err)
		}

		// appending previous and current mods for history processing
		oldMods = append(oldMods, oldDpos...)
		for j := 0; j < len(chunk); j++ {
			newMods = append(newMods, dpos[j])
		}
	}

	return oldMods, newMods, nil
}

func handleMetaDataRules(ctx context.Context, demandPartners map[string]struct{}, tx *sql.Tx) error {
	metaDataQueue, err := prepareDPODataForMetadata(ctx, demandPartners, tx)
	if err != nil {
		return fmt.Errorf("failed to prepare dpo data for metadata table %w", err)
	}

	if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
		return fmt.Errorf("failed to process dpos metadata queue for chunk: %w", err)
	}

	return nil
}

func makeChunksDPO(requests []dto.DPORuleUpdateRequest) ([][]dto.DPORuleUpdateRequest, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]dto.DPORuleUpdateRequest

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

func prepareDPO(chunk []dto.DPORuleUpdateRequest, demandPartners map[string]struct{}) []*models.DpoRule {
	dpos := make([]*models.DpoRule, 0, len(chunk))

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
			Active:        true,
		}

		demandPartners[data.DemandPartner] = struct{}{}
		dpos = append(dpos, DPOOptimizationRule.ToModel())
	}

	return dpos
}

func getOldDPO(ctx context.Context, dpos []*models.DpoRule) ([]any, error) {
	oldDpos := make([]any, 0, len(dpos))
	ids := make([]string, 0, len(dpos))

	for _, dpo := range dpos {
		ids = append(ids, dpo.RuleID)
	}

	oldMods, err := models.DpoRules(models.DpoRuleWhere.RuleID.IN(ids)).All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	m := make(map[string]*models.DpoRule)
	for _, oldMod := range oldMods {
		m[oldMod.RuleID] = oldMod
	}

	for _, dpo := range dpos {
		oldDpo, ok := m[dpo.RuleID]
		if !ok {
			oldDpos = append(oldDpos, nil)
			continue
		}
		oldDpo.Active = dpo.Active
		oldDpos = append(oldDpos, oldDpo)
	}

	return oldDpos, nil
}

func prepareDPODataForMetadata(ctx context.Context, demandPartners map[string]struct{}, tx *sql.Tx) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for demandPartner := range demandPartners {
		modDpos, err := models.DpoRules(models.DpoRuleWhere.DemandPartnerID.EQ(demandPartner)).All(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("cannot get dpo rules for demand partner id [%v]: %w", demandPartner, err)
		}

		dpos := make(core.DemandPartnerOptimizationRuleSlice, 0, len(modDpos))
		dpos.FromModel(modDpos)

		dposRT := core.DpoRT{
			DemandPartnerID: demandPartner,
			IsInclude:       false,
		}

		for _, dpo := range dpos {
			dposRT.Rules = append(dposRT.Rules, dpo.ToRtRule())
		}

		b, err := json.Marshal(dposRT)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON for dpo metadata value: %w", err)
		}

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           utils.DPOMetaDataKeyPrefix + ":" + demandPartner,
			Value:         b,
		})
	}

	return metaDataQueue, nil
}

func bulkInsertDPO(ctx context.Context, tx *sql.Tx, dpos []*models.DpoRule) error {
	req := prepareBulkInsertDPORequest(dpos)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertDPORequest(dpos []*models.DpoRule) *bulkInsertRequest {
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
			models.DpoRuleColumns.Active,
		},
		conflictColumns: []string{
			models.DpoRuleColumns.RuleID,
		},
		updateColumns: []string{
			models.DpoRuleColumns.Factor,
			models.DpoRuleColumns.UpdatedAt,
			models.DpoRuleColumns.Active,
		},
		valueStrings: make([]string, 0, len(dpos)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(dpos)*multiplier)

	for i, dpo := range dpos {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6,
				offset+7, offset+8, offset+9, offset+10, offset+11, offset+12, offset+13),
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
			constant.PostgresCurrentTime,
			constant.PostgresCurrentTime,
			dpo.Active,
		)
	}

	return req
}
