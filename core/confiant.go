package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
	"github.com/volatiletech/sqlboiler/v4/types"
)

var getConfiantQuery = `SELECT * FROM confiant 
        WHERE (publisher_id, domain) IN (%s)`

type ConfiantUpdateRequest struct {
	Publisher string  `json:"publisher_id" validate:"required"`
	Domain    string  `json:"domain"`
	Hash      string  `json:"confiant_key"`
	Rate      float64 `json:"rate"`
}

type ConfiantUpdateRespose struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Confiant struct {
	ConfiantKey *string    `boil:"confiant_key" json:"confiant_key,omitempty" toml:"confiant_key" yaml:"confiant_key"`
	PublisherID string     `boil:"publisher_id" json:"publisher_id,omitempty" toml:"publisher_id" yaml:"publisher_id"`
	Domain      *string    `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Rate        *float64   `boil:"rate" json:"rate,omitempty" toml:"rate" yaml:"rate"`
	CreatedAt   *time.Time `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

type ConfiantSlice []*Confiant

type GetConfiantOptions struct {
	Filter     ConfiantFilter         `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type ConfiantFilter struct {
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	ConfiantID  filter.StringArrayFilter `json:"confiant_key,omitempty"`
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	Rate        filter.StringArrayFilter `json:"rate,omitempty"`
}

func (confiant *Confiant) FromModel(mod *models.Confiant) error {
	confiant.PublisherID = mod.PublisherID
	confiant.CreatedAt = &mod.CreatedAt
	confiant.UpdatedAt = mod.UpdatedAt.Ptr()
	confiant.Domain = &mod.Domain
	confiant.Rate = &mod.Rate
	confiant.ConfiantKey = &mod.ConfiantKey

	return nil
}

func (confiant *Confiant) FromModelToCOnfiantWIthoutDomains(slice models.ConfiantSlice) error {
	for _, mod := range slice {
		if len(mod.Domain) == 0 {
			confiant.PublisherID = mod.PublisherID
			confiant.CreatedAt = &mod.CreatedAt
			confiant.UpdatedAt = mod.UpdatedAt.Ptr()
			confiant.Domain = &mod.Domain
			confiant.Rate = &mod.Rate
			confiant.ConfiantKey = &mod.ConfiantKey
			break
		}

	}
	return nil
}

func (cs *ConfiantSlice) FromModel(slice models.ConfiantSlice) error {

	for _, mod := range slice {
		c := Confiant{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func GetConfiants(ctx context.Context, ops *GetConfiantOptions) (ConfiantSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.ConfiantColumns.PublisherID).AddArray(ops.Pagination.Do())

	if ops.Selector == "id" {
		qmods = qmods.Add(qm.Select("DISTINCT " + models.ConfiantColumns.PublisherID))
	} else {
		qmods = qmods.Add(qm.Select("DISTINCT *"))
		qmods = qmods.Add(qm.Load(models.ConfiantRels.Publisher))

	}
	mods, err := models.Confiants(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve publishers")
	}

	res := make(ConfiantSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (filter *ConfiantFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.ConfiantColumns.PublisherID))
	}

	if len(filter.ConfiantID) > 0 {
		mods = append(mods, filter.ConfiantID.AndIn(models.ConfiantColumns.ConfiantKey))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.ConfiantColumns.Domain))
	}

	if len(filter.Rate) > 0 {
		mods = append(mods, filter.Rate.AndIn(models.ConfiantColumns.Rate))
	}

	return mods
}

func LoadConfiantByPublisherAndDomain(ctx context.Context, pubDom models.PublisherDomainSlice) (map[string]models.Confiant, error) {
	confiantMap := make(map[string]models.Confiant)

	var confiants []models.Confiant

	query := createGetConfiantsQuery(pubDom)
	err := queries.Raw(query).Bind(ctx, bcdb.DB(), &confiants)
	if err != nil {
		return nil, err
	}

	for _, confiant := range confiants {
		confiantMap[confiant.PublisherID+":"+confiant.Domain] = confiant
	}
	return confiantMap, err
}

func UpdateMetaDataQueue(c *fiber.Ctx, data *ConfiantUpdateRequest) error {

	key := buildKey(data)
	value, err := buildValue(c, data)

	mod := models.MetadataQueue{
		Key:           key,
		TransactionID: bcguid.NewFromf(data.Publisher, data.Domain, time.Now()),
		Value:         value,
	}

	err = mod.Insert(c.Context(), bcdb.DB(), boil.Infer())

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to insert metadata update to queue", err)
	}
	return nil
}

func buildKey(data *ConfiantUpdateRequest) string {
	key := "confiant:v2:" + data.Publisher
	if data.Domain != "" {
		key = key + ":" + data.Domain
	}
	return key
}

func buildValue(c *fiber.Ctx, data *ConfiantUpdateRequest) (types.JSON, error) {

	keyRate := keyRate{
		Key:  data.Hash,
		Rate: data.Rate,
	}

	val, err := json.Marshal(keyRate)

	if err != nil {
		return nil, utils.ErrorResponse(c, fiber.StatusBadRequest, "Confiant failed to parse hash value", err)
	}
	return val, err
}

func UpdateConfiant(c *fiber.Ctx, data *ConfiantUpdateRequest) error {
	modConf := models.Confiant{
		PublisherID: data.Publisher,
		ConfiantKey: data.Hash,
		Rate:        data.Rate,
		Domain:      data.Domain,
	}

	return modConf.Upsert(
		c.Context(),
		bcdb.DB(),
		true,
		[]string{models.ConfiantColumns.PublisherID, models.ConfiantColumns.Domain},
		boil.Blacklist(models.ConfiantColumns.CreatedAt),
		boil.Infer(),
	)
}

func createGetConfiantsQuery(pubDom models.PublisherDomainSlice) string {
	tupleCondition := ""
	for i, mod := range pubDom {
		if i > 0 {
			tupleCondition += ","
		}
		tupleCondition += fmt.Sprintf("('%s','%s')", mod.PublisherID, mod.Domain)
	}
	return fmt.Sprintf(getConfiantQuery, tupleCondition)
}

type keyRate struct {
	Key  string  `json:"key"`
	Rate float64 `json:"rate"`
}
