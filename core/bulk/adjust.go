package bulk

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils/constant"
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

func GetFloors(ctx context.Context, data dto.AdjustRequest) (core.FloorSlice, error) {

	domains := make([]interface{}, 0, len(data.Domain))
	for _, value := range data.Domain {
		domains = append(domains, value)
	}

	floors, err := models.Floors(
		qm.WhereIn("domain IN ?", domains...)).All(ctx, bcdb.DB())

	if err != nil {
		return nil, fmt.Errorf("failed to fetch all floors: %w", err)
	}
	res := make(core.FloorSlice, 0)
	res.FromModel(floors)

	return res, nil
}

func (s AdjustService) UpdateFactors(ctx context.Context, factors core.FactorSlice, value float64, service *BulkService) error {
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

func UpdateFloors(ctx context.Context, floors core.FloorSlice, value float64) error {
	var requests []constant.FloorUpdateRequest

	for _, item := range floors {
		floor := constant.FloorUpdateRequest{
			Publisher: item.Publisher,
			Domain:    item.Domain,
			Device:    item.Device,
			Country:   item.Country,
			Floor:     value,
			RuleId:    item.RuleId,
		}
		requests = append(requests, floor)
	}

	err := BulkInsertFloors(ctx, requests)
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
