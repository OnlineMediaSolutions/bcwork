package core

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"
)

type SeatOwnerGetOptions struct {
	Filter     SeatOwnerGetFilter     `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type SeatOwnerGetFilter struct {
	ID                       filter.StringArrayFilter `json:"id,omitempty"`
	SeatOwnerName            filter.StringArrayFilter `json:"seat_owner_name,omitempty"`
	SeatOwnerDomain          filter.StringArrayFilter `json:"seat_owner_domain,omitempty"`
	PublisherAccount         filter.StringArrayFilter `json:"publisher_account,omitempty"`
	CertificationAuthorityID filter.StringArrayFilter `json:"certification_authority_id,omitempty"`
}

func (filter *SeatOwnerGetFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.ID) > 0 {
		mods = append(mods, filter.ID.AndIn(models.SeatOwnerColumns.ID))
	}

	if len(filter.SeatOwnerName) > 0 {
		mods = append(mods, filter.SeatOwnerName.AndIn(models.SeatOwnerColumns.SeatOwnerName))
	}

	if len(filter.SeatOwnerDomain) > 0 {
		mods = append(mods, filter.SeatOwnerDomain.AndIn(models.SeatOwnerColumns.SeatOwnerDomain))
	}

	if len(filter.PublisherAccount) > 0 {
		mods = append(mods, filter.PublisherAccount.AndIn(models.SeatOwnerColumns.PublisherAccount))
	}

	if len(filter.CertificationAuthorityID) > 0 {
		mods = append(mods, filter.CertificationAuthorityID.AndIn(models.SeatOwnerColumns.CertificationAuthorityID))
	}

	return mods
}

func (d *DemandPartnerService) GetSeatOwners(ctx context.Context, ops *SeatOwnerGetOptions) ([]*dto.SeatOwner, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.SeatOwnerColumns.ID).
		AddArray(ops.Pagination.Do()).
		Add(qm.Select("DISTINCT *"))

	mods, err := models.SeatOwners(qmods...).All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to retrieve seat owners: %w", err)
	}

	seatOwners := make([]*dto.SeatOwner, 0, len(mods))
	for _, mod := range mods {
		dp := &dto.SeatOwner{}
		dp.FromModel(mod)
		seatOwners = append(seatOwners, dp)
	}

	return seatOwners, nil
}
