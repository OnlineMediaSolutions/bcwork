package core

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"slices"
	"time"

	"math/rand"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/m6yf/bcwork/bcdb/pagination"
	"github.com/m6yf/bcwork/bcdb/qmods"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rotisserie/eris"
	"github.com/supertokens/supertokens-golang/supertokens"
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

	supertokens.GetUsersNewestFirst(supertokens.DefaultTenantId, nil, nil, nil, nil)

	users := make([]*constant.User, 0, len(mods))
	for _, mod := range mods {
		user := new(constant.User)
		user.FromModel(mod)
		users = append(users, user)
	}

	return users, nil
}

func (u *UserService) CreateUser(ctx context.Context, data *constant.User) error {
	tempPassword := generateTemporaryPassword()

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

	// TODO: uncomment
	// err = sendRegistrationEmail(mod.Email, tempPassword)
	// if err != nil {
	// 	return eris.Wrap(err, "failed to send email with temporary credentials")
	// }

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

func generateTemporaryPassword() string {
	const (
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits  = "0123456789"
		// special = "!@#$%^&*()"

		length     = 10
		useLetters = true
		useSpecial = false
		useNum     = true
	)

	b := make([]byte, length)
	b[0] = letters[rand.Intn(len(letters))] // Ensure at least one letter
	b[1] = digits[rand.Intn(len(digits))]   // Ensure at least one digit

	combined := letters + digits //+ special
	for i := 2; i < length; i++ {
		b[i] = combined[rand.Intn(len(combined))]
	}

	return string(b)
}

func sendRegistrationEmail(email, password string) error {
	registrationTemplate := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Registration Successful</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 20px;
					background-color: #f9f9f9;
				}
				.container {
					max-width: 500px;
					margin: auto;
					padding: 20px;
					background: #fff;
					border: 1px solid #ddd;
					border-radius: 5px;
					box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
				}
				h1 {
					color: #333;
				}
				.credentials {
					margin-top: 20px;
					padding: 10px;
					border: 1px solid #ccc;
					border-radius: 5px;
					background: #f1f1f1;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Welcome to OMS!</h1>
				<p>An account has been created for you in the OMS. Here are your temporary credentials:</p>
				<div class="credentials">
					<p><strong>Email:</strong> {{ .Email }}</p>
					<p><strong>Password:</strong> {{ .Password }}</p>
				</div>
				<p>Please <a href="https://login.nanoook.com/auth/forgot-password">change password</a></p>
			</div>
		</body>
		</html>
	`

	type UserCredentials struct {
		Email    string
		Password string
	}

	credentials := UserCredentials{
		Email:    email,
		Password: password,
	}

	tmpl := template.Must(template.New("registrationTemplate").Parse(registrationTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, credentials); err != nil {
		return err
	}

	return modules.SendEmail(modules.EmailRequest{
		To:      []string{email},
		Subject: "Temporary credentials for OMS",
		Bcc:     email,
		Body:    buf.String(),
		IsHTML:  true,
	})
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