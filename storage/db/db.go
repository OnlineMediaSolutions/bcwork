package db

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type PublisherSyncStorage interface {
	HadLoadingErrorLastTime(ctx context.Context, key string) bool
	UpsertPublisher(ctx context.Context, publisher *models.Publisher, updateColumns boil.Columns) error
	InsertPublisherDomain(ctx context.Context, domain *models.PublisherDomain) error
	SaveResultOfLastSync(ctx context.Context, key string, hasErrors bool) error
}

type DB struct {
	dbClient *sqlx.DB
}

var _ PublisherSyncStorage = (*DB)(nil)

func New(dbClient *sqlx.DB) *DB {
	return &DB{
		dbClient: dbClient,
	}
}

func (d *DB) UpsertPublisher(ctx context.Context, publisher *models.Publisher, updateColumns boil.Columns) error {
	updateOnConflict := true
	conflictColumns := []string{models.PublisherColumns.PublisherID}

	err := publisher.Upsert(ctx, d.dbClient, updateOnConflict, conflictColumns, updateColumns, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) InsertPublisherDomain(ctx context.Context, domain *models.PublisherDomain) error {
	err := domain.Insert(ctx, d.dbClient, boil.Infer())
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil
		}

		return err
	}

	return nil
}

func (d *DB) HadLoadingErrorLastTime(ctx context.Context, key string) bool {
	lastResult, err := models.PublisherSyncs(qm.Where(models.PublisherSyncColumns.Key+" = ?", key)).One(ctx, d.dbClient)
	if err != nil {
		return false
	}

	return lastResult.HadError
}

func (d *DB) SaveResultOfLastSync(ctx context.Context, key string, hasErrors bool) error {
	result := models.PublisherSync{
		Key:      key,
		HadError: hasErrors,
	}

	updateOnConflict := true
	conflictColumns := []string{models.PublisherSyncColumns.Key}

	err := result.Upsert(ctx, d.dbClient, updateOnConflict, conflictColumns, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	const errDuplicateKey = "duplicate key value violates unique constraint"

	return strings.Contains(err.Error(), errDuplicateKey)
}
