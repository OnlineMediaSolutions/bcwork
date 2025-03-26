package adstxt

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
)

type AdsTxtManager interface {
	CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreatePublisherDomainAdsTxtLines(ctx context.Context, tx *sql.Tx, publisherID, domain string) error
	UpdateAdsTxtMetadata(ctx context.Context, data map[string]*dto.AdsTxtGroupedByDPData) error
	UpdateAdsTxtMaterializedViews(ctx context.Context) error
}

type AdsTxtModule struct {
}

func NewAdsTxtModule() *AdsTxtModule {
	// TODO: add views updates
	return &AdsTxtModule{}
}

func (a *AdsTxtModule) CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(demandPartnerConnectionQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(demandPartnerChildQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(seatOwnerQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreatePublisherDomainAdsTxtLines(ctx context.Context, tx *sql.Tx, domain, publisherID string) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesFromPublisherDomainTemplateQuery(), domain, publisherID)
}

func (a *AdsTxtModule) UpdateAdsTxtMaterializedViews(ctx context.Context) error {
	logger.Logger(ctx).Info().Msg("start refreshing ads txt views")

	updateStartTime := time.Now()
	viewsNames := []string{
		models.ViewNames.AdsTXTMainView,
		models.ViewNames.AdsTXTGroupByDPView,
	}
	query := `REFRESH MATERIALIZED VIEW %v;`

	for _, viewName := range viewsNames {
		viewUpdateStartTime := time.Now()
		_, err := bcdb.DB().Exec(fmt.Sprintf(query, viewName))
		if err != nil {
			return fmt.Errorf("cannot refresh view [%v]: %w", viewName, err)
		}

		logger.Logger(ctx).Info().Msgf("view [%v] successfully refreshed in %v", viewName, time.Since(viewUpdateStartTime).String())
	}

	logger.Logger(ctx).Info().Msgf("all ads.txt views successfully refreshed in %v", time.Since(updateStartTime).String())

	return nil
}
