package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService struct {
	supertokenClient *supertokens_module.SuperTokensClient
}

func NewUserService(supertokenClient *supertokens_module.SuperTokensClient) *UserService {
	return &UserService{supertokenClient: supertokenClient}
}

type UserOptions struct {
	Filter     UserFilter             `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type UserFilter struct {
	FirstName        filter.StringArrayFilter `json:"first_name,omitempty"`
	LastName         filter.StringArrayFilter `json:"last_name,omitempty"`
	Email            filter.StringArrayFilter `json:"email,omitempty"`
	Role             filter.StringArrayFilter `json:"role,omitempty"`
	OrganizationName filter.StringArrayFilter `json:"organization_name,omitempty"`
	Address          filter.StringArrayFilter `json:"address,omitempty"`
	Phone            filter.StringArrayFilter `json:"phone,omitempty"`
	Enabled          filter.BoolFilter        `json:"enabled,omitempty"` // TODO: process correct for all
}

func (filter *UserFilter) queryMod() qmods.QueryModsSlice {
	mods := make(qmods.QueryModsSlice, 0)
	if filter == nil {
		return mods
	}

	if len(filter.FirstName) > 0 {
		mods = append(mods, filter.FirstName.AndIn(models.UserColumns.FirstName))
	}

	if len(filter.LastName) > 0 {
		mods = append(mods, filter.LastName.AndIn(models.UserColumns.LastName))
	}

	if len(filter.Email) > 0 {
		mods = append(mods, filter.Email.AndIn(models.UserColumns.Email))
	}

	if len(filter.Role) > 0 {
		mods = append(mods, filter.Role.AndIn(models.UserColumns.Role))
	}

	if len(filter.OrganizationName) > 0 {
		mods = append(mods, filter.OrganizationName.AndIn(models.UserColumns.OrganizationName))
	}

	if len(filter.Address) > 0 {
		mods = append(mods, filter.Address.AndIn(models.UserColumns.Address))
	}

	if len(filter.Phone) > 0 {
		mods = append(mods, filter.Phone.AndIn(models.UserColumns.Phone))
	}

	// if len(filter.Enabled) > 0 {
	// 	mods = append(mods, filter.Enabled.Where(models.UserColumns.Enabled))
	// }

	return mods
}

func (u *UserService) GetUsers(ctx context.Context, ops *UserOptions) ([]*constant.User, error) {
	qmods := ops.Filter.queryMod().
		Order(ops.Order, nil, models.UserColumns.UserID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Users(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve users")
	}

	users := make([]*constant.User, 0, len(mods))
	for _, mod := range mods {
		user := new(constant.User)
		user.FromModel(mod)
		users = append(users, user)
	}

	return users, nil
}

func (u *UserService) CreateUser(ctx context.Context, data *constant.User) error {
	tempPassword := "abcd1234" // TODO: temp password

	userID, err := u.supertokenClient.CreateUser(ctx, data.Email, tempPassword)
	if err != nil {
		return eris.Wrap(err, "failed to create user in supertoken")
	}

	mod := data.ToModel()
	mod.UserID = userID

	err = mod.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to create user in supertoken")
	}

	// TODO: save email to allow use it for third party providers
	// TODO: send email to user with credentials

	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, data *constant.User) error {
	mod, err := models.Users(models.UserWhere.ID.EQ(data.ID)).One(ctx, bcdb.DB())
	if err != nil {
		return eris.Wrap(err, fmt.Sprintf("failed to get user with id [%v] to update", data.ID))
	}

	columns, err := prepareDataForUpdate(data, mod)
	if err != nil {
		return eris.Wrap(err, "error preparing data for update")
	}
	if len(columns) == 0 {
		return errors.New("there are no new values to update user")
	}

	if isEnabledUpdating(columns) && !mod.Enabled {
		err := u.supertokenClient.RevokeAllSessionsForUser(mod.UserID)
		if err != nil {
			return errors.New("error revoking all sessions in supertoken")
		}
	}

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(columns...))
	if err != nil {
		return eris.Wrap(err, "failed to update targeting")
	}

	return nil
}

func prepareDataForUpdate(newData *constant.User, currentData *models.User) ([]string, error) {
	columns := make([]string, 0, 8)

	// first_name
	if newData.FirstName != currentData.FirstName.String {
		currentData.FirstName = null.StringFrom(newData.FirstName)
		columns = append(columns, models.UserColumns.FirstName)
	}
	// last_name
	if newData.LastName != currentData.LastName.String {
		currentData.LastName = null.StringFrom(newData.LastName)
		columns = append(columns, models.UserColumns.LastName)
	}
	// organization_name
	if newData.OrganizationName != currentData.OrganizationName {
		currentData.OrganizationName = newData.OrganizationName
		columns = append(columns, models.UserColumns.OrganizationName)
	}
	// address
	if newData.Address != currentData.Address {
		currentData.Address = newData.Address
		columns = append(columns, models.UserColumns.Address)
	}
	// phone
	if newData.Phone != currentData.Phone {
		currentData.Phone = newData.Phone
		columns = append(columns, models.UserColumns.Phone)
	}
	// role
	if newData.Role != currentData.Role {
		currentData.Role = newData.Role
		columns = append(columns, models.UserColumns.Role)
	}
	// enabled
	if newData.Enabled != currentData.Enabled {
		currentData.Enabled = newData.Enabled
		// if user becomes enabled, clear disabled_at time; else fill disabled_at
		if newData.Enabled {
			currentData.DisabledAt = null.NewTime(time.Time{}, false)
		} else {
			currentData.DisabledAt = null.TimeFrom(time.Now())
		}
		columns = append(columns, models.UserColumns.Enabled)
		columns = append(columns, models.UserColumns.DisabledAt)
	}

	return columns, nil
}

func isRoleUpdating(columns []string) bool {
	return slices.Contains(columns, models.UserColumns.Role)
}

func isEnabledUpdating(columns []string) bool {
	return slices.Contains(columns, models.UserColumns.Enabled)
}

// creating new tenant
// tenantId := "admin"
// emailPasswordEnabled := true
// thirdPartyEnabled := true
// passwordlessEnabled := true
// resp, err := multitenancy.CreateOrUpdateTenant(tenantId, multitenancymodels.TenantConfig{
// 	EmailPasswordEnabled: &emailPasswordEnabled,
// 	ThirdPartyEnabled:    &thirdPartyEnabled,
// 	PasswordlessEnabled:  &passwordlessEnabled,
// })
// if err != nil {
// 	return c.Status(500).JSON(err.Error())
// }
// fmt.Printf("%#v\n", resp)

// add tenant to user
// if id == "c91e28d7-7a74-4229-b11c-9300391a4dfd" {
// 	resp, err := multitenancy.AssociateUserToTenant(tenantId, id)
// 	if err != nil {
// 		return c.Status(500).JSON(err.Error())
// 	}
// 	fmt.Printf("%#v\n", resp)
// }

// get roles
// resp2, err := userroles.GetRolesForUser("public", id, nil)
// if err != nil {
// 	return c.Status(500).JSON(err.Error())
// }
// fmt.Printf("%#v\n", resp2)

// getting users
// p, err := supertokens.GetUsersNewestFirst("public", nil, nil, nil, nil)
// if err != nil {
// 	return nil, err
// }
// var needMetaData bool
// if needMetaData {
// 	for i, user := range p.Users {
// 		id := user.User[auth.SuperTokensIDKey].(string)
// 		metadata, err := usermetadata.GetUserMetadata(id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Printf("%v. %#v", i+1, metadata)
// 	}
// }
// users := make([]constant.User, 0, len(p.Users))
// for _, user := range p.Users {
// 	email, ok := user.User[auth.SuperTokensEmailKey].(string)
// 	if !ok {
// 		log.Printf("error casting [%v] to string", user.User[auth.SuperTokensEmailKey])
// 	}
// 	timeJoined, ok := user.User[auth.SuperTokensTimeJoinedKey].(float64)
// 	if !ok {
// 		log.Printf("error casting [%v] to float64", user.User[auth.SuperTokensTimeJoinedKey])
// 	}

// 	users = append(users, constant.User{
// 		Email:     email,
// 		CreatedAt: time.Unix(0, int64(timeJoined)*int64(time.Millisecond)),
// 	})
// }
