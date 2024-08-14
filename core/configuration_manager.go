package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

var updateQuery = `INSERT INTO configuration ("key","value","description","updated_at","created_at")
VALUES ('%s', '%s', '%s', '%s', '%s')
ON CONFLICT ("key")
DO UPDATE SET value = EXCLUDED.value, description = EXCLUDED.description, updated_at = EXCLUDED.updated_at`

type ConfigurationRequest struct {
	Key         string `json:"key" validate:"required"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}

type ConfigurationPayload struct {
	Filter     ConfigurationFilter    `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type ConfigurationFilter struct {
	Key   filter.StringArrayFilter `json:"key,omitempty"`
	Value filter.StringArrayFilter `json:"value,omitempty"`
}

type ConfigurationSlice []*Configuration

type Configuration struct {
	Key         string     `boil:"key" json:"key" toml:"key" yaml:"key"`
	Value       string     `boil:"value" json:"value,omitempty" toml:"value" yaml:"value,omitempty"`
	Description string     `boil:"description" json:"description" toml:"description" yaml:"description"`
	CreatedAt   time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   *time.Time `boil:"updated_at" json:"updated_at,omitempty" toml:"updated_at" yaml:"updated_at,omitempty"`
}

func UpdateConfiguration(c *fiber.Ctx, data *ConfigurationRequest) error {

	currentTime := time.Now().Format("2006-01-02 15:04")
	query := fmt.Sprintf(updateQuery, data.Key, data.Value, data.Description, currentTime, currentTime)

	_, err := queries.Raw(query).Exec(bcdb.DB())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("%s", err.Error()))
	}
	return nil
}

func GetConfigurations(ctx context.Context, ops *ConfigurationPayload) (ConfigurationSlice, error) {

	qmods := ops.Filter.QueryMod().Order(ops.Order, nil, models.ConfigurationColumns.Key).AddArray(ops.Pagination.Do())

	qmods = qmods.Add(qm.Select("DISTINCT *"))

	mods, err := models.Configurations(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve configurations")
	}

	res := make(ConfigurationSlice, 0)
	res.FromModel(mods)

	return res, nil
}

func (cs *ConfigurationSlice) FromModel(slice models.ConfigurationSlice) error {

	for _, mod := range slice {
		c := Configuration{}
		err := c.FromModel(mod)
		if err != nil {
			return eris.Cause(err)
		}
		*cs = append(*cs, &c)
	}

	return nil
}

func (configuration *Configuration) FromModel(mod *models.Configuration) error {

	configuration.Key = mod.Key
	configuration.Value = mod.Value
	configuration.Description = mod.Description.String
	configuration.UpdatedAt = mod.UpdatedAt.Ptr()
	configuration.CreatedAt = mod.CreatedAt.Time

	return nil
}

func (filter *ConfigurationFilter) QueryMod() qmods.QueryModsSlice {

	mods := make(qmods.QueryModsSlice, 0)

	if filter == nil {
		return mods
	}

	if len(filter.Key) > 0 {
		mods = append(mods, filter.Key.AndIn(models.ConfigurationColumns.Key))
	}

	if len(filter.Value) > 0 {
		mods = append(mods, filter.Value.AndIn(models.ConfigurationColumns.Value))
	}
	return mods
}
