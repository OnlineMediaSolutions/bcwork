package bccron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Next(t *testing.T) {
	t.Parallel()

	type args struct {
		cron string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid_everyDay",
			args: args{
				cron: "0 0 * * *",
			},
			want: func() int {
				now := time.Now().Truncate(time.Second)
				t := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
				return int(t.Sub(now).Seconds())
			}(),
		},
		{
			name: "valid_everyHour",
			args: args{
				cron: "0 * * * *",
			},
			want: func() int {
				now := time.Now().Truncate(time.Second)
				t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
				return int(t.Sub(now).Seconds())
			}(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Next(tt.args.cron)
			assert.Equal(t, tt.want, got)
		})
	}
}
