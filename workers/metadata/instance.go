package metadata

import (
	"context"
	"encoding/json"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

// MetadataInstance is an object representing metadata storage instance (usually redis or rest gateway).
type MetadataInstance struct {
	InstanceID string          `json:"instance_id"`
	Bitwise    int64           `json:"bitwise"`
	Type       string          `json:"type"`
	Config     json.RawMessage `json:"config"`
	Updater    Updater         `json:"-"`
}

type MetadataInstanceSlice []*MetadataInstance

func (mi *MetadataInstance) FromModel(mod *models.MetadataInstance) error {
	mi.InstanceID = mod.InstanceID
	mi.Bitwise = mod.Bitwise
	mi.Type = mod.Type
	mi.Config = mod.Config.JSON

	err := mi.initUpdate()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (mis *MetadataInstanceSlice) FromModel(slice models.MetadataInstanceSlice) error {
	for _, mod := range slice {
		mi := MetadataInstance{}
		err := mi.FromModel(mod)
		if err != nil {
			return errors.WithMessagef(err, "failed to convert metadata instance to core object")
		}
		*mis = append(*mis, &mi)
	}

	return nil
}

func GetMetadataInstances(ctx context.Context) (MetadataInstanceSlice, error) {
	mis, err := models.MetadataInstances().All(ctx, bcdb.DB())
	if err != nil {
		return nil, errors.Wrap(err, "error while to fetch metadata instances from db")
	}

	if mis == nil {
		mis = models.MetadataInstanceSlice{}
	}

	res := make(MetadataInstanceSlice, 0)
	err = res.FromModel(mis)
	if err != nil {
		return nil, errors.Wrap(err, "error while to converting to metadata instances from db to core")
	}

	return res, nil
}

// SetBit register instance update for specific metadata record
func (mi *MetadataInstance) SetBit(ctx context.Context, mod *models.MetadataQueue) error {
	_, err := queries.Raw("UPDATE metadata_queue SET commited_instances = commited_instances | $1 WHERE transaction_id=$2", mi.Bitwise, mod.TransactionID).
		ExecContext(ctx, bcdb.DB())
	if err != nil {
		return errors.Wrapf(err, "failed to turn on metadata instance bit(transaction_id:%s,instance:%s)", mod.TransactionID, mi.InstanceID)
	}

	return nil
}

// initUpdate is a generic function that initiate instance updater and send metadata
func (mi *MetadataInstance) initUpdate() error {
	switch mi.Type {
	case "redis":
		redisUpdater := RedisUpdater{}
		err := json.Unmarshal(mi.Config, &redisUpdater)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal redis updater config(instance:%s)", mi.InstanceID)
		}
		mi.Updater = redisUpdater

	case "http":
		httpUpdater := HttpUpdater{}
		err := json.Unmarshal(mi.Config, &httpUpdater)
		if err != nil {
			return errors.Wrapf(err, "failed to unmarshal http updater config(instance:%s)", mi.InstanceID)
		}
		mi.Updater = httpUpdater

	default:
		return errors.Errorf("unsupported metadata instance type(type:%s)", mi.Type)
	}

	return nil
}
