package bulk

import (
	"bytes"
	"encoding/json"
	"github.com/m6yf/bcwork/api/rest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestFactorBulkPostHandler_InvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/factor/bulk", FactorBulkPostHandler)

	invalidJSON := `{"publisher": "publisher1", "domain": "domain1", "device": "desktop", "factor": 1.23, "country": "US"`

	req := httptest.NewRequest("POST", "/factor/bulk", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var response rest.Response
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "error", response.Status)
	require.Contains(t, response.Message, "error when parsing request body for bulk update")
}
