package core

import (
	"bytes"
	"cmp"
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
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/modules/history"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserService struct {
	supertokenClient      supertokens_module.TokenManagementSystem
	historyModule         history.HistoryModule
	sendRegistrationEmail bool // Temporary, remove after decoupling email sender service
}

func NewUserService(
	supertokenClient supertokens_module.TokenManagementSystem,
	historyModule history.HistoryModule,
	sendRegistrationEmail bool,
) *UserService {
	return &UserService{
		supertokenClient:      supertokenClient,
		historyModule:         historyModule,
		sendRegistrationEmail: sendRegistrationEmail,
	}
}

type UserOptions struct {
	Filter     UserFilter             `json:"filter"`
	Pagination *pagination.Pagination `json:"pagination"`
	Order      order.Sort             `json:"order"`
	Selector   string                 `json:"selector"`
}

type UserFilter struct {
	FirstName        filter.StringArrayFilter   `json:"first_name,omitempty"`
	LastName         filter.StringArrayFilter   `json:"last_name,omitempty"`
	Email            filter.StringArrayFilter   `json:"email,omitempty"`
	Role             filter.StringArrayFilter   `json:"role,omitempty"`
	Types            filter.String2DArrayFilter `json:"types,omitempty"`
	OrganizationName filter.StringArrayFilter   `json:"organization_name,omitempty"`
	Address          filter.StringArrayFilter   `json:"address,omitempty"`
	Phone            filter.StringArrayFilter   `json:"phone,omitempty"`
	Enabled          filter.BoolFilter          `json:"enabled,omitempty"`
}

func (u *UserService) GetUsers(ctx context.Context, ops *UserOptions) ([]*dto.User, error) {
	qmods := ops.Filter.queryMod().
		Order(ops.Order, nil, models.UserColumns.UserID).
		AddArray(ops.Pagination.Do())

	mods, err := models.Users(qmods...).All(ctx, bcdb.DB())
	if err != nil && err != sql.ErrNoRows {
		return nil, eris.Wrap(err, "failed to retrieve users")
	}

	users := make([]*dto.User, 0, len(mods))
	for _, mod := range mods {
		user := new(dto.User)
		user.FromModel(mod)
		users = append(users, user)
	}

	return users, nil
}

func (u *UserService) GetUserInfo(ctx context.Context, userID string) (*dto.User, error) {
	email, err := u.supertokenClient.GetEmailByUserID(userID)
	if err != nil {
		return nil, eris.Wrap(err, "failed to get user email")
	}

	mod, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, "failed to get user from db")
	}

	user := new(dto.User)
	user.FromModel(mod)

	return user, nil
}

func (u *UserService) GetUsersByTypes(ctx context.Context) (*dto.UsersByTypes, error) {
	mods, err := models.Users(models.UserWhere.Types.IsNotNull()).All(ctx, bcdb.DB())
	if err != nil {
		return nil, eris.Wrap(err, "failed to get users with types from db")
	}

	usersByTypes := &dto.UsersByTypes{}
	for _, mod := range mods {
		usersByTypes.Append(mod)
	}

	return usersByTypes, nil
}

func (u *UserService) CreateUser(ctx context.Context, data *dto.User) error {
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

	if u.sendRegistrationEmail {
		err := sendRegistrationEmail(
			mod.Email,
			tempPassword,
			u.supertokenClient.GetWebURL(),
			u.supertokenClient.GetWebURL()+supertokens_module.ChangePasswordSupertokenPath,
		)
		if err != nil {
			return eris.Wrap(err, "failed to send email with temporary credentials")
		}
	}

	u.historyModule.SaveAction(ctx, nil, mod, nil)

	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, data *dto.User) error {
	mod, err := models.Users(models.UserWhere.ID.EQ(data.ID)).One(ctx, bcdb.DB())
	if err != nil {
		return eris.Wrap(err, fmt.Sprintf("failed to get user with id [%v] to update", data.ID))
	}

	oldMod := *mod
	slices.SortStableFunc(data.Types, func(a, b string) int { return cmp.Compare(a, b) })

	columns, err := prepareUserDataForUpdate(data, mod)
	if err != nil {
		return eris.Wrap(err, "error preparing data for update")
	}
	if len(columns) == 0 {
		return errors.New("there are no new values to update user")
	}

	if isEnabledUpdating(columns) && !mod.Enabled {
		err := u.supertokenClient.RevokeAllSessionsForUser(mod.Email)
		if err != nil {
			return errors.New("error revoking all sessions in supertoken")
		}
	}

	_, err = mod.Update(ctx, bcdb.DB(), boil.Whitelist(columns...))
	if err != nil {
		return eris.Wrap(err, "failed to update targeting")
	}

	u.historyModule.SaveAction(ctx, &oldMod, mod, nil)

	return nil
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

	if len(filter.Types) > 0 {
		mods = append(mods, filter.Types.AndIn(models.UserColumns.Types))
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

	if len(filter.Enabled) > 0 {
		mods = append(mods, filter.Enabled.Where(models.UserColumns.Enabled))
	}

	return mods
}

func prepareUserDataForUpdate(newData *dto.User, currentData *models.User) ([]string, error) {
	columns := make([]string, 0, 8)

	if newData.FirstName != currentData.FirstName {
		currentData.FirstName = newData.FirstName
		columns = append(columns, models.UserColumns.FirstName)
	}

	if newData.LastName != currentData.LastName {
		currentData.LastName = newData.LastName
		columns = append(columns, models.UserColumns.LastName)
	}

	if newData.OrganizationName != currentData.OrganizationName {
		currentData.OrganizationName = newData.OrganizationName
		columns = append(columns, models.UserColumns.OrganizationName)
	}

	if newData.Address != currentData.Address.String {
		currentData.Address = helpers.GetNullString(newData.Address)
		columns = append(columns, models.UserColumns.Address)
	}

	if newData.Phone != currentData.Phone.String {
		currentData.Phone = helpers.GetNullString(newData.Phone)
		columns = append(columns, models.UserColumns.Phone)
	}

	if newData.Role != currentData.Role {
		currentData.Role = newData.Role
		columns = append(columns, models.UserColumns.Role)
	}

	if !slices.Equal(newData.Types, currentData.Types) {
		currentData.Types = newData.Types
		columns = append(columns, models.UserColumns.Types)
	}

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

func sendRegistrationEmail(email, password, signInLink, passwordChangelink string) error {
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
				<p>Now you can <a href="{{ .SignInLink }}">sign in</a> using these credentials or third-party providers (e.g. Google, Apple).</p>
				<p>Important: in order to sign in using third-party providers your email must be the same as from above.</p>
				<p>Also please <a href="{{ .PasswordChangeLink }}">change password</a>.</p>
				<p>Temporary password valid for {{ .MaxDaysForTemporaryPassword }} days.</p>
			</div>
		</body>
		</html>
	`

	type UserCredentials struct {
		Email                       string
		Password                    string
		MaxDaysForTemporaryPassword int
		PasswordChangeLink          string
		SignInLink                  string
	}

	credentials := UserCredentials{
		Email:                       email,
		Password:                    password,
		MaxDaysForTemporaryPassword: supertokens_module.MaxDaysForTemporaryPassword,
		PasswordChangeLink:          passwordChangelink,
		SignInLink:                  signInLink,
	}

	tmpl := template.Must(template.New("registrationTemplate").Parse(registrationTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, credentials); err != nil {
		return err
	}

	return modules.SendEmail(modules.EmailRequest{
		To:      []string{email},
		Subject: "Temporary credentials for OMS",
		Bcc:     []string{email},
		Body:    buf.String(),
		IsHTML:  true,
	})
}
