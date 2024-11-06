package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUserGetHandler(t *testing.T) {
	endpoint := "/test/user/get"

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
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
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

func TestUserGetInfoHandler(t *testing.T) {
	endpoint := "/test/user/info"

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name    string
		query   string
		want    want
		wantErr bool
	}{
		{
			name: "validRequest",
			query: func() string {
				mod, _ := models.Users(models.UserWhere.Email.EQ("user_1@oms.com")).One(context.Background(), bcdb.DB())
				return "?id=" + mod.UserID
			}(),
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}`,
			},
		},
		{
			name:  "nothingFound",
			query: "?id=abcd",
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to retrieve user info","error":"failed to get user email: user not found in supertokens"}`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodGet, baseURL+endpoint+tt.query, nil)
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
	endpoint := "/test/user/set"

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
			requestBody: `{"first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Member"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"user successfully created"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Member"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for User. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			// based on results of "validRequest"
			name:        "duplicateUser",
			requestBody: `{"first_name": "John","last_name": "Doe","email": "user_3@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Member"}`,
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
	endpoint := "/test/user/update"

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
			requestBody: `{"id": 2, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Admin", "enabled": false}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"user successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{id": 2, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Admin", "enabled": false}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for User. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "noUserFoundToUpdate",
			requestBody: `{"id": 100, "first_name": "John","last_name": "Doe","email": "user_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Admin", "enabled": false}`,
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
	endpoint := "/test/user/verify/get"

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
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
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
				payload := `{"formFields": [{"id": "email","value": "user_1@oms.com"},{"id": "password","value": "abcd1234"}]}`
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
	endpoint := "/auth/signin"

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
			requestBody: `{"formFields": [{"id": "email","value": "user_1@oms.com"},{"id": "password","value": "abcd1234"}]}`,
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
				response:   "{\"status\": \"SIGN_IN_UP_NOT_ALLOWED\"}\n",
			},
		},
		{
			name:        "userNeedToChangeTemporaryPassword",
			requestBody: `{"formFields": [{"id": "email","value": "user_temp@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   "{\"status\": \"TEMPORARY_PASSWORD_NEEDS_TO_BE_CHANGED\"}\n",
			},
		},
		{
			name:        "userNotFound",
			requestBody: `{"formFields": [{"id": "email","value": "user_100@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   "{\"status\": \"SIGN_IN_UP_NOT_ALLOWED\"}\n",
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
	endpoint := "/test/user/verify/admin/get"

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
			signInRequestBody: `{"formFields": [{"id": "email","value": "user_admin@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
			},
		},
		{
			name:              "developerUser",
			requestBody:       `{"filter": {"email": ["user_1@oms.com"]}}`,
			signInRequestBody: `{"formFields": [{"id": "email","value": "user_developer@oms.com"},{"id": "password","value": "abcd1234"}]}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
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
	endpoint := "/test/user/verify/get"

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
	assert.Equal(t, "{\"status\": \"TEMPORARY_PASSWORD_NEEDS_TO_BE_CHANGED\"}\n", string(signInBody))

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
		`[{"id":1,"first_name":"name_1","last_name":"surname_1","email":"user_1@oms.com","role":"Member","organization_name":"OMS","address":"Israel","phone":"+972559999999","enabled":true,"created_at":"2024-09-01T13:46:41.302Z","disabled_at":null}]`,
		string(getUsersBody),
	)
}

func TestUserUpdate_History(t *testing.T) {
	endpoint := "/user/update"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
		history    dto.History
	}

	tests := []struct {
		name               string
		requestBody        string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "noChanges",
			requestBody:        `{"id": 7, "first_name": "name_history","last_name": "surname_history","email": "user_history@oms.com","organization_name": "Apple","address": "USA","phone": "+66666666666","role": "Member", "enabled": true}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["User"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: false,
			},
		},
		{
			name:               "validRequest",
			requestBody:        `{"id": 7, "first_name": "name_history","last_name": "surname_history","email": "user_history@oms.com","organization_name": "Apple","address": "USA","phone": "+66666666666","role": "Admin", "enabled": true}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["User"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Updated",
					Subject:      "User",
					Item:         "name_history surname_history",
					Changes: []dto.Changes{
						{
							Property: "role",
							OldValue: "Member",
							NewValue: "Admin",
						},
					},
				},
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
			req.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			_, err = http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			historyReq, err := http.NewRequest(fiber.MethodPost, baseURL+historyEndpoint, strings.NewReader(tt.historyRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			historyReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			historyReq.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			historyResp, err := http.DefaultClient.Do(historyReq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, historyResp.StatusCode)

			body, err := io.ReadAll(historyResp.Body)
			assert.NoError(t, err)
			defer historyResp.Body.Close()

			var (
				got   []dto.History
				found bool
			)
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			if !tt.want.hasHistory {
				assert.Equal(t, []dto.History{}, got)
				return
			}

			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}
				if reflect.DeepEqual(tt.want.history, got[i]) {
					found = true
				}
			}

			assert.Equal(t, true, found)
		})
	}
}

func TestUserSet_History(t *testing.T) {
	endpoint := "/user/set"
	historyEndpoint := "/history/get"

	type want struct {
		statusCode int
		hasHistory bool
		history    dto.History
	}

	tests := []struct {
		name               string
		requestBody        string
		historyRequestBody string
		want               want
		wantErr            bool
	}{
		{
			name:               "validRequest",
			requestBody:        `{"first_name": "History_2","last_name": "History_2","email": "user_history_2@oms.com","organization_name": "OMS","address": "Israel","phone": "+972559999999","role": "Member"}`,
			historyRequestBody: `{"filter": {"user_id": [-1],"subject": ["User"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				hasHistory: true,
				history: dto.History{
					UserID:       -1,
					UserFullName: "Internal Worker",
					Action:       "Created",
					Subject:      "User",
					Item:         "History_2 History_2",
				},
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
			req.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			_, err = http.DefaultClient.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			historyReq, err := http.NewRequest(fiber.MethodPost, baseURL+historyEndpoint, strings.NewReader(tt.historyRequestBody))
			if err != nil {
				t.Fatal(err)
			}
			historyReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			historyReq.Header.Set(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

			historyResp, err := http.DefaultClient.Do(historyReq)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, historyResp.StatusCode)

			body, err := io.ReadAll(historyResp.Body)
			assert.NoError(t, err)
			defer historyResp.Body.Close()

			var (
				got   []dto.History
				found bool
			)
			err = json.Unmarshal(body, &got)
			assert.NoError(t, err)
			if !tt.want.hasHistory {
				assert.Equal(t, []dto.History{}, got)
				return
			}

			for i := range got {
				got[i].ID = 0
				got[i].Date = time.Time{}
				for j := range got[i].Changes {
					got[i].Changes[j].ID = ""
				}
				if reflect.DeepEqual(tt.want.history, got[i]) {
					found = true
				}
			}

			assert.Equal(t, true, found)
		})
	}
}
