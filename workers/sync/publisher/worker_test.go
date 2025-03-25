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

const keyValue = "key"

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
											key := keyValue
											return &key
										}(),
										LastModified: &now,
									},
								},
							}, nil).
							GetObjectInputMock.
							Expect("bucket", keyValue).
							Return([]byte(`[{"_id":"id", "accountManager": {"id":"am_id"}, "site": ["1.com"]}]`), nil)
					}(),
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, keyValue).
							Return(false).
							UpsertPublisherAndDomainsMock.
							Expect(
								minimock.AnyContext,
								&models.Publisher{
									PublisherID:      "id",
									AccountManagerID: null.String{Valid: true, String: "am_id"},
								},
								[]*models.PublisherDomain{
									{PublisherID: "id", Domain: "1.com"},
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
							SaveResultOfLastSyncMock.
							Expect(minimock.AnyContext, keyValue, false).
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
											key := keyValue
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
							Expect(minimock.AnyContext, keyValue).
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
											key := keyValue
											return &key
										}(),
										LastModified: &now,
									},
								},
							}, nil).
							GetObjectInputMock.
							Expect("bucket", keyValue).
							Return([]byte(`[{"_id":"id", "accountManager": {"id":"am_id"}, "site": ["1.com"]}]`), nil)
					}(),
					DB: func() db.PublisherSyncStorage {
						return dbmocks.NewPublisherSyncStorageMock(ctrl).
							HadLoadingErrorLastTimeMock.
							Expect(minimock.AnyContext, keyValue).
							Return(false).
							UpsertPublisherAndDomainsMock.
							Expect(
								minimock.AnyContext,
								&models.Publisher{
									PublisherID:      "id",
									AccountManagerID: null.String{Valid: true, String: "am_id"},
								},
								[]*models.PublisherDomain{
									{PublisherID: "id", Domain: "1.com"},
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
							Return(errors.New("error while inserting publisher domain")).
							SaveResultOfLastSyncMock.
							Expect(minimock.AnyContext, keyValue, true).
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
							Expect(minimock.AnyContext, keyValue).
							Return(false)
					}(),
				}
			}(),
			args: args{
				key: keyValue,
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
							Expect(minimock.AnyContext, keyValue).
							Return(true)
					}(),
				}
			}(),
			args: args{
				key: keyValue,
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
							Expect(minimock.AnyContext, keyValue).
							Return(false)
					}(),
				}
			}(),
			args: args{
				key: keyValue,
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
