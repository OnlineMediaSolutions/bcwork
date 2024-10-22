package rest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestUserSetHandler(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/set"

	userService := core.NewUserService(supertokenClient, false)
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
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	createUserTables(bcdb.DB())
	err := createUsersInSupertokens(supertokenClient)
	if err != nil {
		t.Fatalf("can't init users in supertokens: %v", err)
	}

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
			requestBody: `{"first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "user"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"user successfully created"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "user"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for User. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			// based on results of "validRequest"
			name:        "duplicateUser",
			requestBody: `{"first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "user"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to create user","error":"failed to create user in supertoken: error creating user in supertoken: status [FIELD_ERROR] not equal 'OK'"}`,
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

// TODO: tests with using cookies from previous response
// TODO: test password reset

func createUserTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(
		`CREATE TABLE public."user" (` +
			`id serial primary key,` +
			`user_id varchar(256) not null,` +
			`email varchar(256) unique not null,` +
			`first_name varchar(256) not null,` +
			`last_name varchar(256) not null,` +
			`role varchar(64) not null,` +
			`organization_name varchar(128) not null,` +
			`address varchar(128),` +
			`phone varchar(32),` +
			`enabled bool not null default true,` +
			`password_changed bool not null default false,` +
			`reset_token varchar(256),` +
			`created_at timestamp not null,` +
			`disabled_at timestamp` +
			`)`,
	)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES(9, 'user_1@oms.com', 'name_1', 'surname_1', 'user', 'OMS', 'Israel', '+972559999999', TRUE, '2024-10-01 13:46:41.302');`)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES(10, 'user_2@oms.com', 'name_2', 'surname_2', 'admin', 'Google', 'USA', '+11111111', TRUE, '2024-10-01 13:46:41.302');`)
	tx.Commit()
}

func createUsersInSupertokens(client supertokens_module.TokenManagementSystem) error {
	ctx := context.Background()

	_, err := client.CreateUser(ctx, "user_1@oms.com", "abcd1234")
	if err != nil {
		return err
	}

	_, err = client.CreateUser(ctx, "user_2@oms.com", "abcd1234")
	if err != nil {
		return err
	}

	return nil
}
