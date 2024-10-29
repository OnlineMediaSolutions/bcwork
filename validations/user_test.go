package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateUser(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.User
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &dto.User{
					FirstName:        "name",
					LastName:         "surname",
					Email:            "email@email.com",
					OrganizationName: "organization",
					Role:             "Member",
				},
			},
			want: []string{},
		},
		{
			name: "badRole",
			args: args{
				request: &dto.User{
					FirstName:        "name",
					LastName:         "surname",
					Email:            "email@email.com",
					OrganizationName: "organization",
					Address:          "address",
					Phone:            "+972 (55) 999-99-99",
					Role:             "unknown_role",
				},
			},
			want: []string{
				roleValidationErrorMessage,
			},
		},
		{
			name: "badEmail",
			args: args{
				request: &dto.User{
					FirstName:        "name",
					LastName:         "surname",
					Email:            "emailemail.com",
					OrganizationName: "organization",
					Address:          "address",
					Phone:            "+972 (55) 999-99-99",
					Role:             "Member",
				},
			},
			want: []string{
				emailValidationErrorMessage,
			},
		},
		{
			name: "badPhone",
			args: args{
				request: &dto.User{
					FirstName:        "name",
					LastName:         "surname",
					Email:            "email@email.com",
					OrganizationName: "organization",
					Address:          "address",
					Phone:            "+972f2222222",
					Role:             "Member",
				},
			},
			want: []string{
				phoneValidationErrorMessage,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateUser(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
