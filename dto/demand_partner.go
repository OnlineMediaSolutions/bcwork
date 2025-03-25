package dto

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/null/v8"
)

const (
	DefaultDemandPartnerScoreValue = 1000
	// integration types
	ORTBIntergrationType         = "oRTB"
	PrebidServerIntergrationType = "Prebid Server"
	AmazonAPSIntergrationType    = "Amazon APS"
	// approval process
	EmailApprovalProcess                 = "via Email"
	DemandPartnerPlatformApprovalProcess = "via DP Platform"
	GDocApprovalProcess                  = "GDoc"
	OtherApprovalProcess                 = "Other"
	// media types
	WebBannersMediaType = "Web Banners"
	VideoMediaType      = "Video"
	InAppMediaType      = "InApp"
)

type DemandPartner struct {
	DemandPartnerID         string                     `json:"demand_partner_id"`
	DemandPartnerName       string                     `json:"demand_partner_name" validate:"required"`
	Connections             []*DemandPartnerConnection `json:"connections"`
	ApprovalProcess         string                     `json:"approval_process" validate:"approvalProcess"`
	DPBlocks                string                     `json:"dp_blocks" validate:"dpBlocks"`
	POCName                 string                     `json:"poc_name"`
	POCEmail                string                     `json:"poc_email"`
	SeatOwnerID             *int                       `json:"seat_owner_id"`
	SeatOwnerName           string                     `json:"seat_owner_name"`
	ManagerID               *int                       `json:"manager_id" validate:"required"`
	ManagerFullName         string                     `json:"manager_full_name"`
	IntegrationType         []string                   `json:"integration_type" validate:"integrationType"`
	MediaTypeList           []string                   `json:"media_type_list"`
	IsInclude               bool                       `json:"is_include"`
	Active                  bool                       `json:"active"`
	IsApprovalNeeded        bool                       `json:"is_approval_needed"`
	ApprovalBeforeGoingLive bool                       `json:"approval_before_going_live"`
	Automation              bool                       `json:"automation"`
	AutomationName          string                     `json:"automation_name"`
	Threshold               float64                    `json:"threshold" validate:"dpThreshold"`
	Score                   int                        `json:"score"`
	Comments                *string                    `json:"comments"`
	CreatedAt               time.Time                  `json:"created_at"`
	UpdatedAt               *time.Time                 `json:"updated_at"`
}

func (dp *DemandPartner) FromModel(mod *models.Dpo) {
	dp.DemandPartnerID = mod.DemandPartnerID
	dp.DemandPartnerName = mod.DemandPartnerName
	dp.Connections = func() []*DemandPartnerConnection {
		connections := make([]*DemandPartnerConnection, 0, len(mod.R.DemandPartnerDemandPartnerConnections))
		for _, modConnection := range mod.R.DemandPartnerDemandPartnerConnections {
			connection := new(DemandPartnerConnection)
			connection.FromModel(mod.DemandPartnerName, modConnection)
			connections = append(connections, connection)
		}

		return connections
	}()
	dp.ApprovalProcess = mod.ApprovalProcess
	dp.DPBlocks = mod.DPBlocks
	dp.POCName = mod.PocName
	dp.POCEmail = mod.PocEmail
	dp.SeatOwnerID = mod.SeatOwnerID.Ptr()
	dp.SeatOwnerName = func() string {
		if mod.R.SeatOwner != nil {
			return mod.R.SeatOwner.SeatOwnerName
		}

		return ""
	}()
	dp.ManagerID = mod.ManagerID.Ptr()
	dp.ManagerFullName = func() string {
		if mod.R.Manager != nil {
			return mod.R.Manager.FirstName + " " + mod.R.Manager.LastName
		}

		return ""
	}()
	dp.IntegrationType = mod.IntegrationType
	dp.MediaTypeList = func() []string {
		mediaTypes := make([]string, 0, len(mod.R.DemandPartnerDemandPartnerConnections))
		for _, modConnection := range mod.R.DemandPartnerDemandPartnerConnections {
			mediaTypes = append(mediaTypes, modConnection.MediaType...)
		}

		return mediaTypes
	}()
	dp.Active = mod.Active
	dp.IsInclude = mod.IsInclude
	dp.IsApprovalNeeded = mod.IsApprovalNeeded
	dp.ApprovalBeforeGoingLive = mod.ApprovalBeforeGoingLive
	dp.Automation = mod.Automation
	dp.AutomationName = mod.AutomationName.String
	dp.Score = mod.Score
	dp.Comments = mod.Comments.Ptr()
	if mod.Threshold.Valid {
		dp.Threshold = mod.Threshold.Float64
	} else {
		dp.Threshold = 0.0
	}
	dp.CreatedAt = mod.CreatedAt
	dp.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dp *DemandPartner) ToModel(id string) *models.Dpo {
	sort.SliceStable(dp.IntegrationType, func(i, j int) bool { return dp.IntegrationType[i] < dp.IntegrationType[j] })

	return &models.Dpo{
		DemandPartnerID:         id,
		DemandPartnerName:       dp.DemandPartnerName,
		ApprovalProcess:         dp.ApprovalProcess,
		DPBlocks:                dp.DPBlocks,
		PocName:                 dp.POCName,
		PocEmail:                dp.POCEmail,
		SeatOwnerID:             null.IntFromPtr(dp.SeatOwnerID),
		ManagerID:               null.IntFromPtr(dp.ManagerID),
		IntegrationType:         dp.IntegrationType,
		Active:                  dp.Active,
		IsInclude:               dp.IsInclude,
		IsApprovalNeeded:        dp.IsApprovalNeeded,
		ApprovalBeforeGoingLive: dp.ApprovalBeforeGoingLive,
		Automation:              dp.Automation,
		AutomationName: func() null.String {
			if dp.AutomationName == "" {
				return null.String{Valid: false, String: ""}
			}

			return null.StringFrom(dp.AutomationName)
		}(),
		Threshold: func() null.Float64 {
			if dp.Threshold == 0 {
				return null.Float64{Valid: false, Float64: 0}
			}

			return null.Float64From(dp.Threshold)
		}(),
		Score: func() int {
			if dp.Score == 0 {
				return DefaultDemandPartnerScoreValue
			}

			return dp.Score
		}(),
		Comments:  null.StringFromPtr(dp.Comments),
		CreatedAt: time.Now().UTC(),
	}
}

type SeatOwner struct {
	ID                       int        `json:"id"`
	SeatOwnerName            string     `json:"seat_owner_name"`
	SeatOwnerDomain          string     `json:"seat_owner_domain"`
	PublisherAccount         string     `json:"publisher_account"`
	CertificationAuthorityID string     `json:"certification_authority_id"`
	AdsTxtLine               string     `json:"ads_txt_line"`
	LineName                 string     `json:"line_name"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

func (so *SeatOwner) FromModel(mod *models.SeatOwner) {
	so.ID = mod.ID
	so.SeatOwnerName = mod.SeatOwnerName
	so.SeatOwnerDomain = mod.SeatOwnerDomain
	so.PublisherAccount = mod.PublisherAccount
	so.CertificationAuthorityID = mod.CertificationAuthorityID.String
	so.AdsTxtLine = buildAdsTxtLine(mod.SeatOwnerDomain, mod.PublisherAccount, mod.CertificationAuthorityID.String, true)
	so.LineName = buildLineName(mod.SeatOwnerName, "Direct")
	so.CreatedAt = mod.CreatedAt
	so.UpdatedAt = mod.UpdatedAt.Ptr()
}

type DemandPartnerChild struct {
	ID                       int        `json:"id"`
	DPConnectionID           int        `json:"dp_connection_id"`
	DPChildName              string     `json:"dp_child_name" validate:"required"`
	DPDomain                 string     `json:"dp_domain" validate:"required"`
	PublisherAccount         string     `json:"publisher_account" validate:"required"`
	CertificationAuthorityID *string    `json:"certification_authority_id"`
	IsRequiredForAdsTxt      bool       `json:"is_required_for_ads_txt"`
	IsDirect                 bool       `json:"is_direct"`
	AdsTxtLine               string     `json:"ads_txt_line"`
	LineName                 string     `json:"line_name"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

func (dpc *DemandPartnerChild) FromModel(demandPartnerName string, mod *models.DemandPartnerChild) {
	dpc.ID = mod.ID
	dpc.DPConnectionID = mod.DPConnectionID
	dpc.DPChildName = mod.DPChildName
	dpc.DPDomain = mod.DPDomain
	dpc.PublisherAccount = mod.PublisherAccount
	dpc.CertificationAuthorityID = mod.CertificationAuthorityID.Ptr()
	dpc.IsRequiredForAdsTxt = mod.IsRequiredForAdsTXT
	dpc.IsDirect = mod.IsDirect
	dpc.AdsTxtLine = buildAdsTxtLine(mod.DPDomain, mod.PublisherAccount, mod.CertificationAuthorityID.String, mod.IsDirect)
	dpc.LineName = buildLineName(demandPartnerName, mod.DPChildName)
	dpc.CreatedAt = mod.CreatedAt
	dpc.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dpc *DemandPartnerChild) ToModel(connectionID int) *models.DemandPartnerChild {
	return &models.DemandPartnerChild{
		ID:                       dpc.ID,
		DPConnectionID:           connectionID,
		DPChildName:              dpc.DPChildName,
		DPDomain:                 dpc.DPDomain,
		PublisherAccount:         dpc.PublisherAccount,
		CertificationAuthorityID: getCertificationAuthorityIDNullString(dpc.CertificationAuthorityID),
		IsDirect:                 dpc.IsDirect,
		IsRequiredForAdsTXT:      dpc.IsRequiredForAdsTxt,
		CreatedAt:                time.Now().UTC(),
	}
}

type DemandPartnerConnection struct {
	ID                       int                   `json:"id"`
	DemandPartnerID          string                `json:"demand_partner_id"`
	DPDomain                 string                `json:"dp_domain" validate:"required"`
	CertificationAuthorityID *string               `json:"certification_authority_id"`
	PublisherAccount         string                `json:"publisher_account" validate:"required"`
	MediaType                []string              `json:"media_type" validate:"mediaType"`
	IsDirect                 bool                  `json:"is_direct"`
	IsRequiredForAdsTxt      bool                  `json:"is_required_for_ads_txt"`
	Children                 []*DemandPartnerChild `json:"children"`
	AdsTxtLine               string                `json:"ads_txt_line"`
	LineName                 string                `json:"line_name"`
	CreatedAt                time.Time             `json:"created_at"`
	UpdatedAt                *time.Time            `json:"updated_at"`
}

func (dpc *DemandPartnerConnection) FromModel(demandPartnerName string, mod *models.DemandPartnerConnection) {
	mediaTypes := []string{}
	if len(mod.MediaType) > 0 {
		mediaTypes = mod.MediaType
	}

	dpc.ID = mod.ID
	dpc.DemandPartnerID = mod.DemandPartnerID
	dpc.CertificationAuthorityID = mod.CertificationAuthorityID.Ptr()
	dpc.DPDomain = mod.DPDomain
	dpc.PublisherAccount = mod.PublisherAccount
	dpc.MediaType = mediaTypes
	dpc.IsDirect = mod.IsDirect
	dpc.IsRequiredForAdsTxt = mod.IsRequiredForAdsTXT
	dpc.Children = func() []*DemandPartnerChild {
		children := make([]*DemandPartnerChild, 0, len(mod.R.DPConnectionDemandPartnerChildren))
		for _, modChild := range mod.R.DPConnectionDemandPartnerChildren {
			child := new(DemandPartnerChild)
			child.FromModel(demandPartnerName, modChild)
			children = append(children, child)
		}

		return children
	}()
	dpc.AdsTxtLine = buildAdsTxtLine(mod.DPDomain, mod.PublisherAccount, mod.CertificationAuthorityID.String, mod.IsDirect)
	dpc.LineName = buildLineName(demandPartnerName, demandPartnerName)
	dpc.CreatedAt = mod.CreatedAt
	dpc.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dpc *DemandPartnerConnection) ToModel(parentID string) *models.DemandPartnerConnection {
	sort.SliceStable(dpc.MediaType, func(i, j int) bool { return dpc.MediaType[i] < dpc.MediaType[j] })

	return &models.DemandPartnerConnection{
		ID:                       dpc.ID,
		DemandPartnerID:          parentID,
		DPDomain:                 dpc.DPDomain,
		CertificationAuthorityID: getCertificationAuthorityIDNullString(dpc.CertificationAuthorityID),
		PublisherAccount:         dpc.PublisherAccount,
		MediaType:                dpc.MediaType,
		IsDirect:                 dpc.IsDirect,
		IsRequiredForAdsTXT:      dpc.IsRequiredForAdsTxt,
		CreatedAt:                time.Now().UTC(),
	}
}

func buildAdsTxtLine(domain, publisherAccount, certificationAuthorityID string, isDirect bool) string {
	lineType := AdsTxtTypeReseller
	if isDirect {
		lineType = AdsTxtTypeDirect
	}

	adsTxtLine := fmt.Sprintf("%v, %v, %v", domain, strings.ReplaceAll(publisherAccount, "%s", "XXXXX"), lineType)

	if certificationAuthorityID != "" {
		adsTxtLine += fmt.Sprintf(", %v", certificationAuthorityID)
	}

	return adsTxtLine
}

func buildLineName(parent, child string) string {
	return fmt.Sprintf("%v - %v", parent, child)
}

func getCertificationAuthorityIDNullString(certificationAuthorityID *string) null.String {
	s := null.StringFromPtr(certificationAuthorityID)

	if s.Valid && s.String == "" {
		return null.NewString("", false)
	}

	return s
}
