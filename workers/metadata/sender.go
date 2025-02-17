package metadata

import (
	"context"
	"encoding/json"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Send add metadata json to queue to be sent to metadata front storage (redis)
func Send(ctx context.Context, key string, value interface{}, exec boil.ContextExecutor) error {
	now := time.Now()
	valJson, err := json.Marshal(value)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal metadata value into json(key:%s)", key)
	}

	mod := models.MetadataQueue{
		TransactionID:     bcguid.NewFromf(now, key),
		Key:               key,
		Value:             valJson,
		CommitedInstances: 0,
		CreatedAt:         now,
	}

	err = mod.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return errors.Wrapf(err, "failed to insert metadata record into queue(key:%s)", key)
	}

	return nil
}
