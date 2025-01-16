package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Column_GetBooleanReplacementValue(t *testing.T) {
	t.Parallel()

	type args struct {
		isTrue bool
	}

	tests := []struct {
		name   string
		column *Column
		args   args
		want   string
	}{
		{
			name: "valid_true",
			column: &Column{
				BooleanReplacement: &BooleanReplacement{
					True:  "is_bool",
					False: "is_not_bool",
				},
			},
			args: args{
				isTrue: true,
			},
			want: "is_bool",
		},
		{
			name: "valid_false",
			column: &Column{
				BooleanReplacement: &BooleanReplacement{
					True:  "is_bool",
					False: "is_not_bool",
				},
			},
			args: args{
				isTrue: false,
			},
			want: "is_not_bool",
		},
		{
			name:   "noBooleanReplacement",
			column: &Column{},
			args: args{
				isTrue: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.column.GetBooleanReplacementValue(tt.args.isTrue)
			assert.Equal(t, tt.want, got)
		})
	}
}
