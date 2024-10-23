package supertokens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_isPasswordNeedsToBeChanged(t *testing.T) {
	t.Parallel()

	now := time.Now()

	type args struct {
		passwordChanged bool
		createdAt       time.Time
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "passwordNotChanged_temporaryPasswordStillValid",
			args: args{
				passwordChanged: false,
				createdAt:       now.AddDate(0, 0, -15),
			},
			want: false,
		},
		{
			name: "passwordNotChanged_temporaryPasswordNotValid",
			args: args{
				passwordChanged: false,
				createdAt:       now.AddDate(0, 0, -31),
			},
			want: true,
		},
		{
			name: "passwordWasChanged",
			args: args{
				passwordChanged: true,
				createdAt:       now.AddDate(0, 0, -31),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isPasswordNeedsToBeChanged(tt.args.passwordChanged, tt.args.createdAt)
			assert.Equal(t, tt.want, got)
		})
	}
}
