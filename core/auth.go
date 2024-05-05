package core

import (
	"context"
	"database/sql"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"

	"github.com/rotisserie/eris"
)

func AuthToUserID(ctx context.Context, authUserID string, impersonate bool) (string, error) {
	auth, err := models.Auths(models.AuthWhere.AuthUserID.EQ(authUserID)).One(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return "", eris.Wrapf(err, "failed to fetch auth record(auth_user_id:%s,http:500)", auth)
	}
	if auth == nil {
		return "", eris.Errorf("auth record not found(auth_user_id:%s,http:500)", auth)
	}

	if impersonate && auth.ImpersonateAsID.Valid {
		return auth.ImpersonateAsID.String, nil
	}
	return auth.UserID, nil
}
