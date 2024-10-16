package supertokens

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/utils"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionContainer, err := session.GetSession(r, w, nil)
		if err != nil {
			w.Header().Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			w.Write([]byte(fmt.Sprintf(`{"error": %v}`, err.Error())))
			return
		}

		userID := sessionContainer.GetUserID()
		tenantID := supertokens.DefaultTenantId // tenantID := sessionContainer.GetTenantId()
		role, err := getUserRole(userID)
		if err != nil {
			w.Header().Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			w.Write([]byte(fmt.Sprintf(`{"error": %v}`, err.Error())))
			return
		}

		// TODO: delete
		log.Printf("user_id [%v], tenant_id [%v], role [%v]", userID, tenantID, role)

		ctx := context.WithValue(r.Context(), RoleContextKey, role)
		ctx = context.WithValue(ctx, UserIDContextKey, userID)
		ctx = context.WithValue(ctx, TenantIDContextKey, tenantID)

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
