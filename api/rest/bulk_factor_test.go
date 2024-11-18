package rest

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBulkFactorForAutomation(t *testing.T) {
	endpoint := "/test/bulk/factor"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			requestBody: `[{"publisher":"publisher1","domain":"domain1","device":"desktop","factor":1.23,"country":"us"},{"publisher":"publisher2","domain":"domain2","device":"mobile","factor":3,"country":"il"},{"publisher":"publisher2","domain":"domain1","device":"mobile","factor":3,"country":"il"},{"publisher":"publisher1","domain":"domain1","device":"mobile","factor":10,"country":"uk"}]`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"factor bulk update successfully processed"}`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			resp, err := http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			defer resp.Body.Close()
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
