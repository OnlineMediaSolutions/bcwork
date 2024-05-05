package questclean

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strings"
	"time"
)

type Worker struct {
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	return nil
}

func (w *Worker) Do(ctx context.Context) error {

	expire := time.Now().AddDate(0, 0, -5).Format("2006-01-02")

	err := ProcessQuestInstance(ctx, "quest1", expire)
	if err != nil {
		return errors.Wrapf(err, "quest1 cleanup failed")
	}

	err = ProcessQuestInstance(ctx, "quest2", expire)
	if err != nil {
		return errors.Wrapf(err, "quest2 cleanup failed")
	}

	return nil

}

func ProcessQuestInstance(ctx context.Context, env string, expire string) error {

	log.Info().Msgf("cleanup and downsample '%s' expire '%s", env, expire)
	qdb, err := connect(env)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to quest1")
	}
	defer qdb.Close()

	tables := make([]struct {
		Table string `boil:"table"`
	}, 0)

	err = queries.Raw("SHOW TABLES").Bind(ctx, qdb.DB, &tables)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to quest1")
	}

	for _, t := range tables {
		if strings.Contains(t.Table, "telemetry") || strings.Index(t.Table, "sys.") == 0 {
			log.Info().Msgf("skip table '%s'", t.Table)
			continue
		}

		q := fmt.Sprintf("ALTER TABLE %s DROP PARTITION WHERE timestamp < to_timestamp('%s:00:00:00', 'yyyy-MM-dd:HH:mm:ss')", t.Table, expire)
		_, err := qdb.ExecContext(ctx, q)
		if err != nil {
			return errors.Wrapf(err, "failed to expire partitions for table '%s', date '%s'", t.Table, expire)
		}

		log.Info().Msgf("table '%s' on '%s' done", t.Table, env)

	}
	return nil
}

func (w *Worker) GetSleep() int {
	return int(0)
}
