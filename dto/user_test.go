package dto

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
)

func Test_UsersByTypes_Append(t *testing.T) {
	t.Parallel()

	type args struct {
		mod *models.User
	}

	tests := []struct {
		name  string
		users *UsersByTypes
		args  args
		want  *UsersByTypes
	}{
		{
			name: "valid",
			users: &UsersByTypes{
				AM: []*userByType{
					{ID: 1, Fullname: "Name Surname"},
				},
			},
			args: args{
				mod: &models.User{
					ID: 2, FirstName: "First", LastName: "Last", Types: []string{UserTypeMediaBuyer},
				},
			},
			want: &UsersByTypes{
				AM: []*userByType{
					{ID: 1, Fullname: "Name Surname"},
				},
				MB: []*userByType{
					{ID: 2, Fullname: "First Last"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.users.Append(tt.args.mod)
			// object mutates
			assert.Equal(t, tt.want, tt.users)
		})
	}
}
