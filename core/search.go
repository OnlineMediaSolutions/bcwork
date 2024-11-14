package core

import (
	"context"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SearchService struct{}

func NewSearchService(ctx context.Context) *SearchService {
	// refreshing view every N minutes
	go func(ctx context.Context) {
		const minutesToUpdate = 10

		ticker := time.NewTicker(minutesToUpdate * time.Minute)
		viewName := models.ViewNames.GlobalSearchView
		query := `REFRESH MATERIALIZED VIEW` + viewName + `;`

		for {
			select {
			case <-ticker.C:
				_, err := bcdb.DB().Exec(query)
				if err != nil {
					logger.Logger(ctx).Err(err).Msgf("cannot refresh view [%v]", viewName)
				}
			case <-ctx.Done():
				return
			}
			logger.Logger(ctx).Debug().Msgf("view [%v] successfully refreshed", viewName)
		}
	}(ctx)

	return &SearchService{}
}

type SearchRequest struct {
	Query   string
	Section string
}

// func (s *SearchService) Search(ctx context.Context, ops *SearchRequest) ([]*dto.SearchResult, error) {
func (s *SearchService) Search(ctx context.Context, ops *SearchRequest) ([]*models.GlobalSearchView, error) {
	query := "%" + ops.Query + "%"

	mods, err := models.GlobalSearchViews(
		qm.Or2(models.GlobalSearchViewWhere.PublisherID.ILIKE(null.StringFrom(query))),
		qm.Or2(models.GlobalSearchViewWhere.PublisherName.ILIKE(null.StringFrom(query))),
		qm.Or2(models.GlobalSearchViewWhere.Domain.ILIKE(null.StringFrom(query))),
		qm.Or2(models.GlobalSearchViewWhere.DemandPartnerName.ILIKE(null.StringFrom(query))),
	).All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	return mods, nil
}
