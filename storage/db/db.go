package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	adstxt "github.com/m6yf/bcwork/modules/ads_txt"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type PublisherSyncStorage interface {
	HadLoadingErrorLastTime(ctx context.Context, key string) bool
	UpsertPublisherAndDomains(ctx context.Context, publisher *models.Publisher, domains []*models.PublisherDomain, updateColumns boil.Columns) error
	SaveResultOfLastSync(ctx context.Context, key string, hasErrors bool) error
}

type DB struct {
	dbClient                    *sqlx.DB
	adsTxtModule                adstxt.AdsTxtLinesCreater
	isNeededToCreateAdsTxtLines bool
}

var _ PublisherSyncStorage = (*DB)(nil)

func New(dbClient *sqlx.DB, adsTxtModule adstxt.AdsTxtLinesCreater, isNeededToCreateAdsTxtLines bool) *DB {
	return &DB{
		dbClient:                    dbClient,
		adsTxtModule:                adsTxtModule,
		isNeededToCreateAdsTxtLines: isNeededToCreateAdsTxtLines,
	}
}

func (d *DB) UpsertPublisherAndDomains(
	ctx context.Context,
	publisher *models.Publisher,
	domains []*models.PublisherDomain,
	updateColumns boil.Columns,
) error {
	tx, err := d.dbClient.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateOnConflict := true
	conflictColumns := []string{models.PublisherColumns.PublisherID}

	err = publisher.Upsert(ctx, tx, updateOnConflict, conflictColumns, updateColumns, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to upsert row [%v] in publisher table: %w", publisher.PublisherID, err)
	}

	for _, domain := range domains {
		isExisted, err := models.PublisherDomains(
			models.PublisherDomainWhere.Domain.EQ(domain.Domain),
			models.PublisherDomainWhere.PublisherID.EQ(domain.PublisherID),
		).Exists(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to check existance of domain [%v:%v]: %w", domain.PublisherID, domain.Domain, err)
		}

		if isExisted {
			// if there is domain in db already, then updating its mirror publisher id
			domain.UpdatedAt = null.TimeFrom(time.Now())
			_, err := domain.Update(ctx, tx, boil.Whitelist(models.PublisherDomainColumns.MirrorPublisherID, models.PublisherDomainColumns.UpdatedAt))
			if err != nil {
				return fmt.Errorf("failed to update domain [%v] in publisher domain table for publisherId [%v]: %w", domain.Domain, domain.PublisherID, err)
			}
		} else {
			err := domain.Insert(ctx, tx, boil.Infer())
			if err != nil {
				return fmt.Errorf("failed to insert domain [%v] to publisher domain table for publisherId [%v]: %w", domain.Domain, domain.PublisherID, err)
			}

			// if it was new domain, then creating ads txt lines for it
			if d.isNeededToCreateAdsTxtLines {
				err := d.adsTxtModule.CreatePublisherDomainAdsTxtLines(ctx, tx, domain.Domain, domain.PublisherID)
				if err != nil {
					return eris.Wrapf(err, "failed to create ads txt lines for publisher [%v], domain [%v]", domain.PublisherID, domain.Domain)
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
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
