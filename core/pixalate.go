package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
	"strings"
	"time"
)

type PixalateUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate"`
}

type PixalateUpdateRespose struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var deletePixalateQuery = `UPDATE pixalate
SET active = false
WHERE pixalate_key in (%s)`

func UpdatePixalateTable(c *fiber.Ctx, data *PixalateUpdateRequest) error {

	updatedPixalate := models.Pixalate{
		PublisherID: data.Publisher,
		PixalateKey: data.Hash,
		Rate:        data.Rate,
		Domain:      data.Domain,
	}

	return updatedPixalate.Upsert(c.Context(), bcdb.DB(), true, []string{models.PixalateColumns.PublisherID, models.PixalateColumns.Domain}, boil.Infer(), boil.Infer())
}

func SoftDeletePixalateInMetaData(c *fiber.Ctx, data *PixalateUpdateRequest) error {

	val, err := json.Marshal(data.Hash)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Pixalate failed to parse hash value")
	}

	mod := models.MetadataQueue{
		Key:           "pixalate:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         val,
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata_queue with Pixalate")
	}
	return nil
}

func UpdateMetaDataQueueWithPixalate(c *fiber.Ctx, data *PixalateUpdateRequest) error {

	val, err := json.Marshal(data.Hash)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Pixalate failed to parse hash value")
	}
	meta, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ(key), qm.OrderBy("created_by desc")).One(c.Context(), bcdb.DB())

	mod := models.MetadataQueue{
		Key:           "pixalate:" + data.Publisher,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         []byte(strconv.FormatFloat(data.Rate, 'f', 2, 64)),
	}

	if data.Domain != "" {
		mod.Key = mod.Key + ":" + data.Domain
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update metadata_queue with Pixalate")
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
	PixalateID  filter.StringArrayFilter `json:"pixalate_key,omitempty"`
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	Rate        filter.StringArrayFilter `json:"rate,omitempty"`
}

type Pixalate struct {
	PixalateKey string     `boil:"pixalate_key" json:"pixalate_key" toml:"pixalate_key" yaml:"pixalate_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain      string     `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        float64    `boil:"rate" json:"rate" toml:"rate" yaml:"rate"`
	Active      bool       `boil:"active" json:"active" toml:"active" yaml:"active"`
	CreatedAt   time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

type PixalateSlice []*Pixalate

func GetPixalate(ctx context.Context, ops *GetPixalateOptions) (PixalateSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.PixalateColumns.PublisherID).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.PixalateColumns.PublisherID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		qmods = qmods.Add(qm.Load(models.PixalateRels.Publisher))

	}
	mods, err := models.Pixalates(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve publishers")
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
	pixalate.CreatedAt = mod.CreatedAt
	pixalate.Domain = mod.Domain
	pixalate.Rate = mod.Rate
	pixalate.PixalateKey = mod.PixalateKey

	return nil
}

func (filter *PixalateFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PixalateColumns.PublisherID))
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PixalateID.AndIn(models.PixalateColumns.PixalateKey))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PixalateColumns.Domain))
	}

	if len(filter.Rate) > 0 {
		mods = append(mods, filter.Rate.AndIn(models.PixalateColumns.Rate))
	}

	return mods
}

func SoftDeletePixalates(ctx context.Context, keys []string) error {

	var wrappedStrings []string
	for _, pixalateId := range keys {
		wrappedStrings = append(wrappedStrings, fmt.Sprintf(`'%s'`, pixalateId))
	}

	softDelete := fmt.Sprintf(deletePixalateQuery, strings.Join(wrappedStrings, ","))

	_, err := queries.Raw(softDelete).Exec(bcdb.DB())
	if err != nil {
		return eris.Wrap(err, "Failed to remove pixalates by keys")
	}

	return nil

}
