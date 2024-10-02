package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strings"
	"time"
)

const updateActiveColumnQuery = `UPDATE ads_txt
SET active = false
WHERE demand_partner_name = ('%s')`

func InsertDataToAdsTxt(ctx *fiber.Ctx, data core.MetadataUpdateRequest, now time.Time) error {
	demand := strings.Split(data.Key, ":")[1]
	err := updateActiveColumnInDB(ctx, demand)
	if err != nil {
		log.Error().Err(err).Str("body", string(ctx.Body())).Msg("failed to update columns in ads_txt table")
		return fmt.Errorf("failed to update columns in ads_txt table", err)
	}

	chunks := createAdsTxtChunks(&data, demand, ctx)

	tx, err := bcdb.DB().BeginTx(ctx.Context(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for adstxt: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		adsTxt := prepareAdsTxtData(chunk, demand, now)

		if err := bulkInsertAdsTxtQueue(ctx.Context(), tx, adsTxt); err != nil {
			log.Error().Err(err).Msgf("failed to process adstxt chunk %d", i)
			return fmt.Errorf("failed to process ads_txt chunk %d: %w", i, err)
		}

	}
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in ads_txt bulk update")
		return fmt.Errorf("failed to commit transaction in ads_txt bulk update: %w", err)
	}
	return nil
}

func createAdsTxtChunks(data *core.MetadataUpdateRequest, demand string, c *fiber.Ctx) [][]interface{} {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]interface{}

	raw, ok := data.Data.([]interface{})
	if !ok {
		log.Error().Str("body", string(c.Body())).Msg("failed to parse data.Data for demand: " + demand)
	}

	for i := 0; i < len(raw); i += chunkSize {
		end := i + chunkSize
		if end > len(raw) {
			end = len(raw)
		}
		chunk := raw[i:end]
		chunks = append(chunks, chunk)
	}
	return chunks
}

func prepareAdsTxtData(chunk []interface{}, demand string, now time.Time) []models.AdsTXT {

	var pubDomains []models.AdsTXT
	for _, d := range chunk {
		data := d.(map[string]interface{})
		pubDomains = append(pubDomains, models.AdsTXT{
			PublisherID:       data["pubid"].(string),
			Domain:            data["domain"].(string),
			DemandPartnerName: demand,
			Active:            true,
			CreatedAt:         now,
			UpdatedAt:         null.TimeFrom(now),
		})
	}
	return pubDomains
}

func bulkInsertAdsTxtQueue(ctx context.Context, tx *sql.Tx, adsTxt []models.AdsTXT) error {
	req := prepareBulkInsertAdsTxtRequest(adsTxt)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertAdsTxtRequest(adsTxts []models.AdsTXT) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.AdsTXT,
		columns: []string{
			models.AdsTXTColumns.PublisherID,
			models.AdsTXTColumns.Domain,
			models.AdsTXTColumns.DemandPartnerName,
			models.AdsTXTColumns.Active,
			models.AdsTXTColumns.CreatedAt,
			models.AdsTXTColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.AdsTXTColumns.PublisherID,
			models.AdsTXTColumns.Domain,
			models.AdsTXTColumns.DemandPartnerName,
		},
		updateColumns: []string{
			models.AdsTXTColumns.PublisherID,
			models.AdsTXTColumns.Domain,
			models.AdsTXTColumns.DemandPartnerName,
			models.AdsTXTColumns.Active,
			models.AdsTXTColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(adsTxts)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(adsTxts)*multiplier)

	for i, adsTxt := range adsTxts {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6),
		)
		req.args = append(req.args,
			adsTxt.PublisherID,
			adsTxt.Domain,
			adsTxt.DemandPartnerName,
			adsTxt.Active,
			currentTime,
			currentTime,
		)
	}

	return req
}

func updateActiveColumnInDB(c *fiber.Ctx, demand string) error {
	query := fmt.Sprintf(updateActiveColumnQuery, demand)

	_, err := queries.Raw(query).ExecContext(c.Context(), bcdb.DB())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to update the column active in ads_txt")
		return fmt.Errorf("failed to update active columns in ads_txt to false", err)

	}
	return nil
}
