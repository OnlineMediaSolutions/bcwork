package core

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/testutils"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"testing"
)

var (
	pool *dockertest.Pool
	app  *fiber.App

	port    = ":9000"
	baseURL = "http://localhost" + port
)

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

	tx.MustExec("CREATE TABLE IF NOT EXISTS metadata_queue (transaction_id varchar(36) primary key not null, key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.Commit()
}

func TestDPOCreateRulesInMetaDataQueueTableMethod(t *testing.T) {
	pg, ctx := initTest()

	//prints the port for debug purposes
	log.Println("port: " + pg.GetPort("5432/tcp"))

	//run the main method for 2 different demand partners
	sendToRT(ctx, "Finkiel")
	sendToRT(ctx, "onetagbcm")

	//checking that the Rules member is empty
	emptyRules, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("dpo:Finkiel"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	fullRules, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("dpo:onetagbcm"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())

	var dpoEmptyRuleData DPOValueData
	json.Unmarshal(emptyRules.Value, &dpoEmptyRuleData)

	var dpoRuleData DPOValueData
	json.Unmarshal(fullRules.Value, &dpoRuleData)

	//checking that all data is according to expectations
	assert.Len(t, dpoEmptyRuleData.Rules, 0)
	assert.Len(t, dpoRuleData.Rules, 2)

	for _, rule := range dpoRuleData.Rules {
		if rule.RuleID == "1234" {
			assert.Equal(t, "(p=20956__d=docsachhay.net__c=gb__os=.*__dt=mobile__pt=.*__b=.*)", rule.Rule, "Rule is incorrect")
			assert.Equal(t, 10, rule.Factor, "Factor should be 10")
		}
		if rule.RuleID == "5678" {
			assert.Equal(t, "(p=20360__d=finkiel.co.il__c=il__os=android__dt=mobile__pt=.*__b=.*)", rule.Rule, "Rule is incorrect")
			assert.Equal(t, 20, rule.Factor, "Factor should be 20")
		}
	}
	pool.Purge(pg)
}

func initTest() (*dockertest.Resource, context.Context) {
	pool = testutils.SetupDockerTestPool()
	pg := testutils.SetupDB(pool)

	err := bcdb.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	createTables(bcdb.DB())
	return pg, ctx
}
