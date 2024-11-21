package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/history"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

const (
	blockTypeBADV = "badv"
	blockTypeBCAT = "bcat"
)

type BlocksService struct {
	historyModule history.HistoryModule
}

func NewBlocksService(historyModule history.HistoryModule) *BlocksService {
	return &BlocksService{
		historyModule: historyModule,
	}
}

var query = `SELECT metadata_queue.*
FROM metadata_queue,(select key,max(created_at) created_at FROM metadata_queue WHERE key LIKE 'bcat:%' OR key like 'badv:%' group by key) last
WHERE last.created_at=metadata_queue.created_at
    AND last.key=metadata_queue.key `

var sortQuery = ` ORDER by metadata_queue.key`

func (b *BlocksService) GetBlocks(ctx context.Context, request *dto.BlockGetRequest) (models.MetadataQueueSlice, error) {
	key := createKeyForQuery(request)
	records := models.MetadataQueueSlice{}

	err := queries.Raw(query+key+sortQuery).Bind(ctx, bcdb.DB(), &records)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all price factors: %w", err)
	}

	return records, nil
}

func (b *BlocksService) UpdateBlocks(ctx context.Context, data *dto.BlockUpdateRequest) error {
	request := &dto.BlockGetRequest{
		Types:     []string{blockTypeBADV, blockTypeBCAT},
		Publisher: data.Publisher,
	}

	if data.Domain != "" {
		request.Domain = data.Domain
	}

	records, err := b.GetBlocks(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to get previous block records: %w", err)
	}

	var (
		badv []string
		bcat []string
	)
	for _, record := range records {
		if strings.HasPrefix(record.Key, blockTypeBADV) {
			err := json.Unmarshal(record.Value, &badv)
			if err != nil {
				return fmt.Errorf("failed to unmarshal badv previous value: %w", err)
			}
		}

		if strings.HasPrefix(record.Key, blockTypeBCAT) {
			err := json.Unmarshal(record.Value, &bcat)
			if err != nil {
				return fmt.Errorf("failed to unmarshal bcat previous value: %w", err)
			}
		}
	}

	var oldData any
	if len(records) > 0 {
		oldData = &dto.BlockUpdateRequest{
			Publisher: data.Publisher,
			Domain:    data.Domain,
			BADV:      badv,
			BCAT:      bcat,
		}
	}

	if data.BADV != nil {
		if err := updateDB(ctx, blockTypeBADV, data.Publisher, data.Domain, data.BADV); err != nil {
			return fmt.Errorf("failed to insert blocks badv metadata update to queue: %w", err)
		}
	}

	if data.BCAT != nil {
		if err := updateDB(ctx, blockTypeBCAT, data.Publisher, data.Domain, data.BCAT); err != nil {
			return fmt.Errorf("failed to insert blocks bcat metadata update to queue: %w", err)
		}
	}

	b.historyModule.SaveAction(ctx, oldData, data, nil)

	return nil
}

func updateDB(ctx context.Context, businessType, publisher, domain string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	mod := models.MetadataQueue{
		Key:           fmt.Sprintf("%s:%s", businessType, publisher),
		TransactionID: bcguid.NewFromf(publisher, domain, businessType, time.Now()),
		Value:         b,
	}

	if domain != "" {
		mod.Key += ":" + domain
	}

	if err := mod.Insert(ctx, bcdb.DB(), boil.Infer()); err != nil {
		return err
	}

	return nil
}

func createKeyForQuery(request *dto.BlockGetRequest) string {
	types := request.Types
	publisher := request.Publisher
	domain := request.Domain

	var query bytes.Buffer

	//If no publisher or no business types or empty body than return all
	if len(publisher) == 0 || len(types) == 0 {
		query.WriteString(` and 1=1 `)
		return query.String()
	}

	for index, btype := range types {
		if index == 0 {
			query.WriteString("AND (")
		}
		if len(publisher) != 0 && len(domain) != 0 {
			query.WriteString(" (metadata_queue.key = '" + btype + ":" + publisher + ":" + domain + "')")

		} else if len(publisher) != 0 {
			query.WriteString(" (metadata_queue.key = '" + btype + ":" + publisher + "')")
		}
		if index < len(types)-1 {
			query.WriteString(" OR")
		}
	}
	query.WriteString(")")
	return query.String()
}
