package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateDemandPartner(t *testing.T) {
	t.Parallel()

	type args struct {
		request *dto.DemandPartner
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					Threshold:         0.01,
					Automation:        true,
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					Children: []*dto.DemandPartnerChild{
						{
							DPChildName:      "child_name",
							DPChildDomain:    "childdomain.com",
							PublisherAccount: "12345",
						},
					},
					Connections: []*dto.DemandPartnerConnection{
						{PublisherAccount: "abcde"},
					},
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{},
		},
		{
			name: "invalid_noRequiredFields",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID: "id",
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				"DemandPartnerName is mandatory, validation failed",
				"DPDomain is mandatory, validation failed",
				"ManagerID is mandatory, validation failed",
			},
		},
		{
			name: "invalid_approvalProcessNotFromAllowedList",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					ApprovalProcess: "some_approval_process",
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				approvalProcessErrorMessage,
			},
		},
		{
			name: "invalid_dpBlocksNotFromAllowedList",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        "some_dp_blocks",
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				dpBlocksErrorMessage,
			},
		},
		{
			name: "invalid_noRequiredFieldsForChild",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					Children: []*dto.DemandPartnerChild{
						{},
					},
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				"Children: DPChildName is mandatory, validation failed",
				"Children: DPChildDomain is mandatory, validation failed",
				"Children: PublisherAccount is mandatory, validation failed",
			},
		},
		{
			name: "invalid_noRequiredFieldsForConnection",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					Children: []*dto.DemandPartnerChild{
						{
							DPChildName:      "child_name",
							DPChildDomain:    "childdomain.com",
							PublisherAccount: "12345",
						},
					},
					Connections: []*dto.DemandPartnerConnection{
						{},
					},
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				"Connections: PublisherAccount is mandatory, validation failed",
			},
		},
		{
			name: "invalid automation values",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					Threshold:         10,
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					Children: []*dto.DemandPartnerChild{
						{
							DPChildName:      "child_name",
							DPChildDomain:    "childdomain.com",
							PublisherAccount: "12345",
						},
					},
					Connections: []*dto.DemandPartnerConnection{
						{PublisherAccount: "abcde"},
					},
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{
				"dp threshold must be >= 0.00 and <= 0.01",
			},
		},
		{
			name: "valid automation values",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:   "id",
					DemandPartnerName: "name",
					DPDomain:          "domain.com",
					Threshold:         0.001,
					Automation:        true,
					CertificationAuthorityID: func() *string {
						s := "cert_id"
						return &s
					}(),
					Children: []*dto.DemandPartnerChild{
						{
							DPChildName:      "child_name",
							DPChildDomain:    "childdomain.com",
							PublisherAccount: "12345",
						},
					},
					Connections: []*dto.DemandPartnerConnection{
						{PublisherAccount: "abcde"},
					},
					ApprovalProcess: dto.GDocApprovalProcess,
					DPBlocks:        dto.EmailApprovalProcess,
					POCName:         "poc_name",
					POCEmail:        "poc_email",
					SeatOwnerID: func() *int {
						n := 1
						return &n
					}(),
					ManagerID: func() *int {
						n := 1
						return &n
					}(),
					IsInclude:               false,
					Active:                  true,
					IsDirect:                true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					IsRequiredForAdsTxt:     true,
					Score:                   5,
					Comments: func() *string {
						s := "comments"
						return &s
					}(),
				},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateDemandPartner(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
