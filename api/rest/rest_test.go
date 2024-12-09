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
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/m6yf/bcwork/validations"
	"github.com/ory/dockertest"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
	viper.SetDefault(config.APIChunkSizeKey, 100)

	boil.DebugMode = true

	pool = testutils.SetupDockerTestPool()
	pg := testutils.SetupDB(pool)

	err := bcdb.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	var st *dockertest.Resource
	st, supertokenClientTest = testutils.SetupSuperTokens(pool)

	createDBTables(bcdb.DB(), supertokenClientTest)

	historyModule := history.NewHistoryClient()

	omsNPTest = NewOMSNewPlatform(context.Background(), supertokenClientTest, historyModule, false)
	verifySessionMiddleware := adaptor.HTTPMiddleware(supertokenClientTest.VerifySession)

	appTest = fiber.New()
	appTest.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	appTest.Use(LoggingMiddleware)
	// bulk
	appTest.Post("/test/bulk/factor", omsNPTest.FactorBulkPostHandler)
	// appTest.Post("/test/bulk/floor", omsNPTest.FloorBulkPostHandler) // TODO: uncomment after floor refactoring
	appTest.Post("/test/bulk/dpo", omsNPTest.DemandPartnerOptimizationBulkPostHandler)
	// floor
	appTest.Post("/test/floor", omsNPTest.FloorPostHandler)
	appTest.Post("/test/floor/get", omsNPTest.FloorGetAllHandler)
	// bulk
	appTest.Post("/test/global/factor/bulk", omsNPTest.GlobalFactorBulkPostHandler)
	appTest.Post("/test/bulk/factor", omsNPTest.FactorBulkPostHandler)
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
	// search
	appTest.Post("/test/search", omsNPTest.SearchHandler)
	// endpoint to test history saving
	appTest.Post("/bulk/global/factor", verifySessionMiddleware, omsNPTest.GlobalFactorBulkPostHandler)
	appTest.Post("/bulk/factor", verifySessionMiddleware, omsNPTest.FactorBulkPostHandler)
	// appTest.Post("/bulk/floor", verifySessionMiddleware, omsNPTest.FloorBulkPostHandler) // TODO: uncomment after floor refactoring
	appTest.Post("/bulk/dpo", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationBulkPostHandler)
	appTest.Post("/publisher/new", verifySessionMiddleware, omsNPTest.PublisherNewHandler)
	appTest.Post("/publisher/update", verifySessionMiddleware, omsNPTest.PublisherUpdateHandler)
	appTest.Post("/floor", verifySessionMiddleware, omsNPTest.FloorPostHandler)
	appTest.Post("/factor", verifySessionMiddleware, omsNPTest.FactorPostHandler)
	appTest.Post("/global/factor", verifySessionMiddleware, omsNPTest.GlobalFactorPostHandler)
	appTest.Post("/dpo/set", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationSetHandler)
	appTest.Post("/dpo/delete", verifySessionMiddleware, omsNPTest.DemandPartnerOptimizationDeleteHandler)
	appTest.Post("/publisher/domain", verifySessionMiddleware, omsNPTest.PublisherDomainPostHandler)
	appTest.Post("/targeting/set", verifySessionMiddleware, omsNPTest.TargetingSetHandler)
	appTest.Post("/targeting/update", verifySessionMiddleware, omsNPTest.TargetingUpdateHandler)
	appTest.Post("/user/update", verifySessionMiddleware, omsNPTest.UserUpdateHandler)
	appTest.Post("/user/set", verifySessionMiddleware, omsNPTest.UserSetHandler)
	appTest.Post("/block", verifySessionMiddleware, omsNPTest.BlockPostHandler)
	appTest.Post("/pixalate", verifySessionMiddleware, omsNPTest.PixalatePostHandler)
	appTest.Post("/pixalate/delete", verifySessionMiddleware, omsNPTest.PixalateDeleteHandler)
	appTest.Post("/confiant", verifySessionMiddleware, omsNPTest.ConfiantPostHandler)
	//adjust
	appTest.Post("/test/adjust/factor", omsNPTest.FactorAdjusterHandler)
	appTest.Post("/test/adjust/floor", omsNPTest.FloorAdjusterHandler)

	//bid caching
	appTest.Post("/test/bid_caching/set", validations.ValidateBidCaching, omsNPTest.BidCachingSetHandler)
	appTest.Post("/test/bid_caching/update", validations.ValidateUpdateBidCaching, omsNPTest.BidCachingUpdateHandler)
	appTest.Post("/test/bid_caching/delete", omsNPTest.BidCachingDeleteHandler)
	//refresh_cache
	appTest.Post("/test/refresh_cache/set", validations.ValidateRefreshCache, omsNPTest.RefreshCacheSetHandler)
	appTest.Post("/test/refresh_cache/update", validations.ValidateUpdateRefreshCache, omsNPTest.RefreshCacheUpdateHandler)
	appTest.Post("/test/refresh_cache/delete", omsNPTest.RefreshCacheDeleteHandler)

	go appTest.Listen(port)

	code := m.Run()

	pool.Purge(pg)
	pool.Purge(st)
	appTest.Shutdown()

	os.Exit(code)
}

func createDBTables(db *sqlx.DB, client supertokens_module.TokenManagementSystem) {
	createUserTableAndUsersInSupertokens(db, client)
	createPublisherTable(db)
	createTargetingTable(db)
	createMetaDataTable(db)
	createHistoryTable(db)
	createConfiantTable(db)
	createPixalateTable(db)
	createPublisherDomainTable(db)
	createGlobalFactorTable(db)
	createFactorTable(db)
	createFloorTable(db)
	createDPORuleTable(db)
	createPublisherDemandTable(db)
	createDPOTable(db)
	createSearchView(db)
	createRefreshCacheTable(db)
	createBidCachingTable(db)
}

func createUserTableAndUsersInSupertokens(db *sqlx.DB, client supertokens_module.TokenManagementSystem) {
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

func createPublisherTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TYPE public."integration_type" AS ENUM (` +
		`'JS Tags (Compass)',` +
		`'JS Tags (NP)',` +
		`'Prebid.js',` +
		`'Prebid Server',` +
		`'oRTB EP');`,
	)
	tx.MustExec(`CREATE TABLE public.publisher (` +
		`publisher_id varchar(36) NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`"name" varchar(1024) NOT NULL,` +
		`account_manager_id varchar(36) NULL,` +
		`media_buyer_id varchar(36) NULL,` +
		`campaign_manager_id varchar(36) NULL,` +
		`office_location varchar(36) NULL,` +
		`pause_timestamp int8 NULL,` +
		`start_timestamp int8 NULL,` +
		`status varchar(36) NULL,` +
		`reactivate_timestamp int8 NULL,` +
		`"integration_type" public."integration_type"[] NULL,` +
		`CONSTRAINT publisher_name_key UNIQUE (name),` +
		`CONSTRAINT publisher_pkey PRIMARY KEY (publisher_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.publisher ` +
		`(publisher_id, name, status, office_location, created_at)` +
		`VALUES('1111111', 'publisher_1', 'Active', 'LATAM', '2024-10-01 13:46:41.302'),` +
		`('22222222', 'publisher_2', 'Active', 'LATAM', '2024-10-01 13:46:41.302'),` +
		`('333', 'publisher_3', 'Active', 'LATAM', '2024-10-01 13:46:41.302'),` +
		`('999', 'online-media-soluctions', 'Active', 'IL', '2024-10-01 13:46:41.302'),` +
		`('444', 'publisher_4', 'Active', 'IL', '2024-10-01 13:46:41.302');`,
	)
	tx.Commit()
}

func createTargetingTable(db *sqlx.DB) {
	tx := db.MustBegin()
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
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(30, '22222222', '2.com', '300X250', 'top', '{al}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 2.0, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 'Active');`)
	tx.MustExec(`INSERT INTO public.targeting ` +
		`(id, publisher_id, "domain", unit_size, placement_type, country, device_type, browser, kv, price_model, value, created_at, updated_at, status)` +
		`VALUES(40, '999', 'oms.com', '300X250', 'top', '{il}', '{mobile}', '{firefox}', '{"key_1":"value_1","key_2":"value_2","key_3":"value_3"}'::jsonb, 'CPM', 2.0, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 'Active');`)
	tx.Commit()
}

func createMetaDataTable(db *sqlx.DB) {
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
		"demand_partner_id varchar(64)," +
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

func createConfiantTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.confiant (` +
		`confiant_key varchar(256) NOT NULL,` +
		`publisher_id varchar(36) NOT NULL,` +
		`"domain" varchar(256) NOT NULL,` +
		`rate float8 DEFAULT 0 NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`CONSTRAINT pk_confiant_1 PRIMARY KEY (domain, publisher_id)` +
		`);`,
	)
	tx.Commit()
}

func createPixalateTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.pixalate (` +
		`id varchar(256) NOT NULL,` +
		`publisher_id varchar(36) NOT NULL,` +
		`"domain" varchar(256) NOT NULL,` +
		`rate float8 DEFAULT 0 NOT NULL,` +
		`active bool DEFAULT true NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`CONSTRAINT pk_pixalate_1 PRIMARY KEY (domain, publisher_id)` +
		`);`,
	)
	tx.Commit()
}

func createPublisherDomainTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.publisher_domain (` +
		`"domain" varchar(256) NOT NULL,` +
		`publisher_id varchar(36) NOT NULL,` +
		`automation bool DEFAULT false NOT NULL,` +
		`gpp_target float8 NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`"integration_type" public."_integration_type" NULL,` +
		`CONSTRAINT publisher_domain_pkey1 PRIMARY KEY (domain, publisher_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.publisher_domain ` +
		`("domain", publisher_id, automation, gpp_target, created_at, updated_at)` +
		`VALUES('oms.com', '999', TRUE, 0.5, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createGlobalFactorTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.global_factor (` +
		`"key" varchar(36) NOT NULL,` +
		`publisher_id varchar(36) NOT NULL,` +
		`value float8 NULL,` +
		`updated_at timestamp NULL,` +
		`created_at timestamp NULL,` +
		`CONSTRAINT global_factor_pkey PRIMARY KEY (key, publisher_id)` +
		`);`,
	)
	tx.Commit()
}

func createFactorTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.factor (` +
		`publisher varchar(64) NOT NULL,` +
		`"domain" varchar(256) NOT NULL,` +
		`country varchar(64) NULL,` +
		`device varchar(64) NULL,` +
		`factor float8 DEFAULT 0 NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`rule_id varchar(36) DEFAULT ''::character varying NOT NULL,` +
		`demand_partner_id varchar(64) DEFAULT ''::character varying NOT NULL,` +
		`browser varchar(64) NULL,` +
		`os varchar(64) NULL,` +
		`placement_type varchar(64) NULL,` +
		`CONSTRAINT factor_pkey PRIMARY KEY (rule_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.factor ` +
		`(publisher, "domain", factor, rule_id, created_at, updated_at)` +
		`VALUES('999', 'oms.com', 0.5, 'oms-factor-rule-id', '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createFloorTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.floor (` +
		`publisher varchar(64) NOT NULL,` +
		`"domain" varchar(256) NOT NULL,` +
		`country varchar(64) NULL,` +
		`device varchar(64) NULL,` +
		`floor float8 DEFAULT 0 NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`rule_id varchar(36) NOT NULL,` +
		`demand_partner_id varchar(64) DEFAULT ''::character varying NOT NULL,` +
		`browser varchar(64) NULL,` +
		`os varchar(64) NULL,` +
		`placement_type varchar(64) NULL,` +
		`CONSTRAINT floor_pkey PRIMARY KEY (rule_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.floor ` +
		`(publisher, "domain", floor, rule_id, created_at, updated_at)` +
		`VALUES('999', 'oms.com', 0.5, 'oms-factor-rule-id', '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createDPORuleTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.dpo_rule (` +
		`rule_id varchar(36) NOT NULL,` +
		`demand_partner_id varchar(64) NOT NULL,` +
		`publisher varchar(64) NULL,` +
		`"domain" varchar(256) NULL,` +
		`country varchar(64) NULL,` +
		`browser varchar(64) NULL,` +
		`os varchar(64) NULL,` +
		`device_type varchar(64) NULL,` +
		`placement_type varchar(64) NULL,` +
		`factor float8 DEFAULT 0 NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`active bool DEFAULT true NOT NULL,` +
		`CONSTRAINT dpo_rule_pkey PRIMARY KEY (rule_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.dpo_rule ` +
		`(demand_partner_id, publisher, "domain", factor, rule_id, created_at, updated_at)` +
		`VALUES('test_demand_partner', '999', 'oms.com', 0.5, 'oms-factor-rule-id', '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createDPOTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.dpo (` +
		`demand_partner_id varchar(64) NOT NULL,` +
		`is_include bool DEFAULT false NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`demand_partner_name varchar(128) NULL,` +
		`active bool DEFAULT true NOT NULL,` +
		`CONSTRAINT dpo_pkey PRIMARY KEY (demand_partner_id)` +
		`);`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, demand_partner_name, active, created_at, updated_at)` +
		`VALUES('test_demand_partner', TRUE, 'Test Demand Partner', TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createPublisherDemandTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.publisher_demand (` +
		`publisher_id varchar(64) NOT NULL,` +
		`"domain" varchar(256) NOT NULL,` +
		`demand_partner_id varchar(64) NOT NULL,` +
		`ads_txt_status bool DEFAULT false NOT NULL,` +
		`active bool DEFAULT true NOT NULL,` +
		`created_at timestamp NOT NULL,` +
		`updated_at timestamp NULL,` +
		`CONSTRAINT publisher_demand_pkey PRIMARY KEY (publisher_id, domain, demand_partner_id)` +
		`);`)
	tx.MustExec(`INSERT INTO public.publisher_demand ` +
		`(publisher_id, "domain", demand_partner_id, ads_txt_status, active, created_at, updated_at)` +
		`VALUES('999', 'oms.com', 'test_demand_partner', TRUE, TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.Commit()
}

func createSearchView(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`create materialized view search_view as ` +
		`select ` +
		`'Publisher list' as section_type, ` +
		`p.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`null as "domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(p.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(null, '') || ':' || coalesce(null, '') as query ` +
		`from publisher p ` +
		`union ` +
		`select ` +
		`'Publisher / domain list' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(null, '') as query ` +
		`from publisher_domain pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`union ` +
		`select ` +
		`'Publisher / domain - Dashboard' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(null, '') as query ` +
		`from publisher_domain pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`union ` +
		`select ` +
		`'Targeting - Bidder' as section_type, ` +
		`f.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`f."domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') || ':' || coalesce(null, '') as query ` +
		`from factor f ` +
		`join publisher p on p.publisher_id = f.publisher ` +
		`union ` +
		`select ` +
		`'Targeting - JS' as section_type, ` +
		`t.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`t."domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(t.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(t."domain", '') || ':' || coalesce(null, '') as query ` +
		`from targeting t ` +
		`join publisher p on p.publisher_id = t.publisher_id ` +
		`union ` +
		`select ` +
		`'Floors' as section_type, ` +
		`f.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`f."domain", ` +
		`null as demand_partner_name, ` +
		`coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') || ':' || coalesce(null, '') as query ` +
		`from floor f ` +
		`join publisher p on p.publisher_id = f.publisher ` +
		`union ` +
		`select ` +
		`'Publisher / domain - Demand' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`d.demand_partner_name, ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') || ':' || coalesce(d.demand_partner_name, '') as query ` +
		`from publisher_demand pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`join dpo d on pd.demand_partner_id = d.demand_partner_id ` +
		`union ` +
		`select ` +
		`'DPO Rule' as section_type, ` +
		`dr.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`dr."domain", ` +
		`d.demand_partner_name, ` +
		`coalesce(dr.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(dr."domain", '') || ':' || coalesce(d.demand_partner_name, '') as query ` +
		`from dpo_rule dr ` +
		`join dpo d on dr.demand_partner_id = d.demand_partner_id ` +
		`left join publisher p on dr.publisher = p.publisher_id ` +
		`union ` +
		`select ` +
		`'Demand - Demand' as section_type, ` +
		`null as publisher_id, ` +
		`null as publisher_name, ` +
		`null as "domain", ` +
		`d.demand_partner_name, ` +
		`coalesce(null, '') || ':' || coalesce(null, '') || ':' || coalesce(null, '') || ':' || coalesce(d.demand_partner_name, '') as query ` +
		`from dpo d;`)
	tx.Commit()
}

func createRefreshCacheTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.refresh_cache (` +
		`publisher VARCHAR(64) NOT NULL,` +
		`domain VARCHAR(256) NOT NULL,` +
		`country VARCHAR(64),` +
		`active bool DEFAULT true NOT NULL,` +
		`device VARCHAR(64),` +
		`refresh_cache SMALLINT NOT NULL,` +
		`created_at TIMESTAMP NOT NULL,` +
		`updated_at TIMESTAMP,` +
		`rule_id VARCHAR(36) PRIMARY KEY,` +
		`demand_partner_id VARCHAR(64) DEFAULT ''::character varying NOT NULL,` +
		`browser VARCHAR(64),` +
		`os VARCHAR(64),` +
		`placement_type VARCHAR(64)` +
		`);`)

	tx.MustExec(`INSERT INTO public.refresh_cache ` +
		`(rule_id,publisher, domain, demand_partner_id, refresh_cache, active, created_at, updated_at)` +
		`VALUES ('123456','21038', 'oms.com', '', 10, TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)

	tx.Commit()
}

func createBidCachingTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`CREATE TABLE public.bid_caching (` +
		`publisher VARCHAR(64) NOT NULL,` +
		`domain VARCHAR(256) NOT NULL,` +
		`country VARCHAR(64),` +
		`active bool DEFAULT true NOT NULL,` +
		`device VARCHAR(64),` +
		`bid_caching SMALLINT NOT NULL,` +
		`created_at TIMESTAMP NOT NULL,` +
		`updated_at TIMESTAMP,` +
		`rule_id VARCHAR(36) PRIMARY KEY,` +
		`demand_partner_id VARCHAR(64) DEFAULT ''::character varying NOT NULL,` +
		`browser VARCHAR(64),` +
		`os VARCHAR(64),` +
		`placement_type VARCHAR(64)` +
		`);`)

	tx.MustExec(`INSERT INTO public.bid_caching ` +
		`(rule_id,publisher, domain, demand_partner_id, bid_caching, active, created_at, updated_at)` +
		`VALUES ('123456','21038', 'oms.com', '', 10, TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)

	tx.Commit()
}
