package bulk

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strings"
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
