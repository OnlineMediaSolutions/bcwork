package filter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dateFilterTestCase struct {
	Input       []byte
	DatesFilter DatesFilter
	Message     string
}

func TestParse(t *testing.T) {
	testCases := []dateFilterTestCase{
		{
			Input: []byte(`{"from":"2019-11-25","to":"2019-11-27"}`),
			DatesFilter: DatesFilter{
				From: time.Date(2019, 11, 25, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2019, 11, 27, 0, 0, 0, 0, time.UTC),
			},
			Message: "date only",
		},
		{
			Input: []byte(`{"from":"2019-11-25 10:00","to":"2019-11-27 13:10"}`),
			DatesFilter: DatesFilter{
				From: time.Date(2019, 11, 25, 10, 0, 0, 0, time.UTC),
				To:   time.Date(2019, 11, 27, 13, 10, 0, 0, time.UTC),
			},
			Message: "date with hours",
		},
		{
			Input: []byte(`{"from":"2019-11-25T10:00","to":"2019-11-27T13:10"}`),
			DatesFilter: DatesFilter{
				From: time.Date(2019, 11, 25, 10, 0, 0, 0, time.UTC),
				To:   time.Date(2019, 11, 27, 13, 10, 0, 0, time.UTC),
			},
			Message: "date with hours with T seperator",
		},
	}

	for _, testCase := range testCases {
		var input DatesFilter
		err := json.Unmarshal(testCase.Input, &input)
		assert.NoError(t, err, "failed to unmarshal input for date filter tests")
		assert.True(t, input.Equal(&testCase.DatesFilter), "unexpected date filter")
	}

}
