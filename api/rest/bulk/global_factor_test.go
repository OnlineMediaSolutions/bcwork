package bulk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/api/rest"
	"github.com/stretchr/testify/require"
)

func TestGlobalFactorBulkPostHandler_InvalidJSON(t *testing.T) {
	const url = "/global/factor/bulk"

	app := fiber.New()
	app.Post(url, GlobalFactorBulkPostHandler)

	invalidJSON := `{"key": "consultant_fee", "publisher_id": "id", "value": 5`

	req := httptest.NewRequest("POST", url, bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	fmt.Print(string(body))

	var response rest.Response
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)
	fmt.Print(response.Status)
	fmt.Print(response.Message)
	require.Equal(t, "error", response.Status)
	require.Contains(t, response.Message, "error parsing request body for global factor bulk update")
}
