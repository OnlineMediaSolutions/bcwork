package rest

import "testing"

func TestCreateKeyForQuery(t *testing.T) {

	// Test cases
	testCases := []struct {
		name          string
		request       *BlockGetRequest
		expectedQuery string
	}{
		{
			name: "Both Publisher and Domain Provided",
			request: &BlockGetRequest{
				Publisher: "example_publisher",
				Domain:    "example_domain",
			},
			expectedQuery: " and metadata_queue.key = 'example_publisher:example_domain'",
		},
		{
			name: "Only Publisher Provided",
			request: &BlockGetRequest{
				Publisher: "example_publisher",
				Domain:    "",
			},
			expectedQuery: " and last.key = 'example_publisher'",
		},
		{
			name: "Neither Publisher nor Domain Provided",
			request: &BlockGetRequest{
				Publisher: "",
				Domain:    "",
			},
			expectedQuery: " and 1=1 ",
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualQuery := createKeyForQuery(tc.request)
			if actualQuery != tc.expectedQuery {
				t.Errorf("Expected query '%s', but got '%s'", tc.expectedQuery, actualQuery)
			}
		})
	}
}
