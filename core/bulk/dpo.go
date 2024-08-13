package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries"
	strconv "strconv"
	"strings"
	"time"
)

const insert_dpo_rule_query = `INSERT INTO dpo_rule (rule_id, demand_partner_id, publisher, domain, country, browser, os, device_type, placement_type, factor,created_at, updated_at) VALUES `

const on_conflict_query = `ON CONFLICT (rule_id) DO UPDATE SET country = EXCLUDED.country,
	factor = EXCLUDED.factor, device_type = EXCLUDED.device_type, domain = EXCLUDED.domain, placement_type = EXCLUDED.placement_type, updated_at = EXCLUDED.updated_at`

func MakeChunksDPO(requests []core.DPOUpdateRequest) ([][]core.DPOUpdateRequest, error) {
	chunkSize := viper.GetInt("api.chunkSize")
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

type DpoSlice []*core.DPOUpdateRequest

func (dp *DpoSlice) FromModel(slice models.DpoSlice) error {
	for _, mod := range slice {
		dpo := &core.DPOUpdateRequest{
			DemandPartner: mod.DemandPartnerID,
		}
		*dp = append(*dp, dpo)
	}
	return nil
}

func ProcessChunksDPO(c *fiber.Ctx, chunks [][]core.DPOUpdateRequest) error {
	errMessage, done := proceedBulkDpo(c, chunks)
	if done {
		log.Error().Err(errMessage).Msgf("Failed to process dpo_rule")
		return errMessage
	}

	errMessage, done = proceedBulkMetadata(c, chunks)
	if done {
		log.Error().Err(errMessage).Msgf("Failed to process dpo metadata queue")
		return errMessage
	}

	return nil
}

func proceedBulkMetadata(c *fiber.Ctx, chunks [][]core.DPOUpdateRequest) (error, bool) {
	tx, _ := bcdb.DB().BeginTx(c.Context(), nil)

	for i, chunk := range chunks {
		metaDataQueue := prepareMetadataDPO(chunk, context.Background())

		if err := BulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("Failed to process dpo metadata queue for chunk %d", i)
			return err, true
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction in DPO bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Failed to commit DPO transaction",
		}), true
	}
	return nil, false
}

func proceedBulkDpo(c *fiber.Ctx, chunks [][]core.DPOUpdateRequest) (error, bool) {

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Failed to begin transaction for inserting bulk dpo",
		}), false
	}

	for i, chunk := range chunks {
		dpos := prepareDataDPO(chunk)

		if err := bulkInsertDPO(tx, dpos); err != nil {
			log.Error().Err(err).Msgf("Failed to process bulk update for dpos %d", i)
			return err, true
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction in DPO bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit DPO transaction",
		}), true
	}
	return nil, false
}

func prepareDataDPO(chunk []core.DPOUpdateRequest) []models.DpoRule {
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

func prepareMetadataDPO(chunk []core.DPOUpdateRequest, ctx context.Context) []models.MetadataQueue {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		modDpos, _ := models.DpoRules(models.DpoRuleWhere.DemandPartnerID.EQ(data.DemandPartner)).All(ctx, bcdb.DB())

		dpos := make(core.DemandPartnerOptimizationRuleSlice, 0, 0)
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
			fmt.Println("Error marshaling JSON for metadata value:", err)
			continue
		}

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           "dpo:" + data.DemandPartner,
			Value:         b,
		})
	}
	return metaDataQueue
}

func bulkInsertDPO(tx *sql.Tx, dpos []models.DpoRule) error {
	createAt := time.Now().Format("2006-01-02 15") + ":00:00"
	values := make([]string, 0)
	for _, rec := range dpos {
		values = append(values, fmt.Sprintf(`('%s','%s','%s','%s','%s','%s','%s','%s','%s',%s,'%s','%s')`,
			rec.RuleID,
			rec.DemandPartnerID,
			rec.Publisher.String,
			rec.Domain.String,
			rec.Country.String,
			rec.Browser.String,
			rec.Os.String,
			rec.DeviceType.String,
			rec.PlacementType.String,
			[]byte(strconv.FormatFloat(rec.Factor, 'f', 0, 64)),
			createAt,
			createAt))
	}

	query := fmt.Sprint(insert_dpo_rule_query, strings.Join(values, ","))
	query += fmt.Sprintf(on_conflict_query)

	_, err := queries.Raw(query).Exec(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "Failed to update dpo_rule in bulk transaction")
	}

	return nil
}
