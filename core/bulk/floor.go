package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type FloorStruct struct {
	RuleId        string  `boil:"rule_id" json:"rule_id" toml:"rule_id" yaml:"rule_id"`
	Publisher     string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	PublisherName string  `boil:"publisher_name" json:"publisher_name" toml:"publisher_name" yaml:"publisher_name"`
	Domain        string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Country       string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Device        string  `boil:"device" json:"device" toml:"device" yaml:"device"`
	Floor         float64 `boil:"floor" json:"floor" toml:"floor" yaml:"floor"`
	Browser       string  `boil:"browser" json:"browser" toml:"browser" yaml:"browser"`
	OS            string  `boil:"os" json:"os" toml:"os" yaml:"os"`
	PlacementType string  `boil:"placement_type" json:"placement_type" toml:"placement_type" yaml:"placement_type"`
}

func MakeChunksFloor(requests []core.FloorUpdateRequest) ([][]core.FloorUpdateRequest, error) {
	chunkSize := viper.GetInt("api.chunkSize")
	var chunks [][]core.FloorUpdateRequest

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

func ProcessChunksFloor(c *fiber.Ctx, chunks [][]core.FloorUpdateRequest) error {

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to begin transaction",
		})
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		floors, metaDataQueue := prepareDataFloor(chunk)

		if err := bulkInsertFloor(c, tx, floors); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for floor chunk %d", i)
			return err
		}

		if err := InsertRegMetaDataQueue(c, tx, metaDataQueue); err != nil {
			log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in floor bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to commit transaction",
		})
	}

	return nil
}

func prepareDataFloor(chunk []core.FloorUpdateRequest) ([]models.Floor, []models.MetadataQueue) {
	var floors []models.Floor
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {

		if data.RuleId == "" {
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
			data.RuleId = floor.GetRuleID()
		}

		floors = append(floors, models.Floor{
			Publisher:     data.Publisher,
			Domain:        data.Domain,
			Device:        data.Device,
			Floor:         data.Floor,
			Country:       data.Country,
			Os:            null.StringFrom(data.OS),
			PlacementType: null.StringFrom(data.PlacementType),
			Browser:       null.StringFrom(data.Browser),
			RuleID:        data.RuleId,
		})

		metadata, _ := SendFloorToRT(context.Background(), data)
		metaDataQueue = append(metaDataQueue, metadata...)

	}

	return floors, metaDataQueue
}

func SendFloorToRT(c context.Context, updateRequest core.FloorUpdateRequest) ([]models.MetadataQueue, error) {
	const PREFIX string = "price:floor:v2"
	modFloor, err := core.FloorQuery(c, updateRequest)

	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to fetch floors")
	}
	var finalRules []core.FloorRealtimeRecord

	finalRules = core.CreateFloorMetadata(modFloor, finalRules, updateRequest)

	finalOutput := struct {
		Rules []core.FloorRealtimeRecord `json:"rules"`
	}{Rules: finalRules}

	value, err := json.Marshal(finalOutput)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to marshal floorRT to JSON")
	}

	key := utils.GetMetadataObject(updateRequest)
	metadataKey := utils.CreateMetadataKey(key, PREFIX)
	metadataValue := utils.CreateMetadataObject(updateRequest, metadataKey, value)

	return []models.MetadataQueue{metadataValue}, nil
}

func bulkInsertFloor(c *fiber.Ctx, tx *sql.Tx, floors []models.Floor) error {
	columns := []string{"rule_id", "publisher", "domain", "device", "floor", "country", "os", "browser", "placement_type", "created_at", "updated_at"}
	conflictColumns := []string{"publisher", "domain", "device", "country"}
	updateColumns := []string{
		"floor = EXCLUDED.floor",
		"os = EXCLUDED.os",
		"browser = EXCLUDED.browser",
		"placement_type = EXCLUDED.placement_type",
		"created_at = EXCLUDED.created_at",
		"updated_at = EXCLUDED.updated_at",
	}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, floor := range floors {
		values = append(values,
			floor.RuleID,
			floor.Publisher,
			floor.Domain,
			floor.Device,
			floor.Floor,
			floor.Country,
			floor.Os,
			floor.Browser,
			floor.PlacementType,
			currTime, currTime)
	}

	return InsertInBulk(c, tx, "floor", columns, values, conflictColumns, updateColumns)
}
