package core

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	pool *dockertest.Pool
)

func TestDPOCreateRulesInMetaDataQueueTableMethod(t *testing.T) {
	type Rule struct {
		Rule   string `json:"rule"`
		Factor int    `json:"factor"`
		RuleID string `json:"rule_id"`
	}

	type DPOValueData struct {
		Rules           []Rule `json:"rules"`
		IsInclude       bool   `json:"is_include"`
		DemandPartnerID string `json:"demand_partner_id"`
	}

	ctx := context.Background()
	//prints the port for debug purposes
	//log.Println("port: " + pg.GetPort("5432/tcp"))

	//run the main method for 2 different demand partners
	err := sendToRT(ctx, "Finkiel")
	assert.NoError(t, err)
	err = sendToRT(ctx, "onetagbcm")
	assert.NoError(t, err)

	//checking that the Rules member is empty
	emptyRules, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("dpo:Finkiel"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	assert.NoError(t, err)
	fullRules, err := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("dpo:onetagbcm"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	assert.NoError(t, err)

	var dpoEmptyRuleData DPOValueData
	err = json.Unmarshal(emptyRules.Value, &dpoEmptyRuleData)
	assert.NoError(t, err)

	var dpoRuleData DPOValueData
	err = json.Unmarshal(fullRules.Value, &dpoRuleData)
	assert.NoError(t, err)

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
