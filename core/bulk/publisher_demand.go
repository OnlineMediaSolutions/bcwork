package bulk

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

const updateActiveColumnQuery = `UPDATE publisher_demand
SET ads_txt_status = false
WHERE demand_partner_id = ('%s')`

func InsertDataToAdsTxt(ctx *fiber.Ctx, request dto.PublisherDomainRequest, now time.Time) error {
	err := updateActiveColumnInDB(ctx, request.DemandParnerId)
	if err != nil {
		return fmt.Errorf("failed to update columns in publisher demand table: %w", err)
	}

	chunks := createAdsTxtChunks(request)
	tx, err := bcdb.DB().BeginTx(ctx.Context(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for publisher demand: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		publisherDemand := preparePublisherDemandData(chunk, request.DemandParnerId, now)

		if err := bulkInsertPublisherDemandQueue(ctx.Context(), tx, publisherDemand); err != nil {
			return fmt.Errorf("failed to process publisher_demand chunk %d: %w", i, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction in publisher_demand bulk update: %w", err)
	}

	return nil
}

func createAdsTxtChunks(request dto.PublisherDomainRequest) [][]dto.PublisherDomainData {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]dto.PublisherDomainData

	for i := 0; i < len(request.Data); i += chunkSize {
		end := i + chunkSize
		if end > len(request.Data) {
			end = len(request.Data)
		}
		chunk := request.Data[i:end]
		chunks = append(chunks, chunk)
	}

	return chunks
}

func preparePublisherDemandData(chunk []dto.PublisherDomainData, demand string, now time.Time) []models.PublisherDemand {
	var pubDomains []models.PublisherDemand
	for _, data := range chunk {
		pubDomains = append(pubDomains, models.PublisherDemand{
			PublisherID:     data.PubId,
			Domain:          data.Domain,
			DemandPartnerID: demand,
			AdsTXTStatus:    data.AdsTxtStatus,
			Active:          true,
			CreatedAt:       now,
			UpdatedAt:       null.TimeFrom(now),
		})
	}

	return pubDomains
}

func bulkInsertPublisherDemandQueue(ctx context.Context, tx *sql.Tx, adsTxt []models.PublisherDemand) error {
	req := prepareBulkInsertAdsTxtRequest(adsTxt)

	return bulkInsert(ctx, tx, req)
}

func prepareBulkInsertAdsTxtRequest(publisherDemands []models.PublisherDemand) *bulkInsertRequest {
	req := &bulkInsertRequest{
		tableName: models.TableNames.PublisherDemand,
		columns: []string{
			models.PublisherDemandColumns.PublisherID,
			models.PublisherDemandColumns.Domain,
			models.PublisherDemandColumns.DemandPartnerID,
			models.PublisherDemandColumns.Active,
			models.PublisherDemandColumns.AdsTXTStatus,
			models.PublisherDemandColumns.CreatedAt,
			models.PublisherDemandColumns.UpdatedAt,
		},
		conflictColumns: []string{
			models.PublisherDemandColumns.PublisherID,
			models.PublisherDemandColumns.Domain,
			models.PublisherDemandColumns.DemandPartnerID,
		},
		updateColumns: []string{
			models.PublisherDemandColumns.PublisherID,
			models.PublisherDemandColumns.Domain,
			models.PublisherDemandColumns.DemandPartnerID,
			models.PublisherDemandColumns.Active,
			models.PublisherDemandColumns.AdsTXTStatus,
			models.PublisherDemandColumns.UpdatedAt,
		},
		valueStrings: make([]string, 0, len(publisherDemands)),
	}

	multiplier := len(req.columns)
	req.args = make([]interface{}, 0, len(publisherDemands)*multiplier)

	for i, publisherDemand := range publisherDemands {
		offset := i * multiplier
		req.valueStrings = append(req.valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7),
		)
		req.args = append(req.args,
			publisherDemand.PublisherID,
			publisherDemand.Domain,
			publisherDemand.DemandPartnerID,
			publisherDemand.Active,
			publisherDemand.AdsTXTStatus,
			constant.PostgresCurrentTime,
			constant.PostgresCurrentTime,
		)
	}

	return req
}

func updateActiveColumnInDB(c *fiber.Ctx, demand string) error {
	query := fmt.Sprintf(updateActiveColumnQuery, demand)

	_, err := queries.Raw(query).ExecContext(c.Context(), bcdb.DB())
	if err != nil {
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to update the column active in publisher_demand table")

		return fmt.Errorf("failed to update active columns in publisher_demand to false %w", err)
	}

	return nil
}
