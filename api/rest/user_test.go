package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
)

func TestUserSetHandler(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	createUserTables(bcdb.DB())

	endpoint := "/user/set"

	userService := core.NewUserService(supertokenClient)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					validations.ValidateUser, userManagementSystem.UserSetHandler,
				},
			},
		},
	})
	defer app.Shutdown()

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
			requestBody: `{"first_name": "John","last_name": "Doe","email": "johndoe@example.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "user"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully added"}`,
			},
		},
		// {
		// 	name:        "invalidRequest",
		// 	requestBody: `{"publisher_id: "22222222"}`,
		// 	want: want{
		// 		statusCode: fiber.StatusBadRequest,
		// 		response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
		// 	},
		// },
		// {
		// 	name:        "hasDuplicate",
		// 	requestBody: `{"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","ru"],"device_type":["mobile","desktop"],"browser":["firefox","chrome"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"Active"}`,
		// 	want: want{
		// 		statusCode: fiber.StatusBadRequest,
		// 		response:   `{"status":"error","message":"found duplicate while creating targeting","error":"checking for duplicates: found duplicate: there is targeting with such parameters","duplicate":{"id":10,"publisher_id":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":null,"created_at":"2024-10-01T13:51:28.407Z","updated_at":"2024-10-01T13:51:28.407Z","status":"Active"}}`,
		// 	},
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, endpoint, strings.NewReader(tt.requestBody))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			resp, err := app.Test(req, -1)
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

func createUserTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(
		`CREATE TABLE public."user" (` +
			`user_id varchar(36) NOT NULL,` +
			`created_at timestamp NOT NULL,` +
			`email varchar(1024) NOT NULL,` +
			`first_name varchar(256) NULL,` +
			`last_name varchar(256) NULL,` +
			`last_activity_at timestamp NULL,` +
			`invited_at timestamp NULL,` +
			`signedup_at timestamp NULL,` +
			`invited_by varchar(36) NULL,` +
			`"role" varchar(64) NOT NULL,` +
			`organization_name varchar(128) NOT NULL,` +
			`address varchar(128) NOT NULL,` +
			`phone varchar(32) NOT NULL,` +
			`enabled bool DEFAULT true NOT NULL,` +
			`disabled_at timestamp NULL,` +
			`id serial NOT NULL,` +
			`CONSTRAINT user_email_key UNIQUE (email),` +
			`CONSTRAINT user_pkey PRIMARY KEY (user_id)` +
			`)`,
	)
	tx.Commit()
}
