package config

import (
	"fmt"
	"time"
)

type InterfaceMap map[string]interface{}

//Get numeric value with
func (c InterfaceMap) GetIntValue(key string) (int, bool, error) {
	if val, found := c[key]; found {
		valInt, ok := val.(int)
		if !ok {
			return 0, false, fmt.Errorf("failed to convert %s configuration to numeric value (value is not an int)", key)
		}

		return valInt, true, nil
	}

	return 0, false, nil
}

func (c InterfaceMap) GetIntValueWithDefault(key string, def int) (int, error) {
	if val, found := c[key]; found {
		valInt, ok := val.(int)
		if !ok {
			return 0, fmt.Errorf("failed to convert %s configuration to numeric value (value is not an int)", key)
		}

		return valInt, nil
	}

	return def, nil
}

//Get date value
func (c InterfaceMap) GetDateValue(key string) (time.Time, bool, error) {
	if val, found := c[key]; found {
		valDate, ok := val.(time.Time)
		if !ok {
			return time.Now(), false, fmt.Errorf("failed to convert %s configuration to time value (value is not a time)", key)
		}

		return valDate, true, nil
	}

	return time.Now(), false, nil
}

//Get date value
func (c InterfaceMap) GetDateValueWithDefault(key string, def time.Time) (time.Time, error) {
	if val, found := c[key]; found {
		valDate, ok := val.(time.Time)
		if !ok {
			return time.Now(), fmt.Errorf("failed to convert %s configuration to time value (value is not a time)", key)
		}

		return valDate, nil
	}

	return def, nil
}

//Get string  value
func (c InterfaceMap) GetStringValue(key string) (string, bool, error) {
	if val, found := c[key]; found {
		valString, ok := val.(string)
		if !ok {
			return "", false, fmt.Errorf("failed to convert %s configuration to string value (value is not a string)", key)
		}

		return valString, true, nil
	}

	return "", false, nil
}

//Get string value
func (c InterfaceMap) GetStringValueWithDefault(key string, def string) (string, error) {
	if val, found := c[key]; found {
		valString, ok := val.(string)
		if !ok {
			return def, fmt.Errorf("failed to convert %s configuration to string value (value is not a string)", key)
		}

		return valString, nil
	}

	return def, nil
}

//Get string  value
func (c InterfaceMap) GetBoolValue(key string) (bool, bool, error) {
	if val, found := c[key]; found {
		valBool, ok := val.(bool)
		if !ok {
			return false, false, fmt.Errorf("failed to convert %s configuration to boolean value (value is not a boolean)", key)
		}

		return valBool, true, nil
	}

	return false, false, nil
}

//Get string value
func (c InterfaceMap) GetBoolValueWithDefault(key string, def bool) (bool, error) {
	if val, found := c[key]; found {
		valBool, ok := val.(bool)
		if !ok {
			return def, fmt.Errorf("failed to convert %s configuration to boolean value (value is not a boolean)", key)
		}

		return valBool, nil
	}

	return def, nil
}
