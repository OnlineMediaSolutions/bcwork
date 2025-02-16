package omserr

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_extractContext(t *testing.T) {
	cases := []struct {
		Input    string
		Expected []string
		Msg      string
	}{
		{"empty", []string{}, "empty"},
		{"(1)sometext(2)", []string{"1", "2"}, "normal 2 values"},
		{"(1)sometext(2)  (3)sdfsdfs", []string{"1", "2", "3"}, "normal 3 values"},
		{"(1)sometext(2)  (3(nested))sdfsdfs", []string{"1", "2", "3(nested)"}, "nested 3 values"},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.Expected, extractContext(tc.Input), tc.Msg)
	}
}

func Test_extractContextValues(t *testing.T) {
	cases := []struct {
		Input    string
		Expected map[string]string
		Msg      string
	}{
		{"empty", map[string]string{}, "empty"},
		{"key:value", map[string]string{"key": "value"}, "one value"},
		{"key:value:value", map[string]string{"key": "value:value"}, "multiple vals"},
	}

	for _, tc := range cases {
		assert.True(t, reflect.DeepEqual(tc.Expected, extractContextValues(tc.Input)), tc.Msg)
	}
}

func Test_ErrorContext(t *testing.T) {
	cases := []struct {
		Input    error
		Expected ErrorContext
		Msg      string
	}{
		{errors.New("empty"), ErrorContext{Context: map[string]string{}}, "empty"},
		{errors.New("failed to load prediction(user_id:some-user,prediction_id:some_prediction,stat_id:bskt:nba:p:points)"), ErrorContext{Context: map[string]string{"user_id": "some-user", "prediction_id": "some_prediction", "stat_id": "bskt:nba:p:points"}}, "few vals"},
		{errors.Wrap(errors.New("failed to load prediction(user_id:some-user)"), "wrapping(user_id:some-user-wrapped)"), ErrorContext{Context: map[string]string{"user_id": "some-user-wrapped"}}, "wrapping errors"},
	}

	for _, tc := range cases {
		assert.True(t, reflect.DeepEqual(tc.Expected.Context, ExtractErrContext(tc.Input).Context), tc.Msg)
	}
}
