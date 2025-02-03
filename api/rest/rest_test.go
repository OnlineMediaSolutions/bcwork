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
	"github.com/m6yf/bcwork/modules/export"
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
	exportModule := export.NewExportModule()

	omsNPTest = NewOMSNewPlatform(context.Background(), supertokenClientTest, historyModule, exportModule, nil, false)
	verifySessionMiddleware := adaptor.HTTPMiddleware(supertokenClientTest.VerifySession)

	appTest = fiber.New()
	appTest.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	appTest.Use(LoggingMiddleware)
	// download
	appTest.Post("/test/download", omsNPTest.DownloadHandler)
	// bulk
	appTest.Post("/test/bulk/factor", omsNPTest.FactorBulkPostHandler)
	// appTest.Post("/test/bulk/floor", omsNPTest.FloorBulkPostHandler) // TODO: uncomment after floor refactoring
	appTest.Post("/test/bulk/dpo", omsNPTest.DemandPartnerOptimizationBulkPostHandler)
	appTest.Post("/test/bulk/global/factor/", omsNPTest.GlobalFactorBulkPostHandler)
	appTest.Post("/test/bulk/factor", omsNPTest.FactorBulkPostHandler)
	// demand partners
	appTest.Post("/test/dp/seat_owner/get", omsNPTest.DemandPartnerGetSeatOwnersHandler)
	appTest.Post("/test/dp/get", omsNPTest.DemandPartnerGetHandler)
	appTest.Post("/test/dp/set", omsNPTest.DemandPartnerSetHandler)
	appTest.Post("/test/dp/update", omsNPTest.DemandPartnerUpdateHandler)
	// block
	appTest.Post("/test/block/get", omsNPTest.BlockGetAllHandler)
	// targeting
	appTest.Post("/test/targeting/get", omsNPTest.TargetingGetHandler)
	appTest.Post("/test/targeting/set", validations.ValidateTargeting, omsNPTest.TargetingSetHandler)
	appTest.Post("/test/targeting/update", validations.ValidateTargeting, omsNPTest.TargetingUpdateHandler)
	appTest.Post("/test/targeting/tags", omsNPTest.TargetingExportTagsHandler)
	// user
	appTest.Get("/test/user/info", omsNPTest.UserGetInfoHandler)
	appTest.Get("/test/user/by_types", omsNPTest.UserGetByTypesHandler)
	appTest.Post("/test/user/get", omsNPTest.UserGetHandler)
	appTest.Post("/test/user/set", validations.ValidateUser, omsNPTest.UserSetHandler)
	appTest.Post("/test/user/update", validations.ValidateUser, omsNPTest.UserUpdateHandler)
	appTest.Post("/test/user/verify/get", verifySessionMiddleware, omsNPTest.UserGetHandler)
	appTest.Post("/test/user/verify/admin/get", verifySessionMiddleware, supertokenClientTest.AdminRoleRequired, omsNPTest.UserGetHandler)
	// publisher
	appTest.Post("/test/publisher/get", omsNPTest.PublisherGetHandler)
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
	appTest.Post("/floor/delete", verifySessionMiddleware, omsNPTest.FloorDeleteHandler)
	appTest.Post("/factor", verifySessionMiddleware, omsNPTest.FactorPostHandler)
	appTest.Post("/factor/delete", verifySessionMiddleware, omsNPTest.FactorDeleteHandler)
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
	appTest.Post("/test/bid_caching/get", omsNPTest.BidCachingGetAllHandler)
	appTest.Post("/test/bid_caching/set", validations.ValidateBidCaching, omsNPTest.BidCachingSetHandler)
	appTest.Post("/test/bid_caching/update", validations.ValidateUpdateBidCaching, omsNPTest.BidCachingUpdateHandler)
	appTest.Post("/test/bid_caching/delete", omsNPTest.BidCachingDeleteHandler)
	//refresh_cache
	appTest.Post("/test/refresh_cache/get", omsNPTest.RefreshCacheGetAllHandler)
	appTest.Post("/test/refresh_cache/set", validations.ValidateRefreshCache, omsNPTest.RefreshCacheSetHandler)
	appTest.Post("/test/refresh_cache/update", validations.ValidateUpdateRefreshCache, omsNPTest.RefreshCacheUpdateHandler)
	appTest.Post("/test/refresh_cache/delete", omsNPTest.RefreshCacheDeleteHandler)

	// floor
	appTest.Post("/test/floor/get", omsNPTest.FloorGetAllHandler)
	appTest.Post("/test/floor", validations.ValidateFloors, omsNPTest.FloorPostHandler)
	//factor
	appTest.Post("/test/factor/get", omsNPTest.FactorGetAllHandler)
	appTest.Post("/test/factor", validations.ValidateFactor, omsNPTest.FactorPostHandler)

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
	createSeatOwnerTable(db)
	createDPOTable(db)
	createSearchView(db)
	createRefreshCacheTable(db)
	createBidCachingTable(db)
	createDemandPartnerConnectionTable(db)
	createDemandPartnerChildTable(db)
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
			`types varchar(64)[], ` +
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
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user1.User.ID + `', 'user_1@oms.com', 'name_1', 'surname_1', 'Member', '{Account Manager}', 'OMS', 'Israel', '+972559999999', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload2 := `{"email": "user_2@oms.com","password": "abcd1234"}`
	req2, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload2))
	resp2, _ := http.DefaultClient.Do(req2)
	data2, _ := io.ReadAll(resp2.Body)
	defer resp2.Body.Close()
	var user2 supertokens_module.CreateUserResponse
	json.Unmarshal(data2, &user2)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user2.User.ID + `', 'user_2@oms.com', 'name_2', 'surname_2', 'Admin', '{Account Manager}', 'Google', 'USA', '+11111111', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload3 := `{"email": "user_temp@oms.com","password": "abcd1234"}`
	req3, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload3))
	resp3, _ := http.DefaultClient.Do(req3)
	data3, _ := io.ReadAll(resp3.Body)
	defer resp3.Body.Close()
	var user3 supertokens_module.CreateUserResponse
	json.Unmarshal(data3, &user3)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user3.User.ID + `', 'user_temp@oms.com', 'name_temp', 'surname_temp', 'Member', '{Campaign Manager}', 'Google', 'USA', '+77777777777', TRUE, '2024-09-01 13:46:41.302');`)

	payload4 := `{"email": "user_disabled@oms.com","password": "abcd1234"}`
	req4, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload4))
	resp4, _ := http.DefaultClient.Do(req4)
	data4, _ := io.ReadAll(resp4.Body)
	defer resp4.Body.Close()
	var user4 supertokens_module.CreateUserResponse
	json.Unmarshal(data4, &user4)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at) ` +
		`VALUES('` + user4.User.ID + `', 'user_disabled@oms.com', 'name_disabled', 'surname_disabled', 'Member', '{Media Buyer}', 'Google', 'USA', '+88888888888', FALSE, '2024-09-01 13:46:41.302');`)

	payload5 := `{"email": "user_admin@oms.com","password": "abcd1234"}`
	req5, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload5))
	resp5, _ := http.DefaultClient.Do(req5)
	data5, _ := io.ReadAll(resp5.Body)
	defer resp5.Body.Close()
	var user5 supertokens_module.CreateUserResponse
	json.Unmarshal(data5, &user5)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user5.User.ID + `', 'user_admin@oms.com', 'name_disabled', 'surname_disabled', 'Admin', '{Account Manager, Campaign Manager}', 'Google', 'USA', '+88888888888', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload6 := `{"email": "user_developer@oms.com","password": "abcd1234"}`
	req6, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload6))
	resp6, _ := http.DefaultClient.Do(req6)
	data6, _ := io.ReadAll(resp6.Body)
	defer resp6.Body.Close()
	var user6 supertokens_module.CreateUserResponse
	json.Unmarshal(data6, &user6)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user6.User.ID + `', 'user_developer@oms.com', 'name_developer', 'surname_developer', 'Developer', '{Campaign Manager, Media Buyer}', 'Apple', 'USA', '+66666666666', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload7 := `{"email": "user_history@oms.com","password": "abcd1234"}`
	req7, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload7))
	resp7, _ := http.DefaultClient.Do(req7)
	data7, _ := io.ReadAll(resp7.Body)
	defer resp7.Body.Close()
	var user7 supertokens_module.CreateUserResponse
	json.Unmarshal(data7, &user7)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", types, organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user7.User.ID + `', 'user_history@oms.com', 'name_history', 'surname_history', 'Member', '{Account Manager, Campaign Manager, Media Buyer}', 'Apple', 'USA', '+66666666666', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

	payload8 := `{"email": "user_history@oms.com","password": "abcd1234"}`
	req8, _ := http.NewRequest(http.MethodPost, client.GetWebURL()+"/public/recipe/signup", strings.NewReader(payload8))
	resp8, _ := http.DefaultClient.Do(req8)
	data8, _ := io.ReadAll(resp8.Body)
	defer resp8.Body.Close()
	var user8 supertokens_module.CreateUserResponse
	json.Unmarshal(data8, &user8)
	tx.MustExec(`INSERT INTO public.user ` +
		`(user_id, email, first_name, last_name, "role", organization_name, address, phone, enabled, created_at, password_changed) ` +
		`VALUES('` + user8.User.ID + `', 'user_without_types@oms.com', 'name_user_without_types', 'surname_user_without_types', 'Member', 'Apple', 'USA', '+66666666666', TRUE, '2024-09-01 13:46:41.302', TRUE);`)

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
	tx.MustExec(`INSERT INTO public.publisher ` +
		`(publisher_id, name, status, office_location, created_at, account_manager_id, media_buyer_id, campaign_manager_id)` +
		`VALUES('555', 'test_publisher', 'Active', 'IL', '2024-10-01 13:46:41.302', '1', '2', '3');`,
	)
	tx.Commit()
}

func createTargetingTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("create table targeting " +
		"(" +
		"id serial primary key," +
		"rule_id varchar(36) not null default '', " +
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
		`active bool DEFAULT true, ` +
		`CONSTRAINT factor_pkey PRIMARY KEY (rule_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.factor ` +
		`(publisher, "domain", factor, rule_id, created_at, updated_at)` +
		`VALUES('999', 'oms.com', 0.5, 'oms-factor-rule-id', '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.MustExec(`INSERT INTO public.factor ` +
		`(publisher, "domain", factor, rule_id, active,created_at, updated_at)` +
		`VALUES('100', 'brightcom.com', 0.5, 'oms-factor-rule-id1', false, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
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
		`active bool DEFAULT false, ` +
		`CONSTRAINT floor_pkey PRIMARY KEY (rule_id)` +
		`);`,
	)
	tx.MustExec(`INSERT INTO public.floor ` +
		`(publisher, "domain", floor, rule_id, active, created_at, updated_at)` +
		`VALUES('999', 'oms.com', 0.5, 'oms-floor-rule-id',true, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.MustExec(`INSERT INTO public.floor ` +
		`(publisher, "domain", floor, rule_id, active,created_at, updated_at)` +
		`VALUES('100', 'brightcom.com', 0.5, 'oms-floor-rule-id1', false, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
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
		`(demand_partner_id, publisher, "domain", factor, rule_id, active, created_at, updated_at)` +
		`VALUES('test_demand_partner', '999', 'oms.com', 0.5, 'oms-dpo-rule-id-1', TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.MustExec(`INSERT INTO public.dpo_rule ` +
		`(demand_partner_id, publisher, "domain", factor, rule_id, active, created_at, updated_at)` +
		`VALUES('test_demand_partner', '333', 'no_active_rules.com', 0.75, 'oms-dpo-rule-id-2', FALSE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
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
		`dp_domain varchar(128) not null default '',` +
		`certification_authority_id varchar(256),` +
		`seat_owner_id int,` +
		`manager_id int references "user"(id),` +
		`poc_name varchar(128) not null default '',` +
		`poc_email varchar(128) not null default '',` +
		`is_approval_needed bool not null default false,` +
		`approval_before_going_live bool not null default false,` +
		`approval_process varchar(64) not null default 'Other',` +
		`dp_blocks varchar(64) not null default 'Other',` +
		`score int not null default 1000,` +
		`"comments" text,` +
		`automation_name varchar(64), ` +
		`threshold float, ` +
		`automation boolean default false not null, ` +
		`CONSTRAINT dpo_pkey PRIMARY KEY (demand_partner_id)` +
		`);`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, demand_partner_name, active, created_at, updated_at)` +
		`VALUES('test_demand_partner', TRUE, 'Test Demand Partner', TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('rubicon', false, '2024-05-21 09:30:50.000', '2024-06-25 14:51:57.000', 'Rubicon DP', true, 'rubicon.com', 'yikg352gsd1', 7, 1, false, 6, 'Other', NULL, false, 'Other', '', ''); `)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('_test', false, '2024-06-25 14:51:57.000', '2024-06-25 14:51:57.000', '_Test', true, 'test.com', 'hdhr26fcb32', 6, 1, true, 7, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('index', false, '2024-05-07 17:17:11.000', '2024-06-25 14:51:57.000', 'Index', false, 'indexexchange.com', NULL, 6, 1, false, 1000, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('Finkiel', false, '2024-06-25 14:51:57.000', '2024-06-25 14:51:57.000', 'Finkiel DP', true, 'finkiel.com', 'jtfliy6893gfc', 10, 1, true, 3, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('amazon', false, '2024-05-07 17:17:11.000', '2024-06-25 14:51:57.000', 'Amazon', true, 'aps.amazon.com', 'gsrdy5352f5', 10, 1, false, 2, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('dfpdanitom', false, '2024-05-07 17:17:11.000', '2024-06-25 14:51:57.000', 'DFP Danitom', true, 'google.com', 'f08c47fec0942fa0', NULL, 1, true, 1, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('rtbhouse', false, '2024-05-21 09:30:50.000', '2024-06-25 14:51:57.000', 'RTB House', true, 'rtbhouse.com', 'ages32412we', 7, 1, false, 10, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('openx', false, '2024-05-07 17:17:11.000', '2024-06-25 14:51:57.000', 'Open X', true, 'openx.com', '235dg3sfgs3', 7, 1, true, 4, 'Other', NULL, false, 'Other', '', '');`)
	tx.MustExec(`INSERT INTO public.dpo ` +
		`(demand_partner_id, is_include, created_at, updated_at, demand_partner_name, active, dp_domain, certification_authority_id, seat_owner_id, manager_id, is_approval_needed, score, approval_process, "comments", approval_before_going_live, dp_blocks, poc_name, poc_email) ` +
		`VALUES('33across', false, '2024-05-07 17:17:11.000', '2024-06-25 14:51:57.000', '33 Across', true, '33across.com', 'fsgfcxxvge31', 7, 1, false, 5, 'Other', NULL, false, 'Other', '', '');`)
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
		`'Publishers list' as section_type, ` +
		`p.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`null as "domain", ` +
		`coalesce(p.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(null, '') as query ` +
		`from publisher p ` +
		`union ` +
		`select ` +
		`'Domains list' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') as query ` +
		`from publisher_domain pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`union ` +
		`select ` +
		`'Domain - Dashboard' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') as query ` +
		`from publisher_domain pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`union ` +
		`select ` +
		`'Bidder Targetings' as section_type, ` +
		`f.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`f."domain", ` +
		`coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') as query ` +
		`from factor f ` +
		`join publisher p on p.publisher_id = f.publisher ` +
		`where f.active = TRUE ` +
		`union ` +
		`select ` +
		`'JS Targetings' as section_type, ` +
		`t.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`t."domain", ` +
		`coalesce(t.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(t."domain", '') as query ` +
		`from targeting t ` +
		`join publisher p on p.publisher_id = t.publisher_id ` +
		`union ` +
		`select ` +
		`'Floors list' as section_type, ` +
		`f.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`f."domain", ` +
		`coalesce(f.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(f."domain", '') as query ` +
		`from floor f ` +
		`join publisher p on p.publisher_id = f.publisher ` +
		`where f.active = TRUE ` +
		`union ` +
		`select ` +
		`'Domain - Demand' as section_type, ` +
		`pd.publisher_id, ` +
		`p."name" as publisher_name, ` +
		`pd."domain", ` +
		`coalesce(pd.publisher_id, '') || ':' || coalesce(p."name", '') || ':' || coalesce(pd."domain", '') as query ` +
		`from publisher_demand pd ` +
		`join publisher p on p.publisher_id = pd.publisher_id ` +
		`join dpo d on pd.demand_partner_id = d.demand_partner_id ` +
		`union ` +
		`select ` +
		`'DPO Rules' as section_type, ` +
		`dr.publisher as publisher_id, ` +
		`p."name" as publisher_name, ` +
		`dr."domain", ` +
		`coalesce(dr.publisher, '') || ':' || coalesce(p."name", '') || ':' || coalesce(dr."domain", '') as query ` +
		`from dpo_rule dr ` +
		`join dpo d on dr.demand_partner_id = d.demand_partner_id ` +
		`left join publisher p on dr.publisher = p.publisher_id ` +
		`where dr.active = TRUE;`)
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
	tx.MustExec(`INSERT INTO public.refresh_cache ` +
		`(rule_id,publisher, domain, demand_partner_id, refresh_cache, active, created_at, updated_at)` +
		`VALUES ('1234567','21038', 'brightcom.com', '', 10, FALSE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)

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
		`control_percentage float8, ` +
		`created_at TIMESTAMP NOT NULL,` +
		`updated_at TIMESTAMP,` +
		`rule_id VARCHAR(36) PRIMARY KEY,` +
		`demand_partner_id VARCHAR(64) DEFAULT ''::character varying NOT NULL,` +
		`browser VARCHAR(64),` +
		`os VARCHAR(64),` +
		`placement_type VARCHAR(64)` +
		`);`)

	tx.MustExec(`INSERT INTO public.bid_caching ` +
		`(rule_id,publisher, domain, demand_partner_id, bid_caching, active, created_at, updated_at, control_percentage)` +
		`VALUES ('123456','21038', 'oms.com', '', 10, TRUE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407', 0.5);`)
	tx.MustExec(`INSERT INTO public.bid_caching ` +
		`(rule_id,publisher, domain, demand_partner_id, bid_caching, active, created_at, updated_at)` +
		`VALUES ('1234567','21000', 'brightcom.com', '', 10, FALSE, '2024-10-01 13:51:28.407', '2024-10-01 13:51:28.407');`)

	tx.Commit()
}

func createSeatOwnerTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`create table if not exists seat_owner ` +
		`( ` +
		`id serial primary key, ` +
		`seat_owner_name varchar(128) not null default '', ` +
		`seat_owner_domain varchar(128) not null default '', ` +
		`publisher_account varchar(256) not null default '%s', ` +
		`certification_authority_id varchar(256), ` +
		`created_at timestamp not null, ` +
		`updated_at timestamp ` +
		`);`)

	tx.MustExec(`insert into seat_owner (seat_owner_domain, seat_owner_name, publisher_account, created_at) ` +
		`values ` +
		`('adsparc.com', 'Adsparc', '7%s', '2024-10-01 13:51:28.407'), ` +
		`('onomagic.com', 'Onomagic', '%s1', '2024-10-01 13:51:28.407'), ` +
		`('sparcmedia.com', 'Sparcmedia', '3%s', '2024-10-01 13:51:28.407'), ` +
		`('audienciad.com', 'Audienciad', '%s2', '2024-10-01 13:51:28.407'), ` +
		`('limpid.tv', 'Limpid', '9%s', '2024-10-01 13:51:28.407'), ` +
		`('getmediamx.com', 'GetMedia', '12%s', '2024-10-01 13:51:28.407'), ` +
		`('brightcom.com', 'Brightcom', '%s', '2024-10-01 13:51:28.407'), ` +
		`('whildey.com', 'Whildey', '%s5', '2024-10-01 13:51:28.407'), ` +
		`('advibe.media', 'SaharMedia', '8%s', '2024-10-01 13:51:28.407'), ` +
		`('onlinemediasolutions.com', 'OMS', '%s', '2024-10-01 13:51:28.407');`)

	tx.Commit()
}

func createDemandPartnerChildTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`create table if not exists demand_partner_child ` +
		`( ` +
		`id serial primary key, ` +
		`dp_connection_id int not null references demand_partner_connection(id), ` +
		`dp_child_name varchar(128) not null default '', ` +
		`dp_child_domain varchar(128) not null default '', ` +
		`publisher_account varchar(256) not null default '', ` +
		`certification_authority_id varchar(256), ` +
		`is_direct bool not null default false, ` +
		`active bool not null default true, ` +
		`is_required_for_ads_txt bool not null default false, ` +
		`created_at timestamp not null, ` +
		`updated_at timestamp ` +
		`);`)

	tx.MustExec(`INSERT INTO public.demand_partner_child ` +
		`(dp_connection_id, created_at, dp_child_name, dp_child_domain, publisher_account, certification_authority_id, is_required_for_ads_txt) ` +
		`values ` +
		`(4, '2024-10-01 13:51:28.407', 'Open X', 'openx.com', '88888', NULL, false), ` +
		`(11, '2024-10-01 13:51:28.407', 'Pubmatic', 'pubmatic.com', '44444', NULL, false), ` +
		`(5, '2024-10-01 13:51:28.407', 'Appnexus', 'appnexus.com', '55555', NULL, false), ` +
		`(11, '2024-10-01 13:51:28.407', 'Appnexus', 'appnexus.com', '121212', NULL, true), ` +
		`(5, '2024-10-01 13:51:28.407', 'Index', 'indexexchange.com', '131313', NULL, true), ` +
		`(5, '2024-10-01 13:51:28.407', 'AOL', 'adtech.com', '111111', NULL, false), ` +
		`(5, '2024-10-01 13:51:28.407', 'Rubicon DP', 'rubicon.com', '66666', 'srwadcae523', true);`)

	tx.Commit()
}

func createDemandPartnerConnectionTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`create table if not exists demand_partner_connection ` +
		`( ` +
		`id serial primary key, ` +
		`demand_partner_id varchar(64) not null references dpo(demand_partner_id), ` +
		`publisher_account varchar(256) not null default '', ` +
		`integration_type varchar(64)[], ` +
		`active bool not null default true, ` +
		`is_direct bool not null default false,` +
		`is_required_for_ads_txt bool not null default false,` +
		`created_at timestamp not null, ` +
		`updated_at timestamp ` +
		`);`)

	tx.MustExec(`INSERT INTO public.demand_partner_connection ` +
		`(demand_partner_id, created_at, publisher_account, "integration_type", is_direct, is_required_for_ads_txt) ` +
		`values ` +
		`('rubicon', '2024-10-01 13:51:28.407', '66666', '{js, s2s}', false, true), ` +
		`('index', '2024-10-01 13:51:28.407', '181818', '{js, s2s}', false, true), ` +
		`('_test', '2024-10-01 13:51:28.407', '77777', '{js, s2s}', false, true), ` +
		`('Finkiel', '2024-10-01 13:51:28.407', '11111', '{js, s2s}', false, true), ` +
		`('amazon', '2024-10-01 13:51:28.407', 's2s141414', '{s2s}', true, true), ` +
		`('amazon', '2024-10-01 13:51:28.407', '141414', '{js}', true, true), ` +
		`('dfpdanitom', '2024-10-01 13:51:28.407', 'pub-2243508421279209', '{js, s2s}', true, true), ` +
		`('rtbhouse', '2024-10-01 13:51:28.407', '202020', '{js, s2s}', false, false), ` +
		`('openx', '2024-10-01 13:51:28.407', '22222', '{js}', false, true), ` +
		`('openx', '2024-10-01 13:51:28.407', 's2s22222', '{s2s}', false, true), ` +
		`('33across', '2024-10-01 13:51:28.407', '33333', '{js, s2s}', false, true);`)

	tx.Commit()
}
