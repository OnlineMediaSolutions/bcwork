package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildAdsTxtLine(t *testing.T) {
	t.Parallel()

	type args struct {
		domain                   string
		publisherAccount         string
		certificationAuthorityID string
		isDirect                 bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "withCertificationAuthorityID",
			args: args{
				domain:                   "domain.com",
				publisherAccount:         "12345",
				certificationAuthorityID: "abcd",
				isDirect:                 false,
			},
			want: "domain.com, 12345, RESELLER, abcd",
		},
		{
			name: "withoutCertificationAuthorityID",
			args: args{
				domain:           "domain.com",
				publisherAccount: "12345",
				isDirect:         false,
			},
			want: "domain.com, 12345, RESELLER",
		},
		{
			name: "lineForSeatOwner",
			args: args{
				domain:           "seatowner.com",
				publisherAccount: "9%s",
				isDirect:         true,
			},
			want: "seatowner.com, 9XXXXX, DIRECT",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildAdsTxtLine(tt.args.domain, tt.args.publisherAccount, tt.args.certificationAuthorityID, tt.args.isDirect)
			assert.Equal(t, tt.want, got)
		})
	}
}
