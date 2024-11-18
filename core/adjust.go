package core

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AdjustService struct {
	historyModule history.HistoryModule
}

func (s AdjustService) GetFactors(ctx context.Context, data dto.AdjustRequest) (FactorSlice, error) {

	domains := make([]interface{}, 0, len(data.Domain))
	for _, value := range data.Domain {
		domains = append(domains, value)
	}

	factors, err := models.Factors(
		qm.WhereIn("domain IN ?", domains...)).All(ctx, bcdb.DB())

	if err != nil {
		return nil, fmt.Errorf("failed to fetch all factors: %w", err)
	}
	res := make(FactorSlice, 0)
	res.FromModel(factors)

	return res, nil
}

func (s AdjustService) UpdateFactors(ctx context.Context, factors FactorSlice) {
	var requests []bulk.FactorUpdateRequest

	bulk.BulkInsertFactors(ctx)

}

func NewAdjustService(historyModule history.HistoryModule) *AdjustService {
	return &AdjustService{
		historyModule: historyModule,
	}
}
