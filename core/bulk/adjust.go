package bulk

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AdjustService struct {
	historyModule history.HistoryModule
}

func (s AdjustService) GetFactors(ctx context.Context, data dto.AdjustRequest) (core.FactorSlice, error) {

	domains := make([]interface{}, 0, len(data.Domain))
	for _, value := range data.Domain {
		domains = append(domains, value)
	}

	factors, err := models.Factors(
		qm.WhereIn("domain IN ?", domains...)).All(ctx, bcdb.DB())

	if err != nil {
		return nil, fmt.Errorf("failed to fetch all factors: %w", err)
	}
	res := make(core.FactorSlice, 0)
	res.FromModel(factors)

	return res, nil
}

func (s AdjustService) UpdateFactors(ctx context.Context, factors core.FactorSlice, value float64, service *BulkFactorService) error {
	var requests []FactorUpdateRequest

	for _, item := range factors {
		factor := FactorUpdateRequest{
			Publisher: item.Publisher,
			Domain:    item.Domain,
			Device:    item.Device,
			Country:   item.Country,
			Factor:    value,
			RuleID:    item.RuleId,
		}
		requests = append(requests, factor)
	}

	err := service.BulkInsertFactors(ctx, requests)
	if err != nil {
		return err
	}
	return nil
}

func NewAdjustService(historyModule history.HistoryModule) *AdjustService {
	return &AdjustService{
		historyModule: historyModule,
	}
}
