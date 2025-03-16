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
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/net/publicsuffix"
)

func (a *AdsTxtModule) UpdateAdsTxtMetadata(ctx context.Context, data map[string]*dto.AdsTxtGroupedByDPData) error {
	modsMeta, err := createAdsTxtMetaData(data)
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

func createAdsTxtMetaData(data map[string]*dto.AdsTxtGroupedByDPData) ([]*models.MetadataQueue, error) {
	type adstxtRealtimeRecord struct {
		PubID  string `json:"pubid"`
		Domain string `json:"domain"`
	}

	records := make(map[string][]adstxtRealtimeRecord, len(data))
	deduplicationMap := make(map[string]struct{}, len(data))

	for _, row := range data {
		adsTxtLine := row.Parent
		if adsTxtLine == nil {
			adsTxtLine = row.Children[0]
		}

		key := fmt.Sprintf(utils.AdsTxtMetaDataKeyTemplate, adsTxtLine.DemandPartnerID)
		deduplicationKey := fmt.Sprintf("%v:%v:%v", key, adsTxtLine.PublisherID, adsTxtLine.Domain)

		// duplicates could appear because of separation media types
		_, ok := deduplicationMap[deduplicationKey]
		if !ok && adsTxtLine.IsReadyToGoLive {
			deduplicationMap[deduplicationKey] = struct{}{}
			records[key] = append(records[key], adstxtRealtimeRecord{
				PubID:  adsTxtLine.PublisherID,
				Domain: adsTxtLine.Domain,
			})

			// adding subdomains
			subdomain, err := publicsuffix.EffectiveTLDPlusOne(adsTxtLine.Domain)
			if err != nil {
				return nil, err
			}

			if subdomain != adsTxtLine.Domain {
				records[key] = append(records[key], adstxtRealtimeRecord{
					PubID:  adsTxtLine.PublisherID,
					Domain: subdomain,
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
