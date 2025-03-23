package adstxt

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/m6yf/bcwork/dto"
)

type AdsTxtLinesCreater interface {
	CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreatePublisherDomainAdsTxtLines(ctx context.Context, tx *sql.Tx, publisherID, domain string) error
	UpdateAdsTxtMetadata(ctx context.Context, data map[string]*dto.AdsTxtGroupedByDPData) error
}

type AdsTxtModule struct {
}

func NewAdsTxtModule() *AdsTxtModule {
	return &AdsTxtModule{}
}

func (a *AdsTxtModule) CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(demandPartnerConnectionQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(demandPartnerChildQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesTemplateQuery(seatOwnerQueryType), pq.Array(ids))
}

func (a *AdsTxtModule) CreatePublisherDomainAdsTxtLines(ctx context.Context, tx *sql.Tx, domain, publisherID string) error {
	return createAdsTxtLine(ctx, tx, getAdsTxtLinesFromPublisherDomainTemplateQuery(), domain, publisherID)
}
