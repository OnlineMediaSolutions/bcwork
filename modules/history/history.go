package history

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	// subjects
	GlobalFactorSubject       = "Serving Fees"
	UserSubject               = "User"
	DPOSubject                = "DPO"
	PublisherSubject          = "Publisher"
	BlockPublisherSubject     = "Blocks - Publisher"
	ConfiantPublisherSubject  = "Confiant - Publisher"
	PixalatePublisherSubject  = "Pixalate - Publisher"
	DomainSubject             = "Domain"
	BlockDomainSubject        = "Blocks - Domain"
	ConfiantDomainSubject     = "Confiant - Domain"
	PixalateDomainSubject     = "Pixalate - Domain"
	FactorSubject             = "Bidder Targeting"
	JSTargetingSubject        = "JS Targeting"
	FloorSubject              = "Floor"
	FactorAutomationSubject   = "Factor Automation"
	BidCachingSubject         = "Max Creatives in Cache - Publisher"
	BidCachingDomainSubject   = "Max Creatives in Cache - Domain"
	RefreshCacheSubject       = "Max Client Refresh - Publisher"
	RefreshCacheDomainSubject = "Max Client Refresh - Domain"

	// actions
	createdAction = "Created"
	updatedAction = "Updated"
	deletedAction = "Deleted"
	unknownAction = "Unknown"
)

type HistoryModule interface {
	SaveAction(ctx context.Context, oldValue, newValue any, options *HistoryOptions)
}

type HistoryClient struct {
}

var _ HistoryModule = (*HistoryClient)(nil)

func NewHistoryClient() *HistoryClient {
	return &HistoryClient{}
}

func (h *HistoryClient) SaveAction(ctx context.Context, oldValue, newValue any, options *HistoryOptions) {
	var (
		subject          string
		isMultipleValues bool
	)

	if options != nil {
		subject = options.Subject
		isMultipleValues = options.IsMultipleValuesExpected
	} else {
		requestPathValue := ctx.Value(constant.RequestPathContextKey)
		requestPath, ok := requestPathValue.(string)
		if !ok {
			logger.Logger(ctx).Error().Msgf("cannot cast requestPath to string")

			return
		}

		subject = subjectsMap[requestPath]
		if subject == "" {
			logger.Logger(ctx).Error().Msg("no subject found")

			return
		}

		isMultipleValues = isMultipleValuesExpected(requestPath)
	}

	userIDValue := ctx.Value(constant.UserIDContextKey)
	userID, ok := userIDValue.(int)
	if !ok {
		// if there is no user_id then history saving was invoked directly
		userID = supertokens_module.AutomationUserID
	}

	innerCtx := context.WithValue(context.Background(), constant.LoggerContextKey, logger.Logger(ctx))

	go h.saveAction(innerCtx, userID, subject, isMultipleValues, oldValue, newValue)
}

func (h *HistoryClient) saveAction(
	ctx context.Context,
	userID int,
	subject string,
	isMultipleValuesExpected bool,
	oldValue any,
	newValue any,
) {
	var (
		oldValues = []any{oldValue}
		newValues = []any{newValue}
		ok        bool
	)

	if isMultipleValuesExpected {
		oldValues, ok = oldValue.([]any)
		if !ok {
			logger.Logger(ctx).Error().Msgf("cannot cast old value (from bulk) to []any")

			return
		}

		newValues, ok = newValue.([]any)
		if !ok {
			logger.Logger(ctx).Error().Msgf("cannot cast new value (from bulk) to []any")

			return
		}
	}

	if len(oldValues) != len(newValues) {
		logger.Logger(ctx).Error().Msgf("amount of old values [%v] not equal amount of new values [%v]", len(oldValues), len(newValues))

		return
	}

	for i := 0; i < len(oldValues); i++ {
		oldValue := oldValues[i]
		newValue := newValues[i]

		action, err := getAction(oldValue, newValue)
		if err != nil {
			logger.Logger(ctx).Error().Msgf("cannot get action: %v", err.Error())
			continue
		}

		valueForItem := newValue
		if action == deletedAction {
			valueForItem = oldValue
		}

		item, err := getItem(subject, valueForItem)
		if err != nil {
			logger.Logger(ctx).Error().Msgf("cannot get item: %v", err.Error())
			continue
		}

		changes, err := getChanges(action, subject, oldValue, newValue)
		if err != nil {
			logger.Logger(ctx).Error().Msgf("cannot get changes for action [%v]: %v", action, err.Error())
			continue
		}

		var oldValueData []byte
		if oldValue != nil {
			oldValueData, err = json.Marshal(oldValue)
			if err != nil {
				logger.Logger(ctx).Error().Msgf("cannot marshal oldValue: %v", err.Error())
				continue
			}
		}

		var newValueData []byte
		if newValue != nil {
			newValueData, err = json.Marshal(newValue)
			if err != nil {
				logger.Logger(ctx).Error().Msgf("cannot marshal newValue: %v", err.Error())
				continue
			}
		}

		mod := models.History{
			UserID:          userID,
			Subject:         subject,
			Action:          action,
			Item:            item.key,
			PublisherID:     null.StringFromPtr(item.publisherID),
			Domain:          null.StringFromPtr(item.domain),
			DemandPartnerID: null.StringFromPtr(item.demandPartnerID),
			EntityID:        null.StringFromPtr(item.entityID),
			OldValue:        null.JSONFrom(oldValueData),
			NewValue:        null.JSONFrom(newValueData),
			Changes:         null.JSONFrom(changes),
			Date:            time.Now().UTC(),
		}

		err = mod.Insert(context.Background(), bcdb.DB(), boil.Infer())
		if err != nil {
			logger.Logger(ctx).Error().Msgf("cannot insert history data: %v", err.Error())
			continue
		}

		logger.Logger(ctx).Debug().Msgf("history with action [%v] for subject [%v] with id [%v] successfully insert", action, subject, mod.ID)
	}
}

func getAction(oldValue, newValue any) (string, error) {
	switch {
	case newValue != nil && oldValue == nil:
		return createdAction, nil
	case newValue != nil && oldValue != nil:
		return updatedAction, nil
	case newValue == nil && oldValue != nil:
		return deletedAction, nil
	}

	return "", errors.New("unknown action")
}

func isMultipleValuesExpected(requestPath string) bool {
	return strings.Contains(requestPath, "/bulk/") ||
		requestPath == "/dpo/delete" ||
		requestPath == "/pixalate/delete" ||
		requestPath == "/factor/delete" ||
		requestPath == "/floor/delete" ||
		requestPath == "/refresh_cache/delete" ||
		requestPath == "/bid_caching/delete" ||
		requestPath == "/bid_caching/delete?domain=true" ||
		requestPath == "/refresh_cache/delete?domain=true"
}
