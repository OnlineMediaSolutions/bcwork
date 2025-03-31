package dto

import (
	"fmt"
	"strings"

	"github.com/m6yf/bcwork/bcdb/order"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types"
)

const (
	// demand partner statuses
	DPStatusPending         = "pending"
	DPStatusApproved        = "approved"
	DPStatusApprovedPaused  = "approved_paused"
	DPStatusRejected        = "rejected"
	DPStatusRejectedTQ      = "rejected_tq"
	DPStatusDisabledSPO     = "disabled_spo"
	DPStatusDisabledNoImps  = "disabled_no_imps"
	DPStatusHighDiscrepancy = "high_discrepancy"
	DPStatusNotSent         = "not_sent"
	DPStatusNoForm          = "no_form"
	DPStatusWillNotBeSent   = "will_not_be_sent"
	// domain statuses
	DomainStatusActive = "active"
	DomainStatusNew    = "new"
	DomainStatusPaused = "paused"
	// ads.txt lines types
	AdsTxtTypeDirect   = "DIRECT"
	AdsTxtTypeReseller = "RESELLER"
	// ads.txt lines statuses
	AdsTxtStatusAdded      = "added"
	AdsTxtStatusDeleted    = "deleted"
	AdsTxtStatusNotScanned = "not_scanned"
	AdsTxtStatusNo         = "no"
	// ads.txt actions
	AdsTxtActionAdd           = "add"
	AdsTxtActionFix           = "fix"
	AdsTxtActionKeep          = "keep"
	AdsTxtActionLowPerfomance = "low_perfomance"
	AdsTxtActionNone          = "none"
	AdsTxtActionRemove        = "remove"
)

var (
	DPStatusMap = map[string]string{
		DPStatusPending:         "Pending",
		DPStatusApproved:        "Approved",
		DPStatusApprovedPaused:  "Approved - Paused",
		DPStatusRejected:        "Rejected",
		DPStatusRejectedTQ:      "Rejected - TQ",
		DPStatusDisabledSPO:     "Disabled - SPO",
		DPStatusDisabledNoImps:  "Disabled - 0 Sold Imps",
		DPStatusHighDiscrepancy: "High Discrepancy",
		DPStatusNotSent:         "Not sent",
		DPStatusNoForm:          "No form",
		DPStatusWillNotBeSent:   "Will not be sent",
	}

	DomainStatusMap = map[string]string{
		DomainStatusActive: "Active",
		DomainStatusNew:    "New",
		DomainStatusPaused: "Paused",
	}

	StatusMap = map[string]string{
		AdsTxtStatusAdded:      "Added",
		AdsTxtStatusDeleted:    "Deleted",
		AdsTxtStatusNotScanned: "Not Scanned",
		AdsTxtStatusNo:         "No",
	}
)

type AdsTxt struct {
	ID                        int               `json:"id"`
	GroupByDPID               int               `json:"group_by_dp_id"`
	CursorID                  int               `json:"cursor_id"`
	PublisherID               string            `json:"publisher_id"`
	PublisherName             string            `json:"publisher_name"`
	MirrorPublisherID         null.String       `json:"mirror_publisher_id"`
	MirrorPublisherName       null.String       `json:"mirror_publisher_name"`
	AccountManagerID          null.String       `json:"account_manager_id"`
	AccountManagerFullName    null.String       `json:"account_manager_full_name"`
	CampaignManagerID         null.String       `json:"campaign_manager_id"`
	CampaignManagerFullName   null.String       `json:"campaign_manager_full_name"`
	Domain                    string            `json:"domain"`
	DomainStatus              string            `json:"domain_status"`
	DemandPartnerID           string            `json:"demand_partner_id"`
	DemandPartnerName         string            `json:"demand_partner_name"`
	DemandPartnerNameExtended string            `json:"demand_partner_name_extended"` // like Amazon - Amazon or OMS - Direct
	DemandPartnerConnectionID null.Int          `json:"demand_partner_connection_id"`
	MediaType                 types.StringArray `json:"media_type"`
	DemandManagerID           null.String       `json:"demand_manager_id"`
	DemandManagerFullName     null.String       `json:"demand_manager_full_name"`
	DemandStatus              string            `json:"demand_status"`
	IsDemandPartnerActive     bool              `json:"is_demand_partner_active"`
	SeatOwnerName             string            `json:"seat_owner_name"`
	Score                     int               `json:"score"`
	Action                    string            `json:"action"`
	Status                    string            `json:"status"`
	IsRequired                bool              `json:"is_required"`
	AdsTxtLine                string            `json:"ads_txt_line"`
	Added                     int               `json:"added"`      // count of added lines
	Total                     int               `json:"total"`      // total amount of lines
	DPEnabled                 bool              `json:"dp_enabled"` // dp is ready to go live
	LastScannedAt             null.Time         `json:"last_scanned_at"`
	ErrorMessage              null.String       `json:"error_message"`
}

type AdsTxtGroupedByDP struct {
	*AdsTxt
	GroupedLines []*AdsTxt `json:"grouped_lines"`
}

func (a *AdsTxtGroupedByDP) FromAdsTxt(row *AdsTxt) {
	a.ID = row.ID
	a.GroupByDPID = row.GroupByDPID
	a.CursorID = row.CursorID
	a.PublisherID = row.PublisherID
	a.PublisherName = row.PublisherName
	a.MirrorPublisherID = row.MirrorPublisherID
	a.AccountManagerID = row.AccountManagerID
	a.AccountManagerFullName = row.AccountManagerFullName
	a.CampaignManagerID = row.CampaignManagerID
	a.CampaignManagerFullName = row.CampaignManagerFullName
	a.Domain = row.Domain
	a.DomainStatus = row.DomainStatus
	a.DemandPartnerID = row.DemandPartnerID
	a.DemandPartnerName = row.DemandPartnerName
	a.DemandPartnerNameExtended = row.DemandPartnerNameExtended
	a.DemandPartnerConnectionID = row.DemandPartnerConnectionID
	a.MediaType = row.MediaType
	a.DemandManagerID = row.DemandManagerID
	a.DemandManagerFullName = row.DemandManagerFullName
	a.DemandStatus = row.DemandStatus
	a.IsDemandPartnerActive = row.IsDemandPartnerActive
	a.SeatOwnerName = row.SeatOwnerName
	a.Score = row.Score
	a.Action = row.Action
	a.Status = row.Status
	a.IsRequired = row.IsRequired
	a.AdsTxtLine = row.AdsTxtLine
	a.Added = row.Added
	a.Total = row.Total
	a.DPEnabled = row.DPEnabled
	a.LastScannedAt = row.LastScannedAt
	a.ErrorMessage = row.ErrorMessage
}

// ProcessParentRow processing parent row of group by dp ads.txt table in priority:
// 1. Main Line (Amazon - Amazon);
// 2. Seat Owner Line (OMS - Direct);
// 3. Any other line (EBDA - OpenX);
func (a *AdsTxtGroupedByDP) ProcessParentRow(row *AdsTxt) {
	const seatOwnerLineSuffix = "- Direct"

	isMainLine := fmt.Sprintf("%v - %v", row.DemandPartnerName, row.DemandPartnerName) == row.DemandPartnerNameExtended
	isSeatOwnerLine := strings.HasSuffix(row.DemandPartnerNameExtended, seatOwnerLineSuffix)

	var isParentRowAlreadySet bool
	if a != nil && a.AdsTxt != nil {
		isParentRowAlreadySet = fmt.Sprintf("%v - %v", row.DemandPartnerName, row.DemandPartnerName) == a.DemandPartnerNameExtended
	}

	if isParentRowAlreadySet {
		return
	}

	if isMainLine {
		a.FromAdsTxt(row)
		return
	}

	if isSeatOwnerLine {
		a.FromAdsTxt(row)
		return
	}
}

type AdsTxtUpdateRequest struct {
	Domain          []string `json:"domain" validate:"required,min=1"`
	DemandPartnerID *string  `json:"demand_partner_id,omitempty"`
	DomainStatus    *string  `json:"domain_status,omitempty" validate:"adsTxtDomainStatus"`
	DemandStatus    *string  `json:"demand_status,omitempty" validate:"adsTxtDemandStatus"`
}

type AdsTxtResponse struct {
	Data  []*AdsTxt `json:"data"`
	Total int64     `json:"total"`
}

type AdsTxtGroupByDPResponse struct {
	Data  []*AdsTxtGroupedByDP `json:"data"`
	Total int64                `json:"total"`
	Order order.Sort           `json:"-"`
}

func (a *AdsTxtGroupByDPResponse) Len() int {
	return len(a.Data)
}

func (a *AdsTxtGroupByDPResponse) Swap(i, j int) {
	a.Data[i], a.Data[j] = a.Data[j], a.Data[i]
}

func (a *AdsTxtGroupByDPResponse) Less(i, j int) bool {
	for _, order := range a.Order {
		var (
			compared bool
			result   bool
		)

		switch order.Name {
		case "publisher_id":
			compared = a.Data[i].PublisherID != a.Data[j].PublisherID
			result = a.Data[i].PublisherID < a.Data[j].PublisherID
			if order.Desc {
				result = a.Data[i].PublisherID > a.Data[j].PublisherID
			}
		case "mirror_publisher_id":
			compared = a.Data[i].MirrorPublisherID.String != a.Data[j].MirrorPublisherID.String
			result = a.Data[i].MirrorPublisherID.String < a.Data[j].MirrorPublisherID.String
			if order.Desc {
				result = a.Data[i].MirrorPublisherID.String > a.Data[j].MirrorPublisherID.String
			}
		case "account_manager_id":
			compared = a.Data[i].AccountManagerID.String != a.Data[j].AccountManagerID.String
			result = a.Data[i].AccountManagerID.String < a.Data[j].AccountManagerID.String
			if order.Desc {
				result = a.Data[i].AccountManagerID.String > a.Data[j].AccountManagerID.String
			}
		case "publisher_name":
			compared = a.Data[i].PublisherName != a.Data[j].PublisherName
			result = a.Data[i].PublisherName < a.Data[j].PublisherName
			if order.Desc {
				result = a.Data[i].PublisherName > a.Data[j].PublisherName
			}
		case "campaign_manager_id":
			compared = a.Data[i].CampaignManagerID.String != a.Data[j].CampaignManagerID.String
			result = a.Data[i].CampaignManagerID.String < a.Data[j].CampaignManagerID.String
			if order.Desc {
				result = a.Data[i].CampaignManagerID.String > a.Data[j].CampaignManagerID.String
			}
		case "domain":
			compared = a.Data[i].Domain != a.Data[j].Domain
			result = a.Data[i].Domain < a.Data[j].Domain
			if order.Desc {
				result = a.Data[i].Domain > a.Data[j].Domain
			}
		case "demand_status":
			compared = a.Data[i].DemandStatus != a.Data[j].DemandStatus
			result = a.Data[i].DemandStatus < a.Data[j].DemandStatus
			if order.Desc {
				result = a.Data[i].DemandStatus > a.Data[j].DemandStatus
			}
		case "domain_status":
			compared = a.Data[i].DomainStatus != a.Data[j].DomainStatus
			result = a.Data[i].DomainStatus < a.Data[j].DomainStatus
			if order.Desc {
				result = a.Data[i].DomainStatus > a.Data[j].DomainStatus
			}
		case "demand_manager_id":
			compared = a.Data[i].DemandManagerID.String != a.Data[j].DemandManagerID.String
			result = a.Data[i].DemandManagerID.String < a.Data[j].DemandManagerID.String
			if order.Desc {
				result = a.Data[i].DemandManagerID.String > a.Data[j].DemandManagerID.String
			}
		case "demand_partner_name":
			compared = a.Data[i].DemandPartnerName != a.Data[j].DemandPartnerName
			result = a.Data[i].DemandPartnerName < a.Data[j].DemandPartnerName
			if order.Desc {
				result = a.Data[i].DemandPartnerName > a.Data[j].DemandPartnerName
			}
		}

		if compared {
			return result
		}
	}

	return a.Data[i].CursorID < a.Data[j].CursorID
}
