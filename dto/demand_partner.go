package dto

import (
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/volatiletech/null/v8"
)

const DefaultDemandPartnerScoreValue = 1000

// TODO: add validations
type DemandPartner struct {
	DemandPartnerID          string                     `json:"demand_partner_id"`
	DemandPartnerName        string                     `json:"demand_partner_name"`
	DPDomain                 string                     `json:"dp_domain"`
	Children                 []*DemandPartnerChild      `json:"children"`
	Connections              []*DemandPartnerConnection `json:"connection"`
	CertificationAuthorityID *string                    `json:"certification_authority_id"`
	ApprovalProcess          string                     `json:"approval_process"`
	DPBlocks                 string                     `json:"dp_blocks"`
	POCName                  string                     `json:"poc_name"`
	POCEmail                 string                     `json:"poc_email"`
	SeatOwnerID              *int                       `json:"seat_owner_id"`
	ManagerID                *int                       `json:"manager_id"`
	IsInclude                bool                       `json:"is_include"`
	Active                   bool                       `json:"active"`
	IsDirect                 bool                       `json:"is_direct"`
	IsApprovalNeeded         bool                       `json:"is_approval_needed"`
	ApprovalBeforeGoingLive  bool                       `json:"approval_before_going_live"`
	IsRequiredForAdsTxt      bool                       `json:"is_required_for_ads_txt"`
	Score                    int                        `json:"score"`
	Comments                 *string                    `json:"comments"`
	CreatedAt                time.Time                  `json:"created_at"`
	UpdatedAt                *time.Time                 `json:"updated_at"`
}

func (dp *DemandPartner) FromModel(mod *models.Dpo) {
	dp.DemandPartnerID = mod.DemandPartnerID
	dp.DemandPartnerName = mod.DemandPartnerName
	dp.DPDomain = mod.DPDomain
	dp.Children = func() []*DemandPartnerChild {
		children := make([]*DemandPartnerChild, 0, len(mod.R.DPParentDemandPartnerChildren))
		for _, modChild := range mod.R.DPParentDemandPartnerChildren {
			child := new(DemandPartnerChild)
			child.FromModel(modChild)
			children = append(children, child)
		}

		return children
	}()
	dp.Connections = func() []*DemandPartnerConnection {
		connections := make([]*DemandPartnerConnection, 0, len(mod.R.DemandPartnerDemandPartnerConnections))
		for _, modConnection := range mod.R.DemandPartnerDemandPartnerConnections {
			connection := new(DemandPartnerConnection)
			connection.FromModel(modConnection)
			connections = append(connections, connection)
		}

		return connections
	}()
	dp.CertificationAuthorityID = mod.CertificationAuthorityID.Ptr()
	dp.ApprovalProcess = mod.ApprovalProcess
	dp.DPBlocks = mod.DPBlocks
	dp.POCName = mod.PocName
	dp.POCEmail = mod.PocEmail
	dp.SeatOwnerID = mod.SeatOwnerID.Ptr()
	dp.ManagerID = mod.ManagerID.Ptr()
	dp.Active = mod.Active
	dp.IsInclude = mod.IsInclude
	dp.IsDirect = mod.IsDirect
	dp.IsApprovalNeeded = mod.IsApprovalNeeded
	dp.ApprovalBeforeGoingLive = mod.ApprovalBeforeGoingLive
	dp.IsRequiredForAdsTxt = mod.IsRequiredForAdsTXT
	dp.Score = mod.Score
	dp.Comments = mod.Comments.Ptr()
	dp.CreatedAt = mod.CreatedAt
	dp.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dp *DemandPartner) ToModel(id string) *models.Dpo {
	return &models.Dpo{
		DemandPartnerID:          id,
		DemandPartnerName:        dp.DemandPartnerName,
		DPDomain:                 dp.DPDomain,
		CertificationAuthorityID: null.StringFromPtr(dp.CertificationAuthorityID),
		ApprovalProcess:          dp.ApprovalProcess,
		DPBlocks:                 dp.DPBlocks,
		PocName:                  dp.POCName,
		PocEmail:                 dp.POCEmail,
		SeatOwnerID:              null.IntFromPtr(dp.SeatOwnerID),
		ManagerID:                null.IntFromPtr(dp.ManagerID),
		Active:                   dp.Active,
		IsInclude:                dp.IsInclude,
		IsDirect:                 dp.IsDirect,
		IsApprovalNeeded:         dp.IsApprovalNeeded,
		ApprovalBeforeGoingLive:  dp.ApprovalBeforeGoingLive,
		IsRequiredForAdsTXT:      dp.IsRequiredForAdsTxt,
		Score: func(s int) int {
			if s == 0 {
				return DefaultDemandPartnerScoreValue
			}
			return s
		}(dp.Score),
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
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

func (so *SeatOwner) FromModel(mod *models.SeatOwner) {
	so.ID = mod.ID
	so.SeatOwnerName = mod.SeatOwnerName
	so.SeatOwnerDomain = mod.SeatOwnerDomain
	so.PublisherAccount = mod.PublisherAccount
	so.CertificationAuthorityID = mod.CertificationAuthorityID.String
	so.CreatedAt = mod.CreatedAt
	so.UpdatedAt = mod.UpdatedAt.Ptr()
}

type DemandPartnerChild struct {
	ID                       int        `json:"id"`
	ParentID                 string     `json:"parent_id"`
	DPChildName              string     `json:"dp_child_name"`
	DPChildDomain            string     `json:"dp_child_domain"`
	PublisherAccount         string     `json:"publisher_account"`
	CertificationAuthorityID *string    `json:"certification_authority_id"`
	IsRequiredForAdsTxt      bool       `json:"is_required_for_ads_txt"`
	Active                   bool       `json:"active"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`
}

func (dpc *DemandPartnerChild) FromModel(mod *models.DemandPartnerChild) {
	dpc.ID = mod.ID
	dpc.ParentID = mod.DPParentID
	dpc.DPChildName = mod.DPChildName
	dpc.DPChildDomain = mod.DPChildDomain
	dpc.PublisherAccount = mod.PublisherAccount
	dpc.CertificationAuthorityID = mod.CertificationAuthorityID.Ptr()
	dpc.IsRequiredForAdsTxt = mod.IsRequiredForAdsTXT
	dpc.Active = mod.Active
	dpc.CreatedAt = mod.CreatedAt
	dpc.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dpc *DemandPartnerChild) ToModel(parentID string) *models.DemandPartnerChild {
	return &models.DemandPartnerChild{
		ID:                       dpc.ID,
		DPParentID:               parentID,
		DPChildName:              dpc.DPChildName,
		DPChildDomain:            dpc.DPChildDomain,
		PublisherAccount:         dpc.PublisherAccount,
		CertificationAuthorityID: null.StringFromPtr(dpc.CertificationAuthorityID),
		Active:                   true,
		IsRequiredForAdsTXT:      dpc.IsRequiredForAdsTxt,
		CreatedAt:                time.Now().UTC(),
	}
}

type DemandPartnerConnection struct {
	ID               int        `json:"id"`
	DemandPartnerID  string     `json:"demand_partner_id"`
	PublisherAccount string     `json:"publisher_account"`
	IntegrationType  []string   `json:"integration_type"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

func (dpc *DemandPartnerConnection) FromModel(mod *models.DemandPartnerConnection) {
	dpc.ID = mod.ID
	dpc.DemandPartnerID = mod.DemandPartnerID
	dpc.PublisherAccount = mod.PublisherAccount
	dpc.IntegrationType = mod.IntegrationType
	dpc.Active = mod.Active
	dpc.CreatedAt = mod.CreatedAt
	dpc.UpdatedAt = mod.UpdatedAt.Ptr()
}

func (dpc *DemandPartnerConnection) ToModel(parentID string) *models.DemandPartnerConnection {
	return &models.DemandPartnerConnection{
		ID:               dpc.ID,
		DemandPartnerID:  parentID,
		PublisherAccount: dpc.PublisherAccount,
		IntegrationType:  dpc.IntegrationType,
		Active:           true,
		CreatedAt:        time.Now().UTC(),
	}
}
