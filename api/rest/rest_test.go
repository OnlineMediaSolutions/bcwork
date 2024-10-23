package rest

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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
	"github.com/ory/dockertest"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	pool                 *dockertest.Pool
	app                  *fiber.App
	supertokenClient     supertokens_module.TokenManagementSystem
	userManagementSystem *UserManagementSystem

	port    = ":9000"
	baseURL = "http://localhost" + port
)

func TestMain(m *testing.M) {
	pool = testutils.SetupDockerTestPool()
	pg := testutils.SetupDB(pool)

	err := bcdb.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	var st *dockertest.Resource
	st, supertokenClient = testutils.SetupSuperTokens(pool, false)

	createUserTablesAndUsersInSupertokens(bcdb.DB(), supertokenClient)
	createTargetingTables(bcdb.DB())
	createBlockTables(bcdb.DB())

	userService := core.NewUserService(supertokenClient, false)
	userManagementSystem = NewUserManagementSystem(userService)

	app = fiber.New()
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	// block
	app.Post("/block/get", BlockGetAllHandler)
	// targeting
	app.Post("/targeting/get", TargetingGetHandler)
	app.Post("/targeting/set", validations.ValidateTargeting, TargetingSetHandler)
	app.Post("/targeting/update", validations.ValidateTargeting, TargetingUpdateHandler)
	app.Post("/targeting/tags", TargetingExportTagsHandler)
	// user
	app.Post("/user/get", userManagementSystem.UserGetHandler)
	app.Post("/user/set", validations.ValidateUser, userManagementSystem.UserSetHandler)
	app.Post("/user/update", validations.ValidateUser, userManagementSystem.UserUpdateHandler)
	app.Post("/user/verify/get", adaptor.HTTPMiddleware(supertokenClient.VerifySession), userManagementSystem.UserGetHandler)
	app.Post("/user/verify/admin/get", adaptor.HTTPMiddleware(supertokenClient.VerifySession), supertokenClient.AdminRoleRequired, userManagementSystem.UserGetHandler)
	go app.Listen(port)

	code := m.Run()

	pool.Purge(pg)
	pool.Purge(st)
	app.Shutdown()

	os.Exit(code)
}

func createUserTablesAndUsersInSupertokens(db *sqlx.DB, client supertokens_module.TokenManagementSystem) {
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

	payload1 := `{"email": "user_1@oms.com","password": "abcd1234"}`
	req1, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload1))
	resp1, _ := http.DefaultClient.Do(req1)
	data1, _ := io.ReadAll(resp1.Body)
	defer resp1.Body.Close()
	var user1 supertokens_module.CreateUserResponse
	json.Unmarshal(data1, &user1)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user1.User.ID + `', 'user_1@oms.com', 'name_1', 'surname_1', 'user', 'OMS', 'Israel', '+972559999999', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload2 := `{"email": "user_2@oms.com","password": "abcd1234"}`
	req2, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload2))
	resp2, _ := http.DefaultClient.Do(req2)
	data2, _ := io.ReadAll(resp2.Body)
	defer resp2.Body.Close()
	var user2 supertokens_module.CreateUserResponse
	json.Unmarshal(data2, &user2)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user2.User.ID + `', 'user_2@oms.com', 'name_2', 'surname_2', 'admin', 'Google', 'USA', '+11111111', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload3 := `{"email": "user_temp@oms.com","password": "abcd1234"}`
	req3, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload3))
	resp3, _ := http.DefaultClient.Do(req3)
	data3, _ := io.ReadAll(resp3.Body)
	defer resp3.Body.Close()
	var user3 supertokens_module.CreateUserResponse
	json.Unmarshal(data3, &user3)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user3.User.ID + `', 'user_temp@oms.com', 'name_temp', 'surname_temp', 'user', 'Google', 'USA', '+77777777777', TRUE, '2024-09-01 13:46:41.302');`)

	payload4 := `{"email": "user_disabled@oms.com","password": "abcd1234"}`
	req4, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload4))
	resp4, _ := http.DefaultClient.Do(req4)
	data4, _ := io.ReadAll(resp4.Body)
	defer resp4.Body.Close()
	var user4 supertokens_module.CreateUserResponse
	json.Unmarshal(data4, &user4)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user4.User.ID + `', 'user_disabled@oms.com', 'name_disabled', 'surname_disabled', 'user', 'Google', 'USA', '+88888888888', FALSE, '2024-09-01 13:46:41.302');`)

	payload5 := `{"email": "user_admin@oms.com","password": "abcd1234"}`
	req5, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload5))
	resp5, _ := http.DefaultClient.Do(req5)
	data5, _ := io.ReadAll(resp5.Body)
	defer resp5.Body.Close()
	var user5 supertokens_module.CreateUserResponse
	json.Unmarshal(data5, &user5)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user5.User.ID + `', 'user_admin@oms.com', 'name_disabled', 'surname_disabled', 'admin', 'Google', 'USA', '+88888888888', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	tx.Commit()
}

func createTargetingTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("create table IF NOT EXISTS publisher " +
		"(" +
		"publisher_id varchar(64) primary key," +
		"name varchar(64) not null" +
		")",
	)
	tx.MustExec(`INSERT INTO public.publisher ` +
		`(publisher_id, name)` +
		`VALUES('1111111', 'publisher_1'),('22222222', 'publisher_2'),('333', 'publisher_3');`)
	tx.MustExec("create table targeting " +
		"(" +
		"id serial primary key," +
		"publisher_id varchar(64) not null references publisher(publisher_id)," +
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
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(9, '1111111', '2.com', '300X250', 'top', '{ru,us}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, '', 0.0, '2024-10-01 13:46:41.302', '2024-10-01 13:46:41.302', 'Active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(10, '22222222', '2.com', '300X250', 'top', '{il,us}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 'Active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(11, '1111111', '2.com', '300X250', 'top', '{ru}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 1.0, '2024-10-01 13:57:05.542', '2024-10-01 13:57:05.542', 'Active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status, daily_cap)` +
		`VALUES(20, '333', '2.com', '300X250', 'top', '{ru,us}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, '', 0.0, '2024-10-01 13:46:41.302', '2024-10-01 13:46:41.302', 'Active', 1000);`)
	tx.MustExec("CREATE TABLE IF NOT EXISTS metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value jsonb,commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.Commit()
}

func createBlockTables(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("CREATE TABLE IF NOT EXISTS metadata_queue (transaction_id varchar(36), key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"f2b8833e-e0e4-57e0-a68b-6792e337ab4d", "badv:20223:realgm.com", nil, "[\"safesysdefender.xyz\"]", 0, "2024-09-20T10:10:10.100", "2024-09-26T10:10:10.100")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"c53c4dd2-6f68-5b62-b613-999a5239ad36", "badv:20356:playpilot.com", nil, "[\"fraction-content.com\"]", 0, "2024-09-20T10:10:10.100", "2024-09-26T10:10:10.100")
	tx.Commit()
}
