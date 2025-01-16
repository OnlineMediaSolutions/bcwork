package export

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExportModule_ExportCSV(t *testing.T) {
	t.Parallel()

	type args struct {
		srcs []json.RawMessage
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid_sameObjects",
			args: args{
				srcs: []json.RawMessage{
					json.RawMessage(`{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`),
					json.RawMessage(`{"id": 2, "name": "publisher_2", "active": false, "domain": "2.com", "factor": 0.05}`),
					json.RawMessage(`{"id": 3, "name": "publisher_3", "active": true, "domain": "3.com", "factor": 0.1}`),
				},
			},
			want: []byte(
				"id,name,active,domain,factor\n" +
					"1,publisher_1,true,1.com,0.01\n" +
					"2,publisher_2,false,2.com,0.05\n" +
					"3,publisher_3,true,3.com,0.1\n",
			),
		},
		{
			name: "valid_differentPositionsOfKeys",
			args: args{
				srcs: []json.RawMessage{
					json.RawMessage(`{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`),
					json.RawMessage(`{"name": "publisher_2", "factor": 0.05, "active": false, "id": 2, "domain": "2.com"}`),
					json.RawMessage(`{"domain": "3.com", "id": 3, "active": true, "factor": 0.1, "name": "publisher_3"}`),
				},
			},
			want: []byte(
				"id,name,active,domain,factor\n" +
					"1,publisher_1,true,1.com,0.01\n" +
					"2,publisher_2,false,2.com,0.05\n" +
					"3,publisher_3,true,3.com,0.1\n",
			),
		},
		{
			name: "valid_oneRow",
			args: args{
				srcs: []json.RawMessage{
					json.RawMessage(`{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`),
				},
			},
			want: []byte(
				"id,name,active,domain,factor\n" +
					"1,publisher_1,true,1.com,0.01\n",
			),
		},
		{
			name: "whenDifferentObjects_thenReturnError",
			args: args{
				srcs: []json.RawMessage{
					json.RawMessage(`{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`),
					json.RawMessage(`{"id": 2, "name": "publisher_2", "active": false, "domain": "2.com"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "whenBadJSON_thenReturnError",
			args: args{
				srcs: []json.RawMessage{
					json.RawMessage(`{"id": 1, "name": "publisher_1", "active": true, "domain": "1.com", "factor": 0.01}`),
					json.RawMessage(`{"id": 2, "name": "publisher_2", "active": false:: "domain": "2.com"`),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewExportModule()

			got, err := e.ExportCSV(context.Background(), tt.args.srcs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
