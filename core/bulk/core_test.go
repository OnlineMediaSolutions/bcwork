package bulk

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

var (
	pool *dockertest.Pool
)

func TestMain(m *testing.M) {
	pg, _ := initTest()
	code := m.Run()

	//Run all core tests
	pool.Purge(pg)
	os.Exit(code)

}

func initTest() (*dockertest.Resource, context.Context) {
	pool = testutils.SetupDockerTestPool()
	pg := testutils.SetupDB(pool)

	err := bcdb.DB().Ping()
	if err != nil {
		log.Fatal()
	}
	ctx := context.Background()
	createTables(bcdb.DB())
	return pg, ctx
}

func createTables(db *sqlx.DB) {
	tx := db.MustBegin()

	tx.MustExec("CREATE TABLE IF NOT EXISTS dpo_rule (rule_id varchar(36) not null primary key,demand_partner_id varchar(64) not null, publisher varchar(64),domain varchar(256),country varchar(64),browser varchar(64),os varchar(64),  device_type varchar(64), placement_type varchar(64), factor float8 not null default 0, created_at timestamp not null,updated_at timestamp, active bool not null default true)")
	tx.MustExec("INSERT INTO dpo_rule (rule_id, demand_partner_id, publisher, domain, country, browser, os, device_type, placement_type, factor, created_at, updated_at, active) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		"1111", "Finkiel", "20360", "mako.co.il", "jp", nil, nil, "mobile", nil, 20, "2024-12-01 14:24:33.100", "2024-12-01 14:24:33.100", false)
	tx.MustExec("INSERT INTO dpo_rule (rule_id, demand_partner_id, publisher, domain, country, browser, os, device_type, placement_type, factor, created_at, updated_at, active) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		"1234", "onetagbcm", "20956", "docsachhay.net", "gb", nil, nil, "mobile", nil, 10, "2024-12-01 14:24:33.100", "2024-12-01 14:24:33.100", true)
	tx.MustExec("INSERT INTO dpo_rule (rule_id, demand_partner_id, publisher, domain, country, browser, os, device_type, placement_type, factor, created_at, updated_at, active) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		"5678", "onetagbcm", "20360", "finkiel.co.il", "il", nil, "android", "mobile", nil, 20, "2024-12-01 14:24:33.100", "2024-12-01 14:24:33.100", true)

	tx.MustExec("CREATE TABLE history" +
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
	tx.MustExec("CREATE TABLE IF NOT EXISTS metadata_queue (transaction_id varchar(36) primary key not null, key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")

	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"426cba39-7d1c-59fd-ad61-36a03a92415b", "price:factor:v2:1234:finkiel.com", nil, "{\"rules\":[{\"rule\":\"(p=1234__d=finkiel.com__c=gb__os=.*__dt=mobile__pt=.*__b=.*)\",\"factor\":2,\"rule_id\":\"e81337e9-983c-50f9-9fca-e1f2131c5ed8\"},{\"rule\":\"(p=1234__d=finkiel.com__c=il__os=.*__dt=desktop__pt=.*__b=.*)\",\"factor\":4,\"rule_id\":\"80ecfa53-2a28-548b-a371-743dbb22c437\"}]}", 0, "2024-09-20T10:10:10.100", "2024-09-26T10:10:10.100")

	tx.MustExec("CREATE TABLE IF NOT EXISTS factor (publisher varchar(64), domain varchar(256), country varchar(64), device varchar(64), factor float8 not null default 0, created_at timestamp not null, updated_at timestamp, rule_id varchar(36) not null default '',demand_partner_id varchar(64) not null default '',browser varchar(64), os varchar(64),placement_type varchar(64), active bool not null default true)")
	tx.MustExec("INSERT INTO factor (publisher, domain, country, device, factor, created_at, updated_at,rule_id,demand_partner_id,browser,os,placement_type, active) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		"1234", "finkiel.com", "il", "desktop", 4, "2024-12-01 14:24:33.100", "2024-12-01 14:24:33.100", "80ecfa53-2a28-548b-a371-743dbb22c437", "", nil, nil, nil, true)
	tx.MustExec("INSERT INTO factor (publisher, domain, country, device, factor, created_at, updated_at,rule_id,demand_partner_id,browser,os,placement_type, active) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		"1234", "finkiel.com", "gb", "mobile", 2, "2024-12-01 14:24:33.100", "2024-12-01 14:24:33.100", "e81337e9-983c-50f9-9fca-e1f2131c5ed8", "", nil, nil, nil, true)
	tx.Commit()
}
