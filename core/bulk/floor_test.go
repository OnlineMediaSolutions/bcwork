package bulk

import (
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"testing"
)

func Test_prepareBulkInsertFloorsRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		floors []models.Floor
	}

	tests := []struct {
		name string
		args args
		want *bulkInsertRequest
	}{
		{
			name: "valid",
			args: args{
				floors: []models.Floor{
					{
						RuleID:        "rule_1",
						Publisher:     "publisher_1",
						Domain:        "1.com",
						Country:       null.StringFrom("IL"),
						Browser:       null.StringFrom("firefox"),
						Os:            null.StringFrom("linux"),
						Device:        null.StringFrom("mobile"),
						PlacementType: null.StringFrom("top"),
						Floor:         0.1,
					},
					{
						RuleID:        "rule_2",
						Publisher:     "publisher_2",
						Domain:        "2.com",
						Country:       null.StringFrom("US"),
						Browser:       null.StringFrom("chrome"),
						Os:            null.StringFrom("macos"),
						Device:        null.StringFrom("tablet"),
						PlacementType: null.StringFrom("bottom"),
						Floor:         0.05,
					},
				},
			},
			want: &bulkInsertRequest{
				tableName: models.TableNames.Floor,
				columns: []string{
					models.FloorColumns.RuleID,
					models.FloorColumns.Publisher,
					models.FloorColumns.Domain,
					models.FloorColumns.Country,
					models.FloorColumns.Browser,
					models.FloorColumns.Os,
					models.FloorColumns.Device,
					models.FloorColumns.PlacementType,
					models.FloorColumns.Floor,
					models.FloorColumns.CreatedAt,
					models.FloorColumns.UpdatedAt,
				},
				conflictColumns: []string{
					models.FloorColumns.RuleID,
				},
				updateColumns: []string{
					models.FloorColumns.Floor,
					models.FloorColumns.UpdatedAt,
				},
				valueStrings: []string{
					"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
					"($12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)",
				},
				args: []interface{}{
					"rule_1", "publisher_1", "1.com", null.String{String: "IL", Valid: true}, null.String{String: "firefox", Valid: true}, null.String{String: "linux", Valid: true}, null.String{String: "mobile", Valid: true}, null.String{String: "top", Valid: true}, 0.1, currentTime, currentTime,
					"rule_2", "publisher_2", "2.com", null.String{String: "US", Valid: true}, null.String{String: "chrome", Valid: true}, null.String{String: "macos", Valid: true}, null.String{String: "tablet", Valid: true}, null.String{String: "bottom", Valid: true}, 0.05, currentTime, currentTime,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareBulkInsertFloorsRequest(tt.args.floors)
			assert.Equal(t, tt.want, got)
		})
	}
}

//func Test_prepareFloors(t *testing.T) {
//	t.Parallel()
//
//	type args struct {
//		chunk []constant.FloorUpdateRequest
//	}
//
//	tests := []struct {
//		name string
//		args args
//		want []models.Floor
//	}{
//		{
//			name: "valid",
//			args: args{
//				chunk: []constant.FloorUpdateRequest{
//					{
//						Publisher:     "publisher_1",
//						Domain:        "1.com",
//						Country:       "IL",
//						Browser:       "firefox",
//						OS:            "linux",
//						Device:        "mobile",
//						PlacementType: "top",
//						Floor:         0.1,
//					},
//				},
//			},
//			want: []models.Floor{
//				{
//					RuleID:        "09baffdc-c450-5491-adae-aefae59e28cc",
//					Publisher:     "publisher_1",
//					Domain:        "1.com",
//					Country:       null.String{String: "us", Valid: true},
//					Browser:       null.StringFrom("firefox"),
//					Os:            null.StringFrom("linux"),
//					Device:        null.String{String: "mobile", Valid: true},
//					PlacementType: null.StringFrom("top"),
//					Floor:         0.1,
//				},
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		tt := tt
//		t.Run(tt.name, func(t *testing.T) {
//			t.Parallel()
//
//			got := prepareFloors(tt.args.chunk)
//			assert.Equal(t, tt.want, got)
//		})
//	}
//}

func Test_prepareFloors(t *testing.T) {
	t.Parallel()

	type args struct {
		chunk []constant.FloorUpdateRequest
	}

	tests := []struct {
		name string
		args args
		want []models.Floor
	}{
		{
			name: "valid",
			args: args{
				chunk: []constant.FloorUpdateRequest{
					{
						Publisher:     "publisher_1",
						Domain:        "1.com",
						Country:       "IL",
						Browser:       "firefox",
						OS:            "linux",
						Device:        "mobile",
						PlacementType: "top",
						Floor:         0.1,
					},
				},
			},
			want: []models.Floor{
				{
					RuleID:        "expected-rule-id",
					Publisher:     "publisher_1",
					Domain:        "1.com",
					Country:       null.StringFrom("IL"),
					Device:        null.StringFrom("mobile"),
					Floor:         0.1,
					Browser:       null.StringFrom("firefox"),
					Os:            null.StringFrom("linux"),
					PlacementType: null.StringFrom("top"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareFloors(tt.args.chunk)
			assert.Equal(t, tt.want, got)
		})
	}
}
