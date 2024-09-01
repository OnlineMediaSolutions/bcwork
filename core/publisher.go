package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"strings"
	"time"
)

// Publisher is an object representing the database table.
type Publisher struct {
	PublisherID         string    `json:"publisher_id"`
	CreatedAt           time.Time `json:"created_at"`
	Name                string    `json:"name"`
	AccountManagerID    string    `json:"account_manager_id,omitempty"`
	MediaBuyerID        string    `json:"media_buyer_id,omitempty"`
	CampaignManagerID   string    `json:"campaign_manager_id,omitempty"`
	OfficeLocation      string    `json:"office_location,omitempty"`
	PauseTimestamp      int64     `json:"pause_timestamp,omitempty"`
	StartTimestamp      int64     `json:"start_timestamp,omitempty"`
	ReactivateTimestamp int64     `json:"reactivate_timestamp,omitempty"`
	Domains             []string  `json:"domains,omitempty"`
	IntegrationType     []string  `json:"integrationType"`
	Status              string    `json:"status"`
	Confiant            Confiant  `json:"confiant,omitempty"`
	Pixalate            Pixalate  `json:"pixalate,omitempty"`
	LatestTimestamp     int64     `json:"latest_timestamp,omitempty"`
}

type PublisherSlice []*Publisher

func (pub *Publisher) FromModel(mod *models.Publisher) error {

	pub.PublisherID = mod.PublisherID
	pub.CreatedAt = mod.CreatedAt
	pub.Name = mod.Name
	pub.Status = mod.Status.String
	pub.AccountManagerID = mod.AccountManagerID.String
	pub.MediaBuyerID = mod.MediaBuyerID.String
	pub.CampaignManagerID = mod.CampaignManagerID.String
	pub.OfficeLocation = mod.OfficeLocation.String
	pub.PauseTimestamp = mod.PauseTimestamp.Int64
	pub.StartTimestamp = mod.StartTimestamp.Int64
	pub.ReactivateTimestamp = mod.ReactivateTimestamp.Int64
	pub.LatestTimestamp = max(pub.StartTimestamp, pub.ReactivateTimestamp)

	if len(mod.IntegrationType) == 0 {
		pub.IntegrationType = []string{}
	} else {
		pub.IntegrationType = mod.IntegrationType
	}

	if mod.R != nil {
		if len(mod.R.PublisherDomains) > 0 {
			for _, dom := range mod.R.PublisherDomains {
				pub.Domains = append(pub.Domains, dom.Domain)
			}
		}
		if len(mod.R.Confiants) > 0 {
			pub.Confiant = Confiant{}
			err := pub.Confiant.FromModelToCOnfiantWIthoutDomains(mod.R.Confiants)
			if err != nil {
				return eris.Wrap(err, "failed to add Confiant data for publisher")
			}
		}
		if len(mod.R.Pixalates) > 0 {
			pub.Pixalate = Pixalate{}
			err := pub.Pixalate.FromModelToPixalateWIthoutDomains(mod.R.Pixalates)
			if err != nil {
				return eris.Wrap(err, "failed to add Pixalate data for publisher")
			}
		}
	}

	return nil
}

func (cs *PublisherSlice) FromModel(slice models.PublisherSlice) error {

	for _, mod := range slice {
		c := Publisher{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

type PublisherFilter struct {
	PublisherID       filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Name              filter.StringArrayFilter `json:"name,omitempty"`
	OfficeLocation    filter.StringArrayFilter `json:"office_location,omitempty"`
	AccountManagerID  filter.StringArrayFilter `json:"account_manager_id,omitempty"`
	MediaBuyerID      filter.StringArrayFilter `json:"media_buyer_id,omitempty"`
	CampaignManagerID filter.StringArrayFilter `json:"campaign_manager_id,omitempty"`
	Search            string                   `json:"search,omitempty"`
	CreatedAt         *filter.DatesFilter      `json:"created_at,omitempty"`
}

func (filter *PublisherFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PublisherColumns.PublisherID))
	}

	if len(filter.Name) > 0 {
		mods = append(mods, filter.Name.AndIn(models.PublisherColumns.Name))
	}

	if len(filter.OfficeLocation) > 0 {
		mods = append(mods, filter.OfficeLocation.AndIn(models.PublisherColumns.OfficeLocation))
	}

	if len(filter.AccountManagerID) > 0 {
		mods = append(mods, filter.AccountManagerID.AndIn(models.PublisherColumns.AccountManagerID))
	}

	if len(filter.MediaBuyerID) > 0 {
		mods = append(mods, filter.MediaBuyerID.AndIn(models.PublisherColumns.MediaBuyerID))
	}

	if len(filter.CampaignManagerID) > 0 {
		mods = append(mods, filter.CampaignManagerID.AndIn(models.PublisherColumns.CampaignManagerID))
	}

	if filter.CreatedAt != nil {
		mods = append(mods, filter.CreatedAt.AndIn(models.PublisherColumns.CreatedAt))
	}

	if filter.Search != "" {
		mods = append(mods,
			qm.And("(LOWER(CAST ("+models.PublisherColumns.Name+" AS TEXT)) LIKE '%"+strings.ToLower(string(filter.Search))+"%'"),
			qm.Or("LOWER(CAST ("+models.PublisherColumns.PublisherID+" AS TEXT)) LIKE '%"+strings.ToLower(string(filter.Search))+"%'"),
			qm.Or("LOWER(CAST ("+models.PublisherColumns.OfficeLocation+" AS TEXT)) LIKE '%"+strings.ToLower(string(filter.Search))+"%')"),
		)
	}

	return mods
}

type GetPublisherOptions struct {
	Filter     PublisherFilter        `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

func GetPublisher(ctx context.Context, ops *GetPublisherOptions) (PublisherSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.PublisherColumns.PublisherID).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.PublisherColumns.PublisherID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		qmods = qmods.Add(qm.Load(models.PublisherRels.PublisherDomains))
		qmods = qmods.Add(qm.Load(models.PublisherRels.Confiants))
		qmods = qmods.Add(qm.Load(models.PublisherRels.Pixalates))

	}
	mods, err := models.Publishers(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publishers")
	}

	res := make(PublisherSlice, 0)
	res.FromModel(mods)

	return res, nil
}

type UpdatePublisherValues struct {
	Name                *string   `json:"name"`
	AccountManagerID    *string   `json:"account_manager_id,omitempty"`
	MediaBuyerID        *string   `json:"media_buyer_id,omitempty"`
	CampaignManagerID   *string   `json:"campaign_manager_id,omitempty"`
	OfficeLocation      *string   `json:"office_location,omitempty"`
	PauseTimestamp      *int64    `json:"pause_timestamp,omitempty"`
	StartTimestamp      *int64    `json:"start_timestamp,omitempty"`
	ReactivateTimestamp *int64    `json:"reactivate_timestamp,omitempty"`
	Status              *string   `json:"status,omitempty"`
	IntegrationType     *[]string `json:"integration_type,omitempty"`
}

func UpdatePublisher(ctx context.Context, publisherID string, vals UpdatePublisherValues) error {

	if publisherID == "" {
		return fmt.Errorf("publisher_id is mandatory when updating a publishe")
	}

	modPublisher := &models.Publisher{
		PublisherID: publisherID,
	}

	//whitelist
	cols := []string{}

	if vals.Name != nil {
		modPublisher.Name = *vals.Name
		cols = append(cols, models.PublisherColumns.Name)
	}

	if vals.AccountManagerID != nil {
		modPublisher.AccountManagerID = null.StringFromPtr(vals.AccountManagerID)
		cols = append(cols, models.PublisherColumns.AccountManagerID)
	}

	if vals.MediaBuyerID != nil {
		modPublisher.MediaBuyerID = null.StringFromPtr(vals.MediaBuyerID)
		cols = append(cols, models.PublisherColumns.MediaBuyerID)
	}

	if vals.CampaignManagerID != nil {
		modPublisher.CampaignManagerID = null.StringFromPtr(vals.CampaignManagerID)
		cols = append(cols, models.PublisherColumns.CampaignManagerID)
	}

	if vals.OfficeLocation != nil {
		modPublisher.OfficeLocation = null.StringFromPtr(vals.OfficeLocation)
		cols = append(cols, models.PublisherColumns.OfficeLocation)
	}

	if vals.PauseTimestamp != nil {
		modPublisher.PauseTimestamp = null.Int64FromPtr(vals.PauseTimestamp)
		cols = append(cols, models.PublisherColumns.PauseTimestamp)
	}

	if vals.StartTimestamp != nil {
		modPublisher.StartTimestamp = null.Int64FromPtr(vals.StartTimestamp)
		cols = append(cols, models.PublisherColumns.StartTimestamp)
	}

	if vals.ReactivateTimestamp != nil {
		modPublisher.ReactivateTimestamp = null.Int64FromPtr(vals.ReactivateTimestamp)
		cols = append(cols, models.PublisherColumns.ReactivateTimestamp)
	}
	if vals.Status != nil {
		modPublisher.Status = null.StringFromPtr(vals.Status)
		cols = append(cols, models.PublisherColumns.Status)
	}
	if vals.IntegrationType != nil {
		modPublisher.IntegrationType = types.StringArray(*vals.IntegrationType)
		cols = append(cols, models.PublisherColumns.IntegrationType)
	}
	if len(cols) == 0 {
		return fmt.Errorf("applicaiton payload contains no vals for update (publisher_id:%s)", modPublisher.PublisherID)
	}

	count, err := modPublisher.Update(ctx, bcdb.DB(), boil.Whitelist(cols...))
	if err != nil {
		return eris.Wrap(err, fmt.Sprintf("failed to update publisher (publisher_id:%s)", modPublisher.PublisherID))
	}
	if count == 0 {
		return eris.Wrap(err, fmt.Sprintf("wrong publishe_id when updating publisher,verify publisher_id really exists (unit_id:%s)", modPublisher.PublisherID))
	}

	return nil

}

type PublisherCreateValues struct {
	Name              string   `json:"name"`
	AccountManagerID  string   `json:"account_manager_id"`
	MediaBuyerID      string   `json:"media_buyer_id"`
	CampaignManagerID string   `json:"campaign_manager_id"`
	OfficeLocation    string   `json:"office_location"`
	Status            string   `json:"status"`
	IntegrationType   []string `json:"integration_type"`
}

func CreatePublisher(ctx context.Context, vals PublisherCreateValues) (string, error) {

	maxAge, err := calculatePublisherKey()

	modPublisher := models.Publisher{
		PublisherID:       maxAge,
		Name:              vals.Name,
		AccountManagerID:  null.StringFrom(vals.AccountManagerID),
		MediaBuyerID:      null.StringFrom(vals.MediaBuyerID),
		CampaignManagerID: null.StringFrom(vals.CampaignManagerID),
		OfficeLocation:    null.StringFrom(vals.OfficeLocation),
		Status:            null.StringFrom(vals.Status),
		IntegrationType:   vals.IntegrationType,
	}

	err = modPublisher.Insert(ctx, bcdb.DB(), boil.Infer())

	if err != nil {
		return "", eris.Wrapf(err, "failed to insert publisher")
	}

	return modPublisher.PublisherID, nil

}

func calculatePublisherKey() (string, error) {
	var maxAge int
	err := models.Publishers(qm.Select("MAX(publisher_id)")).QueryRow(bcdb.DB()).Scan(&maxAge)
	if err != nil {
		eris.Wrapf(err, "failed to calculate max publisher id")
	}

	return fmt.Sprintf("%d", maxAge+1), err
}

func PublisherCount(ctx context.Context, filter *PublisherFilter) (int64, error) {

	c, err := models.Publishers(filter.QueryMod()...).Count(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return 0, eris.Wrapf(err, "failed to fetch all publishers")
	}

	return c, nil
}
