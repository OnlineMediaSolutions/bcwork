package bccron

import (
	"time"
)

func Next(cron string) int {
	now := time.Now().Truncate(time.Second)
	diff := MustParse(cron).Next(now).Sub(now).Seconds()
	if diff < 0 {
		return 0
	}
	return int(diff)
}
