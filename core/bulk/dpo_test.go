package bulk

import (
	"testing"

	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_prepareBulkInsertDPORequest(t *testing.T) {
	t.Parallel()

	type args struct {
		dpos []*models.DpoRule
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				dpos: []*models.DpoRule{
					{
						RuleID:          "rule_1",
						DemandPartnerID: "dp_1",
						Publisher:       null.String{Valid: true, String: "publisher_1"},
						Domain:          null.String{Valid: true, String: "1.com"},
						Country:         null.String{Valid: true, String: "IL"},
						Browser:         null.String{Valid: true, String: "firefox"},
						Os:              null.String{Valid: true, String: "linux"},
						DeviceType:      null.String{Valid: true, String: "mobile"},
						PlacementType:   null.String{Valid: true, String: "top"},
						Factor:          0.1,
						Active:          true,
					},
					{
						RuleID:          "rule_2",
						DemandPartnerID: "dp_2",
						Publisher:       null.String{Valid: true, String: "publisher_2"},
						Domain:          null.String{Valid: true, String: "2.com"},
						Country:         null.String{Valid: true, String: "US"},
						Browser:         null.String{Valid: true, String: "chrome"},
						Os:              null.String{Valid: true, String: "macos"},
						DeviceType:      null.String{Valid: true, String: "mobile"},
						PlacementType:   null.String{Valid: true, String: "bottom"},
						Factor:          0.05,
						Active:          true,
					},
					{
						RuleID:          "rule_3",
						DemandPartnerID: "dp_3",
						Publisher:       null.String{Valid: true, String: "publisher_3"},
						Domain:          null.String{Valid: true, String: "3.com"},
						Country:         null.String{Valid: true, String: "RU"},
						Browser:         null.String{Valid: true, String: "opera"},
						Os:              null.String{Valid: true, String: "windows"},
						DeviceType:      null.String{Valid: true, String: "mobile"},
						PlacementType:   null.String{Valid: true, String: "side"},
						Factor:          0.15,
						Active:          true,
					},
				},
			},
			want: &bulkInsertRequest{
				tableName: models.TableNames.DpoRule,
				columns: []string{
					models.DpoRuleColumns.RuleID,
					models.DpoRuleColumns.DemandPartnerID,
					models.DpoRuleColumns.Publisher,
					models.DpoRuleColumns.Domain,
					models.DpoRuleColumns.Country,
					models.DpoRuleColumns.Browser,
					models.DpoRuleColumns.Os,
					models.DpoRuleColumns.DeviceType,
					models.DpoRuleColumns.PlacementType,
					models.DpoRuleColumns.Factor,
					models.DpoRuleColumns.CreatedAt,
					models.DpoRuleColumns.UpdatedAt,
					models.DpoRuleColumns.Active,
				},
				conflictColumns: []string{
					models.DpoRuleColumns.RuleID,
				},
				updateColumns: []string{
					models.DpoRuleColumns.Factor,
					models.DpoRuleColumns.UpdatedAt,
					models.DpoRuleColumns.Active,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
					"($14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)",
					"($27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39)",
				},
				args: []interface{}{
					"rule_1", "dp_1",
					null.String{Valid: true, String: "publisher_1"},
					null.String{Valid: true, String: "1.com"},
					null.String{Valid: true, String: "IL"},
					null.String{Valid: true, String: "firefox"},
					null.String{Valid: true, String: "linux"},
					null.String{Valid: true, String: "mobile"},
					null.String{Valid: true, String: "top"},
					0.1, constant.PostgresCurrentTime, constant.PostgresCurrentTime, true,

					"rule_2", "dp_2",
					null.String{Valid: true, String: "publisher_2"},
					null.String{Valid: true, String: "2.com"},
					null.String{Valid: true, String: "US"},
					null.String{Valid: true, String: "chrome"},
					null.String{Valid: true, String: "macos"},
					null.String{Valid: true, String: "mobile"},
					null.String{Valid: true, String: "bottom"},
					0.05, constant.PostgresCurrentTime, constant.PostgresCurrentTime, true,

					"rule_3", "dp_3",
					null.String{Valid: true, String: "publisher_3"},
					null.String{Valid: true, String: "3.com"},
					null.String{Valid: true, String: "RU"},
					null.String{Valid: true, String: "opera"},
					null.String{Valid: true, String: "windows"},
					null.String{Valid: true, String: "mobile"},
					null.String{Valid: true, String: "side"},
					0.15, constant.PostgresCurrentTime, constant.PostgresCurrentTime, true,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertDPORequest(tt.args.dpos)
			assert.Equal(t, tt.want, got)
		})
	}
}
