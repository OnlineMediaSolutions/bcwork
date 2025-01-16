package helpers

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func Test_getQuarter(t *testing.T) {
	t.Parallel()

	type args struct {
		month time.Month
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Q1_January",
			args: args{
				month: time.January,
			},
			want: "Q1",
		},
		{
			name: "Q1_February",
			args: args{
				month: time.February,
			},
			want: "Q1",
		},
		{
			name: "Q1_March",
			args: args{
				month: time.March,
			},
			want: "Q1",
		},
		{
			name: "Q2_April",
			args: args{
				month: time.April,
			},
			want: "Q2",
		},
		{
			name: "Q2_May",
			args: args{
				month: time.May,
			},
			want: "Q2",
		},
		{
			name: "Q2_June",
			args: args{
				month: time.June,
			},
			want: "Q2",
		},
		{
			name: "Q3_July",
			args: args{
				month: time.July,
			},
			want: "Q3",
		},
		{
			name: "Q3_August",
			args: args{
				month: time.August,
			},
			want: "Q3",
		},
		{
			name: "Q3_September",
			args: args{
				month: time.September,
			},
			want: "Q3",
		},
		{
			name: "Q4_October",
			args: args{
				month: time.October,
			},
			want: "Q4",
		},
		{
			name: "Q4_November",
			args: args{
				month: time.November,
			},
			want: "Q4",
		},
		{
			name: "Q4_December",
			args: args{
				month: time.December,
			},
			want: "Q4",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetQuarter(tt.args.month)
			assert.Equal(t, tt.want, got)
		})
	}
}
