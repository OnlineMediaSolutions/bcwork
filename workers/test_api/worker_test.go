package testapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prepareData(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				data: []byte(`{"key":"tech_fee","created_at":"2024-09-17T09:16:53.587236Z","value":5.01,"updated_at":"2024-09-24T13:41:43.7Z"}`),
			},
			want: `{"key":"tech_fee","value":5.01}`,
		},
		{
			name: "nothingToRemove",
			args: args{
				data: []byte(`{"key":"tech_fee","value":5.01}`),
			},
			want: `{"key":"tech_fee","value":5.01}`,
		},
		{
			name: "removedFromBeginningAndEnding",
			args: args{
				data: []byte(`{"created_at":"2024-09-17T09:16:53.587236Z","key":"tech_fee","value":5.01,"updated_at":"2024-09-24T13:41:43.7Z"}`),
			},
			want: `{"key":"tech_fee","value":5.01}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareData(tt.args.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_prepareMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		report [][]string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				report: [][]string{
					{"endpoint_1", "error_1"},
					{"endpoint_2", "error_2"},
					{"endpoint_3", "error_3"},
				},
			},
			want: "Test API worker. Failed tests:\n1. endpoint_1: error_1\n2. endpoint_2: error_2\n3. endpoint_3: error_3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := prepareMessage(tt.args.report)
			assert.Equal(t, tt.want, got)
		})
	}
}
