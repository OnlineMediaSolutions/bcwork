package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/dto"
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

var deletePixalateQuery = `UPDATE pixalate
SET active = false
WHERE id in (%s)`

func (p *PixalateService) UpdatePixalateTable(ctx context.Context, data *dto.PixalateUpdateRequest) error {
	var oldModPointer any
	mod, err := models.Pixalates(
		models.PixalateWhere.PublisherID.EQ(data.Publisher),
		models.PixalateWhere.Domain.EQ(data.Domain),
	).One(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

	p.historyModule.SaveAction(ctx, oldModPointer, mod, nil)

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

func (p *PixalateService) UpdateMetaDataQueueWithPixalate(ctx context.Context, data *dto.PixalateUpdateRequest) error {
	mod := models.MetadataQueue{
		Key:           "pixalate:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(strconv.FormatFloat(data.Rate, 'f', 2, 64)),
	}

	if !data.Active {
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
	Active      *filter.BoolFilter       `json:"active,omitempty"`
}

func (p *PixalateService) GetPixalate(ctx context.Context, ops *GetPixalateOptions) (dto.PixalateSlice, error) {
	qmods := ops.Filter.QueryMod().
		Order(ops.Order, nil, models.PixalateColumns.PublisherID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Pixalates(qmods...).All(ctx, bcdb.DB())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "Failed to retrieve Pixalates")
	}

	res := make(dto.PixalateSlice, 0)
	err = res.FromModel(mods)
	if err != nil {
		return nil, fmt.Errorf("failed to map pixalate: %w", err)
	}

	return res, nil
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

	if filter.Active != nil {
		mods = append(mods, filter.Active.Where(models.PixalateColumns.Active))
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

	p.historyModule.SaveAction(ctx, oldMods, newMods, nil)

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
