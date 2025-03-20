package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProcessParentLine(t *testing.T) {
	t.Parallel()

	type args struct {
		row *AdsTxt
	}

	tests := []struct {
		name string
		data *AdsTxtGroupedByDPData
		args args
		want *AdsTxt
	}{
		{
			name: "alreadyHasMainLine",
			data: &AdsTxtGroupedByDPData{
				Parent: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
			args: args{
				row: &AdsTxt{
					DemandPartnerName:         "Secondary",
					DemandPartnerNameExtended: "Primary - Secondary",
				},
			},
			want: &AdsTxt{
				DemandPartnerName:         "Primary",
				DemandPartnerNameExtended: "Primary - Primary",
			},
		},
		{
			name: "secondaryLineRewritedWithMainLine",
			data: &AdsTxtGroupedByDPData{
				Parent: &AdsTxt{
					DemandPartnerName:         "Secondary",
					DemandPartnerNameExtended: "Primary - Secondary",
				},
			},
			args: args{
				row: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
			want: &AdsTxt{
				DemandPartnerName:         "Primary",
				DemandPartnerNameExtended: "Primary - Primary",
			},
		},
		{
			name: "secondaryLineRewritedBySeatOwner",
			data: &AdsTxtGroupedByDPData{
				Parent: &AdsTxt{
					DemandPartnerName:         "Secondary",
					DemandPartnerNameExtended: "Primary - Secondary",
				},
			},
			args: args{
				row: &AdsTxt{
					DemandPartnerName:         "SeatOwner",
					DemandPartnerNameExtended: "SeatOwner - Direct",
				},
			},
			want: &AdsTxt{
				DemandPartnerName:         "SeatOwner",
				DemandPartnerNameExtended: "SeatOwner - Direct",
			},
		},
		{
			name: "seatOwnerLineRewritedByMainLine",
			data: &AdsTxtGroupedByDPData{
				Parent: &AdsTxt{
					DemandPartnerName:         "SeatOwner",
					DemandPartnerNameExtended: "SeatOwner - Direct",
				},
			},
			args: args{
				row: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
			want: &AdsTxt{
				DemandPartnerName:         "Primary",
				DemandPartnerNameExtended: "Primary - Primary",
			},
		},
		{
			name: "secondaryLineNotChangedWithOtherSecondaryLine",
			data: &AdsTxtGroupedByDPData{
				Parent: &AdsTxt{
					DemandPartnerName:         "Secondary",
					DemandPartnerNameExtended: "Primary - Secondary",
				},
			},
			args: args{
				row: &AdsTxt{
					DemandPartnerName:         "Tertiary",
					DemandPartnerNameExtended: "Primary - Tertiary",
				},
			},
			want: &AdsTxt{
				DemandPartnerName:         "Secondary",
				DemandPartnerNameExtended: "Primary - Secondary",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.data.ProcessParentRow(tt.args.row)
			// object mutates
			assert.Equal(t, tt.want, tt.data.Parent)
		})
	}
}
