package dto

import (
	"github.com/volatiletech/null/v8"
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
	// ads.txt lines statuses
	AdsTxtStatusAdded      = "added"
	AdsTxtStatusDeleted    = "deleted"
	AdsTxtStatusNotScanned = "not_scanned"
	AdsTxtStatusNo         = "no"
	// domain statuses
	DomainStatusActive = "active"
	DomainStatusNew    = "new"
	DomainStatusPaused = "paused"
)

type AdsTxt struct {
	ID                        int         `json:"id"`
	PublisherID               string      `json:"publisher_id"`
	PublisherName             string      `json:"publisher_name"`
	AccountManagerID          null.String `json:"account_manager_id"`
	AccountManagerFullName    string      `json:"account_manager_full_name"`
	CampaignManagerID         null.String `json:"campaign_manager_id"`
	CampaignManagerFullName   string      `json:"campaign_manager_full_name"`
	Domain                    string      `json:"domain"`
	DomainStatus              string      `json:"domain_status"`
	DomainIntegrationType     string      `json:"domain_intergration_type"`
	DemandPartnerName         string      `json:"demand_partner_name"`
	DemandPartnerNameExtended string      `json:"demand_partner_name_extended"` // like Amazon - Amazon or OMS - Direct
	DPIntegrationType         string      `json:"dp_intergration_type"`
	DemandManagerID           null.String `json:"demand_manager_id"`
	DemandManagerFullName     string      `json:"demand_manager_full_name"`
	DemandStatus              string      `json:"demand_status"`
	SeatOwnerName             string      `json:"seat_owner_name"`
	Score                     int         `json:"score"`
	Status                    string      `json:"status"`
	IsRequired                bool        `json:"is_required"`
	AdsTxtLine                string      `json:"ads_txt_line"`
	Added                     int         `json:"added"` // count of added lines
	Total                     int         `json:"total"` // total amount of lines
	IsReadyToWork             bool        `json:"is_ready_to_go_live"`
	LastScannedAt             null.Time   `json:"last_scanned_at"`
	ErrorMessage              null.String `json:"error_message"`
}

type AdsTxtGroupedByDPData struct {
	Parent   *AdsTxt   `json:"parent"`
	Children []*AdsTxt `json:"children"`
}

type AdsTxtUpdateRequest struct {
	ID           int    `json:"id"`
	DomainStatus string `json:"domain_status"`
	DemandStatus string `json:"demand_status"`
}
