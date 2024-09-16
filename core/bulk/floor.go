package bulk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strconv"
	"strings"
	"time"
)

type FloorUpdateRequest struct {
	RuleId        string  `json:"rule_id"`
	Publisher     string  `json:"publisher"`
	Domain        string  `json:"domain"`
	Device        string  `json:"device"`
	Floor         float64 `json:"floor"`
	Country       string  `json:"country"`
	Browser       string  `json:"browser"`
	OS            string  `json:"os"`
	PlacementType string  `json:"placement_type"`
}

func (f FloorUpdateRequest) GetPublisher() string     { return f.Publisher }
func (f FloorUpdateRequest) GetDomain() string        { return f.Domain }
func (f FloorUpdateRequest) GetDevice() string        { return f.Device }
func (f FloorUpdateRequest) GetCountry() string       { return f.Country }
func (f FloorUpdateRequest) GetBrowser() string       { return f.Browser }
func (f FloorUpdateRequest) GetOS() string            { return f.OS }
func (f FloorUpdateRequest) GetPlacementType() string { return f.PlacementType }

const insert_floor_rule_query = `INSERT INTO floor (rule_id, publisher, domain, country, browser, os, device, placement_type, floor,created_at, updated_at) VALUES `

const floor_on_conflict_query = `ON CONFLICT (rule_id) DO UPDATE SET floor = EXCLUDED.floor, updated_at = NOW()`

func MakeChunksFloor(requests []FloorUpdateRequest) ([][]FloorUpdateRequest, error) {
	chunkSize := viper.GetInt("api.chunkSize")
	var chunks [][]FloorUpdateRequest

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

func ProcessChunksFloor(c *fiber.Ctx, chunks [][]FloorUpdateRequest) error {
	if err := processFloorChunks(c, chunks); err != nil {
		return err
	}

	if err := processMetadataChunks(c, chunks); err != nil {
		return err
	}

	return nil
}

func processFloorChunks(c *fiber.Ctx, chunks [][]FloorUpdateRequest) error {
	for i, chunk := range chunks {
		tx, err := bcdb.DB().BeginTx(c.Context(), nil) // Start a new transaction for each chunk
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to begin transaction",
			})
		}

		floors := prepareDataFloor(chunk)

		if err := bulkInsertFloor(tx, floors); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for floor chunk %d", i)
			tx.Rollback() // Rollback if insertion fails
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msgf("failed to commit transaction for floor chunk %d", i)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to commit floor data",
			})
		}

		log.Info().Msgf("Successfully processed floor chunk %d", i)
	}

	return nil
}

func processMetadataChunks(c *fiber.Ctx, chunks [][]FloorUpdateRequest) error {
	for i, chunk := range chunks {
		tx, err := bcdb.DB().BeginTx(c.Context(), nil) // Start a new transaction for metadata insertion
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to begin transaction for metadata",
			})
		}

		metaDataValue := prepareMetadataFloor(chunk, context.Background())

		if err := BulkInsertMetaDataQueue(c, tx, metaDataValue); err != nil {
			log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
			tx.Rollback() // Rollback if metadata insertion fails
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msgf("failed to commit transaction for metadata chunk %d", i)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to commit metadata data",
			})
		}

		log.Info().Msgf("Successfully processed metadata chunk %d", i)
	}

	return nil
}

func prepareMetadataFloor(chunk []FloorUpdateRequest, ctx context.Context) []models.MetadataQueue {
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		const PREFIX string = "price:floor:v2"
		modFloor, _ := core.FloorQuery(ctx, core.FloorUpdateRequest(data))

		var finalRules []core.FloorRealtimeRecord

		finalRules = core.CreateFloorMetadata(modFloor, finalRules)

		finalOutput := struct {
			Rules []core.FloorRealtimeRecord `json:"rules"`
		}{Rules: finalRules}

		value, _ := json.Marshal(finalOutput)

		key := utils.GetMetadataObject(data)
		metadataKey := utils.CreateMetadataKey(key, PREFIX)
		metadataValue := utils.CreateMetadataObject(data, metadataKey, value)

		fmt.Println("modFloor:", modFloor)
		fmt.Println("finalRules:", finalRules)
		fmt.Println("metadataKey:", metadataKey)
		fmt.Println("metadataValue:", metadataValue)

		metaDataQueue = append(metaDataQueue, metadataValue)

	}

	return metaDataQueue
}

func prepareDataFloor(chunk []FloorUpdateRequest) []models.Floor {
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
			Device:        floor.Device,
			Floor:         floor.Floor,
			Country:       floor.Country,
			Os:            floor.OS,
			Browser:       floor.Browser,
			PlacementType: floor.PlacementType,
			RuleID:        floor.GetRuleID(),
		})
	}

	return floors
}

func bulkInsertFloor(tx *sql.Tx, floors []models.Floor) error {
	createAt := time.Now().Format("2006-01-02 15") + ":00:00"
	values := make([]string, 0)
	for _, rec := range floors {
		values = append(values, fmt.Sprintf(`('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')`,
			rec.RuleID,
			rec.Publisher,
			rec.Domain,
			rec.Country,
			rec.Browser,
			rec.Os,
			rec.Device,
			rec.PlacementType,
			[]byte(strconv.FormatFloat(rec.Floor, 'f', 0, 64)),
			createAt,
			createAt))
	}

	query := fmt.Sprint(insert_floor_rule_query, strings.Join(values, ","))
	query += fmt.Sprintf(floor_on_conflict_query)

	_, err := queries.Raw(query).Exec(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
