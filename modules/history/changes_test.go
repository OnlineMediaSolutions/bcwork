package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValueContainsData(t *testing.T) {
	t.Parallel()

	type args struct {
		value any
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "noData",
			args: args{
				value: nil,
			},
			want: false,
		},
		{
			name: "dataEqualEmptyString",
			args: args{
				value: func() any {
					s := ""
					return s
				}(),
			},
			want: false,
		},
		{
			name: "hasStringData",
			args: args{
				value: func() any {
					s := "data"
					return s
				}(),
			},
			want: true,
		},
		{
			name: "hasDifferentData",
			args: args{
				value: func() any {
					s := 5
					return s
				}(),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isValueContainsData(tt.args.value)
			assert.Equal(t, tt.want, got)
		})
	}
}
