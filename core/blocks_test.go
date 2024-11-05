package core

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
)

func TestCreateKeyForQuery(t *testing.T) {
	testCases := []struct {
		name     string
		request  dto.BlockGetRequest
		expected string
	}{
		{
			name:     "Empty request",
			request:  dto.BlockGetRequest{},
			expected: " and 1=1 ",
		},
		{
			name: "Publisher and types provided",
			request: dto.BlockGetRequest{
				Types:     []string{"badv", "bcat"},
				Publisher: "publisher",
				Domain:    "domain",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher:domain') OR (metadata_queue.key = 'bcat:publisher:domain'))",
		},
		{
			name: "Publisher provided, no domain",
			request: dto.BlockGetRequest{
				Types:     []string{"badv", "cbat"},
				Publisher: "publisher",
			},
			expected: "AND ( (metadata_queue.key = 'badv:publisher') OR (metadata_queue.key = 'cbat:publisher'))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := createKeyForQuery(&tc.request)
			if result != tc.expected {
				t.Errorf("Test %s failed: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
