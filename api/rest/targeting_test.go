package rest

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/stretchr/testify/assert"
)

func TestTargetingGetHandler(t *testing.T) {
	endpoint := "/targeting/get"

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					TargetingGetHandler,
				},
			},
		},
	})
	defer app.Shutdown()

	db, pool, pg := testutils.SetupDB(t)
	defer func() {
		db.Close()
		pool.Purge(pg)
	}()

	createTargetingTables(db)

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
			requestBody: `{"filter": {"publisher": ["22222222"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[{"id":10,"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":null,"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":0,"status":"active"}]`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"filter": {"publisher: ["22222222"]}}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"status":"error","message":"failed to parse request for getting targeting data","error":"invalid character '2' after object key"}`,
			},
		},
		{
			name:        "nothingFound",
			requestBody: `{"filter": {"publisher": ["xxxxxxxx"]}}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `[]`,
			},
		},
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

func TestTargetingSetHandler(t *testing.T) {
	endpoint := "/targeting/set"

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					validations.ValidateTargeting, TargetingSetHandler,
				},
			},
		},
	})
	defer app.Shutdown()

	db, pool, pg := testutils.SetupDB(t)
	defer func() {
		db.Close()
		pool.Purge(pg)
	}()

	createTargetingTables(db)

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
			requestBody: `{"publisher":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully added"}`,
			},
		},
		{
			name:        "invalidRequest",
			requestBody: `{"publisher: "22222222"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "hasDuplicate",
			requestBody: `{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","ru"],"device_type":["mobile","desktop"],"browser":["firefox","chrome"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to create targeting","error":"checking for duplicates: there is same targeting (id=10) with such parameters [country=[il us],device_type=[mobile],browser=[firefox],os=[],kv={\"key_1\": \"value_1\", \"key_2\": \"value_2\", \"key_3\": \"value_3\"}]"}`,
			},
		},
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

func TestTargetingUpdateHandler(t *testing.T) {
	endpoint := "/targeting/update"

	app := testutils.SetupApp(&testutils.AppSetup{
		Endpoints: []testutils.EndpointSetup{
			{
				Method: fiber.MethodPost,
				Path:   endpoint,
				Handlers: []fiber.Handler{
					validations.ValidateTargeting, TargetingUpdateHandler,
				},
			},
		},
	})
	defer app.Shutdown()

	db, pool, pg := testutils.SetupDB(t)
	defer func() {
		db.Close()
		pool.Purge(pg)
	}()

	createTargetingTables(db)

	type want struct {
		statusCode int
		response   string
	}

	tests := []struct {
		name        string
		endpoint    string
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			endpoint:    endpoint + "?id=10",
			requestBody: `{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully updated"}`,
			},
		},
		{
			name:        "invalidRequest",
			endpoint:    endpoint,
			requestBody: `{"publisher: "22222222"}`,
			want: want{
				statusCode: fiber.StatusBadRequest,
				response:   `{"message":"Invalid request body for Targeting. Please ensure it's a valid JSON.","status":"error"}`,
			},
		},
		{
			name:        "noTargetingFoundToUpdate",
			endpoint:    endpoint + "?id=12",
			requestBody: `{"publisher":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"failed to get targeting with id [12] to update: sql: no rows in result set"}`,
			},
		},
		{
			// based on results of "validRequest"
			name:        "nothingToUpdate",
			endpoint:    endpoint + "?id=10",
			requestBody: `{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"there are no new values to update targeting"}`,
			},
		},
		{
			name:        "duplicateConflictOnUpdatedEntity",
			endpoint:    endpoint + "?id=11",
			requestBody: `{"publisher":"1111111","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["us"],"device_type":["mobile"],"browser":["firefox"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"error checking for duplicates: there is same targeting (id=9) with such parameters [country=[ru us],device_type=[mobile],browser=[firefox],os=[],kv={\"key_1\": \"value_1\", \"key_2\": \"value_2\", \"key_3\": \"value_3\"}]"}`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(fiber.MethodPost, tt.endpoint, strings.NewReader(tt.requestBody))
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

func createTargetingTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("create table targeting " +
		"(" +
		"id serial primary key," +
		"publisher varchar(64) not null," +
		"domain varchar(256) not null," +
		"unit_size varchar(64) not null," +
		"placement_type varchar(64)," +
		"country text[]," +
		"device_type text[]," +
		"browser text[]," +
		"os text[]," +
		"kv jsonb," +
		"price_model varchar(64) not null," +
		"value float8 not null," +
		"daily_cap int," +
		"created_at timestamp not null," +
		"updated_at timestamp," +
		"status  varchar(64) not null" +
		")",
	)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(9, '1111111', '2.com', '300X250', 'top', '{ru,us}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, '', 0.0, '2024-10-01 13:46:41.302', '2024-10-01 13:46:41.302', 'active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(10, '22222222', '2.com', '300X250', 'top', '{il,us}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 'active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(11, '1111111', '2.com', '300X250', 'top', '{ru}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:57:05.542', '2024-10-01 13:57:05.542', 'active');`)
	tx.MustExec("CREATE TABLE metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value jsonb,commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.Commit()
}
