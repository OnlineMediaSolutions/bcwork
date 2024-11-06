package history

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/m6yf/bcwork/utils/constant"
)

var subjectsMap = map[string]string{
	// TODO: these bulks require different flow
	// "/bulk/dpo":    dpoSubject,
	// "/bulk/floor":  floorSubject,
	// "/bulk/factor": factorSubject,

	"/bulk/global/factor":               globalFactorSubject,
	"/publisher/new":                    publisherSubject,
	"/publisher/update":                 publisherSubject,
	"/floor":                            floorSubject,
	"/factor":                           factorSubject,
	"/global/factor":                    globalFactorSubject,
	"/dpo/set":                          dpoSubject,
	"/dpo/delete":                       dpoSubject,
	"/dpo/update":                       dpoSubject,
	"/publisher/domain":                 domainSubject,
	"/publisher/domain?automation=true": factorAutomationSubject,
	"/targeting/set":                    jsTargetingSubject,
	"/targeting/update":                 jsTargetingSubject,
	"/user/set":                         userSubject,
	"/user/update":                      userSubject,
	"/block":                            blockPublisherSubject,
	"/pixalate":                         pixalatePublisherSubject,
	"/confiant":                         confiantPublisherSubject,
	"/block?domain=true":                blockDomainSubject,
	"/pixalate?domain=true":             pixalateDomainSubject,
	"/confiant?domain=true":             confiantDomainSubject,
}

func (h *HistoryClient) HistoryMiddleware(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		return err
	}

	ctx := c.Context()
	requestPath := string(c.Request().RequestURI())
	logger.Logger(ctx).Debug().Msgf("requestPath - %v", requestPath)

	subject := subjectsMap[requestPath]
	if subject == "" {
		logger.Logger(ctx).Error().Msg("no subject found")
		return nil
	}
	logger.Logger(ctx).Debug().Msg(subject)

	requestIDValue := ctx.Value(constant.RequestIDContextKey)
	requestID, ok := requestIDValue.(string)
	if !ok {
		logger.Logger(ctx).Error().Msgf("cannot cast requestID to string")
		return nil
	}

	userIDValue := ctx.Value(constant.UserIDContextKey)
	userID, ok := userIDValue.(int)
	if !ok {
		logger.Logger(ctx).Error().Msgf("cannot cast userID to int")
		return nil
	}

	logger.Logger(ctx).Debug().Msgf("[HistoryClient] requestID - %v", requestID)

	innerCtx := context.WithValue(context.Background(), constant.LoggerContextKey, logger.Logger(ctx))
	go h.saveAction(innerCtx, userID, requestID, subject, requestPath)

	return nil
}
