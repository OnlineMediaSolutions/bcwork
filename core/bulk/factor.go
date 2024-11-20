package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"

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

func (b *BulkService) BulkInsertFactors(ctx context.Context, requests []FactorUpdateRequest) error {
	chunks, err := makeChunksFactor(requests)
	if err != nil {
		return fmt.Errorf("failed to create chunks for factor updates: %w", err)
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	pubDomains := make(map[string]struct{})
	oldMods, newMods, err := handleBulkFactor(ctx, tx, chunks, pubDomains, len(requests))
	if err != nil {
		return err
	}

	err = handleMetaDataFactorRules(ctx, pubDomains, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in factor bulk update")
		return fmt.Errorf("failed to commit transaction in factor bulk update: %w", err)
	}

	b.historyModule.SaveAction(ctx, oldMods, newMods, &history.HistoryOptions{Subject: history.FactorSubject, IsMultipleValuesExpected: true})

	return nil
}

func handleMetaDataFactorRules(ctx context.Context, pubDomains map[string]struct{}, tx *sql.Tx) error {

	metaDataQueue, err := prepareMetaDataWithFactors(ctx, pubDomains, tx)
	if err != nil {
		log.Error().Err(err).Msgf("failed to prepare factor data for metadata table")
		return fmt.Errorf("failed to prepare factor data for metadata table %w", err)
	}

	if err := bulkInsertMetaDataQueue(ctx, tx, metaDataQueue); err != nil {
		log.Error().Err(err).Msgf("failed to process factor metadata queue for chunk")
		return fmt.Errorf("failed to process factor metadata queue for chunk: %w", err)
	}
	return nil
}

func prepareMetaDataWithFactors(ctx context.Context, pubDomains map[string]struct{}, tx *sql.Tx) ([]models.MetadataQueue, error) {
	var metaDataQueue []models.MetadataQueue

	for pubDomain := range pubDomains {
		pubDomainSplit := strings.Split(pubDomain, ":")

		modFactor, err := models.Factors(models.FactorWhere.Publisher.EQ(pubDomainSplit[0]), models.FactorWhere.Domain.EQ(pubDomainSplit[1])).All(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("cannot get factor rules for publisher + demand partner id [%v]: %w", pubDomainSplit, err)
		}

		key := utils.FactorMetaDataKeyPrefix + ":" + pubDomainSplit[0] + ":" + pubDomainSplit[1]

		var finalRules []core.FactorRealtimeRecord
		finalRules = core.CreateFactorMetadata(modFactor, finalRules)

		finalOutput := struct {
			Rules []core.FactorRealtimeRecord `json:"rules"`
		}{Rules: finalRules}

		value, err := json.Marshal(finalOutput)
		if err != nil {
			return nil, eris.Wrap(err, "failed to marshal bulk factorRT to JSON")
		}

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           key,
			Value:         value,
		})
	}

	return metaDataQueue, nil
}

func handleBulkFactor(
	ctx context.Context,
	tx *sql.Tx,
	chunks [][]FactorUpdateRequest,
	pubDomains map[string]struct{},
	amountOfRequests int,
) ([]any, []any, error) {
	oldMods := make([]any, 0, amountOfRequests)
	newMods := make([]any, 0, amountOfRequests)

	for i, chunk := range chunks {
		factors, oldFactors := createFactorsData(ctx, chunk, pubDomains)

		if err := bulkInsertFactors(ctx, tx, factors); err != nil {
			log.Error().Err(err).Msgf("failed to process factor bulk update for chunk %d", i)
			return nil, nil, fmt.Errorf("failed to process factor bulk update for chunk %d: %w", i, err)
		}

		oldMods = append(oldMods, oldFactors...)
		for j := 0; j < len(chunk); j++ {
			newMods = append(newMods, factors[j])
		}
	}

	return oldMods, newMods, nil
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

func createFactorsData(ctx context.Context, chunk []FactorUpdateRequest, pubDomain map[string]struct{}) ([]*models.Factor, []any) {
	factors := make([]*models.Factor, 0, len(chunk))
	oldFactors := make([]any, 0, len(chunk))
	ids := make([]string, 0, len(chunk))

	for _, data := range chunk {
		factor := &core.Factor{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.Device,
			Factor:    data.Factor,
			Country:   data.Country,
		}

		factor.RuleId = factor.GetRuleID()
		ids = append(ids, factor.RuleId)
		factors = append(factors, factor.ToModel())
		pubDomain[data.Publisher+":"+data.Domain] = struct{}{}
	}

	oldMods, err := models.Factors(models.FactorWhere.RuleID.IN(ids)).All(ctx, bcdb.DB())
	if err != nil {
		return nil, nil
	}

	m := make(map[string]*models.Factor)
	for _, oldMod := range oldMods {
		m[oldMod.RuleID] = oldMod
	}

	for _, factor := range factors {
		oldFactor, ok := m[factor.RuleID]
		if !ok {
			oldFactors = append(oldFactors, nil)
			continue
		}
		oldFactors = append(oldFactors, oldFactor)
	}

	return factors, oldFactors
}

func bulkInsertFactors(ctx context.Context, tx *sql.Tx, factors []*models.Factor) error {
	req := prepareBulkInsertFactorsRequest(factors)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertFactorsRequest(factors []*models.Factor) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.Factor,
		columns: []string{
			models.FactorColumns.Publisher,
			models.FactorColumns.Domain,
			models.FactorColumns.Device,
			models.FactorColumns.Country,
			models.FactorColumns.Factor,
			models.FactorColumns.RuleID,
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
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8),
		)
		req.args = append(req.args,
			factor.Publisher,
			factor.Domain,
			factor.Device,
			factor.Country,
			factor.Factor,
			factor.RuleID,
			currentTime,
			currentTime,
		)
	}

	return req
}
