package core

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
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
	"time"
)

type PublisherDomainUpdateRequest struct {
	PublisherID string   `json:"publisher_id" validate:"required"`
	Domain      string   `json:"domain"`
	GppTarget   *float64 `json:"gpp_target,omitempty"`
	Automation  bool     `json:"automation"`
}

type GetPublisherDomainOptions struct {
	Filter     PublisherDomainFilter  `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type PublisherDomainFilter struct {
	Domain      filter.StringArrayFilter `json:"domain,omitempty"`
	PublisherID filter.StringArrayFilter `json:"publisher_id,omitempty"`
	Automation  filter.StringArrayFilter `json:"automation,omitempty"`
	GppTarget   filter.StringArrayFilter `json:"gpp_target,omitempty"`
}

func GetPublisherDomain(ctx context.Context, ops *GetPublisherDomainOptions) (PublisherDomainSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.PublisherDomainColumns.PublisherID).AddArray(ops.Pagination.Do())
	qmods = qmods.Add(qm.Select("DISTINCT *"))

	mods, err := models.PublisherDomains(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "Failed to retrieve publisher domains values")
	}

	confiantMap, err := LoadConfiantByPublisherAndDomain(ctx, mods)
	pixalateMap, err := LoadPixalateByPublisherAndDomain(ctx, mods)

	if err != nil {
		return nil, eris.Wrap(err, "Error while retreving confiants for publisher domains values")
	}
	res := make(PublisherDomainSlice, 0)
	res.FromModel(mods, confiantMap, pixalateMap)

	return res, nil
}

type PublisherDomainSlice []*PublisherDomain

func (cs *PublisherDomainSlice) FromModel(slice models.PublisherDomainSlice, confiantMap map[string]models.Confiant, pixalateMap map[string]models.Pixalate) error {

	for _, mod := range slice {
		c := PublisherDomain{}
		key := mod.PublisherID + ":" + mod.Domain
		confiant := confiantMap[key]
		pixalate := pixalateMap[key]
		err := c.FromModel(mod, confiant, pixalate)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

type PublisherDomain struct {
	PublisherID     string     `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain          string     `boil:"domain" json:"domain,omitempty" toml:"domain" yaml:"domain,omitempty"`
	Automation      bool       `boil:"automation" json:"automation" toml:"automation" yaml:"automation"`
	GppTarget       float64    `boil:"gpp_target" json:"gpp_target" toml:"gpp_target" yaml:"gpp_target"`
	IntegrationType []string   `boil:"integration_type" json:"integration_type" toml:"integration_type" yaml:"integration_type"`
	CreatedAt       time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Confiant        Confiant   `boil:"confiant" json:"confiant,omitempty" toml:"confiant" yaml:"confiant"`
	Pixalate        Pixalate   `boil:"pixalate" json:"pixalate,omitempty" toml:"pixalate" yaml:"pixalate"`
	UpdatedAt       *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func (filter *PublisherDomainFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.PublisherID) > 0 {
		mods = append(mods, filter.PublisherID.AndIn(models.PublisherDomainColumns.PublisherID))
	}

	if len(filter.Domain) > 0 {
		mods = append(mods, filter.Domain.AndIn(models.PublisherDomainColumns.Domain))
	}

	if len(filter.GppTarget) > 0 {
		mods = append(mods, filter.GppTarget.AndIn(models.PublisherDomainColumns.GPPTarget))
	}

	if len(filter.Automation) > 0 {
		mods = append(mods, filter.Automation.AndIn(models.PublisherDomainColumns.Automation))
	}

	return mods
}

func (pubDom *PublisherDomain) FromModel(mod *models.PublisherDomain, confiant models.Confiant, pixalate models.Pixalate) error {

	pubDom.PublisherID = mod.PublisherID
	pubDom.CreatedAt = mod.CreatedAt
	pubDom.UpdatedAt = mod.UpdatedAt.Ptr()
	pubDom.Domain = mod.Domain
	pubDom.GppTarget = mod.GPPTarget.Float64
	pubDom.Automation = mod.Automation
	pubDom.IntegrationType = mod.IntegrationType
	pubDom.Confiant = Confiant{}
	pubDom.Pixalate = Pixalate{}
	if len(confiant.ConfiantKey) > 0 {
		pubDom.Confiant.createConfiant(confiant)
	}
	if len(pixalate.ID) > 0 {
		pubDom.Pixalate.createPixalate(pixalate)
	}
	return nil
}

func (newConfiant *Confiant) createConfiant(confiant models.Confiant) {
	newConfiant.PublisherID = confiant.PublisherID
	newConfiant.CreatedAt = &confiant.CreatedAt
	newConfiant.UpdatedAt = confiant.UpdatedAt.Ptr()
	newConfiant.Domain = confiant.Domain
	newConfiant.Rate = confiant.Rate
	newConfiant.ConfiantKey = confiant.ConfiantKey
}

func UpdatePublisherDomain(c *fiber.Ctx, data *PublisherDomainUpdateRequest) error {

	var gppTarget sql.NullFloat64
	if data.GppTarget == nil {
		gppTarget = sql.NullFloat64{Float64: 0, Valid: false}
	} else {
		gppTarget = sql.NullFloat64{Float64: *data.GppTarget, Valid: true}
	}
	modConf := models.PublisherDomain{
		Domain:      data.Domain,
		PublisherID: data.PublisherID,
		Automation:  data.Automation,
		GPPTarget:   null.Float64(gppTarget),
	}

	return modConf.Upsert(c.Context(), bcdb.DB(), true, []string{models.PublisherDomainColumns.PublisherID, models.PublisherDomainColumns.Domain}, boil.Infer(), boil.Infer())
}
