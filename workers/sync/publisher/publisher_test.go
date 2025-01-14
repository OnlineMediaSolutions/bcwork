package publisher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gojuno/minimock/v3"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/storage/db"
	dbmocks "github.com/m6yf/bcwork/storage/db/mocks"
	s3storage "github.com/m6yf/bcwork/storage/s3_storage"
	s3mocks "github.com/m6yf/bcwork/storage/s3_storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Test_Worker_Do(t *testing.T) {
	t.Parallel()

	now := time.Now()

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
					Bucket:     "bucket",
					Prefix:     "prefix",
					DaysBefore: -2,
					S3: func() s3storage.S3 {
						return s3mocks.NewS3Mock(ctrl).
							ListS3ObjectsMock.
							Expect("bucket", "prefix").
							Return(&s3.ListObjectsV2Output{
								Contents: []*s3.Object{
									{
										Key: func() *string {
											key := "key"
											return &key
										}(),
										LastModified: &now,
									},
								},
							}, nil).
							GetObjectInputMock.
							Expect("bucket", "key").
							Return([]byte(`[{"_id":"id", "accountManager": {"id":"am_id"}, "site": ["1.com"]}]`), nil)
					}(),
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(false).
							UpsertPublisherMock.
							Expect(
								minimock.AnyContext,
								&models.Publisher{
									PublisherID:      "id",
									AccountManagerID: null.String{Valid: true, String: "am_id"},
								},
								boil.Blacklist(
									models.PublisherColumns.CreatedAt,
									models.PublisherColumns.Status,
									models.PublisherColumns.IntegrationType,
									models.PublisherColumns.CampaignManagerID,
									models.PublisherColumns.OfficeLocation,
									models.PublisherColumns.ReactivateTimestamp,
									models.PublisherColumns.StartTimestamp,
								),
							).
							Return(nil).
							InsertPublisherDomainMock.
							Expect(
								minimock.AnyContext,
								&models.PublisherDomain{
									PublisherID: "id",
									Domain:      "1.com",
								},
							).
							Return(nil).
							SaveResultOfLastSyncMock.
							Expect(minimock.AnyContext, "key", false).
							Return(nil)
					}(),
				}
			}(),
			wantErr: false,
		},
		{
			name: "nothingToUpdate",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)
				return &Worker{
					Bucket:     "bucket",
					Prefix:     "prefix",
					DaysBefore: -2,
					S3: func() s3storage.S3 {
						return s3mocks.NewS3Mock(ctrl).
							ListS3ObjectsMock.
							Expect("bucket", "prefix").
							Return(&s3.ListObjectsV2Output{
								Contents: []*s3.Object{
									{
										Key: func() *string {
											key := "key"
											return &key
										}(),
										LastModified: func() *time.Time {
											t := now.AddDate(0, 0, -5)
											return &t
										}(),
									},
								},
							}, nil)
					}(),
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(false)
					}(),
				}
			}(),
			wantErr: false,
		},
		{
			name: "errorWhileProcessingObject",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)
				return &Worker{
					Bucket:     "bucket",
					Prefix:     "prefix",
					DaysBefore: -2,
					S3: func() s3storage.S3 {
						return s3mocks.NewS3Mock(ctrl).
							ListS3ObjectsMock.
							Expect("bucket", "prefix").
							Return(&s3.ListObjectsV2Output{
								Contents: []*s3.Object{
									{
										Key: func() *string {
											key := "key"
											return &key
										}(),
										LastModified: &now,
									},
								},
							}, nil).
							GetObjectInputMock.
							Expect("bucket", "key").
							Return([]byte(`[{"_id":"id", "accountManager": {"id":"am_id"}, "site": ["1.com"]}]`), nil)
					}(),
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(false).
							UpsertPublisherMock.
							Expect(
								minimock.AnyContext,
								&models.Publisher{
									PublisherID:      "id",
									AccountManagerID: null.String{Valid: true, String: "am_id"},
								},
								boil.Blacklist(
									models.PublisherColumns.CreatedAt,
									models.PublisherColumns.Status,
									models.PublisherColumns.IntegrationType,
									models.PublisherColumns.CampaignManagerID,
									models.PublisherColumns.OfficeLocation,
									models.PublisherColumns.ReactivateTimestamp,
									models.PublisherColumns.StartTimestamp,
								),
							).
							Return(nil).
							InsertPublisherDomainMock.
							Expect(
								minimock.AnyContext,
								&models.PublisherDomain{
									PublisherID: "id",
									Domain:      "1.com",
								},
							).
							Return(errors.New("error while inserting publisher domain")).
							SaveResultOfLastSyncMock.
							Expect(minimock.AnyContext, "key", true).
							Return(nil)
					}(),
				}
			}(),
			wantErr: false,
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

func Test_Worker_isNeededToUpdate(t *testing.T) {
	t.Parallel()

	now := time.Now()

	type args struct {
		key          string
		lastModified *time.Time
	}

	tests := []struct {
		name   string
		worker *Worker
		args   args
		want   bool
	}{
		{
			name: "needToUpdate_becauseOfLastModifiedTime",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)
				return &Worker{
					DaysBefore: -2,
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(false)
					}(),
				}
			}(),
			args: args{
				key: "key",
				lastModified: func() *time.Time {
					t := now.AddDate(0, 0, -1)
					return &t
				}(),
			},
			want: true,
		},
		{
			name: "needToUpdate_becauseOfHadErrorLastTime",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)
				return &Worker{
					DaysBefore: -2,
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(true)
					}(),
				}
			}(),
			args: args{
				key: "key",
				lastModified: func() *time.Time {
					t := now.AddDate(0, 0, -3)
					return &t
				}(),
			},
			want: true,
		},
		{
			name: "noNeedToUpdate",
			worker: func() *Worker {
				ctrl := minimock.NewController(t)
				return &Worker{
					DaysBefore: -2,
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, "key").
							Return(false)
					}(),
				}
			}(),
			args: args{
				key: "key",
				lastModified: func() *time.Time {
					t := now.AddDate(0, 0, -3)
					return &t
				}(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.worker.isNeededToUpdate(context.Background(), tt.args.key, tt.args.lastModified)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_LoadedPublisher_ToModel(t *testing.T) {
	t.Parallel()

	type want struct {
		publisher *models.Publisher
		domains   models.PublisherDomainSlice
		blacklist boil.Columns
	}

	tests := []struct {
		name            string
		loadedPublisher *LoadedPublisher
		want            want
	}{
		{
			name: "maxBlacklistColumnLength",
			loadedPublisher: &LoadedPublisher{
				Id:         "1",
				Name:       "publisher",
				MediaBuyer: &field{Id: "media_buyer_id"},
				PausedDate: 1000,
				Site:       []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:    "1",
					Name:           "publisher",
					MediaBuyerID:   null.String{Valid: true, String: "media_buyer_id"},
					PauseTimestamp: null.Int64{Valid: true, Int64: 1000},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
						models.PublisherColumns.AccountManagerID,
						models.PublisherColumns.CampaignManagerID,
						models.PublisherColumns.OfficeLocation,
						models.PublisherColumns.ReactivateTimestamp,
						models.PublisherColumns.StartTimestamp,
					},
				},
			},
		},
		{
			name: "minBlacklistColumnLength",
			loadedPublisher: &LoadedPublisher{
				Id:              "1",
				Name:            "publisher",
				MediaBuyer:      &field{Id: "media_buyer_id"},
				StartDate:       500,
				PausedDate:      1000,
				ReactivatedDate: 2000,
				AccountManager:  &field{Id: "account_manager_id"},
				CampaignManager: &field{Id: "campaign_manager_id"},
				OfficeLocation:  "office",
				Site:            []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:         "1",
					Name:                "publisher",
					MediaBuyerID:        null.String{Valid: true, String: "media_buyer_id"},
					StartTimestamp:      null.Int64{Valid: true, Int64: 500},
					PauseTimestamp:      null.Int64{Valid: true, Int64: 1000},
					ReactivateTimestamp: null.Int64{Valid: true, Int64: 2000},
					AccountManagerID:    null.String{Valid: true, String: "account_manager_id"},
					CampaignManagerID:   null.String{Valid: true, String: "campaign_manager_id"},
					OfficeLocation:      null.String{Valid: true, String: "office"},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
					},
				},
			},
		},
		{
			name: "managerIDFromMap",
			loadedPublisher: &LoadedPublisher{
				Id:         "1",
				Name:       "publisher",
				MediaBuyer: &field{Id: "62de259de6e2871c098001e9"},
				PausedDate: 1000,
				Site:       []string{"1.com", "2.com", "3.com"},
			},
			want: want{
				publisher: &models.Publisher{
					PublisherID:    "1",
					Name:           "publisher",
					MediaBuyerID:   null.String{Valid: true, String: "18"},
					PauseTimestamp: null.Int64{Valid: true, Int64: 1000},
				},
				domains: models.PublisherDomainSlice{
					{PublisherID: "1", Domain: "1.com"},
					{PublisherID: "1", Domain: "2.com"},
					{PublisherID: "1", Domain: "3.com"},
				},
				blacklist: boil.Columns{
					Kind: 4,
					Cols: []string{
						models.PublisherColumns.CreatedAt,
						models.PublisherColumns.Status,
						models.PublisherColumns.IntegrationType,
						models.PublisherColumns.AccountManagerID,
						models.PublisherColumns.CampaignManagerID,
						models.PublisherColumns.OfficeLocation,
						models.PublisherColumns.ReactivateTimestamp,
						models.PublisherColumns.StartTimestamp,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			publisher, domains, blacklist := tt.loadedPublisher.ToModel(getMockManagersMap())
			assert.Equal(t, tt.want.publisher, publisher)
			assert.Equal(t, tt.want.domains, domains)
			assert.Equal(t, tt.want.blacklist, blacklist)
		})
	}
}

func Test_getManagerID(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "idFromMap",
			args: args{
				id: "62de259de6e2871c098001e9",
			},
			want: "18",
		},
		{
			name: "initialId",
			args: args{
				id: "someunknownid",
			},
			want: "someunknownid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getManagerID(tt.args.id, getMockManagersMap())
			assert.Equal(t, tt.want, got)
		})
	}
}

func getMockManagersMap() map[string]string {
	return map[string]string{
		"62de259de6e2871c098001e9": "18",
	}
}
