package rest

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGlobalFactorBulkPostHandler_InvalidJSON(t *testing.T) {
	endpoint := "/global/factor/bulk"

	invalidJSON := `{"key": "consultant_fee", "publisher_id": "id", "value": 5`

	req, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(invalidJSON))
	assert.NoError(t, err)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, `{"status":"error","message":"error parsing request body for global factor bulk update","error":"unexpected end of JSON input"}`, string(body))
}
