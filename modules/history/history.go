package history

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/m6yf/bcwork/storage/cache"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	// subjects
	globalFactorSubject      = "Serving Fees"
	userSubject              = "User"
	dpoSubject               = "DPO"
	publisherSubject         = "Publisher"
	blockPublisherSubject    = "Blocks - Publisher"
	confiantPublisherSubject = "Confiant - Publisher"
	pixalatePublisherSubject = "Pixalate - Publisher"
	domainSubject            = "Domain"
	blockDomainSubject       = "Blocks - Domain"
	confiantDomainSubject    = "Confiant - Domain"
	pixalateDomainSubject    = "Pixalate - Domain"
	factorSubject            = "Bidder Targeting"
	jsTargetingSubject       = "JS Targeting"
	floorSubject             = "Floor"
	factorAutomationSubject  = "Factor Automation"

	// actions
	createdAction = "Created"
	updatedAction = "Updated"
	deletedAction = "Deleted"
	unknownAction = "Unknown"
)

type HistoryModule interface {
	SaveOldAndNewValuesToCache(ctx context.Context, oldValue, newValue any)
}

type HistoryClient struct {
	cache cache.Cache
}

var _ HistoryModule = (*HistoryClient)(nil)

func NewHistoryClient(cache cache.Cache) *HistoryClient {
	return &HistoryClient{
		cache: cache,
	}
}

func (h *HistoryClient) SaveOldAndNewValuesToCache(ctx context.Context, oldValue, newValue any) {
	requestIDValue := ctx.Value(constant.RequestIDContextKey)
	requestID, ok := requestIDValue.(string)
	if !ok {
		logger.Logger(ctx).Debug().Msgf("cannot cast requestID to string")
		return
	}

	h.cache.Set(requestID+":"+cache.HistoryOldValueCacheKey, oldValue)
	h.cache.Set(requestID+":"+cache.HistoryNewValueCacheKey, newValue)
}

func (h *HistoryClient) saveAction(ctx context.Context, userID int, requestID, subject string) {
	oldValue, ok := h.cache.Get(requestID + ":" + cache.HistoryOldValueCacheKey)
	if !ok {
		logger.Logger(ctx).Debug().Msgf("old value not ok")
	}
	h.cache.Delete(requestID + cache.HistoryOldValueCacheKey)

	newValue, ok := h.cache.Get(requestID + ":" + cache.HistoryNewValueCacheKey)
	if !ok {
		logger.Logger(ctx).Debug().Msgf("new value not ok")
	}
	h.cache.Delete(requestID + cache.HistoryNewValueCacheKey)

	logger.Logger(ctx).Debug().Msgf("old value - %v", oldValue)
	logger.Logger(ctx).Debug().Msgf("new value - %v", newValue)

	action, err := getAction(oldValue, newValue)
	if err != nil {
		logger.Logger(ctx).Debug().Msgf("cannot get action: %v", err.Error())
		return
	}
	logger.Logger(ctx).Debug().Msg(action)

	valueForItem := newValue
	if action == deletedAction {
		valueForItem = oldValue
	}
	item, err := getItem(subject, valueForItem)
	if err != nil {
		logger.Logger(ctx).Debug().Msgf("cannot get item: %v", err.Error())
		return
	}
	logger.Logger(ctx).Debug().Msgf("item - %#v", item)

	var changes []byte
	if action == updatedAction {
		changes, err = getChanges(oldValue, newValue)
		if err != nil {
			logger.Logger(ctx).Debug().Msgf("cannot get changes: %v", err.Error())
			return
		}
	}

	var oldValueData []byte
	if oldValue != nil {
		oldValueData, err = json.Marshal(oldValue)
		if err != nil {
			logger.Logger(ctx).Debug().Msgf("cannot marshal oldValue: %v", err.Error())
		}
	}

	var newValueData []byte
	if newValue != nil {
		newValueData, err = json.Marshal(newValue)
		if err != nil {
			logger.Logger(ctx).Debug().Msgf("cannot marshal newValue: %v", err.Error())
		}
	}

	mod := models.History{
		UserID:      userID,
		Subject:     subject,
		Action:      action,
		Item:        item.key,
		PublisherID: null.StringFromPtr(item.publisherID),
		Domain:      null.StringFromPtr(item.domain),
		EntityID:    null.StringFromPtr(item.entityID),
		OldValue:    null.JSONFrom(oldValueData),
		NewValue:    null.JSONFrom(newValueData),
		Changes:     null.JSONFrom(changes),
		Date:        time.Now().UTC(),
	}

	err = mod.Insert(context.Background(), bcdb.DB(), boil.Infer())
	if err != nil {
		logger.Logger(ctx).Debug().Msgf("cannot insert history data: %v", err.Error())
	}

	logger.Logger(ctx).Debug().Msgf("history for subject [%v] successfully insert", subject)
}

func getSubject(ctx context.Context, requestPath string) string {
	logger.Logger(ctx).Debug().Msgf("requestPath - %v", requestPath)
	return subjectsMap[requestPath]
}

func getAction(oldValue, newValue any) (string, error) {
	switch {
	case newValue != nil && oldValue == nil:
		return createdAction, nil
	case newValue != nil && oldValue != nil:
		return updatedAction, nil
	case newValue == nil && oldValue != nil:
		return deletedAction, nil
	default:
		return "", errors.New("unknown action")
	}
}
