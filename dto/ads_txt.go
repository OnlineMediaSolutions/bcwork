package dto

import (
	"encoding/json"

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
	ID                        int               `boil:"id" json:"id"`
	CursorID                  int               `boil:"cursor_id" json:"cursor_id"`
	PublisherID               string            `boil:"publisher_id" json:"publisher_id"`
	PublisherName             string            `boil:"publisher_name" json:"publisher_name"`
	MirrorPublisherID         null.String       `boil:"mirror_publisher_id" json:"mirror_publisher_id"`
	MirrorPublisherName       null.String       `boil:"mirror_publisher_name" json:"mirror_publisher_name"`
	AccountManagerID          null.String       `boil:"account_manager_id" json:"account_manager_id"`
	AccountManagerFullName    null.String       `boil:"account_manager_full_name" json:"account_manager_full_name"`
	CampaignManagerID         null.String       `boil:"campaign_manager_id" json:"campaign_manager_id"`
	CampaignManagerFullName   null.String       `boil:"campaign_manager_full_name" json:"campaign_manager_full_name"`
	Domain                    string            `boil:"domain" json:"domain"`
	DomainStatus              string            `boil:"domain_status" json:"domain_status"`
	DemandPartnerID           string            `boil:"demand_partner_id" json:"demand_partner_id"`
	DemandPartnerName         string            `boil:"demand_partner_name" json:"demand_partner_name"`
	DemandPartnerNameExtended string            `boil:"demand_partner_name_extended" json:"demand_partner_name_extended"` // like Amazon - Amazon or OMS - Direct
	DemandPartnerConnectionID null.Int          `boil:"demand_partner_connection_id" json:"demand_partner_connection_id"`
	MediaType                 types.StringArray `boil:"media_type" json:"media_type"`
	DemandManagerID           null.String       `boil:"demand_manager_id" json:"demand_manager_id"`
	DemandManagerFullName     null.String       `boil:"demand_manager_full_name" json:"demand_manager_full_name"`
	DemandStatus              string            `boil:"demand_status" json:"demand_status"`
	IsDemandPartnerActive     bool              `boil:"is_demand_partner_active" json:"is_demand_partner_active"`
	SeatOwnerName             string            `boil:"seat_owner_name" json:"seat_owner_name"`
	Score                     int               `boil:"score" json:"score"`
	Action                    string            `boil:"action" json:"action"`
	Status                    string            `boil:"status" json:"status"`
	IsRequired                bool              `boil:"is_required" json:"is_required"`
	AdsTxtLine                string            `boil:"ads_txt_line" json:"ads_txt_line"`
	Added                     int               `boil:"added" json:"added"`           // count of added lines
	Total                     int               `boil:"total" json:"total"`           // total amount of lines
	DPEnabled                 bool              `boil:"dp_enabled" json:"dp_enabled"` // dp is ready to go live
	LastScannedAt             null.Time         `boil:"last_scanned_at" json:"last_scanned_at"`
	ErrorMessage              null.String       `boil:"error_message" json:"error_message"`
}

type AdsTxtGroupedByDP struct {
	*AdsTxt         `boil:",bind"`
	GroupedLinesRaw json.RawMessage `boil:"grouped_lines_raw" json:"-"`
	GroupedLines    []*AdsTxt       `boil:"_" json:"grouped_lines"`
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
}
