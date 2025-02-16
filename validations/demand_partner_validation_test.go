package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_validateDemandPartner(t *testing.T) {
	t.Parallel()

	const (
		comments = "comments"
		certID   = "cert_id"
	)

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
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                0.001,
					Automation:               true,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
				},
			},
			want: []string{},
		},
		{
			name: "invalid_noRequiredFields",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:          "id",
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					ApprovalProcess:          dto.GDocApprovalProcess,
					DPBlocks:                 dto.EmailApprovalProcess,
					POCName:                  "poc_name",
					POCEmail:                 "poc_email",
					SeatOwnerID:              helpers.GetPointerToInt(1),
					IsInclude:                false,
					Active:                   true,
					IsApprovalNeeded:         true,
					ApprovalBeforeGoingLive:  true,
					Score:                    5,
					Comments:                 helpers.GetPointerToString(certID),
				},
			},
			want: []string{
				"DemandPartnerName is mandatory, validation failed",
				"DPDomain is mandatory, validation failed",
				"ManagerID is mandatory, validation failed",
				"integration type must be in allowed list: oRTB,Prebid Server,Amazon APS",
			},
		},
		{
			name: "invalid_approvalProcessNotFromAllowedList",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					Threshold:                0.001,
					DPDomain:                 "domain.com",
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					ApprovalProcess:          "some_approval_process",
					DPBlocks:                 dto.EmailApprovalProcess,
					POCName:                  "poc_name",
					POCEmail:                 "poc_email",
					SeatOwnerID:              helpers.GetPointerToInt(1),
					ManagerID:                helpers.GetPointerToInt(1),
					IsInclude:                false,
					Active:                   true,
					IsApprovalNeeded:         true,
					ApprovalBeforeGoingLive:  true,
					Score:                    5,
					Comments:                 helpers.GetPointerToString(comments),
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
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					Threshold:                0.001,
					DPDomain:                 "domain.com",
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					ApprovalProcess:          dto.GDocApprovalProcess,
					DPBlocks:                 "some_dp_blocks",
					POCName:                  "poc_name",
					POCEmail:                 "poc_email",
					SeatOwnerID:              helpers.GetPointerToInt(1),
					ManagerID:                helpers.GetPointerToInt(1),
					IsInclude:                false,
					Active:                   true,
					IsApprovalNeeded:         true,
					ApprovalBeforeGoingLive:  true,
					Score:                    5,
					Comments:                 helpers.GetPointerToString(comments),
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
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                0.001,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
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
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					Threshold:                0.001,
					DPDomain:                 "domain.com",
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
				},
			},
			want: []string{
				"Connections: PublisherAccount is mandatory, validation failed",
				"media type must be in allowed list: Web Banners,Video,InApp",
			},
		},
		{
			name: "invalid automation values",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                10,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
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
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                0.001,
					Automation:               true,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
				},
			},
			want: []string{},
		},
		{
			name: "test valid min threshold",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                0,
					Automation:               true,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
				},
			},
			want: []string{},
		},
		{
			name: "test valid max threshold",
			args: args{
				request: &dto.DemandPartner{
					DemandPartnerID:          "id",
					DemandPartnerName:        "name",
					DPDomain:                 "domain.com",
					Threshold:                0.010,
					Automation:               true,
					CertificationAuthorityID: helpers.GetPointerToString(certID),
					IntegrationType:          []string{dto.PrebidServerIntergrationType},
					Connections: []*dto.DemandPartnerConnection{
						{
							PublisherAccount: "abcde",
							Children: []*dto.DemandPartnerChild{
								{
									DPChildName:      "child_name",
									DPChildDomain:    "childdomain.com",
									PublisherAccount: "12345",
								},
							},
							MediaType: []string{dto.WebBannersMediaType},
						},
					},
					ApprovalProcess:         dto.GDocApprovalProcess,
					DPBlocks:                dto.EmailApprovalProcess,
					POCName:                 "poc_name",
					POCEmail:                "poc_email",
					SeatOwnerID:             helpers.GetPointerToInt(1),
					ManagerID:               helpers.GetPointerToInt(1),
					IsInclude:               false,
					Active:                  true,
					IsApprovalNeeded:        true,
					ApprovalBeforeGoingLive: true,
					Score:                   5,
					Comments:                helpers.GetPointerToString(comments),
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
