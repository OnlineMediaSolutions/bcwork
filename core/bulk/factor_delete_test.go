package bulk

import (
	"context"
	"encoding/json"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"testing"
)

func TestDeleteFactorMechanism(t *testing.T) {
	historyModule := history.NewHistoryClient()
	service := BulkService{historyModule: historyModule}
	ctx := context.Background()

	//checking that we have 2 factors in metadata_queue
	metaDataFactors, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:factor:v2:1234:finkiel.com")).All(ctx, bcdb.DB())
	value := metaDataFactors[0].Value
	var rules MetaDataRules
	json.Unmarshal(value, &rules)

	factors, _ := models.Factors().All(ctx, bcdb.DB())
	assert.Equal(t, len(rules.Rules), 2)
	assert.Equal(t, len(factors), 2)

	//delete 1 factor
	ids := []string{"80ecfa53-2a28-548b-a371-743dbb22c437"}
	service.BulkDeleteFactor(ctx, ids)

	//checking that we have only 1 factor in metadata_queue
	metaDataFactor, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:factor:v2:1234:finkiel.com"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	json.Unmarshal(metaDataFactor.Value, &rules)
	factors, _ = models.Factors().All(ctx, bcdb.DB())

	assert.Equal(t, len(rules.Rules), 1)
	assert.Equal(t, rules.Rules[0].RuleID, "e81337e9-983c-50f9-9fca-e1f2131c5ed8")
	assert.Equal(t, len(factors), 2)

	//check that one factor is active and one is not
	for _, factor := range factors {
		if factor.RuleID == "80ecfa53-2a28-548b-a371-743dbb22c437" {
			assert.False(t, factor.Active)
		}
		if factor.RuleID == "e81337e9-983c-50f9-9fca-e1f2131c5ed8" {
			assert.True(t, factor.Active)
		}
	}
}

type MetaDataRules struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Rule   string `json:"rule"`
	Factor int    `json:"factor"`
	RuleID string `json:"rule_id"`
}
