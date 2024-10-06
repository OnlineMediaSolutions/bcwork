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
				response:   `[{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":["linux"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"daily_cap":0,"status":"active"}]`,
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
			requestBody: `{"publisher":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":["linux"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":1,"status":"active"}`,
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
			requestBody: `{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","ru"],"device_type":["mobile","desktop"],"browser":["firefox","chrome"],"os":["linux","windows"],"kv":{"key_1":"value_1","key_4":"value_4","key_5":"value_5"},"price_model":"CPM","value":1,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to create targeting","error":"could not create targeting: there is targeting with such parameters"}`,
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
		requestBody string
		want        want
		wantErr     bool
	}{
		{
			name:        "validRequest",
			requestBody: `{"publisher":"22222222","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":["linux"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusOK,
				response:   `{"status":"success","message":"targeting successfully updated"}`,
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
			name:        "noTargetingFoundToUpdate",
			requestBody: `{"publisher":"33333333","domain":"2.com","unit_size":"300X250","placement_type":"top","country":["il","us"],"device_type":["mobile"],"browser":["firefox"],"os":["linux"],"kv":{"key_1":"value_1","key_2":"value_2","key_3":"value_3"},"price_model":"CPM","value":2,"status":"active"}`,
			want: want{
				statusCode: fiber.StatusInternalServerError,
				response:   `{"status":"error","message":"failed to update targeting","error":"no targeting found to update"}`,
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

func createTargetingTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("create table targeting " +
		"(" +
		"id serial primary key," +
		"hash varchar(36) not null," +
		"rule_id varchar(36) not null," +
		"publisher varchar(64)," +
		"domain varchar(256)," +
		"unit_size varchar(64)," +
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
		`(id, hash, rule_id, publisher, "domain", unit_size, placement_type, country, device_type, browser, os, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(9, '3ec92550-e7b1-55ce-a175-5822a129632f', '2bc8edcf-c4ab-5fe4-809c-ef54b167372c', '1111111', '3.com', '300X250', 'top', '{ru,us}', '{mobile}', '{firefox}', '{linux}', '{"key_1": "value_1", "key_2": "value_2", "key_3": "value_3"}'::jsonb, '', 0.0, '2024-10-01 13:46:41.302', '2024-10-01 13:46:41.302', 'active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, hash, rule_id, publisher, "domain", unit_size, placement_type, country, device_type, browser, os, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(10, '7e41d20c-f624-511c-88b8-baa28281a303', '029af331-ed85-5284-a551-fed9f8f8f63a', '22222222', '2.com', '300X250', 'top', '{il,us}', '{mobile}', '{firefox}', '{linux}', '{"key_1": "value_1", "key_2": "value_2", "key_3": "value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 'active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, hash, rule_id, publisher, "domain", unit_size, placement_type, country, device_type, browser, os, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(11, '6f7ed004-7791-50ae-847b-6e61194f9669', '454a6636-60c8-5f09-903a-a6924cbbad3d', '1111111', '2.com', '300X250', 'top', '{ru}', '{mobile}', '{firefox}', '{linux}', '{"key_1": "value_1", "key_2": "value_2", "key_3": "value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:57:05.542', '2024-10-01 13:57:05.542', 'active');`)
	tx.MustExec("CREATE TABLE metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.Commit()
}
