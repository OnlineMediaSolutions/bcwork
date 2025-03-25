package dto

import (
	"testing"

	"github.com/m6yf/bcwork/utils/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
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

func Test_getCertificationAuthorityIDNullString(t *testing.T) {
	t.Parallel()

	type args struct {
		certificationAuthorityID *string
	}

	tests := []struct {
		name string
		args args
		want null.String
	}{
		{
			name: "usualCertificationAuthorityID",
			args: args{
				certificationAuthorityID: helpers.GetPointerToString("a1b2c3d4"),
			},
			want: null.String{String: "a1b2c3d4", Valid: true},
		},
		{
			name: "emptyString",
			args: args{
				certificationAuthorityID: helpers.GetPointerToString(""),
			},
			want: null.String{String: "", Valid: false},
		},
		{
			name: "nilPointer",
			args: args{
				certificationAuthorityID: nil,
			},
			want: null.String{String: "", Valid: false},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getCertificationAuthorityIDNullString(tt.args.certificationAuthorityID)
			assert.Equal(t, tt.want, got)
		})
	}
}
