package supertokens

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog/log"
	"github.com/supertokens/supertokens-golang/recipe/session"
	session_errors "github.com/supertokens/supertokens-golang/recipe/session/errors"
)

func (c *SuperTokensClient) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		if c.isAllowedAPIKey(r) {
			ctx := context.WithValue(r.Context(), constant.UserIDContextKey, WorkerUserID)
			ctx = context.WithValue(ctx, constant.RoleContextKey, DeveloperRoleName)
			next.ServeHTTP(w, r.WithContext(ctx))

			return
		}

		sessionContainer, err := session.GetSession(r, w, nil)
		if err != nil {
			var (
				message []byte
				status  int
			)
			switch err.(type) {
			case session_errors.TryRefreshTokenError:
				message = []byte(`{"error": "Session expired. Try to refresh."}`)
				status = http.StatusForbidden

			case session_errors.UnauthorizedError:
				message = []byte(`{"error": "unauthorized"}`)
				status = http.StatusUnauthorized
			default:
				message = []byte(fmt.Sprintf(`{"error": "can't get session: %v"}`, err.Error()))
				status = http.StatusInternalServerError
			}

			_, err := w.Write(message)
			if err != nil {
				log.Error().Err(err).Msg("failed writing response body")
			}
			w.WriteHeader(status)

			return
		}

		userID := sessionContainer.GetUserID()
		email, err := c.GetEmailByUserID(userID)
		if err != nil {
			_, err := w.Write([]byte(fmt.Sprintf(`{"error": "can't get user email: %v"}`, err.Error())))
			if err != nil {
				log.Error().Err(err).Msg("failed writing response body while getting email by user id")
			}
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		user, err := getUserByEmail(r.Context(), email)
		if err != nil {
			_, err = w.Write([]byte(fmt.Sprintf(`{"error": "can't get user: %v"}`, err.Error())))
			if err != nil {
				log.Error().Err(err).Msg("failed writing response body while get user by email")
			}
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if !user.Enabled {
			_, err = w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errNotAllowed.Error())))
			if err != nil {
				log.Error().Err(err).Msg("failed writing response body while checking is user disabled")
			}
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), constant.UserIDContextKey, user.ID)
		ctx = context.WithValue(ctx, constant.UserEmailContextKey, user.Email)
		ctx = context.WithValue(ctx, constant.RoleContextKey, user.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sc *SuperTokensClient) AdminRoleRequired(c *fiber.Ctx) error {
	role := c.Context().Value(constant.RoleContextKey)

	if role == nil || !slices.Contains([]string{AdminRoleName, DeveloperRoleName}, role.(string)) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin role required", errors.New("current user doesn't have admin role"))
	}

	return c.Next()
}

func (sc *SuperTokensClient) isAllowedAPIKey(r *http.Request) bool {
	return slices.Contains(sc.apiKeys, r.Header.Get(constant.HeaderOMSWorkerAPIKey))
}

func getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	mod, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, bcdb.DB())
	if err != nil {
		return nil, fmt.Errorf("can't get user by email: %w", err)
	}

	return mod, nil
}
