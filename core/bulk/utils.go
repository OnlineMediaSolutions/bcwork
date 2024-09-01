package bulk

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"time"
)

func InsertInBulk(c *fiber.Ctx, tx *sql.Tx, tableName string, columns []string, values []interface{}, conflictColumns, updateColumns []string) error {
	columnCount := len(columns)
	valueStrings := make([]string, 0, len(values)/columnCount)
	valueArgs := make([]interface{}, 0, len(values))

	for i := 0; i < len(values)/columnCount; i++ {
		placeholders := make([]string, columnCount)
		for j := 0; j < columnCount; j++ {
			placeholders[j] = fmt.Sprintf("$%d", i*columnCount+j+1)
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		valueArgs = append(valueArgs, values[i*columnCount:(i+1)*columnCount]...)
	}

	columnNames := strings.Join(columns, ", ")
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, columnNames, strings.Join(valueStrings, ","))

	if conflictColumns != nil && updateColumns != nil {
		stmt += fmt.Sprintf(" ON CONFLICT (%s) DO UPDATE SET %s", strings.Join(conflictColumns, ", "), strings.Join(updateColumns, ", "))
	}

	log.Info().Msgf("Executing bulk insert for %s: %s", tableName, stmt)
	if _, err := tx.ExecContext(c.Context(), stmt, valueArgs...); err != nil {
		log.Error().Err(err).Msgf("Failed to execute bulk insert for %s: %s", tableName, stmt)
		return fmt.Errorf("failed to insert into %s in bulk: %w", tableName, err)
	}

	return nil
}

func generatePlaceholders(rowCount, colCount int) string {
	rowTemplate := strings.Builder{}
	for i := 1; i <= colCount; i++ {
		if i > 1 {
			rowTemplate.WriteString(", ")
		}
		rowTemplate.WriteString(fmt.Sprintf("$%d", i))
	}

	rowTemplateStr := fmt.Sprintf("(%s)", rowTemplate.String())

	placeholders := strings.Repeat(rowTemplateStr+", ", rowCount)

	return placeholders[:len(placeholders)-2]
}

func InsertRegMetaDataQueue(c *fiber.Ctx, tx *sql.Tx, metaDataQueue []models.MetadataQueue) error {
	columns := []string{"key", "transaction_id", "value", "version", "commited_instances", "created_at", "updated_at"}
	currTime := time.Now().In(boil.GetLocation())

	var values []interface{}

	for _, metaData := range metaDataQueue {
		values = append(values, metaData.Key, metaData.TransactionID, metaData.Value, nil, 0, currTime, currTime)
	}

	numRows := len(metaDataQueue)

	if numRows == 0 {
		return nil
	}

	placeholderStr := generatePlaceholders(numRows, len(columns))

	query := fmt.Sprintf(
		"INSERT INTO metadata_queue (%s) VALUES %s",
		strings.Join(columns, ", "),
		placeholderStr,
	)

	_, err := tx.ExecContext(c.Context(), query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert into metadata_queue in bulk: %w", err)
	}

	return nil
}

func BulkInsertMetaDataQueue(c *fiber.Ctx, tx *sql.Tx, metaDataQueue []models.MetadataQueue) error {
	columns := []string{"key", "transaction_id", "value", "commited_instances", "created_at", "updated_at"}

	var values []interface{}
	currTime := time.Now().In(boil.GetLocation())
	for _, metaData := range metaDataQueue {
		values = append(values, metaData.Key, metaData.TransactionID, metaData.Value, 0, currTime, currTime)
	}

	return InsertInBulk(c, tx, "metadata_queue", columns, values, nil, nil)
}
