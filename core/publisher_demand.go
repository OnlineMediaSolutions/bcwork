package core

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/valyala/fasthttp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type MetadataUpdateRequest struct {
	Key     string      `json:"key"`
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

func InsertDataToMetaData(c *fiber.Ctx, data MetadataUpdateRequest, value []byte, now time.Time) error {
	mod := models.MetadataQueue{
		TransactionID: bcguid.NewFromf(data.Key, now),
		Key:           data.Key,
		Value:         value,
		CreatedAt:     now,
	}

	err := mod.Insert(c.Context(), bcdb.DB(), boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to insert metadata update to queue %w", err)
	}

	return nil
}

type PublisherDemandResponseSlice []*PublisherDemandResponse

func GetPublisherDemandData(ctx *fasthttp.RequestCtx, ops *GetPublisherDemandOptions) (PublisherDemandResponseSlice, error) {
	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.PublisherDemandColumns.PublisherID).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.PublisherDemandColumns.PublisherID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		if ops.Filter.Active != nil {
			qmods = qmods.Add(qm.Load(models.PublisherDemandRels.DemandPartner, qm.Where("active = ?", *ops.Filter.Active)))
		} else {
			qmods = qmods.Add(qm.Load(models.PublisherDemandRels.DemandPartner))
		}
	}

	mods, err := models.PublisherDemands(qmods...).All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "Failed to retrieve PublisherDemandResponse")
	}

	res := make(PublisherDemandResponseSlice, 0)
	for _, mod := range mods {
		if mod.R.DemandPartner != nil {
			var publisherDemandResponse PublisherDemandResponse
			if err := publisherDemandResponse.FromModel(mod, mod.R.DemandPartner); err == nil {
				res = append(res, &publisherDemandResponse)
			}
		}
	}

	return res, nil
}

type GetPublisherDemandOptions struct {
	Filter     PublisherDemandFilter  `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PublisherDemandFilter struct {
	Publisher    filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Domain       filter.StringArrayFilter `json:"domain,omitempty"`
	Demand       filter.StringArrayFilter `json:"demand,omitempty"`
	Active       *filter.BoolFilter       `json:"active,omitempty"`
	AdsTxtStatus *filter.BoolFilter       `json:"ads_txt_status,omitempty"`
}

type PublisherDemandResponse struct {
	PublisherID       string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain            *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	DemandPartnerName *string    `boil:"demand_partner_name" json:"demand_partner_name,omitempty" toml:"demand_partner_name" yaml:"demand_partner_name"`
	DemandPartnerID   *string    `boil:"demand_partner_id" json:"demand_partner_id,omitempty" toml:"demand_partner_id" yaml:"demand_partner_id"`
	Active            *bool      `boil:"active" json:"active,omitempty" toml:"active" yaml:"active"`
	AdsTxtStatus      *bool      `boil:"ads_txt_status" json:"ads_txt_status,omitempty" toml:"ads_txt_status" yaml:"active"`
	CreatedAt         *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt         *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (filter *PublisherDemandFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Publisher) > 0 {
		mods = append(mods, filter.Publisher.AndIn(models.PublisherDemandColumns.PublisherID))
	}

	if len(filter.Demand) > 0 {
		mods = append(mods, filter.Demand.AndIn(models.PublisherDemandColumns.DemandPartnerID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDemandColumns.Domain))
	}

	if filter.AdsTxtStatus != nil {
		mods = append(mods, filter.AdsTxtStatus.Where(models.PublisherDemandColumns.AdsTXTStatus))
	}

	return mods
}

func (cs *PublisherDemandResponseSlice) FromModel(slice models.PublisherDemandSlice) error {
	for _, mod := range slice {
		c := PublisherDemandResponse{}
		demandPartner := mod.R.DemandPartner

		err := c.FromModel(mod, demandPartner)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (publisherDemandResponse *PublisherDemandResponse) FromModel(mod *models.PublisherDemand, demandPartner *models.Dpo) error {
	publisherDemandResponse.Domain = &mod.Domain
	publisherDemandResponse.PublisherID = mod.PublisherID
	publisherDemandResponse.CreatedAt = &mod.CreatedAt
	publisherDemandResponse.UpdatedAt = mod.UpdatedAt.Ptr()
	publisherDemandResponse.AdsTxtStatus = &mod.AdsTXTStatus
	publisherDemandResponse.DemandPartnerName = &demandPartner.DemandPartnerName
	publisherDemandResponse.DemandPartnerID = &demandPartner.DemandPartnerID
	publisherDemandResponse.Active = &demandPartner.Active

	return nil
}
