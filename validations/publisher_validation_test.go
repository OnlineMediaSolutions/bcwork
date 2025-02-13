package validations

import (
	"testing"

	"github.com/m6yf/bcwork/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validatePublisher(t *testing.T) {
	t.Parallel()

	type args struct {
		request interface{}
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "valid_createRequest",
			args: args{
				request: &dto.PublisherCreateValues{
					Name:            "publisher",
					IntegrationType: []string{dto.ORTBIntergrationType},
					MediaType:       []string{dto.VideoMediaType},
				},
			},
			want: []string{},
		},
		{
			name: "invalid_createRequest",
			args: args{
				request: &dto.PublisherCreateValues{},
			},
			want: []string{
				"Name is mandatory, validation failed",
				// "integration type must be in allowed list: oRTB,Prebid Server,Amazon APS",
				// "media type must be in allowed list: Web Banners,Video,InApp",
			},
		},
		{
			name: "valid_updateRequest",
			args: args{
				request: &dto.UpdatePublisherValues{
					Name: func() *string {
						s := "new_publisher"
						return &s
					}(),
					IntegrationType: []string{dto.ORTBIntergrationType},
					MediaType:       []string{dto.VideoMediaType},
				},
			},
			want: []string{},
		},
		{
			name: "valid_emptyUpdateRequest",
			args: args{
				request: &dto.UpdatePublisherValues{},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validatePublisher(tt.args.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
