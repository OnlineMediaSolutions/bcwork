package core

import (
	"context"
	"encoding/json"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"testing"
)

var (
	pool *dockertest.Pool
)

func TestDPOCreateRulesInMetaDataQueueTableMethod(t *testing.T) {
	ctx := context.Background()
	//prints the port for debug purposes
	//log.Println("port: " + pg.GetPort("5432/tcp"))

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
}
