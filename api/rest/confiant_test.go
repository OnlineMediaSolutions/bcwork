package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"testing"
)

func TestApiTest(t *testing.T) {
	c := ConfiantGetAllHandler
	print(c)
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		// First test case
		{
			description:  "get HTTP status 200",
			route:        "/test",
			expectedCode: 200,
		},
	}
	log.Println("IM here")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Debug entpoint")
	})
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, -1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

}
