package adstxt

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type adsTxtLineTemplate struct {
	PublisherID               string   `boil:"publisher_id"`
	Domain                    string   `boil:"domain"`
	DemandPartnerConnectionID null.Int `boil:"demand_partner_connection_id"`
	DemandPartnerChildID      null.Int `boil:"demand_partner_child_id"`
	SeatOwnerID               null.Int `boil:"seat_owner_id"`
	DemandStatus              string   `boil:"demand_status"`
}

func createAdsTxtLine(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) error {
	var lines []*adsTxtLineTemplate
	err := queries.Raw(query, args...).
		Bind(ctx, tx, &lines)
	if err != nil {
		return err
	}

	chunks, err := makeAdsTxtChunks(lines)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		err := insertAdsTxtLines(ctx, tx, chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func makeAdsTxtChunks(lines []*adsTxtLineTemplate) ([][]*adsTxtLineTemplate, error) {
	chunkSize := viper.GetInt(config.APIChunkSizeKey)
	var chunks [][]*adsTxtLineTemplate

	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunk := lines[i:end]
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func insertAdsTxtLines(ctx context.Context, tx *sql.Tx, lines []*adsTxtLineTemplate) error {
	columns := []string{
		models.AdsTXTColumns.PublisherID,
		models.AdsTXTColumns.Domain,
		models.AdsTXTColumns.DemandPartnerConnectionID,
		models.AdsTXTColumns.DemandPartnerChildID,
		models.AdsTXTColumns.SeatOwnerID,
		models.AdsTXTColumns.DemandStatus,
		models.AdsTXTColumns.CreatedAt,
	}
	valueStrings := make([]string, 0, len(lines))
	multiplier := len(columns)
	args := make([]interface{}, 0, len(lines)*multiplier)

	for i, line := range lines {
		offset := i * multiplier
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%v, $%v, $%v, $%v, $%v, $%v, $%v)",
				offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7),
		)
		args = append(args,
			line.PublisherID,
			line.Domain,
			line.DemandPartnerConnectionID.Ptr(),
			line.DemandPartnerChildID.Ptr(),
			line.SeatOwnerID.Ptr(),
			line.DemandStatus,
			constant.PostgresCurrentTime,
		)
	}

	columnNames := strings.Join(columns, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", models.TableNames.AdsTXT, columnNames, strings.Join(valueStrings, ","))

	_, err := queries.Raw(query, args...).ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
