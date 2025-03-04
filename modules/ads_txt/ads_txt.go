package adstxt

import (
	"context"
	"database/sql"

	"github.com/m6yf/bcwork/dto"
)

type AdsTxtLinesCreater interface {
	CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error
	// TODO: needs publisher management
	UpdateAdsTxtMetadata(ctx context.Context, data map[string]*dto.AdsTxtGroupedByDPData) error
}

type AdsTxtModule struct {
}

func NewAdsTxtModule() *AdsTxtModule {
	return &AdsTxtModule{}
}

func (a *AdsTxtModule) CreateDemandPartnerConnectionAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, ids, demandPartnerConnectionQueryType)
}

func (a *AdsTxtModule) CreateDemandPartnerChildAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, ids, demandPartnerChildQueryType)
}

func (a *AdsTxtModule) CreateSeatOwnerAdsTxtLines(ctx context.Context, tx *sql.Tx, ids []int) error {
	return createAdsTxtLine(ctx, tx, ids, seatOwnerQueryType)
}
