package core

import (
	"context"
	"database/sql"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
)

type HistoryService struct{}

func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

type HistoryOptions struct {
	Filter     HistoryFilter          `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type HistoryFilter struct {
	UserID          filter.IntArrayFilter    `json:"user_id,omitempty"`
	Action          filter.StringArrayFilter `json:"action,omitempty"`
	Subject         filter.StringArrayFilter `json:"subject,omitempty"`
	Item            filter.StringArrayFilter `json:"item,omitempty"`
	PublisherID     filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Domain          filter.StringArrayFilter `json:"domain,omitempty"`
	DemandPartnerID filter.StringArrayFilter `json:"demand_partner_id,omitempty"`
	EntityID        filter.StringArrayFilter `json:"entity_id,omitempty"`
}

func (h *HistoryService) GetHistory(ctx context.Context, ops *HistoryOptions) ([]*dto.History, error) {
	userQueryMod := ops.Filter.userQueryMod()
	users, err := models.Users(userQueryMod...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve users")
	}

	usersMap := make(map[int]string, len(users))
	for _, user := range users {
		usersMap[user.ID] = user.FirstName + " " + user.LastName
	}

	qmods := ops.Filter.queryMod().
		Order(ops.Order, nil, models.HistoryColumns.ID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Histories(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve history")
	}

	historyData := make([]*dto.History, 0, len(mods))
	for _, mod := range mods {
		history := new(dto.History)
		err := history.FromModel(mod, usersMap)
		if err != nil {
			return nil, eris.Wrap(err, "failed to map history to dto")
		}
		historyData = append(historyData, history)
	}

	return historyData, nil
}

func (filter *HistoryFilter) queryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.UserID) > 0 {
		mods = append(mods, filter.UserID.AndIn(models.HistoryColumns.UserID))
	}

	if len(filter.Action) > 0 {
		mods = append(mods, filter.Action.AndIn(models.HistoryColumns.Action))
	}

	if len(filter.Subject) > 0 {
		mods = append(mods, filter.Subject.AndIn(models.HistoryColumns.Subject))
	}

	if len(filter.Item) > 0 {
		mods = append(mods, filter.Item.AndIn(models.HistoryColumns.Item))
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.HistoryColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.HistoryColumns.Domain))
	}

	if len(filter.DemandPartnerID) > 0 {
		mods = append(mods, filter.DemandPartnerID.AndIn(models.HistoryColumns.DemandPartnerID))
	}

	if len(filter.EntityID) > 0 {
		mods = append(mods, filter.EntityID.AndIn(models.HistoryColumns.EntityID))
	}

	return mods
}

func (filter *HistoryFilter) userQueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0, 2)
	if filter == nil {
		return mods
	}

	if len(filter.UserID) > 0 {
		mods = append(mods, filter.UserID.AndIn(models.UserColumns.ID))
	}

	return mods
}
