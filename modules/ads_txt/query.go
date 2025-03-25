package adstxt

import (
	"fmt"

	"github.com/m6yf/bcwork/models"
)

const (
	demandPartnerConnectionQueryType = iota
	demandPartnerChildQueryType
	seatOwnerQueryType

	demandPartnerDemandStatusStatement = "case when d.is_approval_needed then 'not_sent' else 'approved' end"
	seatOwnerDemandStatusStatement     = "'approved'"
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
		where t.id = ANY($1);
	`
	var fieldName, tableName, demandStatus, joinClause string
	switch queryType {
	case demandPartnerConnectionQueryType:
		fieldName = models.AdsTXTColumns.DemandPartnerConnectionID
		tableName = models.TableNames.DemandPartnerConnection
		demandStatus = demandPartnerDemandStatusStatement
		joinClause = "join dpo d on d.demand_partner_id = t.demand_partner_id "
	case demandPartnerChildQueryType:
		fieldName = models.AdsTXTColumns.DemandPartnerChildID
		tableName = models.TableNames.DemandPartnerChild
		demandStatus = demandPartnerDemandStatusStatement
		joinClause = `
			join demand_partner_connection dpc on t.dp_connection_id = dpc.id 
			join dpo d on d.demand_partner_id = dpc.demand_partner_id 
		`
	case seatOwnerQueryType:
		fieldName = models.AdsTXTColumns.SeatOwnerID
		tableName = models.TableNames.SeatOwner
		demandStatus = seatOwnerDemandStatusStatement
		joinClause = "join dpo d on d.seat_owner_id = t.id "
	}

	return fmt.Sprintf(baseQuery, fieldName, demandStatus, tableName, joinClause)
}

func getAdsTxtLinesFromPublisherDomainTemplateQuery() string {
	baseQuery := `
		select 
			pd.publisher_id,
			pd."domain",
			dpc.id as demand_partner_connection_id,
			null::int as demand_partner_child_id,
			null::int as seat_owner_id,
			%s as demand_status
		from publisher_domain pd
		cross join demand_partner_connection dpc 
		join dpo d on dpc.demand_partner_id = d.demand_partner_id 
		where pd."domain" = $1 and publisher_id = $2
		union
		select 
			pd.publisher_id,
			pd."domain",
			null::int as demand_partner_connection_id,
			dpc.id as demand_partner_child_id,
			null::int as seat_owner_id,
			%s as demand_status
		from publisher_domain pd
		cross join demand_partner_child dpc 
		join demand_partner_connection dpc2 on dpc2.id = dpc.dp_connection_id
		join dpo d on dpc2.demand_partner_id = d.demand_partner_id 
		where pd."domain" = $1 and publisher_id = $2
		union
		select 
			pd.publisher_id,
			pd."domain",
			null::int as demand_partner_connection_id,
			null::int as demand_partner_child_id,
			so.id as seat_owner_id,
			%s as demand_status
		from publisher_domain pd
		cross join seat_owner so 
		join dpo d on so.id = d.seat_owner_id 
		where pd."domain" = $1 and publisher_id = $2;
	`

	return fmt.Sprintf(baseQuery, demandPartnerDemandStatusStatement, demandPartnerDemandStatusStatement, seatOwnerDemandStatusStatement)
}
