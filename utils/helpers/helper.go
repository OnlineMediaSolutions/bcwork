package helpers

import (
	"github.com/rs/zerolog/log"
	"math"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"
)

func GetStringWithDefaultValue(str, defaultValue string) string {
	if str == "" {
		return defaultValue
	}
	return str
}

func GetStringFromSliceWithDefaultValue(elems []string, sep, defaultValue string) string {
	if len(elems) == 0 {
		return defaultValue
	}

	return "(" + strings.Join(elems, sep) + ")"
}

func GetNullString(s string) null.String {
	if s == "" {
		return null.NewString("", false)
	}
	return null.StringFrom(s)
}

func RoundFloat(value float64) float64 {
	return math.Round(value*100) / 100
}

func FormatDate(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse timestamp")
		return ""
	}
	return t.Format("2006-01-02")
}
