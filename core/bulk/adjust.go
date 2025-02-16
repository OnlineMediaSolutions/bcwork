package bulk

import (
	"context"
	"fmt"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type AdjustService struct {
	historyModule history.HistoryModule
}

func (b *BulkService) AdjustFactors(ctx context.Context, data dto.AdjustRequest) error {
	domains := make([]interface{}, 0, len(data.Domain))
	for _, value := range data.Domain {
		domains = append(domains, value)
	}

	factors, err := models.Factors(
		qm.WhereIn("domain IN ?", domains...)).All(ctx, bcdb.DB())

	if err != nil || len(factors) == 0 {
		return fmt.Errorf("failed to fetch factors: %w", err)
	}
	res := make(dto.FactorSlice, 0)
	err = res.FromModel(factors)
	if err != nil {
		return fmt.Errorf("failed to map factors: %w", err)
	}

	var requests []FactorUpdateRequest

	for _, item := range res {
		factor := FactorUpdateRequest{
			Publisher: item.Publisher,
			Domain:    item.Domain,
			Device:    item.Device,
			Country:   item.Country,
			Factor:    data.Value,
			RuleID:    item.RuleId,
		}
		requests = append(requests, factor)
	}

	err = b.BulkInsertFactors(ctx, requests)
	if err != nil {
		return err
	}

	return nil
}

func (b *BulkService) AdjustFloors(ctx context.Context, data dto.AdjustRequest) error {
	domains := make([]interface{}, 0, len(data.Domain))
	for _, value := range data.Domain {
		domains = append(domains, value)
	}

	floors, err := models.Floors(
		qm.WhereIn("domain IN ?", domains...)).All(ctx, bcdb.DB())

	if err != nil || len(floors) == 0 {
		return fmt.Errorf("failed to fetch floors: %w", err)
	}

	res := make(dto.FloorSlice, 0)
	err = res.FromModel(floors)
	if err != nil {
		return fmt.Errorf("failed to map floors: %w", err)
	}

	var requests []dto.FloorUpdateRequest

	for _, item := range res {
		floor := dto.FloorUpdateRequest{
			Publisher: item.Publisher,
			Domain:    item.Domain,
			Device:    item.Device,
			Country:   item.Country,
			Floor:     data.Value,
			RuleId:    item.RuleId,
		}
		requests = append(requests, floor)
	}

	err = BulkInsertFloors(ctx, requests)
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
