package testapi

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	httpmocks "github.com/m6yf/bcwork/modules/http_client/mocks"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/modules/messager/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Worker_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		worker  *Worker
		wantErr bool
	}{
		{
			name: "valid",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)

				return &Worker{
					BaseURL: "localhost",
					cases: []testCase{
						{
							Name:     "TestPing",
							Endpoint: "/ping",
							Method:   http.MethodGet,
							Payload:  "",
							Want:     `{"status":"OK"}`,
						},
					},
					httpClient: func() httpclient.Doer {
						return httpmocks.NewHttpClientMock(ctrl).
							DoMock.
							Expect(
								minimock.AnyContext,
								http.MethodGet,
								"localhost/ping",
								strings.NewReader(""),
							).
							Return([]byte(`{"status":"OK"}`), 200, nil)
					}(),
				}
			}(),
			wantErr: false,
		},
		{
			name: "wrongResponse",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)

				return &Worker{
					BaseURL: "localhost",
					cases: []testCase{
						{
							Name:     "TestPublisherGet",
							Endpoint: "/publisher/get",
							Method:   http.MethodPost,
							Payload:  `{"filter":{"publisher_id":["9995"]}}`,
							Want:     `[{"publisher_id":"9995"}]`,
						},
					},
					httpClient: func() httpclient.Doer {
						return httpmocks.NewHttpClientMock(ctrl).
							DoMock.
							Expect(
								minimock.AnyContext,
								http.MethodPost,
								"localhost/publisher/get",
								strings.NewReader(`{"filter":{"publisher_id":["9995"]}}`),
							).
							Return([]byte(`{"status":"error"}`), 500, nil)
					}(),
					messager: func() messager.Messager {
						return mocks.NewMessagerMock(ctrl).
							SendMessageMock.
							Expect("*Test API worker. Failed tests:*\n1. _TestPublisherGet [POST /publisher/get]_: ```not equal:\ngot  = {\"status\":\"error\"}\nwant = [{\"publisher_id\":\"9995\"}]```").
							Return(nil)
					}(),
				}
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.worker.Do(context.Background())
			if tt.wantErr {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)
		})
	}
}

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
			want: "*Test API worker. Failed tests:*\n1. _endpoint_1_: ```error_1```\n2. _endpoint_2_: ```error_2```\n3. _endpoint_3_: ```error_3```",
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
