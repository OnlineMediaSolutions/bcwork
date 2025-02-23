package adstxt

import (
	"fmt"

	"github.com/m6yf/bcwork/models"
)

const (
	demandPartnerConnectionQueryType = iota
	demandPartnerChildQueryType
	seatOwnerQueryType
)

func getAdsTxtLinesTemplateQuery(queryType int) string {
	baseQuery := `
		select 
			pd."domain", 
			pd.publisher_id,
			t.id as %s,
			%s as demand_status
		from publisher_domain as pd
		cross join %s as t
		%s
		where d.active and t.id = ANY($1);
	`
	var fieldName, tableName, demandStatus, joinClause string
	switch queryType {
	case demandPartnerConnectionQueryType:
		fieldName = models.AdsTXTColumns.DemandPartnerConnectionID
		tableName = models.TableNames.DemandPartnerConnection
		demandStatus = "case when d.is_approval_needed then 'not_sent' else 'approved' end"
		joinClause = "join dpo d on d.demand_partner_id = t.demand_partner_id "
	case demandPartnerChildQueryType:
		fieldName = models.AdsTXTColumns.DemandPartnerChildID
		tableName = models.TableNames.DemandPartnerChild
		demandStatus = "case when d.is_approval_needed then 'not_sent' else 'approved' end"
		joinClause = `
			join demand_partner_connection dpc on t.dp_connection_id = dpc.id 
			join dpo d on d.demand_partner_id = dpc.demand_partner_id 
		`
	case seatOwnerQueryType:
		fieldName = models.AdsTXTColumns.SeatOwnerID
		tableName = models.TableNames.SeatOwner
		demandStatus = "'approved'"
		joinClause = "join dpo d on d.seat_owner_id = t.id "
	}

	return fmt.Sprintf(baseQuery, fieldName, demandStatus, tableName, joinClause)
}

func getDynamicColumnName(queryType int) string {
	var columnName string
	switch queryType {
	case demandPartnerConnectionQueryType:
		columnName = models.AdsTXTColumns.DemandPartnerConnectionID
	case demandPartnerChildQueryType:
		columnName = models.AdsTXTColumns.DemandPartnerChildID
	case seatOwnerQueryType:
		columnName = models.AdsTXTColumns.SeatOwnerID
	}

	return columnName
}

func getDynamicColumnValue(queryType int, line *adsTxtLineTemplate) *int {
	var columnValue *int
	switch queryType {
	case demandPartnerConnectionQueryType:
		columnValue = line.DemandPartnerConnectionID.Ptr()
	case demandPartnerChildQueryType:
		columnValue = line.DemandPartnerChildID.Ptr()
	case seatOwnerQueryType:
		columnValue = line.SeatOwnerID.Ptr()
	}

	return columnValue
}
