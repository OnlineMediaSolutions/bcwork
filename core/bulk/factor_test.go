package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestPrepareFactorsData(t *testing.T) {
	tests := []struct {
		name     string
		input    []constant.FactorUpdateRequest
		expected []models.Factor
	}{
		{
			name: "valid input",
			input: []constant.FactorUpdateRequest{
				{
					Publisher:     "Publisher1",
					Domain:        "example.com",
					Country:       "US",
					Device:        "Mobile",
					Factor:        1.0,
					Browser:       "Chrome",
					OS:            "Android",
					PlacementType: "Banner",
				},
				{
					Publisher:     "Publisher2",
					Domain:        "example.org",
					Country:       "",
					Device:        "Desktop",
					Factor:        0.5,
					Browser:       "Firefox",
					OS:            "",
					PlacementType: "Interstitial",
				},
			},
			expected: []models.Factor{
				{
					Publisher:     "Publisher1",
					Domain:        "example.com",
					Country:       null.NewString("US", true),
					Device:        null.NewString("Mobile", true),
					Factor:        1.0,
					Browser:       null.NewString("Chrome", true),
					Os:            null.NewString("Android", true),
					PlacementType: null.NewString("Banner", true),
					RuleID:        "21a3ef12-ce56-5b4a-b2d3-1d3f25b20aba",
				},
				{
					Publisher:     "Publisher2",
					Domain:        "example.org",
					Country:       null.NewString("", false),
					Device:        null.NewString("Desktop", true),
					Factor:        0.5,
					Browser:       null.NewString("Firefox", true),
					Os:            null.NewString("", false),
					PlacementType: null.NewString("Interstitial", true),
					RuleID:        "7088da9b-adaa-5c51-a2ac-d8d3703bff89",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := prepareFactorsData(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
