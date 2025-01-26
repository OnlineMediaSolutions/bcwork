package dto

import (
	"github.com/volatiletech/null/v8"
)

type AdsTxt struct {
	ID                      int         `json:"id"`
	PublisherID             string      `json:"publisher_id"`
	PublisherName           string      `json:"publisher_name"`
	AccountManagerID        null.String `json:"account_manager_id"`
	AccountManagerFullName  string      `json:"account_manager_full_name"`
	CampaignManagerID       null.String `json:"campaign_manager_id"`
	CampaignManagerFullName string      `json:"campaign_manager_full_name"`
	Domain                  string      `json:"domain"`
	DomainStatus            string      `json:"domain_status"`
	DomainIntegrationType   string      `json:"domain_intergration_type"`
	DemandPartnerName       string      `json:"demand_partner_name"`
	DPIntegrationType       string      `json:"dp_intergration_type"`
	DemandManagerID         null.String `json:"demand_manager_id"`
	DemandManagerFullName   string      `json:"demand_manager_full_name"`
	DemandStatus            string      `json:"demand_status"`
	SeatOwnerName           string      `json:"seat_owner_name"`
	Score                   int         `json:"score"`
	Status                  string      `json:"status"`
	IsRequired              bool        `json:"is_required"`
	AdsTxtLine              string      `json:"ads_txt_line"`
	LastScannedAt           null.Time   `json:"last_scanned_at"`
	ErrorMessage            null.String `json:"error_message"`
}
