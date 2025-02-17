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
				value: "",
			},
			want: false,
		},
		{
			name: "hasStringData",
			args: args{
				value: "data",
			},
			want: true,
		},
		{
			name: "hasDifferentData",
			args: args{
				value: 5,
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
