package core

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_prepareUserDataForUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		newData     *dto.User
		currentData *models.User
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "updateAllFields",
			args: args{
				newData: &dto.User{
					FirstName:        "Joe",
					LastName:         "Smith",
					OrganizationName: "OMS",
					Address:          "Israel",
					Phone:            "+1111",
					Role:             "user",
					Types:            []string{"type_1"},
					Enabled:          false,
				},
				currentData: &models.User{
					FirstName:        "Ivan",
					LastName:         "Ivanov",
					OrganizationName: "Google",
					Address:          null.StringFrom("USA"),
					Phone:            null.StringFrom("+2222"),
					Role:             "admin",
					Types:            []string{"type_1", "type_2"},
					Enabled:          true,
				},
			},
			want: []string{
				models.UserColumns.FirstName,
				models.UserColumns.LastName,
				models.UserColumns.OrganizationName,
				models.UserColumns.Address,
				models.UserColumns.Phone,
				models.UserColumns.Role,
				models.UserColumns.Types,
				models.UserColumns.Enabled,
				models.UserColumns.DisabledAt,
			},
		},
		{
			name: "updatePartialFields",
			args: args{
				newData: &dto.User{
					FirstName:        "Joe",
					LastName:         "Smith",
					OrganizationName: "OMS",
					Address:          "Israel",
					Phone:            "+1111",
					Role:             "user",
					Enabled:          true,
				},
				currentData: &models.User{
					FirstName:        "Joe",
					LastName:         "Smith",
					OrganizationName: "Google",
					Address:          null.StringFrom("USA"),
					Phone:            null.StringFrom("+2222"),
					Role:             "admin",
					Enabled:          true,
				},
			},
			want: []string{
				models.UserColumns.OrganizationName,
				models.UserColumns.Address,
				models.UserColumns.Phone,
				models.UserColumns.Role,
			},
		},
		{
			name: "nothingToUpdate",
			args: args{
				newData: &dto.User{
					FirstName:        "Joe",
					LastName:         "Smith",
					OrganizationName: "OMS",
					Address:          "Israel",
					Phone:            "+1111",
					Role:             "user",
					Types:            []string{"type_1", "type_2"},
					Enabled:          false,
				},
				currentData: &models.User{
					FirstName:        "Joe",
					LastName:         "Smith",
					OrganizationName: "OMS",
					Address:          null.StringFrom("Israel"),
					Phone:            null.StringFrom("+1111"),
					Role:             "user",
					Types:            []string{"type_1", "type_2"},
					Enabled:          false,
				},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := prepareUserDataForUpdate(tt.args.newData, tt.args.currentData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
