package supertokens

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/supertokens/supertokens-golang/recipe/session"
	session_errors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
)

func (c *SuperTokensClient) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		sessionContainer, err := session.GetSession(r, w, nil)
		if err != nil {
			switch err.(type) {
			case session_errors.TryRefreshTokenError:
				w.Write([]byte(`{"error": "Session expired. Try to refresh."}`))
				w.WriteHeader(http.StatusForbidden)
				return
			case session_errors.UnauthorizedError:
				isLocalHost, err := isRemoteAddrEqualLocalHost(r)
				if err != nil {
					w.Write([]byte(fmt.Sprintf(`{"error": "can't check request address: %v"}`, err.Error())))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if c.skipSessionVerificationForLocalHost && isLocalHost {
					ctx := context.WithValue(r.Context(), RoleContextKey, AdminRoleName)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				} else {
					w.Write([]byte(`{"error": "unauthorized"}`))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			default:
				w.Write([]byte(fmt.Sprintf(`{"error": "can't get session: %v"}`, err.Error())))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		userID := sessionContainer.GetUserID()
		email, err := getEmailByUserID(userID)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "can't get user email: %v"}`, err.Error())))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := getUserByEmail(r.Context(), email)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "can't get user: %v"}`, err.Error())))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.Enabled {
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errUserDisabled.Error())))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailContextKey, user.Email)
		ctx = context.WithValue(ctx, RoleContextKey, user.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (sc *SuperTokensClient) AdminRoleRequired(c *fiber.Ctx) error {
	role := c.Context().Value(RoleContextKey)

	if role == nil || role.(string) != AdminRoleName {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin role required", errors.New("current user doesn't have admin role"))
	}

	return c.Next()
}

func isRemoteAddrEqualLocalHost(r *http.Request) (bool, error) {
	ip := r.RemoteAddr

	ip, _, err := net.SplitHostPort(ip)
	if err != nil {
		return false, err
	}

	return ip == "127.0.0.1", nil
}

func getEmailByUserID(userID string) (string, error) {
	user, err := thirdpartyemailpassword.GetUserById(userID)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}

func getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	mod, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, bcdb.DB())
	if err != nil {
		return nil, err
	}

	return mod, nil
}
