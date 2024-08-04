package bulk

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
	"time"
)

type FloorUpdateRequest struct {
	utils.MetadataKey
	Publisher string  `json:"publisher"`
	Domain    string  `json:"domain"`
	Device    string  `json:"device"`
	Floor     float64 `json:"floor"`
	Country   string  `json:"country"`
}

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
		floors, metaDataQueue := prepareDataFloor(chunk)

		if err := bulkInsertFloor(c, tx, floors); err != nil {
			log.Error().Err(err).Msgf("failed to process bulk update for floor chunk %d", i)
			return err
		}

		if err := BulkInsertMetaDataQueue(c, tx, metaDataQueue); err != nil {
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

func prepareDataFloor(chunk []FloorUpdateRequest) ([]models.Floor, []models.MetadataQueue) {
	var floors []models.Floor
	var metaDataQueue []models.MetadataQueue

	for _, data := range chunk {
		floors = append(floors, models.Floor{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.Device,
			Floor:     data.Floor,
			Country:   data.Country,
		})

		metadataKey := utils.MetadataKey{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			Device:    data.Device,
			Country:   data.Country,
		}
		key := utils.CreateMetadataKey(metadataKey, "price:floor")

		metaDataQueue = append(metaDataQueue, models.MetadataQueue{
			Key:           key,
			TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(data.Floor, 'f', 2, 64)),
		})
	}

	return floors, metaDataQueue
}

func bulkInsertFloor(c *fiber.Ctx, tx *sql.Tx, floors []models.Floor) error {
	columns := []string{"publisher", "domain", "device", "floor", "country", "created_at", "updated_at"}
	conflictColumns := []string{"publisher", "domain", "device", "country"}
	updateColumns := []string{"floor = EXCLUDED.floor", "created_at = EXCLUDED.created_at", "updated_at = EXCLUDED.updated_at"}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, floor := range floors {
		values = append(values, floor.Publisher, floor.Domain, floor.Device, floor.Floor, floor.Country, currTime, currTime)
	}

	return InsertInBulk(c, tx, "floor", columns, values, conflictColumns, updateColumns)
}
