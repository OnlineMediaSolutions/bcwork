package constant

import (
	"time"

	"github.com/m6yf/bcwork/models"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/volatiletech/null/v8"
)

type User struct {
	ID               int        `json:"id"`
	UserID           string     `json:"-"` // supertokens ID
	FirstName        string     `json:"first_name" validate:"required"`
	LastName         string     `json:"last_name" validate:"required"`
	Email            string     `json:"email" validate:"email"`
	Role             string     `json:"role" validate:"role"`
	OrganizationName string     `json:"organization_name" validate:"required"`
	Address          string     `json:"address" validate:"required"`
	Phone            string     `json:"phone" validate:"phone"`
	Enabled          bool       `json:"enabled"`
	CreatedAt        time.Time  `json:"created_at"`
	DisabledAt       *time.Time `json:"disabled_at"`
}

func (u *User) FromModel(mod *models.User) {
	u.ID = mod.ID
	u.UserID = mod.UserID
	u.FirstName = mod.FirstName.String
	u.LastName = mod.LastName.String
	u.Email = mod.Email
	u.Role = mod.Role
	u.OrganizationName = mod.OrganizationName
	u.Address = mod.Address
	u.Phone = mod.Phone
	u.Enabled = mod.Enabled
	u.CreatedAt = mod.CreatedAt
	u.DisabledAt = mod.DisabledAt.Ptr()
}

func (u User) ToModel() *models.User {
	return &models.User{
		UserID:           u.UserID,
		FirstName:        null.StringFrom(u.FirstName),
		LastName:         null.StringFrom(u.LastName),
		Email:            u.Email,
		Role:             u.Role,
		OrganizationName: u.OrganizationName,
		Address:          u.Address,
		Phone:            u.Phone,
		Enabled:          u.Enabled,
		CreatedAt:        time.Now(),
	}
}

// func (u User) ExtractMetaData() map[string]interface{} {
// 	return map[string]interface{}{
// 		supertokens_module.SuperTokensMetaDataFirstNameKey:        u.FirstName,
// 		supertokens_module.SuperTokensMetaDataLastNameKey:         u.LastName,
// 		supertokens_module.SuperTokensMetaDataOrganizationNameKey: u.OrganizationName,
// 		supertokens_module.SuperTokensMetaDataAddressKey:          u.Address,
// 		supertokens_module.SuperTokensMetaDataPhoneKey:            u.Phone,
// 		supertokens_module.SuperTokensMetaDataEnabledKey:          u.Enabled,
// 		supertokens_module.SuperTokensMetaDataDisabledAtKey:       u.DisabledAt,
// 	}
// }

func ExtractUserMetaDataFromModel(mod *models.User) map[string]interface{} {
	var disabledAt *time.Time
	if mod.DisabledAt.Valid {
		disabledAt = &mod.DisabledAt.Time
	}

	return map[string]interface{}{
		supertokens_module.SuperTokensMetaDataFirstNameKey:        mod.FirstName.String,
		supertokens_module.SuperTokensMetaDataLastNameKey:         mod.LastName.String,
		supertokens_module.SuperTokensMetaDataOrganizationNameKey: mod.OrganizationName,
		supertokens_module.SuperTokensMetaDataAddressKey:          mod.Address,
		supertokens_module.SuperTokensMetaDataPhoneKey:            mod.Phone,
		supertokens_module.SuperTokensMetaDataEnabledKey:          mod.Enabled,
		supertokens_module.SuperTokensMetaDataDisabledAtKey:       disabledAt,
	}
}
