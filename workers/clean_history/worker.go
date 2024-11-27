package clean_history

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strconv"
	"time"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
	skipInitRun bool
}

const delete_query = `DELETE from history where date <= ('%s') and user_id = ('%s');`
const automationUser = -2

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local_prod")
	w.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}
	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	if w.skipInitRun {
		fmt.Println("Skipping work as per the skip_init_run flag.")
		w.skipInitRun = false
		return nil
	}

	log.Info().Msg("Start to Remove old rows from history table")

	query := w.createQuery()
	err := deleteDataFromHistoryTable(ctx, query)
	if err != nil {
		return err
	}
	log.Info().Msg("Finished History table clean Process")
	return nil
}

func (w *Worker) createQuery() string {
	start := time.Now().UTC().Add(-72 * time.Hour).Truncate(24 * time.Hour).Format("2006-01-02")
	return fmt.Sprintf(delete_query, start, strconv.Itoa(automationUser))
}

func deleteDataFromHistoryTable(ctx context.Context, query string) error {

	_, err := queries.Raw(query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("error deleting data from history table: %w", err)
	}
	return nil
}

func (w *Worker) GetSleep() int {
	log.Info().Msg(fmt.Sprintf("next run in: %d seconds. V1.3.4", bccron.Next(w.Cron)))
	if w.Cron != "" {
		return bccron.Next(w.Cron)
	}
	return 0
}
