package core

import (
	"testing"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_getDemandPartnerColumnsToUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		newData *models.Dpo
		oldData *models.Dpo
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "updateAllFields",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                false,
					UpdatedAt:                null.TimeFrom(time.Now()),
					DemandPartnerName:        "demand_partner_new_name",
					Active:                   true,
					DPDomain:                 "demandnew.com",
					IsDirect:                 false,
					CertificationAuthorityID: null.StringFrom("abcd"),
					SeatOwnerID:              null.IntFrom(1),
					ManagerID:                null.IntFrom(2),
					IsApprovalNeeded:         false,
					Score:                    1,
					IsRequiredForAdsTXT:      false,
				},
				oldData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                true,
					UpdatedAt:                null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName:        "demand_partner_name",
					Active:                   false,
					DPDomain:                 "demand.com",
					IsDirect:                 true,
					CertificationAuthorityID: null.StringFrom("efgh"),
					SeatOwnerID:              null.IntFrom(2),
					ManagerID:                null.IntFrom(3),
					IsApprovalNeeded:         true,
					Score:                    2,
					IsRequiredForAdsTXT:      true,
				},
			},
			want: []string{
				models.DpoColumns.IsInclude,
				models.DpoColumns.UpdatedAt,
				models.DpoColumns.DemandPartnerName,
				models.DpoColumns.Active,
				models.DpoColumns.DPDomain,
				models.DpoColumns.IsDirect,
				models.DpoColumns.CertificationAuthorityID,
				models.DpoColumns.SeatOwnerID,
				models.DpoColumns.ManagerID,
				models.DpoColumns.IsApprovalNeeded,
				models.DpoColumns.Score,
				models.DpoColumns.IsRequiredForAdsTXT,
			},
		},
		{
			name: "updatePartialFields",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                false,
					UpdatedAt:                null.TimeFrom(time.Now()),
					DemandPartnerName:        "demand_partner_name",
					Active:                   true,
					DPDomain:                 "demand.com",
					IsDirect:                 false,
					CertificationAuthorityID: null.StringFrom("abcd"),
					SeatOwnerID:              null.IntFrom(1),
					ManagerID:                null.IntFrom(2),
					IsApprovalNeeded:         true,
					Score:                    1,
					IsRequiredForAdsTXT:      true,
				},
				oldData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                false,
					UpdatedAt:                null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName:        "demand_partner_name",
					Active:                   true,
					DPDomain:                 "demand.com",
					IsDirect:                 true,
					CertificationAuthorityID: null.StringFrom("abcd"),
					SeatOwnerID:              null.IntFrom(2),
					ManagerID:                null.IntFrom(3),
					IsApprovalNeeded:         true,
					Score:                    2,
					IsRequiredForAdsTXT:      true,
				},
			},
			want: []string{
				models.DpoColumns.UpdatedAt,
				models.DpoColumns.IsDirect,
				models.DpoColumns.SeatOwnerID,
				models.DpoColumns.ManagerID,
				models.DpoColumns.Score,
			},
		},
		{
			name: "noNewFieldsToUpdate",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                false,
					UpdatedAt:                null.TimeFrom(time.Now()),
					DemandPartnerName:        "demand_partner_new_name",
					Active:                   true,
					DPDomain:                 "demandnew.com",
					IsDirect:                 false,
					CertificationAuthorityID: null.StringFrom("abcd"),
					SeatOwnerID:              null.IntFrom(1),
					ManagerID:                null.IntFrom(2),
					IsApprovalNeeded:         false,
					Score:                    1,
					IsRequiredForAdsTXT:      false,
				},
				oldData: &models.Dpo{
					DemandPartnerID:          "id",
					IsInclude:                false,
					UpdatedAt:                null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName:        "demand_partner_new_name",
					Active:                   true,
					DPDomain:                 "demandnew.com",
					IsDirect:                 false,
					CertificationAuthorityID: null.StringFrom("abcd"),
					SeatOwnerID:              null.IntFrom(1),
					ManagerID:                null.IntFrom(2),
					IsApprovalNeeded:         false,
					Score:                    1,
					IsRequiredForAdsTXT:      false,
				},
			},
			want: []string{
				models.DpoColumns.UpdatedAt,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := getDemandPartnerColumnsToUpdate(tt.args.newData, tt.args.oldData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
