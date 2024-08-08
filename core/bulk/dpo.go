package bulk

import (
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
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

func ProcessChunksDPO(c *fiber.Ctx, chunks [][]core.DPOUpdateRequest) error {

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		dpos, metaDataQueue := prepareDataDPO(chunk)
		print(metaDataQueue)

		if err := bulkInsertDPO(c, tx, dpos); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for dpos %d", i)
			return err
		}

		//TODO need to write to the METADATA_QUEUE
		//if err := BulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
		//	log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
		//	return err
		//}
	}
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in DPO bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit DPO transaction",
		})
	}
	return nil
}

func prepareDataDPO(chunk []core.DPOUpdateRequest) ([]models.DpoRule, []models.MetadataQueue) {
	var dpos []models.DpoRule
	var metaDataQueue []models.MetadataQueue

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

		metadataKey := utils.MetadataKey{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.DeviceType,
			Country:   data.Country,
		}
		key := utils.CreateMetadataKey(metadataKey, "price:floor")

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			Key:           key,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(data.Factor, 'f', 2, 64)),
		})
	}

	return dpos, metaDataQueue
}

func bulkInsertDPO(c *fiber.Ctx, tx *sql.Tx, dpos []models.DpoRule) error {
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
