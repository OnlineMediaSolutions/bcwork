package core

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type PixalateService struct {
	historyModule history.HistoryModule
}

func NewPixalateService(historyModule history.HistoryModule) *PixalateService {
	return &PixalateService{
		historyModule: historyModule,
	}
}

var getPixalateQuery = `SELECT * FROM pixalate 
        WHERE (publisher_id, domain) IN (%s)`

type PixalateUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Rate      float64 `json:"rate"`
	Active    bool    `json:"active"`
}

type PixalateUpdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var deletePixalateQuery = `UPDATE pixalate
SET active = false
WHERE pixalate_key in (%s)`

func (p *PixalateService) UpdatePixalateTable(ctx context.Context, data *PixalateUpdateRequest) error {
	var oldModPointer any
	mod, err := models.Pixalates(
		models.PixalateWhere.PublisherID.EQ(data.Publisher),
		models.PixalateWhere.Domain.EQ(data.Domain),
	).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if mod == nil {
		mod = &models.Pixalate{
			PublisherID: data.Publisher,
			ID:          bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
			Rate:        data.Rate,
			Domain:      data.Domain,
			Active:      data.Active,
			CreatedAt:   time.Now().UTC(),
		}

		err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	} else {
		oldMod := *mod
		oldModPointer = &oldMod

		mod.Rate = data.Rate
		mod.Active = data.Active
		mod.UpdatedAt = null.TimeFrom(time.Now().UTC())

		_, err := mod.Update(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return err
		}
	}

	p.historyModule.SaveOldAndNewValuesToCache(ctx, oldModPointer, mod)

	return nil
}

func (p *PixalateService) SoftDeletePixalateInMetaData(ctx context.Context, keys []string) error {
	metas, err := models.Pixalates(models.PixalateWhere.ID.IN(keys)).All(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to fetch metadata_queue for Pixalate: %w", err)
	}

	for _, meta := range metas {
		mod := models.MetadataQueue{
			Key:           "pixalate:" + meta.PublisherID,
			TransactionID: bcguid.NewFromf(meta.PublisherID, meta.Domain, time.Now()),
			Value:         []byte(strconv.FormatFloat(0, 'f', 2, 64)),
		}

		if meta.Domain != "" {
			mod.Key = mod.Key + ":" + meta.Domain
		}

		err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
		if err != nil {
			return fmt.Errorf("failed to update metadata_queue with Pixalate: %w", err)
		}
	}

	return nil
}

func (p *PixalateService) UpdateMetaDataQueueWithPixalate(ctx context.Context, data *PixalateUpdateRequest) error {
	mod := models.MetadataQueue{
		Key:           "pixalate:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(strconv.FormatFloat(data.Rate, 'f', 2, 64)),
	}

	if data.Active == false {
		mod.Value = []byte("0")
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err := mod.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to update metadata_queue with Pixalate: %w", err)
	}

	return nil
}

type GetPixalateOptions struct {
	Filter     PixalateFilter         `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PixalateFilter struct {
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	PixalateKey filter.StringArrayFilter `json:"pixalate_key,omitempty"`
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	Rate        filter.StringArrayFilter `json:"rate,omitempty"`
	Active      filter.StringArrayFilter `json:"active,omitempty"`
}

type Pixalate struct {
	PixalateKey string     `boil:"pixalate_key" json:"pixalate_key,omitempty" toml:"pixalate_key" yaml:"pixalate_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain      *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        *float64   `boil:"rate" json:"rate,omitempty" toml:"rate" yaml:"rate"`
	Active      *bool      `boil:"active" json:"active,omitempty" toml:"active" yaml:"active"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

type PixalateSlice []*Pixalate

func (p *PixalateService) GetPixalate(ctx context.Context, ops *GetPixalateOptions) (PixalateSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.PixalateColumns.PublisherID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Pixalates(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve Pixalates")
	}

	res := make(PixalateSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (cs *PixalateSlice) FromModel(slice models.PixalateSlice) error {
	for _, mod := range slice {
		c := Pixalate{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (pixalate *Pixalate) FromModel(mod *models.Pixalate) error {
	pixalate.PublisherID = mod.PublisherID
	pixalate.CreatedAt = &mod.CreatedAt
	pixalate.UpdatedAt = mod.UpdatedAt.Ptr()
	pixalate.Domain = &mod.Domain
	pixalate.Rate = &mod.Rate
	pixalate.PixalateKey = mod.ID
	pixalate.Active = &mod.Active

	return nil
}

func (pixalate *Pixalate) FromModelToPixalateWIthoutDomains(slice models.PixalateSlice) error {
	for _, mod := range slice {
		if len(mod.Domain) == 0 {
			pixalate.PublisherID = mod.PublisherID
			pixalate.CreatedAt = &mod.CreatedAt
			pixalate.UpdatedAt = mod.UpdatedAt.Ptr()
			pixalate.Domain = &mod.Domain
			pixalate.Rate = &mod.Rate
			pixalate.PixalateKey = mod.ID
			pixalate.Active = &mod.Active
			break
		}
	}

	return nil
}

func (newPixalate *Pixalate) createPixalate(pixalate models.Pixalate) {
	newPixalate.PublisherID = pixalate.PublisherID
	newPixalate.CreatedAt = &pixalate.CreatedAt
	newPixalate.UpdatedAt = pixalate.UpdatedAt.Ptr()
	newPixalate.Domain = &pixalate.Domain
	newPixalate.Rate = &pixalate.Rate
	newPixalate.PixalateKey = pixalate.ID
	newPixalate.Active = &pixalate.Active
}

func (filter *PixalateFilter) QueryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PixalateColumns.PublisherID))
	}

	if len(filter.PixalateKey) > 0 {
		mods = append(mods, filter.PixalateKey.AndIn(models.PixalateColumns.ID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PixalateColumns.Domain))
	}

	if len(filter.Rate) > 0 {
		mods = append(mods, filter.Rate.AndIn(models.PixalateColumns.Rate))
	}

	if len(filter.Active) > 0 {
		mods = append(mods, filter.Active.AndIn(models.PixalateColumns.Active))
	}

	return mods
}

func (p *PixalateService) SoftDeletePixalates(ctx context.Context, keys []string) error {
	var wrappedStrings []string
	for _, pixalateId := range keys {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, pixalateId))
	}

	mods, err := models.Pixalates(models.PixalateWhere.ID.IN(keys)).All(ctx, bcdb.DB())
	if err != nil {
		return eris.Wrap(err, "Failed to get pixalates by keys")
	}

	softDelete := fmt.Sprintf(deletePixalateQuery, strings.Join(wrappedStrings, ","))

	_, err = queries.Raw(softDelete).Exec(bcdb.DB())
	if err != nil {
		return eris.Wrap(err, "Failed to remove pixalates by keys")
	}

	oldMods := make([]any, 0, len(mods))
	newMods := make([]any, 0, len(mods))

	for i := range mods {
		oldMods = append(oldMods, mods[i])
		newMods = append(newMods, &models.Pixalate{
			ID:          mods[i].ID,
			PublisherID: mods[i].PublisherID,
			Domain:      mods[i].Domain,
			Rate:        mods[i].Rate,
			Active:      false,
			CreatedAt:   mods[i].CreatedAt,
			UpdatedAt:   mods[i].UpdatedAt,
		})
	}

	p.historyModule.SaveOldAndNewValuesToCache(ctx, oldMods, newMods)

	return nil
}

func LoadPixalateByPublisherAndDomain(ctx context.Context, pubDom models.PublisherDomainSlice) (map[string]models.Pixalate, error) {
	pixalateMap := make(map[string]models.Pixalate)

	var pixalate []models.Pixalate
	query := createGetPixalatesQuery(pubDom)
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &pixalate)
	if err != nil {
		return nil, err
	}

	for _, pixalate := range pixalate {
		pixalateMap[pixalate.PublisherID+":"+pixalate.Domain] = pixalate
	}

	return pixalateMap, err
}

func createGetPixalatesQuery(pubDom models.PublisherDomainSlice) string {
	tupleCondition := ""
	for i, mod := range pubDom {
		if i > 0 {
			tupleCondition += ","
		}
		tupleCondition += fmt.Sprintf("('%s','%s')", mod.PublisherID, mod.Domain)
	}

	return fmt.Sprintf(getPixalateQuery, tupleCondition)
}
