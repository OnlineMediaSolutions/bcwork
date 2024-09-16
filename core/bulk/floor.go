package bulk

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
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

	tx, err := bcdb.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to begin transaction",
		})
	}
	defer tx.Rollback()

	for i, chunk := range chunks {
		floors := prepareDataFloor(chunk)

		if err := bulkInsertFloor(tx, floors); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for floor chunk %d", i)
			return err
		}

		//if err := BulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
		//	log.Error().Err(err).Msgf("failed to process metadata queue for chunk %d", i)
		//	return err
		//}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in floor bulk update")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to commit transaction",
		})
	}

	return nil
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
