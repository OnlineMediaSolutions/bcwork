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
	tx.MustExec("CREATE TABLE IF NOT EXISTS metadata_queue (transaction_id varchar(36) primary key not null, key varchar(256), version varchar(16),value varchar(512),commited_instances integer, created_at timestamp, updated_at timestamp)")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"0db4298c-84f0-5b13-8147-3e38a2526d15", "dpo:onetagbcm", nil, "{\"rules\": [{\"rule\": \"(p=20972__d=juno.com__c=ca__os=android__dt=.*__pt=.*__b=.*)\", \"factor\": 90, \"rule_id\": \"1234\"}, {\"rule\": \"(p=20360__d=coinmarketcap.com__c=gb__os=android__dt=.*__pt=.*__b=.*)\", \"factor\": 90, \"rule_id\": \"6543\"}], \"is_include\": false, \"demand_partner_id\": \"onetagbcm\"}", 0, "2024-10-20T10:10:10.100", "2024-10-26T10:10:10.100")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"62423155-7a24-590e-9868-70228c206d42", "dpo:onetagbcm", nil, "{\"rules\":[{\"rule\":\"(p=20972__d=juno.com__c=ca__os=android__dt=.*__pt=.*__b=.*)\",\"factor\":90,\"rule_id\":\"24e18928-bba8-5344-a1cc-6c7ca59ea4d0\"},{\"rule\":\"(p=20360__d=coinmarketcap.com__c=gb__os=android__dt=.*__pt=.*__b=.*)\",\"factor\":90,\"rule_id\":\"621f4990-ccd8-50b4-8dfa-a685c6681e52\"}],\"is_include\":false,\"demand_partner_id\":\"onetagbcm\"}", 0, "2024-10-20T10:10:10.100", "2024-10-26T10:10:10.100")
	tx.MustExec("INSERT INTO metadata_queue (transaction_id, key, version, value, commited_instances, created_at, updated_at) "+
		"VALUES ($1,$2, $3, $4, $5, $6, $7)",
		"64b6cf0e-48ae-5401-87ff-d67a877be8c8", "dpo:onetagbcm", nil, "{\"rules\":[{\"rule\":\"(p=20972__d=juno.com__c=ca__os=android__dt=.*__pt=.*__b=.*)\",\"factor\":90,\"rule_id\":\"1234\"},{\"rule\":\"(p=20360__d=coinmarketcap.com__c=gb__os=android__dt=.*__pt=.*__b=.*)\",\"factor\":90,\"rule_id\":\"5678\"},{\"rule\":\"(p=20360__d=coinmarketcap.com__c=gb__os=android__dt=.*__pt=.*__b=.*)\",\"factor\":90,\"rule_id\":\"9877\"}],\"is_include\":false,\"demand_partner_id\":\"onetagbcm\"}", 0, "2024-11-20T10:10:10.100", "2024-11-26T10:10:10.100")
	tx.Commit()
}

func TestDeleteDpoRuleId(t *testing.T) {
	pg, ctx := initTest()

	//prints the port for debug purposes
	log.Println("port: " + pg.GetPort("5432/tcp"))

	request := DPODeleteRequest{
		DemandPartner: "onetagbcm",
		RuleId:        "9877",
	}

	//run the main method
	DeleteDpoRuleId(ctx, request)

	//checking that the ruleId was removed
	rule, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("dpo:onetagbcm"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())

	var dpoValueData DPOValueData
	json.Unmarshal(rule.Value, &dpoValueData)

	for _, rule := range dpoValueData.Rules {
		if rule.RuleID == "9877" {
			assert.Equal(t, "5678", rule.RuleID, "Rule with rule_id 9877 should have been removed")
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
