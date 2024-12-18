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

func TestDeleteFloorMechanism(t *testing.T) {

	historyModule := history.NewHistoryClient()
	service := BulkService{historyModule: historyModule}
	ctx := context.Background()

	//checking that we have 2 floors in metadata_queue
	metaDataFloors, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:floor:v2:1234:finkiel.com")).All(ctx, bcdb.DB())
	value := metaDataFloors[0].Value
	var rules MetaDataFloorRules
	json.Unmarshal(value, &rules)

	floors, _ := models.Floors().All(ctx, bcdb.DB())
	assert.Equal(t, len(rules.Rule), 2)
	assert.Equal(t, len(floors), 2)

	//delete 1 floor
	ids := []string{"80ecfa53-2a28-548b-a371-743dbb22c439"}
	service.BulkDeleteFloor(ctx, ids)

	//checking that we have only 1 floor in metadata_queue
	metaDataFloor, _ := models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:floor:v2:1234:finkiel.com"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	json.Unmarshal(metaDataFloor.Value, &rules)
	floors, _ = models.Floors().All(ctx, bcdb.DB())

	assert.Equal(t, len(rules.Rule), 1)
	assert.Equal(t, rules.Rule[0].RuleID, "e81337e9-983c-50f9-9fca-e1f2131c5ed0")
	assert.Equal(t, len(floors), 2)

	//check that one floor is active and one is not
	for _, floor := range floors {
		if floor.RuleID == "80ecfa53-2a28-548b-a371-743dbb22c439" {
			assert.False(t, floor.Active)
		}
		if floor.RuleID == "e81337e9-983c-50f9-9fca-e1f2131c5ed0" {
			assert.True(t, floor.Active)
		}
	}

	//delete 2nd floor
	ids = []string{"e81337e9-983c-50f9-9fca-e1f2131c5ed0"}
	service.BulkDeleteFloor(ctx, ids)

	//checking that metadata_queue rules array is empty
	metaDataFloor, _ = models.MetadataQueues(models.MetadataQueueWhere.Key.EQ("price:floor:v2:1234:finkiel.com"), qm.OrderBy("updated_at desc")).One(ctx, bcdb.DB())
	json.Unmarshal(metaDataFloor.Value, &rules)
	floors, _ = models.Floors().All(ctx, bcdb.DB())

	assert.Equal(t, len(rules.Rule), 0)
}

type MetaDataFloorRules struct {
	Rule []FloorRule `json:"rules"`
}

type FloorRule struct {
	Rule   string `json:"rule"`
	Floor  int    `json:"floor"`
	RuleID string `json:"rule_id"`
}
