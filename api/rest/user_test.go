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
	"github.com/m6yf/bcwork/models"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestUserGetHandler(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, true)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/get"

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					userManagementSystem.UserGetHandler,
				},
			},
		},
	})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
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
			requestBody: `{"filter": {"email": ["user_1@oms.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"user","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{filter": {"email": ["user_1@oms.com"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting users data","error":"invalid character 'f' looking for beginning of object key string"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"email": ["user_4@oms.com"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[]`,
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

func TestUserSetHandler(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, true)
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

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
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

func TestUserUpdateHandler(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, true)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/update"

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					validations.ValidateUser, userManagementSystem.UserUpdateHandler,
				},
			},
		},
	})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
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
			requestBody: `{"id": 2, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "admin", "enabled": false}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"user successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{id": 2, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "admin", "enabled": false}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for User. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "noUserFoundToUpdate",
			requestBody: `{"id": 100, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "admin", "enabled": false}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update user","error":"failed to get user with id [100] to update: sql: no rows in result set"}`,
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

func TestVerifySession(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, false)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/get"

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					adaptor.HTTPMiddleware(supertokenClient.VerifySession), userManagementSystem.UserGetHandler,
				},
			},
		},
	})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
	}

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name         string
		requestBody  string
		needToSignIn bool
		want         want
		wantErr      bool
	}{
		{
			name:         "unauthorized",
			requestBody:  `{"filter": {"email": ["user_1@oms.com"]}}`,
			needToSignIn: false,
			want: want{
				statusCode: fiber.StatusUnauthorized,
				response:   `{"error": "unauthorized"}`,
			},
		},
		{
			name:         "validRequest",
			requestBody:  `{"filter": {"email": ["user_1@oms.com"]}}`,
			needToSignIn: true,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"user","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
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

			if tt.needToSignIn {
				payload := `{"formFields": [{"id": "email","value": "user_2@oms.com"},{"id": "password","value": "abcd1234"}]}`
				signInReq, err := http.NewRequest(fiber.MethodPost, baseURL+"/auth/signin", strings.NewReader(payload))
				if err != nil {
					t.Fatal(err)
				}
				signInReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

				signInResp, err := http.DefaultClient.Do(signInReq)
				if tt.wantErr {
					t.Fatal(err)
				}

				req.Header.Set(fiber.HeaderCookie, "sAccessToken="+signInResp.Header.Get("St-Access-Token"))
			}

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

func TestSignIn(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, false)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/auth/signin"

	app := testutils.SetupApp(&testutils.AppSetup{})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
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
			name:        "valid",
			requestBody: `{"formFields": [{"id": "email","value": "user_2@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `"status":"OK"`,
			},
		},
		{
			name:        "userDisabled",
			requestBody: `{"formFields": [{"id": "email","value": "user_disabled@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   "user disabled\n",
			},
		},
		{
			name:        "userNeedToChangeTemporaryPassword",
			requestBody: `{"formFields": [{"id": "email","value": "user_temp@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   "temporary password needs to be changed\n",
			},
		},
		{
			name:        "userNotFound",
			requestBody: `{"formFields": [{"id": "email","value": "user_100@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   "provided email not allowed to sign in/up using third-party providers\n",
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
			assert.Contains(t, string(body), tt.want.response) // to surpass dynamic user_id and time creation
		})
	}
}

func TestAdminRoleRequired(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, false)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/get"

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					adaptor.HTTPMiddleware(supertokenClient.VerifySession),
					supertokenClient.AdminRoleRequired,
					userManagementSystem.UserGetHandler,
				},
			},
		},
	})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	if err != nil {
		t.Fatalf("can't create tables and init users in supertokens: %v", err)
	}

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name              string
		requestBody       string
		signInRequestBody string
		want              want
		wantErr           bool
	}{
		{
			name:              "notAdminUser",
			requestBody:       `{"filter": {"email": ["user_1@oms.com"]}}`,
			signInRequestBody: `{"formFields": [{"id": "email","value": "user_1@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusForbidden,
				response:   `{"status":"error","message":"admin role required","error":"current user doesn't have admin role"}`,
			},
		},
		{
			name:              "adminUser",
			requestBody:       `{"filter": {"email": ["user_1@oms.com"]}}`,
			signInRequestBody: `{"formFields": [{"id": "email","value": "user_2@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"user","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
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

			signInReq, err := http.NewRequest(fiber.MethodPost, baseURL+"/auth/signin", strings.NewReader(tt.signInRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			signInReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			signInResp, err := http.DefaultClient.Do(signInReq)
			if tt.wantErr {
				t.Fatal(err)
			}

			req.Header.Set(fiber.HeaderCookie, "sAccessToken="+signInResp.Header.Get("St-Access-Token"))

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

func TestResetTemporaryPasswordFlow(t *testing.T) {
	pool := testutils.SetupDockerTestPool(t)
	pg := testutils.SetupDB(t, pool)
	st, supertokenClient := testutils.SetupSuperTokens(t, pool, false)
	defer func() {
		pool.Purge(pg)
		pool.Purge(st)
	}()

	port := ":8000"
	baseURL := "http://localhost" + port
	endpoint := "/user/get"

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem := NewUserManagementSystem(userService)

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					adaptor.HTTPMiddleware(supertokenClient.VerifySession),
					userManagementSystem.UserGetHandler,
				},
			},
		},
	})
	defer app.Shutdown()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	go app.Listen(port)

	err := createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	assert.NoError(t, err)

	// trying to sign in
	signInPayload := `{"formFields": [{"id": "email","value": "user_temp@oms.com"},{"id": "password","value": "abcd1234"}]}`
	signInReq, err := http.NewRequest(fiber.MethodPost, baseURL+"/auth/signin", strings.NewReader(signInPayload))
	assert.NoError(t, err)
	signInReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	// need to change temporary password
	signInResp, err := http.DefaultClient.Do(signInReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, signInResp.StatusCode)
	signInBody, err := io.ReadAll(signInResp.Body)
	assert.NoError(t, err)
	defer signInResp.Body.Close()
	assert.Equal(t, "temporary password needs to be changed\n", string(signInBody))

	// getting token for temporary password reset
	getTokenPayload := `{"formFields":[{"id":"email","value":"user_temp@oms.com"}]}`
	getTokenReq, err := http.NewRequest(fiber.MethodPost, baseURL+"/auth/user/password/reset/token", strings.NewReader(getTokenPayload))
	assert.NoError(t, err)
	getTokenReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getTokenResp, err := http.DefaultClient.Do(getTokenReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, getTokenResp.StatusCode)
	user, err := models.Users(models.UserWhere.Email.EQ("user_temp@oms.com")).One(context.Background(), bcdb.DB())
	assert.NoError(t, err)

	// reseting password
	resetPasswordPayload := `{"formFields":[{"id":"password","value":"abcd1234"}],"token":"` + user.ResetToken.String + `","method":"token"}`
	resetPasswordReq, err := http.NewRequest(fiber.MethodPost, baseURL+"/auth/user/password/reset", strings.NewReader(resetPasswordPayload))
	assert.NoError(t, err)
	resetPasswordReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resetPasswordResp, err := http.DefaultClient.Do(resetPasswordReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resetPasswordResp.StatusCode)

	// trying to sign in after password changing
	signInReq, err = http.NewRequest(fiber.MethodPost, baseURL+"/auth/signin", strings.NewReader(signInPayload))
	assert.NoError(t, err)
	signInResp, err = http.DefaultClient.Do(signInReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, signInResp.StatusCode)

	// doing get users request
	getUsersPayload := `{"filter": {"email": ["user_1@oms.com"]}}`
	getUsersReq, err := http.NewRequest(fiber.MethodPost, baseURL+endpoint, strings.NewReader(getUsersPayload))
	assert.NoError(t, err)
	getUsersReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	getUsersReq.Header.Set(fiber.HeaderCookie, "sAccessToken="+signInResp.Header.Get("St-Access-Token"))

	getUsersResp, err := http.DefaultClient.Do(getUsersReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, getUsersResp.StatusCode)
	getUsersBody, err := io.ReadAll(getUsersResp.Body)
	assert.NoError(t, err)
	defer getUsersResp.Body.Close()
	assert.Equal(
		t,
		`[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"user","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
		string(getUsersBody),
	)
}

func createUserTablesAndUsersInSupertokens(db *sqlx.DB, client supertokens_module.TokenManagementSystem) error {
	ctx := context.Background()
	tx := db.MustBeginTx(ctx, nil)
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

	id1, err := client.CreateUser(ctx, "user_1@oms.com", "abcd1234")
	if err != nil {
		return err
	}
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + id1 + `', 'user_1@oms.com', 'name_1', 'surname_1', 'user', 'OMS', 'Israel', '+972559999999', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	id2, err := client.CreateUser(ctx, "user_2@oms.com", "abcd1234")
	if err != nil {
		return err
	}
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + id2 + `', 'user_2@oms.com', 'name_2', 'surname_2', 'admin', 'Google', 'USA', '+11111111', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	id3, err := client.CreateUser(ctx, "user_temp@oms.com", "abcd1234")
	if err != nil {
		return err
	}
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + id3 + `', 'user_temp@oms.com', 'name_temp', 'surname_temp', 'user', 'Google', 'USA', '+77777777777', TRUE, '2024-09-01 13:46:41.302');`)

	id4, err := client.CreateUser(ctx, "user_disabled@oms.com", "abcd1234")
	if err != nil {
		return err
	}
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + id4 + `', 'user_disabled@oms.com', 'name_disabled', 'surname_disabled', 'user', 'Google', 'USA', '+88888888888', FALSE, '2024-09-01 13:46:41.302');`)

	tx.Commit()

	return nil
}
