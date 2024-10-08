package core

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
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
		log.Error().Err(err).Str("body", string(c.Body())).Msg("failed to insert metadata update to queue")
		return fmt.Errorf("failed to insert metadata update to queue", err)
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
		qmods = qmods.Add(qm.Load(models.PublisherDemandRels.DemandPartner))
	}

	//
	//qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.DpoColumns.DemandPartnerID).AddArray(ops.Pagination.Do())
	//if ops.Selector == "id" {
	//	qmods = qmods.Add(qm.Select("DISTINCT " + models.DpoColumns.DemandPartnerID))
	//} else {
	//	qmods = qmods.Add(qm.Select("DISTINCT *"))
	//	qmods = qmods.Add(qm.Load(models.DpoRels.DemandPartnerPublisherDemands))
	//}

	mods, err := models.Dpos(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve PublisherDemandResponse")
	}
	res := make(PublisherDemandResponseSlice, 0)
	res.FromModel(mods)

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
	Active       filter.StringArrayFilter `json:"active,omitempty"`
	AdsTxtStatus filter.StringArrayFilter `json:"ads_txt_status,omitempty"`
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

		mods = append(mods, filter.Publisher.(models.DPO.DemandPartnerPublisherDemandsPublisherID))
		mods = append(mods, filter.Publisher.AndIn(models.DPO.DemandPartnerPublisherDemandsPublisherID))
	}

	if len(filter.Demand) > 0 {
		mods = append(mods, filter.Demand.AndIn(models.PublisherDemandColumns.DemandPartnerID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDemandColumns.Domain))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.PublisherDemandColumns.Active))
	}

	return mods
}

func (cs *PublisherDemandResponseSlice) FromModel(slice models.DpoSlice) error {

	for _, mod := range slice {

		for _, pubDemand := range mod.R.GetDemandPartnerPublisherDemands() {
			c := PublisherDemandResponse{}
			err := c.FromModel(mod, pubDemand)
			if err != nil {
				return eris.Cause(err)
			}
			*cs = append(*cs, &c)
		}

	}
	return nil
}
func (publisherDemandResponse *PublisherDemandResponse) FromModel(mod *models.Dpo, pubDemand *models.PublisherDemand) error {
	publisherDemandResponse.PublisherID = pubDemand.PublisherID
	publisherDemandResponse.CreatedAt = &pubDemand.CreatedAt
	publisherDemandResponse.UpdatedAt = pubDemand.UpdatedAt.Ptr()
	publisherDemandResponse.Domain = &pubDemand.Domain
	publisherDemandResponse.DemandPartnerName = &mod.DemandPartnerName.String
	publisherDemandResponse.DemandPartnerID = &mod.DemandPartnerID
	publisherDemandResponse.AdsTxtStatus = &pubDemand.AdsTXTStatus
	publisherDemandResponse.Active = &mod.Active

	return nil
}
