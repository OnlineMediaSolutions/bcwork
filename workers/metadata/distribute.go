package metadata

import (
	"context"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Worker struct {
	Sleep       int
	Limit       int
	DatabaseEnv string `json:"dbenv"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	var err error
	w.Sleep, err = conf.GetIntValueWithDefault("sleep", 60)
	if err != nil {
		return errors.Wrapf(err, "failed to read 'sleep'")
	}

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("do")

	instances, err := GetMetadataInstances(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to pull metadata instances")
	}

	if len(instances) == 0 {
		log.Info().Msg("no metadata instances were found, exiting")
		return nil
	}

	bitwise := int64(0)
	for _, inst := range instances {
		bitwise |= inst.Bitwise
	}

	queue, err := models.MetadataQueues(models.MetadataQueueWhere.CommitedInstances.LT(bitwise),
		qm.OrderBy(models.MetadataQueueColumns.CreatedAt)).All(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to pull metadata queue records for update")
	}

	for _, rec := range queue {
		log.Info().Str("key", rec.Key).Str("transaction_id", rec.TransactionID).Msg("distributing metadata")
		for _, mi := range instances {
			err := mi.Updater.Update(ctx, rec)
			if err != nil {
				log.Error().Err(err).Msgf("failed to update record in metadata instance, retry in 1 minute(transaction_id:%s,instance:%s)", rec.TransactionID, mi.InstanceID)
				continue
			}
			err = mi.SetBit(ctx, rec)
			if err != nil {
				log.Error().Err(err).Msgf("failed to set instance bit, metadata was successful but is not registered and will retry(transaction_id:%s,instance:%s)", rec.TransactionID, mi.InstanceID)
				continue
			}

			log.Info().Msgf("metadata update transaction successfully completed(transaction_id:%s,instance:%s)", rec.TransactionID, mi.InstanceID)
		}
	}

	return nil

}

func (w *Worker) GetSleep() int {
	return w.Sleep
}
