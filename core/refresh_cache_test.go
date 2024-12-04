package core

import (
	"testing"
)

func TestCreateSoftDeleteQueryRefreshCache(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{"rule1", "rule2"},
			expected: "UPDATE refresh_cache SET active = false WHERE rule_id IN ('rule1','rule2');",
		},
		{
			input:    []string{"rule3"},
			expected: "UPDATE refresh_cache SET active = false WHERE rule_id IN ('rule3');",
		},
		{
			input:    []string{},
			expected: "UPDATE refresh_cache SET active = false WHERE rule_id IN ();",
		},
	}

	for _, test := range tests {
		result := createSoftDeleteQueryRefreshCache(test.input)
		if result != test.expected {
			t.Errorf("For input %v, expected %q but got %q", test.input, test.expected, result)
		}
	}
}
