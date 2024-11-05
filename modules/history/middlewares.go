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

	"/bulk/global/factor":               globalFactorSubject, // TODO: save action
	"/publisher/new":                    publisherSubject,
	"/publisher/update":                 publisherSubject,
	"/floor":                            floorSubject,  // TODO: save action
	"/factor":                           factorSubject, // TODO: save action
	"/global/factor":                    globalFactorSubject,
	"/dpo/set":                          dpoSubject, // TODO: save action
	"/dpo/delete":                       dpoSubject, // TODO: save action
	"/dpo/update":                       dpoSubject, // TODO: save action
	"/publisher/domain":                 domainSubject,
	"/publisher/domain?automation=true": factorAutomationSubject,
	"/targeting/set":                    jsTargetingSubject,
	"/targeting/update":                 jsTargetingSubject,
	"/user/set":                         userSubject,
	"/user/update":                      userSubject,
	"/block":                            blockPublisherSubject,    // TODO: save action
	"/pixalate":                         pixalatePublisherSubject, // TODO: save action
	"/confiant":                         confiantPublisherSubject, // TODO: save action
	"/block?domain=true":                blockDomainSubject,       // TODO: save action
	"/pixalate?domain=true":             pixalateDomainSubject,    // TODO: save action
	"/confiant?domain=true":             confiantDomainSubject,    // TODO: save action
}

func (h *HistoryClient) HistoryMiddleware(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		return err
	}

	ctx := c.Context()

	subject := getSubject(ctx, string(c.Request().RequestURI()))
	if subject == "" {
		logger.Logger(ctx).Debug().Msg("no subject found")
		return nil
	}
	logger.Logger(ctx).Debug().Msg(subject)

	requestIDValue := ctx.Value(constant.RequestIDContextKey)
	requestID, ok := requestIDValue.(string)
	if !ok {
		logger.Logger(ctx).Debug().Msgf("cannot cast requestID to string")
		return nil
	}

	userIDValue := ctx.Value(constant.UserIDContextKey)
	userID, ok := userIDValue.(int)
	if !ok {
		logger.Logger(ctx).Debug().Msgf("cannot cast userID to int")
		return nil
	}

	logger.Logger(ctx).Debug().Msgf("[HistoryClient] requestID - %v", requestID)

	innerCtx := context.WithValue(context.Background(), constant.LoggerContextKey, logger.Logger(ctx))
	go h.saveAction(innerCtx, userID, requestID, subject)

	return nil
}
