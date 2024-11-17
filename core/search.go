package core

import (
	"context"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
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
		var start time.Time

		ticker := time.NewTicker(minutesToUpdate * time.Minute)
		viewName := models.ViewNames.GlobalSearchView
		query := `REFRESH MATERIALIZED VIEW ` + viewName + `;`

		for {
			select {
			case <-ticker.C:
				start = time.Now()
				_, err := bcdb.DB().Exec(query)
				if err != nil {
					logger.Logger(ctx).Err(err).Msgf("cannot refresh view [%v]", viewName)
					continue
				}
			case <-ctx.Done():
				return
			}
			logger.Logger(ctx).Debug().Msgf("view [%v] successfully refreshed in %v", viewName, time.Since(start).String())
		}
	}(ctx)

	return &SearchService{}
}

type SearchRequest struct {
	Query       string `json:"query"`
	SectionType string `json:"section_type"`
}

// func (s *SearchService) Search(ctx context.Context, ops *SearchRequest) ([]*dto.SearchResult, error) {
func (s *SearchService) Search(ctx context.Context, req *SearchRequest) ([]*models.GlobalSearchView, error) {
	query := null.StringFrom("%" + req.Query + "%")

	qmods := make([]qm.QueryMod, 0, 2)
	qmods = append(qmods, qm.Where(
		models.GlobalSearchViewColumns.PublisherID+" ILIKE $1"+" OR "+
			models.GlobalSearchViewColumns.PublisherName+" ILIKE $1"+" OR "+
			models.GlobalSearchViewColumns.Domain+" ILIKE $1"+" OR "+
			models.GlobalSearchViewColumns.DemandPartnerName+" ILIKE $1",
		query,
	))

	if req.SectionType != "" && req.SectionType != dto.AllSectionType {
		qmods = append(
			qmods,
			qm.Where(models.GlobalSearchViewColumns.SectionType+" = $2", null.StringFrom(req.SectionType)),
		)
	}

	mods, err := models.GlobalSearchViews(qmods...).All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	return mods, nil
}
