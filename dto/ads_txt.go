package dto

import (
	"fmt"
	"strings"

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

type AdsTxt struct {
	CursorID                  int               `json:"cursor_id"`
	ID                        int               `json:"id"`
	PublisherID               string            `json:"publisher_id"`
	PublisherName             string            `json:"publisher_name"`
	AccountManagerID          null.String       `json:"account_manager_id"`
	AccountManagerFullName    string            `json:"account_manager_full_name"`
	CampaignManagerID         null.String       `json:"campaign_manager_id"`
	CampaignManagerFullName   string            `json:"campaign_manager_full_name"`
	Domain                    string            `json:"domain"`
	DomainStatus              string            `json:"domain_status" validate:"adsTxtDomainStatus"`
	DemandPartnerID           string            `json:"demand_partner_id"`
	DemandPartnerName         string            `json:"demand_partner_name"`
	DemandPartnerNameExtended string            `json:"demand_partner_name_extended"` // like Amazon - Amazon or OMS - Direct
	DemandPartnerConnectionID null.Int          `json:"demand_partner_connection_id"`
	MediaType                 types.StringArray `json:"media_type"`
	DemandManagerID           null.String       `json:"demand_manager_id"`
	DemandManagerFullName     string            `json:"demand_manager_full_name"`
	DemandStatus              string            `json:"demand_status" validate:"adsTxtDemandStatus"`
	IsDemandPartnerActive     bool              `json:"is_demand_partner_active"`
	SeatOwnerName             string            `json:"seat_owner_name"`
	Score                     int               `json:"score"`
	Action                    string            `json:"action"`
	Status                    string            `json:"status"`
	IsRequired                bool              `json:"is_required"`
	AdsTxtLine                string            `json:"ads_txt_line"`
	Added                     int               `json:"added"` // count of added lines
	Total                     int               `json:"total"` // total amount of lines
	IsReadyToGoLive           bool              `json:"is_ready_to_go_live"`
	LastScannedAt             null.Time         `json:"last_scanned_at"`
	ErrorMessage              null.String       `json:"error_message"`
	IsMirrorUsed              bool              `json:"is_mirror_used"`
}

func (a *AdsTxt) Mirror(source *AdsTxt) {
	a.AdsTxtLine = source.AdsTxtLine
	a.Status = source.Status
	a.DomainStatus = source.DomainStatus
	a.DemandStatus = source.DemandStatus
	a.Added = source.Added
	a.Total = source.Total
	a.IsReadyToGoLive = source.IsReadyToGoLive
	a.IsMirrorUsed = true
}

type AdsTxtGroupedByDPData struct {
	Parent   *AdsTxt   `json:"parent"`
	Children []*AdsTxt `json:"children"`
}

// ProcessParentRow processing parent row of group by dp ads.txt table in priority:
// 1. Main Line (Amazon - Amazon);
// 2. Seat Owner Line (OMS - Direct);
// 3. Any other line (EBDA - OpenX);
func (a *AdsTxtGroupedByDPData) ProcessParentRow(row *AdsTxt) {
	const seatOwnerLineSuffix = "- Direct"

	isMainLine := fmt.Sprintf("%v - %v", row.DemandPartnerName, row.DemandPartnerName) == row.DemandPartnerNameExtended
	isSeatOwnerLine := strings.HasSuffix(row.DemandPartnerNameExtended, seatOwnerLineSuffix)

	var isParentRowAlreadySet bool
	if a.Parent != nil {
		isParentRowAlreadySet = fmt.Sprintf("%v - %v", row.DemandPartnerName, row.DemandPartnerName) == a.Parent.DemandPartnerNameExtended
	}

	if isParentRowAlreadySet {
		return
	}

	if isMainLine {
		a.Parent = row
		return
	}

	if isSeatOwnerLine {
		a.Parent = row
		return
	}
}

type AdsTxtUpdateRequest struct {
	ID           int    `json:"id"`
	DomainStatus string `json:"domain_status"`
	DemandStatus string `json:"demand_status"`
}
