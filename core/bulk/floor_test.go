package bulk

//
//import (
//	"github.com/m6yf/bcwork/models"
//	"github.com/m6yf/bcwork/utils/constant"
//	"testing"
//)
//
////func Test_prepareBulkInsertFloorsRequest(t *testing.T) {
////	t.Parallel()
////
////	type args struct {
////		floors []models.Floor
////	}
////
////	tests := []struct {
////		name string
////		args args
////		want *bulkInsertRequest
////	}{
////		{
////			name: "valid",
////			args: args{
////				floors: []models.Floor{
////					{
////						RuleID:        "rule_1",
////						Publisher:     "publisher_1",
////						Domain:        "1.com",
////						Country:       null.StringFrom("IL"),
////						Browser:       null.StringFrom("firefox"),
////						Os:            null.StringFrom("linux"),
////						Device:        null.StringFrom("mobile"),
////						PlacementType: null.StringFrom("top"),
////						Floor:         0.1,
////					},
////				},
////			},
////			want: &bulkInsertRequest{
////				tableName: models.TableNames.Floor,
////				columns: []string{
////					models.FloorColumns.RuleID,
////					models.FloorColumns.Publisher,
////					models.FloorColumns.Domain,
////					models.FloorColumns.Country,
////					models.FloorColumns.Browser,
////					models.FloorColumns.Os,
////					models.FloorColumns.Device,
////					models.FloorColumns.PlacementType,
////					models.FloorColumns.Floor,
////					models.FloorColumns.CreatedAt,
////					models.FloorColumns.UpdatedAt,
////				},
////				conflictColumns: []string{
////					models.FloorColumns.RuleID,
////				},
////				updateColumns: []string{
////					models.FloorColumns.Floor,
////					models.FloorColumns.UpdatedAt,
////				},
////				valueStrings: []string{
////					"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
////				},
////				args: []interface{}{
////					"rule_1", "publisher_1", "1.com", "IL", "firefox", "linux", "mobile", "top", 0.1, currentTime, currentTime,
////				},
////			},
////		},
////	}
////
////	for _, tt := range tests {
////		tt := tt
////		t.Run(tt.name, func(t *testing.T) {
////			t.Parallel()
////
////			got := prepareBulkInsertFloorsRequest(tt.args.floors)
////			assert.Equal(t, tt.want, got)
////		})
////	}
////}
//
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
//					{
//						Publisher:     "publisher_2",
//						Domain:        "2.com",
//						Country:       "US",
//						Browser:       "chrome",
//						OS:            "macos",
//						Device:        "mobile",
//						PlacementType: "bottom",
//						Floor:         0.05,
//					},
//					{
//						Publisher:     "publisher_3",
//						Domain:        "3.com",
//						Country:       "RU",
//						Browser:       "opera",
//						OS:            "windows",
//						Device:        "mobile",
//						PlacementType: "side",
//						Floor:         0.15,
//					},
//				},
//			},
//			want: []models.Floor{
//				{
//					RuleID:        "09baffdc-c450-5491-adae-aefae59e28cc",
//					Publisher:     "publisher_1",
//					Domain:        "1.com",
//					Country:       "IL",
//					Browser:       "firefox",
//					Os:            "linux",
//					Device:        "mobile",
//					PlacementType: "top",
//					Floor:         0.1,
//				},
//				{
//					RuleID:        "95a177ad-80f0-5c11-87f9-ac73da58dc14",
//					Publisher:     "publisher_2",
//					Domain:        "2.com",
//					Country:       "US",
//					Browser:       "chrome",
//					Os:            "macos",
//					Device:        "mobile",
//					PlacementType: "bottom",
//					Floor:         0.05,
//				},
//				{
//					RuleID:        "a8b46eac-8d71-5108-959a-0cb6826d49ec",
//					Publisher:     "publisher_3",
//					Domain:        "3.com",
//					Country:       "RU",
//					Browser:       "opera",
//					Os:            "windows",
//					Device:        "mobile",
//					PlacementType: "side",
//					Floor:         0.15,
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
