package dpo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDemandPartners_DesiredOutput(t *testing.T) {
	demandData := []byte(`[
		{
			"automation_name": "index-pbs",
			"api_name": "indexs2s",
			"threshold": 0.001
		},
		{
			"automation_name": "onetag-bcm",
			"api_name": "onetagbcm",
			"threshold": 0.001
		},
		{
			"automation_name": "pubmatic-pbs",
			"api_name": "pubmaticbcm",
			"threshold": 0.001
		},
		{
			"automation_name": "sovrn",
			"api_name": "sovrnbcm",
			"threshold": 0.001
		},
		{
			"automation_name": "yieldmo-audienciad",
			"api_name": "yieldmo",
			"threshold": 0.001
		}
	]`)

	var err error
	stringErrors := []string{}

	worker := &Worker{
		Demands: map[string]*DemandSetup{},
	}

	demands, err := worker.getDemandPartners(demandData, err, stringErrors)
	expected := map[string]*DemandSetup{
		"index-pbs": {
			Name:      "index-pbs",
			ApiName:   "indexs2s",
			Threshold: 0.001,
		},
		"onetag-bcm": {
			Name:      "onetag-bcm",
			ApiName:   "onetagbcm",
			Threshold: 0.001,
		},
		"pubmatic-pbs": {
			Name:      "pubmatic-pbs",
			ApiName:   "pubmaticbcm",
			Threshold: 0.001,
		},
		"sovrn": {
			Name:      "sovrn",
			ApiName:   "sovrnbcm",
			Threshold: 0.001,
		},
		"yieldmo-audienciad": {
			Name:      "yieldmo-audienciad",
			ApiName:   "yieldmo",
			Threshold: 0.001,
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, demands)
	assert.Empty(t, stringErrors)
}
