package rest

import (
	"context"

	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/core/bulk"
	adstxt "github.com/m6yf/bcwork/modules/ads_txt"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/export"
	"github.com/m6yf/bcwork/modules/history"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
)

type OMSNewPlatform struct {
	userService          *core.UserService
	targetingService     *core.TargetingService
	domainService        *core.DomainService
	historyService       *core.HistoryService
	publisherService     *core.PublisherService
	globalFactorService  *core.GlobalFactorService
	bulkService          bulk.Bulker
	confiantService      *core.ConfiantService
	pixalateService      *core.PixalateService
	blocksService        *core.BlocksService
	floorService         *core.FloorService
	factorService        *core.FactorService
	demandPartnerService *core.DemandPartnerService
	dpoService           *core.DPOService
	adjustService        bulk.Adjuster
	searchService        *core.SearchService
	bidCachingService    *core.BidCachingService
	refreshCacheService  *core.RefreshCacheService
	emailService         *core.EmailService
	downloadService      *core.DownloadService
	adsTxtService        *core.AdsTxtService
	dpApiService         *core.DpAPIService
}

func NewOMSNewPlatform(
	ctx context.Context,
	supertokenClient supertokens_module.TokenManagementSystem,
	historyModule history.HistoryModule,
	exportModule export.Exporter,
	compassModule compass.CompassModule,
	adstxtModule adstxt.AdsTxtLinesCreater,
	sendRegistrationEmail bool, // Temporary, remove after decoupling email sender service
) *OMSNewPlatform {
	userService := core.NewUserService(supertokenClient, historyModule, sendRegistrationEmail)
	targetingService := core.NewTargetingService(historyModule)
	domainService := core.NewDomainService(historyModule)
	historyService := core.NewHistoryService()
	publisherService := core.NewPublisherService(historyModule, compassModule)
	globalFactorService := core.NewGlobalFactorService(historyModule)
	bulkService := bulk.NewBulkService(historyModule)
	confiantService := core.NewConfiantService(historyModule)
	pixalateService := core.NewPixalateService(historyModule)
	blocksService := core.NewBlocksService(historyModule)
	floorService := core.NewFloorService(historyModule)
	factorService := core.NewFactorService(historyModule)
	demandPartnerService := core.NewDemandPartnerService(historyModule, adstxtModule)
	dpoService := core.NewDPOService(historyModule)
	bidCachingService := core.NewBidCachingService(historyModule)
	refreshCacheService := core.NewRefreshCacheService(historyModule)
	searchService := core.NewSearchService(ctx)
	emailService := core.NewEmailService(ctx)
	downloadService := core.NewDownloadService(exportModule)
	adsTxtService := core.NewAdsTxtService(historyModule, compassModule)

	return &OMSNewPlatform{
		userService:          userService,
		targetingService:     targetingService,
		domainService:        domainService,
		historyService:       historyService,
		publisherService:     publisherService,
		globalFactorService:  globalFactorService,
		bulkService:          bulkService,
		confiantService:      confiantService,
		pixalateService:      pixalateService,
		blocksService:        blocksService,
		floorService:         floorService,
		factorService:        factorService,
		demandPartnerService: demandPartnerService,
		dpoService:           dpoService,
		searchService:        searchService,
		bidCachingService:    bidCachingService,
		refreshCacheService:  refreshCacheService,
		adjustService:        bulkService,
		emailService:         emailService,
		downloadService:      downloadService,
		adsTxtService:        adsTxtService,
	}
}
