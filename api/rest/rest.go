package rest

import (
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/modules/history"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
)

type OMSNewPlatform struct {
	userService         *core.UserService
	targetingService    *core.TargetingService
	domainService       *core.DomainService
	historyService      *core.HistoryService
	publisherService    *core.PublisherService
	globalFactorService *core.GlobalFactorService
	bulkService         *bulk.BulkService
	bulkFactorService   *bulk.BulkFactorService
	confiantService     *core.ConfiantService
	pixalateService     *core.PixalateService
	blocksService       *core.BlocksService
	floorService        *core.FloorService
	factorService       *core.FactorService
	dpoService          *core.DPOService
}

func NewOMSNewPlatform(
	supertokenClient supertokens_module.TokenManagementSystem,
	historyModule history.HistoryModule,
	sendRegistrationEmail bool, // Temporary, remove after decoupling email sender service
) *OMSNewPlatform {
	userService := core.NewUserService(supertokenClient, historyModule, sendRegistrationEmail)
	targetingService := core.NewTargetingService(historyModule)
	domainService := core.NewDomainService(historyModule)
	historyService := core.NewHistoryService()
	publisherService := core.NewPublisherService(historyModule)
	globalFactorService := core.NewGlobalFactorService(historyModule)
	bulkService := bulk.NewBulkService(historyModule)
	bulkFactorService := bulk.NewBulkFactorService(historyModule)
	confiantService := core.NewConfiantService(historyModule)
	pixalateService := core.NewPixalateService(historyModule)
	blocksService := core.NewBlocksService(historyModule)
	floorService := core.NewFloorService(historyModule)
	factorService := core.NewFactorService(historyModule)
	dpoService := core.NewDPOService(historyModule)

	return &OMSNewPlatform{
		userService:         userService,
		targetingService:    targetingService,
		domainService:       domainService,
		historyService:      historyService,
		publisherService:    publisherService,
		globalFactorService: globalFactorService,
		bulkService:         bulkService,
		bulkFactorService:   bulkFactorService,
		confiantService:     confiantService,
		pixalateService:     pixalateService,
		blocksService:       blocksService,
		floorService:        floorService,
		factorService:       factorService,
		dpoService:          dpoService,
	}
}
