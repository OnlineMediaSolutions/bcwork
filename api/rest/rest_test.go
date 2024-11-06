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
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/history"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/storage/cache"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/ory/dockertest"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	pool                 *dockertest.Pool
	appTest              *fiber.App
	supertokenClientTest supertokens_module.TokenManagementSystem
	omsNPTest            *OMSNewPlatform

	port    = ":9000"
	baseURL = "http://localhost" + port
)

func TestMain(m *testing.M) {
	viper.SetDefault(config.AWSWorkerAPIKeyKey, "aws_worker_api_key")
	viper.SetDefault(config.CronWorkerAPIKeyKey, "cron_worker_api_key")

	pool = testutils.SetupDockerTestPool()
	pg := testutils.SetupDB(pool)

	err := bcdb.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	var st *dockertest.Resource
	st, supertokenClientTest = testutils.SetupSuperTokens(pool)

	createDBTables(bcdb.DB(), supertokenClientTest)

	cache := cache.NewInMemoryCache()
	historyModule := history.NewHistoryClient(cache)

	omsNPTest = NewOMSNewPlatform(supertokenClientTest, historyModule, false)
	verifySessionMiddleware := adaptor.HTTPMiddleware(supertokenClientTest.VerifySession)

	appTest = fiber.New()
	appTest.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	appTest.Use(LoggingMiddleware)
	appTest.Use(historyModule.HistoryMiddleware)
	// floor
	appTest.Post("/test/floor", omsNPTest.FloorPostHandler)
	appTest.Post("/test/floor/get", omsNPTest.FloorGetAllHandler)
	// bulk
	appTest.Post("/test/global/factor/bulk", omsNPTest.GlobalFactorBulkPostHandler)
	// block
	appTest.Post("/test/block/get", omsNPTest.BlockGetAllHandler)
	// targeting
	appTest.Post("/test/targeting/get", omsNPTest.TargetingGetHandler)
	appTest.Post("/test/targeting/set", validations.ValidateTargeting, omsNPTest.TargetingSetHandler)
	appTest.Post("/test/targeting/update", validations.ValidateTargeting, omsNPTest.TargetingUpdateHandler)
	appTest.Post("/test/targeting/tags", omsNPTest.TargetingExportTagsHandler)
	// user
	appTest.Get("/test/user/info", omsNPTest.UserGetInfoHandler)
	appTest.Post("/test/user/get", omsNPTest.UserGetHandler)
	appTest.Post("/test/user/set", validations.ValidateUser, omsNPTest.UserSetHandler)
	appTest.Post("/test/user/update", validations.ValidateUser, omsNPTest.UserUpdateHandler)
	appTest.Post("/test/user/verify/get", verifySessionMiddleware, omsNPTest.UserGetHandler)
	appTest.Post("/test/user/verify/admin/get", verifySessionMiddleware, supertokenClientTest.AdminRoleRequired, omsNPTest.UserGetHandler)
	// history
	appTest.Post("/history/get", omsNPTest.HistoryGetHandler)
	// endpoint to test history saving
	appTest.Post("/bulk/global/factor", verifySessionMiddleware, omsNPTest.GlobalFactorBulkPostHandler)    // TODO: add test
	appTest.Post("/publisher/new", verifySessionMiddleware, omsNPTest.PublisherNewHandler)                 // TODO: add test
	appTest.Post("/publisher/update", verifySessionMiddleware, omsNPTest.PublisherUpdateHandler)           // TODO: add test
	appTest.Post("/floor", verifySessionMiddleware, omsNPTest.FloorPostHandler)                            // TODO: add test
	appTest.Post("/factor", verifySessionMiddleware, omsNPTest.FactorPostHandler)                          // TODO: add test
	appTest.Post("/global/factor", verifySessionMiddleware, omsNPTest.GlobalFactorPostHandler)             // TODO: add test
	appTest.Post("/dpo/set", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationSetHandler)       // TODO: add test
	appTest.Post("/dpo/delete", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationDeleteHandler) // TODO: add test
	appTest.Post("/dpo/update", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationUpdateHandler) // TODO: add test
	appTest.Post("/publisher/domain", verifySessionMiddleware, omsNPTest.PublisherDomainPostHandler)       // TODO: add test + ?automation=true
	appTest.Post("/targeting/set", verifySessionMiddleware, omsNPTest.TargetingSetHandler)
	appTest.Post("/targeting/update", verifySessionMiddleware, omsNPTest.TargetingUpdateHandler)
	appTest.Post("/user/update", verifySessionMiddleware, omsNPTest.UserUpdateHandler)
	appTest.Post("/user/set", verifySessionMiddleware, omsNPTest.UserSetHandler)
	appTest.Post("/block", verifySessionMiddleware, omsNPTest.BlockPostHandler)                // TODO: add test + ?automation=true
	appTest.Post("/pixalate", verifySessionMiddleware, omsNPTest.PixalatePostHandler)          // TODO: add test + ?automation=true
	appTest.Post("/pixalate/delete", verifySessionMiddleware, omsNPTest.PixalateDeleteHandler) // TODO: add test
	appTest.Post("/confiant", verifySessionMiddleware, omsNPTest.ConfiantPostHandler)          // TODO: add test + ?automation=true

	go appTest.Listen(port)

	code := m.Run()

	pool.Purge(pg)
	pool.Purge(st)
	appTest.Shutdown()

	os.Exit(code)
}

func createDBTables(db *sqlx.DB, client supertokens_module.TokenManagementSystem) {
	createUserTablesAndUsersInSupertokens(db, client)
	createTargetingTables(db)
	createBlockTables(db)
	createHistoryTable(db)
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
		`VALUES('` + user1.User.ID + `', 'user_1@oms.com', 'name_1', 'surname_1', 'Member', 'OMS', 'Israel', '+972559999999', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload2 := `{"email": "user_2@oms.com","password": "abcd1234"}`
	req2, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload2))
	resp2, _ := http.DefaultClient.Do(req2)
	data2, _ := io.ReadAll(resp2.Body)
	defer resp2.Body.Close()
	var user2 supertokens_module.CreateUserResponse
	json.Unmarshal(data2, &user2)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user2.User.ID + `', 'user_2@oms.com', 'name_2', 'surname_2', 'Admin', 'Google', 'USA', '+11111111', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload3 := `{"email": "user_temp@oms.com","password": "abcd1234"}`
	req3, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload3))
	resp3, _ := http.DefaultClient.Do(req3)
	data3, _ := io.ReadAll(resp3.Body)
	defer resp3.Body.Close()
	var user3 supertokens_module.CreateUserResponse
	json.Unmarshal(data3, &user3)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user3.User.ID + `', 'user_temp@oms.com', 'name_temp', 'surname_temp', 'Member', 'Google', 'USA', '+77777777777', TRUE, '2024-09-01 13:46:41.302');`)

	payload4 := `{"email": "user_disabled@oms.com","password": "abcd1234"}`
	req4, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload4))
	resp4, _ := http.DefaultClient.Do(req4)
	data4, _ := io.ReadAll(resp4.Body)
	defer resp4.Body.Close()
	var user4 supertokens_module.CreateUserResponse
	json.Unmarshal(data4, &user4)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user4.User.ID + `', 'user_disabled@oms.com', 'name_disabled', 'surname_disabled', 'Member', 'Google', 'USA', '+88888888888', FALSE, '2024-09-01 13:46:41.302');`)

	payload5 := `{"email": "user_admin@oms.com","password": "abcd1234"}`
	req5, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload5))
	resp5, _ := http.DefaultClient.Do(req5)
	data5, _ := io.ReadAll(resp5.Body)
	defer resp5.Body.Close()
	var user5 supertokens_module.CreateUserResponse
	json.Unmarshal(data5, &user5)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user5.User.ID + `', 'user_admin@oms.com', 'name_disabled', 'surname_disabled', 'Admin', 'Google', 'USA', '+88888888888', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload6 := `{"email": "user_developer@oms.com","password": "abcd1234"}`
	req6, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload6))
	resp6, _ := http.DefaultClient.Do(req6)
	data6, _ := io.ReadAll(resp6.Body)
	defer resp6.Body.Close()
	var user6 supertokens_module.CreateUserResponse
	json.Unmarshal(data6, &user6)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user6.User.ID + `', 'user_developer@oms.com', 'name_developer', 'surname_developer', 'Developer', 'Apple', 'USA', '+66666666666', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload7 := `{"email": "user_history@oms.com","password": "abcd1234"}`
	req7, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload7))
	resp7, _ := http.DefaultClient.Do(req7)
	data7, _ := io.ReadAll(resp7.Body)
	defer resp7.Body.Close()
	var user7 supertokens_module.CreateUserResponse
	json.Unmarshal(data7, &user7)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user7.User.ID + `', 'user_history@oms.com', 'name_history', 'surname_history', 'Member', 'Apple', 'USA', '+66666666666', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

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

func createHistoryTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("create table history" +
		"(" +
		"id serial primary key," +
		"user_id int not null," +
		"subject varchar(64) not null," +
		"item text not null," +
		"publisher_id varchar(64)," +
		"domain varchar(64)," +
		"entity_id varchar(64)," +
		"action varchar(64) not null," +
		"old_value jsonb," +
		"new_value jsonb," +
		"changes jsonb," +
		"date timestamp not null" +
		");",
	)
	tx.MustExec(`INSERT INTO public.history ` +
		`(id, user_id, subject, item, "action", old_value, new_value, changes, "date", publisher_id, "domain", entity_id) ` +
		`VALUES(198, 26, 'JS Targeting', '[de us]_120x600_[desktop]_[windows]_[safari]_{leaderboard true}_CPM_5.02_{0 false}_[]', 'Updated', '{"id": 18, "kv": null, "os": ["windows"], "value": 5.01, "domain": "finkiel.com", "status": "Active", "browser": ["safari"], "country": ["de", "us"], "daily_cap": null, "unit_size": "120x600", "created_at": "2024-11-03T15:15:56.996018Z", "updated_at": "2024-11-04T07:31:37.65011Z", "device_type": ["desktop"], "price_model": "CPM", "publisher_id": "1111111", "placement_type": "leaderboard"}'::jsonb, '{"id": 18, "kv": null, "os": ["windows"], "value": 5.02, "domain": "finkiel.com", "status": "Active", "browser": ["safari"], "country": ["de", "us"], "daily_cap": null, "unit_size": "120x600", "created_at": "2024-11-03T15:15:56.996018Z", "updated_at": "2024-11-04T08:05:50.319973Z", "device_type": ["desktop"], "price_model": "CPM", "publisher_id": "1111111", "placement_type": "leaderboard"}'::jsonb, '[{"property": "value", "new_value": 5.02, "old_value": 5.01}]'::jsonb, '2024-11-04 08:05:50.333', '1111111', 'finkiel.com', '18');`)
	tx.MustExec(`INSERT INTO public.history ` +
		`(id, user_id, subject, item, "action", old_value, new_value, changes, "date", publisher_id, "domain", entity_id) ` +
		`VALUES(199, 26, 'Domain', 'finkiel.com (1111111)', 'Updated', '{"domain": "finkiel.com", "automation": false, "created_at": "2024-09-17T09:16:51.949465Z", "gpp_target": 20, "updated_at": "2024-10-31T10:38:28.877837Z", "publisher_id": "1111111"}'::jsonb, '{"domain": "finkiel.com", "automation": true, "created_at": "2024-09-17T09:16:51.949465Z", "gpp_target": 20, "updated_at": "2024-11-04T08:20:18.80533Z", "publisher_id": "1111111"}'::jsonb, '[{"property": "automation", "new_value": true, "old_value": false}]'::jsonb, '2024-11-04 08:20:18.812', '1111111', 'finkiel.com', NULL);`)
	tx.MustExec(`INSERT INTO public.history ` +
		`(id, user_id, subject, item, "action", old_value, new_value, changes, "date", publisher_id, "domain", entity_id) ` +
		`VALUES(200, 26, 'Domain', 'finkiel.com (1111111)', 'Updated', '{"domain": "finkiel.com", "automation": true, "created_at": "2024-09-17T09:16:51.949465Z", "gpp_target": 20, "updated_at": "2024-11-04T08:20:18.80533Z", "publisher_id": "1111111"}'::jsonb, '{"domain": "finkiel.com", "automation": false, "created_at": "2024-09-17T09:16:51.949465Z", "gpp_target": 20, "updated_at": "2024-11-04T08:20:32.079207Z", "publisher_id": "1111111"}'::jsonb, '[{"property": "automation", "new_value": false, "old_value": true}]'::jsonb, '2024-11-04 08:20:32.085', '1111111', 'finkiel.com', NULL);`)
	tx.MustExec(`INSERT INTO public.history ` +
		`(id, user_id, subject, item, "action", old_value, new_value, changes, "date", publisher_id, "domain", entity_id) ` +
		`VALUES(201, -1, 'Factor Automation', 'online-image-editor1.com (1111111)', 'Updated', '{"domain": "online-image-editor1.com", "automation": false, "created_at": "2024-10-31T14:36:25.731208Z", "gpp_target": 20, "updated_at": "2024-11-04T08:21:15.787426Z", "publisher_id": "1111111"}'::jsonb, '{"domain": "online-image-editor1.com", "automation": false, "created_at": "2024-10-31T14:36:25.731208Z", "gpp_target": 25, "updated_at": "2024-11-04T08:21:27.254635Z", "publisher_id": "1111111"}'::jsonb, '[{"property": "gpp_target", "new_value": 25, "old_value": 20}]'::jsonb, '2024-11-04 08:21:27.262', '1111111', 'online-image-editor1.com', NULL);`)
	tx.Commit()
}
