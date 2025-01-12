package core

import (
	"github.com/m6yf/bcwork/dto"
	"testing"
)

func TestBuildResultMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []dto.Result
		expected map[string]map[string]dto.ActivityStatus
	}{
		{
			name: "Test with mixed PubImps values",
			input: []dto.Result{
				{Domain: "finkiel.com", PublisherId: 111, PubImps: 6000},
				{Domain: "finkiel.com", PublisherId: 222, PubImps: 100},
				{Domain: "finkiel.com", PublisherId: 333, PubImps: 10},
				{Domain: "anotherDomain.com", PublisherId: 111, PubImps: 5000},
			},
			expected: map[string]map[string]dto.ActivityStatus{
				"finkiel.com": {
					"111": dto.ActivityStatus(2), // PubImps >= 5000 ACTIVE
					"222": dto.ActivityStatus(1), // PubImps >= 20 and < 5000 LOW
					"333": dto.ActivityStatus(0), // PubImps < 20 PAUSED
				},
				"anotherDomain.com": {
					"111": dto.ActivityStatus(2), // PubImps >= 5000 ACTIVE
				},
			},
		},
		{
			name:     "All domain are paused",
			input:    []dto.Result{},
			expected: map[string]map[string]dto.ActivityStatus{},
		},
		{
			name: "Test with all PubImps less than 20",
			input: []dto.Result{
				{Domain: "finkiel.com", PublisherId: 111, PubImps: 10},
				{Domain: "finkiel.com", PublisherId: 222, PubImps: 5},
			},
			expected: map[string]map[string]dto.ActivityStatus{
				"finkiel.com": {
					"111": dto.ActivityStatus(0), // PAUSED
					"222": dto.ActivityStatus(0), // PAUSED
				},
			},
		},
		{
			name:     "No Input",
			input:    []dto.Result{},
			expected: map[string]map[string]dto.ActivityStatus{},
		},
		{
			name:     "No input to method",
			input:    nil,
			expected: map[string]map[string]dto.ActivityStatus{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildResultMap(tt.input)

			for domain, publishers := range tt.expected {
				for pubId, expectedStatus := range publishers {
					actualStatus, exists := result[domain][pubId]
					if !exists {
						t.Errorf("Expected to find domain %s, publisher %s, but it was not found", domain, pubId)
					}
					if actualStatus != expectedStatus {
						t.Errorf("For domain %s, publisher %s, expected %d but got %d", domain, pubId, expectedStatus, actualStatus)
					}
				}
			}
		})
	}
}
