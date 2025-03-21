package core

import (
	"testing"
	"time"

	"github.com/m6yf/bcwork/models"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func Test_getModelsColumnsToUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		oldData          any
		newData          any
		blacklistColumns []string
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "dpo_updateAllFields",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         false,
					UpdatedAt:         null.TimeFrom(time.Now()),
					DemandPartnerName: "demand_partner_new_name",
					Active:            true,
					SeatOwnerID:       null.IntFrom(1),
					ManagerID:         null.IntFrom(2),
					IsApprovalNeeded:  false,
					Score:             1,
				},
				oldData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         true,
					UpdatedAt:         null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName: "demand_partner_name",
					Active:            false,
					SeatOwnerID:       null.IntFrom(2),
					ManagerID:         null.IntFrom(3),
					IsApprovalNeeded:  true,
					Score:             2,
				},
				blacklistColumns: []string{
					models.DpoColumns.DemandPartnerID,
					models.DpoColumns.CreatedAt,
				},
			},
			want: []string{
				models.DpoColumns.IsInclude,
				models.DpoColumns.UpdatedAt,
				models.DpoColumns.DemandPartnerName,
				models.DpoColumns.Active,
				models.DpoColumns.SeatOwnerID,
				models.DpoColumns.ManagerID,
				models.DpoColumns.IsApprovalNeeded,
				models.DpoColumns.Score,
			},
		},
		{
			name: "dpo_updatePartialFields",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         false,
					UpdatedAt:         null.TimeFrom(time.Now()),
					DemandPartnerName: "demand_partner_name",
					Active:            true,
					SeatOwnerID:       null.IntFrom(1),
					ManagerID:         null.IntFrom(2),
					IsApprovalNeeded:  true,
					Score:             1,
				},
				oldData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         false,
					UpdatedAt:         null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName: "demand_partner_name",
					Active:            true,
					SeatOwnerID:       null.IntFrom(2),
					ManagerID:         null.IntFrom(3),
					IsApprovalNeeded:  true,
					Score:             2,
				},
				blacklistColumns: []string{
					models.DpoColumns.DemandPartnerID,
					models.DpoColumns.CreatedAt,
				},
			},
			want: []string{
				models.DpoColumns.UpdatedAt,
				models.DpoColumns.SeatOwnerID,
				models.DpoColumns.ManagerID,
				models.DpoColumns.Score,
			},
		},
		{
			name: "dponoNewFieldsToUpdate",
			args: args{
				newData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         false,
					UpdatedAt:         null.TimeFrom(time.Now()),
					DemandPartnerName: "demand_partner_new_name",
					Active:            true,
					SeatOwnerID:       null.IntFrom(1),
					ManagerID:         null.IntFrom(2),
					IsApprovalNeeded:  false,
					Score:             1,
				},
				oldData: &models.Dpo{
					DemandPartnerID:   "id",
					IsInclude:         false,
					UpdatedAt:         null.TimeFrom(time.Now().Add(time.Hour * -24)),
					DemandPartnerName: "demand_partner_new_name",
					Active:            true,
					SeatOwnerID:       null.IntFrom(1),
					ManagerID:         null.IntFrom(2),
					IsApprovalNeeded:  false,
					Score:             1,
				},
				blacklistColumns: []string{
					models.DpoColumns.DemandPartnerID,
					models.DpoColumns.CreatedAt,
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

			got, err := getModelsColumnsToUpdate(tt.args.oldData, tt.args.newData, tt.args.blacklistColumns)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
