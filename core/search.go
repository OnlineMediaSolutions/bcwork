package core

import (
	"context"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/spf13/viper"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type SearchService struct{}

func NewSearchService(ctx context.Context) *SearchService {
	// refreshing view every N minutes
	go func(ctx context.Context) {
		var (
			start                  time.Time
			defaultMinutesToUpdate time.Duration = 10 * time.Minute
		)

		minutesToUpdate := time.Duration(viper.GetInt(config.SearchViewUpdateRateKey)) * time.Minute
		if minutesToUpdate == 0 {
			minutesToUpdate = defaultMinutesToUpdate
		}

		ticker := time.NewTicker(minutesToUpdate)
		viewName := models.ViewNames.SearchView
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
				ticker.Stop()
				return
			}
			logger.Logger(ctx).Debug().Msgf("view [%v] successfully refreshed in %v", viewName, time.Since(start).String())
		}
	}(ctx)

	return &SearchService{}
}

func (s *SearchService) Search(ctx context.Context, req *dto.SearchRequest) (map[string][]dto.SearchResult, error) {
	query := null.StringFrom("%" + req.Query + "%")

	qmods := make([]qm.QueryMod, 0, 2)
	qmods = append(qmods, qm.Where(models.SearchViewColumns.Query+" ILIKE $1", query))

	if req.SectionType != "" {
		qmods = append(
			qmods,
			qm.Where(models.SearchViewColumns.SectionType+" = $2", null.StringFrom(req.SectionType)),
		)
	}

	mods, err := models.SearchViews(qmods...).All(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	return dto.PrepareSearchResults(mods, req.SectionType), nil
}
