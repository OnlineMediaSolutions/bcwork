package adstxt

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (a *AdsTxtModule) UpdateAdsTxtMetadata(ctx context.Context, data map[string]*dto.AdsTxtGroupedByDPData) error {
	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	modsMeta, err := createAdsTxtMetaData(data)
	if err != nil {
		return fmt.Errorf("failed to create ads txt metadata: %w", err)
	}

	for _, modMeta := range modsMeta {
		err = modMeta.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return fmt.Errorf("failed to insert metadata record: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to make commit for updating ads txt metadata: %w", err)
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
		key := fmt.Sprintf(utils.AdsTxtMetaDataKeyTemplate, row.Parent.DemandPartnerID)
		deduplicationKey := fmt.Sprintf("%v:%v:%v", key, row.Parent.PublisherID, row.Parent.Domain)

		// duplicates could appear because of separation media types
		_, ok := deduplicationMap[deduplicationKey]
		if !ok && row.Parent.IsReadyToWork {
			deduplicationMap[deduplicationKey] = struct{}{}
			records[key] = append(records[key], adstxtRealtimeRecord{
				PubID:  row.Parent.PublisherID,
				Domain: row.Parent.Domain,
			})
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
