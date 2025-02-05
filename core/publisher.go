package core

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type PublisherService struct {
	historyModule history.HistoryModule
	compassModule compass.CompassModule
}

func NewPublisherService(historyModule history.HistoryModule, compassModule compass.CompassModule) *PublisherService {
	return &PublisherService{
		historyModule: historyModule,
		compassModule: compassModule,
	}
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

func (p *PublisherService) GetPublisher(ctx context.Context, ops *GetPublisherOptions) (dto.PublisherSlice, error) {
	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.PublisherColumns.PublisherID).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.PublisherColumns.PublisherID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		qmods = qmods.Add(qm.Load(models.PublisherRels.PublisherDomains))
		qmods = qmods.Add(qm.Load(models.PublisherRels.Confiants))
		qmods = qmods.Add(qm.Load(models.PublisherRels.Pixalates))
		qmods = qmods.Add(qm.Load(models.PublisherRels.BidCachings))
		qmods = qmods.Add(qm.Load(models.PublisherRels.RefreshCaches))

	}
	mods, err := models.Publishers(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publishers")
	}

	users, err := models.Users().All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve users")
	}

	usersMap := make(map[string]string, len(users))
	for _, user := range users {
		userID := strconv.Itoa(user.ID)
		usersMap[userID] = user.FirstName + " " + user.LastName
	}

	res := make(dto.PublisherSlice, 0)
	res.FromModel(mods, usersMap)

	return res, nil
}

func (p *PublisherService) UpdatePublisher(ctx context.Context, publisherID string, vals dto.UpdatePublisherValues) error {
	if publisherID == "" {
		return fmt.Errorf("publisher_id is mandatory when updating a publisher")
	}

	modPublisher, err := models.Publishers(models.PublisherWhere.PublisherID.EQ(publisherID)).One(ctx, bcdb.DB())
	if err != nil {
		return eris.Wrap(err, fmt.Sprintf("failed to get publisher with id [%v] to update", publisherID))
	}

	oldModPublisher := *modPublisher

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
		return eris.Wrap(err, fmt.Sprintf("wrong publisher_id when updating publisher,verify publisher_id really exists (unit_id:%s)", modPublisher.PublisherID))
	}

	p.historyModule.SaveAction(ctx, &oldModPublisher, modPublisher, nil)

	return nil
}

func (p *PublisherService) CreatePublisher(ctx context.Context, vals dto.PublisherCreateValues) (string, error) {
	maxAge, err := calculatePublisherKey()

	modPublisher := &models.Publisher{
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

	p.historyModule.SaveAction(ctx, nil, modPublisher, nil)

	return modPublisher.PublisherID, nil
}

func (p *PublisherService) PublisherCount(ctx context.Context, filter *PublisherFilter) (int64, error) {
	c, err := models.Publishers(filter.QueryMod()...).Count(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return 0, eris.Wrapf(err, "failed to fetch all publishers")
	}

	return c, nil
}

func calculatePublisherKey() (string, error) {
	var maxPublisherIdValue int

	err := queries.Raw("select max(CAST(publisher_id AS NUMERIC))\nfrom publisher").QueryRow(bcdb.DB()).Scan(&maxPublisherIdValue)
	if err != nil {
		eris.Wrapf(err, "failed to calculate max publisher id")
	}

	return fmt.Sprintf("%d", maxPublisherIdValue+1), err
}
