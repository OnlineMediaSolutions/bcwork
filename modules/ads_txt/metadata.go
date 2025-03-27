package adstxt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/logger"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/net/publicsuffix"
)

func (a *AdsTxtModule) UpdateAdsTxtMetadata(ctx context.Context, resp *dto.AdsTxtGroupByDPResponse) error {
	modsMeta, err := createAdsTxtMetaData(ctx, resp)
	if err != nil {
		return fmt.Errorf("failed to create ads txt metadata: %w", err)
	}

	if len(modsMeta) == 0 {
		return errors.New("no data to update ads txt metadata")
	}

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, modMeta := range modsMeta {
		err := modMeta.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return fmt.Errorf("failed to insert metadata record: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to make commit updates for ads txt metadata: %w", err)
	}

	return nil
}

func createAdsTxtMetaData(ctx context.Context, resp *dto.AdsTxtGroupByDPResponse) ([]*models.MetadataQueue, error) {
	type adstxtRealtimeRecord struct {
		PubID  string `json:"pubid"`
		Domain string `json:"domain"`
	}

	records := make(map[string][]adstxtRealtimeRecord, len(resp.Data))
	deduplicationMap := make(map[string]struct{}, len(resp.Data))

	for _, row := range resp.Data {
		key := fmt.Sprintf(utils.AdsTxtMetaDataKeyTemplate, row.Parent.DemandPartnerID)
		deduplicationKey := fmt.Sprintf("%v:%v:%v", key, row.Parent.PublisherID, row.Parent.Domain)

		// duplicates could appear because of multiple connections for same demand partner
		_, ok := deduplicationMap[deduplicationKey]
		if !ok && row.Parent.IsReadyToGoLive {
			deduplicationMap[deduplicationKey] = struct{}{}
			records[key] = append(records[key], adstxtRealtimeRecord{
				PubID:  row.Parent.PublisherID,
				Domain: row.Parent.Domain,
			})

			// adding top level domain
			topLevelDomain, err := publicsuffix.EffectiveTLDPlusOne(row.Parent.Domain)
			if err != nil {
				logger.Logger(ctx).Err(err).Msgf("cannot extract top level domain for %v", row.Parent.Domain)
				continue
			}

			if topLevelDomain != row.Parent.Domain {
				records[key] = append(records[key], adstxtRealtimeRecord{
					PubID:  row.Parent.PublisherID,
					Domain: topLevelDomain,
				})
			}
		}
	}

	modsMeta := make([]*models.MetadataQueue, 0, len(records))
	for key, record := range records {
		value, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}

		modsMeta = append(modsMeta, &models.MetadataQueue{
			TransactionID: bcguid.NewFromf(time.Now()),
			Key:           key,
			Value:         value,
		})
	}

	return modsMeta, nil
}
