package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/helpers"
)

type User struct {
	ID               int        `json:"id"`
	UserID           string     `json:"-"` // supertokens ID
	FirstName        string     `json:"first_name" validate:"required"`
	LastName         string     `json:"last_name" validate:"required"`
	Email            string     `json:"email" validate:"email"`
	Role             string     `json:"role" validate:"role"`
	OrganizationName string     `json:"organization_name" validate:"required"`
	Address          string     `json:"address"`
	Phone            string     `json:"phone" validate:"phone"`
	Enabled          bool       `json:"enabled"`
	CreatedAt        time.Time  `json:"created_at"`
	DisabledAt       *time.Time `json:"disabled_at"`
}

func (u *User) FromModel(mod *models.User) {
	u.ID = mod.ID
	u.UserID = mod.UserID
	u.FirstName = mod.FirstName
	u.LastName = mod.LastName
	u.Email = mod.Email
	u.Role = mod.Role
	u.OrganizationName = mod.OrganizationName
	u.Address = mod.Address.String
	u.Phone = mod.Phone.String
	u.Enabled = mod.Enabled
	u.CreatedAt = mod.CreatedAt
	u.DisabledAt = mod.DisabledAt.Ptr()
}

func (u User) ToModel() *models.User {
	return &models.User{
		UserID:           u.UserID,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Email:            u.Email,
		Role:             u.Role,
		OrganizationName: u.OrganizationName,
		Address:          helpers.GetNullString(u.Address),
		Phone:            helpers.GetNullString(u.Phone),
		Enabled:          true,
		PasswordChanged:  false,
		CreatedAt:        time.Now(),
	}
}
