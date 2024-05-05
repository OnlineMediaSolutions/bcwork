package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type StringMap map[string]string

// Get numeric value with
func (c StringMap) GetIntValue(key string) (int, bool, error) {
	if val, found := c[key]; found {
		if val == "-" || val == "" {
			return 0, false, nil
		}
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return 0, false, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to numeric value (%s is not numeric)", key, val))
		}

		return valInt, true, nil
	}

	return 0, false, nil
}

func (c StringMap) GetIntValueWithDefault(key string, def int) (int, error) {
	if val, found := c[key]; found {
		if val == "-" || val == "" {
			return def, nil
		}
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return def, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to numeric value (%s is not numeric)", key, val))
		}

		return valInt, nil
	}

	return def, nil
}

// Get date value
func (c StringMap) GetDateValue(key string) (time.Time, bool, error) {
	if val, found := c[key]; found {
		if val == "-" || val == "" {
			return time.Time{}, false, nil
		}
		valDate, err := time.Parse("2006-01-02", val)
		if err != nil {
			return time.Now(), false, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to date value (%s is not date in the format YYYY-MM-DD)", key, val))
		}

		return valDate, true, nil
	}

	return time.Now(), false, nil
}

// Get date value
func (c StringMap) GetDateValueWithDefault(key string, def time.Time) (time.Time, error) {
	if val, found := c[key]; found {
		valDate, err := time.Parse("2006-01-02", val)
		if err != nil {
			return def, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to date value (%s is not date in the format YYYY-MM-DD)", key, val))
		}

		return valDate, nil
	}

	return def, nil
}

// Get date value
func (c StringMap) GetDateHourValue(key string) (time.Time, bool, error) {
	if val, found := c[key]; found {
		if val == "-" || val == "" {
			return time.Time{}, false, nil
		}
		valDate, err := time.Parse("2006-01-02-15", val)
		if err != nil {
			return time.Now(), false, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to date hour value (%s is not date in the format YYYY-MM-DD)", key, val))
		}

		return valDate, true, nil
	}

	return time.Now(), false, nil
}

// Get date value
func (c StringMap) GetDateHourValueWithDefault(key string, def time.Time) (time.Time, error) {
	if val, found := c[key]; found {
		valDate, err := time.Parse("2006-01-02-15", val)
		if err != nil {
			return def, errors.Wrap(err, fmt.Sprintf("failed to convert %s configuration to date hour value (%s is not date in the format YYYY-MM-DD)", key, val))
		}

		return valDate, nil
	}

	return def, nil
}

// Get string  value
func (c StringMap) GetStringValue(key string) (string, bool) {
	if val, found := c[key]; found {
		return val, true
	}

	return "", false
}

// Get string value
func (c StringMap) GetStringValueWithDefault(key string, def string) string {
	if val, found := c[key]; found {
		return val
	}

	return def
}

// Get string  value
func (c StringMap) GetStringSlice(key string, sep string) ([]string, bool) {
	if val, found := c[key]; found {
		return strings.Split(val, sep), true
	}

	return []string{}, false
}

// Get string  value
func (c StringMap) GetBoolValue(key string) (bool, bool) {
	if val, found := c[key]; found {
		if strings.ToLower(val) == "true" || strings.ToLower(val) == "on" {
			return true, true
		}

		return false, true
	}

	return false, false
}

// Get string value
func (c StringMap) GetBoolValueWithDefault(key string, def bool) bool {
	if val, found := c[key]; found {
		if strings.ToLower(val) == "true" || strings.ToLower(val) == "on" || val == "1" {
			return true
		} else if strings.ToLower(val) == "false" || strings.ToLower(val) == "off" || val == "0" {
			return false
		}

		return def
	}

	return def
}

// Get duration  value
func (c StringMap) GetDurationValue(key string) (time.Duration, error) {
	if val, found := c[key]; found {
		d, err := time.ParseDuration(val)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to parse duration")
		}
		return d, nil
	}

	return 0, nil
}

// Get duration  value
func (c StringMap) GetDurationValueWithDefault(key string, def time.Duration) (time.Duration, error) {
	if val, found := c[key]; found {
		d, err := time.ParseDuration(val)
		if err != nil {
			return def, errors.Wrapf(err, "failed to parse duration")
		}
		return d, nil
	}

	return def, nil
}

func FromQuery(values map[string][]string) StringMap {
	res := StringMap{}

	for k, v := range values {
		if len(v) > 0 {
			res[k] = v[0]
		}
	}

	return res
}
