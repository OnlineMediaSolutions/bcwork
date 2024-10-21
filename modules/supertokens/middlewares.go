package supertokens

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
)

func VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		sessionContainer, err := session.GetSession(r, w, nil)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "can't get session: %v"}`, err.Error())))
			return
		}

		// TODO: redirect POST http://localhost:8000/auth/session/refresh

		userID := sessionContainer.GetUserID()
		email, err := getEmailByUserID(userID)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "can't get user email: %v"}`, err.Error())))
			return
		}

		user, err := getUserByEmail(r.Context(), email)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"error": "can't get user: %v"}`, err.Error())))
			return
		}

		if !user.Enabled {
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errUserDisabled.Error())))
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailContextKey, user.Email)
		ctx = context.WithValue(ctx, RoleContextKey, user.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminRoleRequired(c *fiber.Ctx) error {
	role := c.Context().Value(RoleContextKey)

	if role == nil || role.(string) != AdminRoleName {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin role required", errors.New("current user doesn't have admin role"))
	}

	return c.Next()
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
