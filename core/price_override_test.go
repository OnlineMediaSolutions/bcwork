package core

import (
	"encoding/json"
	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/types"
	"testing"
)

func TestBuildPriceOverrideValue(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		data := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{
				{IP: "192.168.1.1", Price: 100.50},
				{IP: "10.0.0.1", Price: 200.75},
			},
		}

		result, err := buildPriceOvverideValue(data)

		assert.NoError(t, err)

		var ips []dto.Ips
		err = json.Unmarshal(result, &ips)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(ips))
		assert.Equal(t, "192.168.1.1", ips[0].IP)
		assert.Equal(t, 100.50, ips[0].Price)
		assert.Equal(t, "10.0.0.1", ips[1].IP)
		assert.Equal(t, 200.75, ips[1].Price)
	})

	t.Run("Empty Input", func(t *testing.T) {
		data := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{},
		}

		result, err := buildPriceOvverideValue(data)

		assert.NoError(t, err)

		var ips []dto.Ips
		err = json.Unmarshal(result, &ips)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(ips))
	})

	t.Run("Special Characters in IP", func(t *testing.T) {
		data := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{
				{IP: "255.255.255.255", Price: 0.0},
				{IP: "localhost", Price: 300.99},
			},
		}

		result, err := buildPriceOvverideValue(data)
		assert.NoError(t, err)

		var ips []dto.Ips
		err = json.Unmarshal(result, &ips)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(ips))
		assert.Equal(t, "255.255.255.255", ips[0].IP)
		assert.Equal(t, 0.0, ips[0].Price)
		assert.Equal(t, "localhost", ips[1].IP)
		assert.Equal(t, 300.99, ips[1].Price)
	})
}

func TestAddNewIpToValue(t *testing.T) {
	t.Run("Success Case", func(t *testing.T) {
		// Prepare initial JSON value and request data
		initialValue := `[{"ip":"192.168.0.1","date":"2025-01-01T12:00:00Z","price":100}]`
		value := types.JSON(initialValue)

		requestData := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{
				{IP: "10.0.0.1", Price: 200.50},
				{IP: "172.16.0.1", Price: 300.75},
			},
		}

		result, err := addNewIpToValue(value, requestData)

		assert.NoError(t, err)

		var resultData []dto.Ips
		err = json.Unmarshal(result, &resultData)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(resultData))
		assert.Equal(t, "192.168.0.1", resultData[0].IP)
		assert.Equal(t, "10.0.0.1", resultData[1].IP)
		assert.Equal(t, 200.50, resultData[1].Price)
		assert.Equal(t, "172.16.0.1", resultData[2].IP)
		assert.Equal(t, 300.75, resultData[2].Price)
	})

	t.Run("Unmarshal Error", func(t *testing.T) {

		invalidValue := `invalid json`
		value := types.JSON(invalidValue)

		requestData := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{
				{IP: "10.0.0.1", Price: 200.50},
			},
		}

		result, err := addNewIpToValue(value, requestData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to unmarshal metadata value")
		assert.Nil(t, result)
	})
}
