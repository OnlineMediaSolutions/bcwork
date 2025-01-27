package core

import (
	"encoding/json"
	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
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

		var ips []ipPriceDate
		err = json.Unmarshal(result, &ips)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(ips))
		assert.Equal(t, "192.168.1.1", ips[0].IP)
		assert.Equal(t, 100.50, ips[0].Price)
		assert.Equal(t, "10.0.0.1", ips[1].IP)
		assert.Equal(t, 200.75, ips[1].Price)
	})

	t.Run("Empty Input", func(t *testing.T) {
		// Prepare empty input
		data := &dto.PriceOverrideRequest{
			Ips: []dto.Ips{},
		}

		result, err := buildPriceOvverideValue(data)

		assert.NoError(t, err)

		var ips []ipPriceDate
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

		var ips []ipPriceDate
		err = json.Unmarshal(result, &ips)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(ips))
		assert.Equal(t, "255.255.255.255", ips[0].IP)
		assert.Equal(t, 0.0, ips[0].Price)
		assert.Equal(t, "localhost", ips[1].IP)
		assert.Equal(t, 300.99, ips[1].Price)
	})
}
