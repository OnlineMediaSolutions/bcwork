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
		data *AdsTxtGroupedByDP
		args args
		want *AdsTxtGroupedByDP
	}{
		{
			name: "alreadyHasMainLine",
			data: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
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
			want: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
		},
		{
			name: "secondaryLineRewritedWithMainLine",
			data: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
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
			want: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
		},
		{
			name: "secondaryLineRewritedBySeatOwner",
			data: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
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
			want: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
					DemandPartnerName:         "SeatOwner",
					DemandPartnerNameExtended: "SeatOwner - Direct",
				},
			},
		},
		{
			name: "seatOwnerLineRewritedByMainLine",
			data: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
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
			want: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
					DemandPartnerName:         "Primary",
					DemandPartnerNameExtended: "Primary - Primary",
				},
			},
		},
		{
			name: "secondaryLineNotChangedWithOtherSecondaryLine",
			data: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
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
			want: &AdsTxtGroupedByDP{
				AdsTxt: &AdsTxt{
					DemandPartnerName:         "Secondary",
					DemandPartnerNameExtended: "Primary - Secondary",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.data.ProcessParentRow(tt.args.row)
			// object mutates
			assert.Equal(t, tt.want, tt.data)
		})
	}
}
