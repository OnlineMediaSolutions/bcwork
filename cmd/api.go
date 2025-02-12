package cmd

import (
	"context"
	"net/http/pprof"
	"strings"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/m6yf/bcwork/api/rest"
	"github.com/m6yf/bcwork/api/rest/bulk"
	"github.com/m6yf/bcwork/api/rest/report"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/export"
	"github.com/m6yf/bcwork/modules/history"
	supertokens_module "github.com/m6yf/bcwork/modules/supertokens"
	"github.com/m6yf/bcwork/validations"
	"github.com/m6yf/bcwork/validations/dpo"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/supertokens"

	_ "github.com/m6yf/bcwork/api/rest/docs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// @title Swagger OMS API
// @version 1.0
// @description API for OMS.

// @contact.name Brightcom Support
// @contact.url http://www.nanoook.com/support
// @contact.email support@gutsy.me

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name oms-worker-api-key

// apiCmd represents api server command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "OMS API Server",
	Long:  ``,

	Run: ApiCmd,
}

func ApiCmd(cmd *cobra.Command, args []string) {
	ctx := context.TODO()

	dbEnv := viper.GetString("database.env")
	if dbEnv != "" {
		err := bcdb.InitDB(dbEnv)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect init DB")
		}
	}

	// log sql
	if viper.GetBool("db." + dbEnv + ".debug") {
		boil.DebugMode = true
	}

	historyModule := history.NewHistoryClient()
	exportModule := export.NewExportModule()
	compassModule := compass.NewCompass()

	apiURL, webURL, initFunc := supertokens_module.GetSuperTokensConfig()
	supertokenClient, err := supertokens_module.NewSuperTokensClient(apiURL, webURL, initFunc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to supertokens")
	}

	omsNP := rest.NewOMSNewPlatform(ctx, supertokenClient, historyModule, exportModule, compassModule, true)

	app := fiber.New(fiber.Config{ErrorHandler: rest.ErrorHandler})
	allowedHeaders := append([]string{"Content-Type", "x-amz-acl"}, supertokens.GetAllCORSHeaders()...)
	allowedHeadersInCommaSeparetedStringFormat := strings.Join(allowedHeaders, ", ")

	app.Use(cors.New(cors.Config{
		Next:         nil,
		AllowOrigins: "http://localhost:3000,https://app-dev.nanoook.com,https://app.nanoook.com,https://login.nanoook.com,https://loginstg.nanoook.com, https://admin.nanoook.com,https://api.nanoook.com",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     allowedHeadersInCommaSeparetedStringFormat,
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	// logging basic information about all requests
	app.Use(rest.LoggingMiddleware)
	// collect profile for each request
	// app.Use(rest.ProfileMiddleware)

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("UP") })
	app.Get("/ping", rest.PingPong)
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	app.Post("/download", validations.ValidateDownload, omsNP.DownloadHandler)

	users := app.Group("/user")
	// unsecured endpoints for internal use
	users.Get("/info", omsNP.UserGetInfoHandler)
	users.Get("/by_types", omsNP.UserGetByTypesHandler)

	// supertokens middleware + session verification
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	app.Use(adaptor.HTTPMiddleware(supertokenClient.VerifySession))

	// debug endpoints for profiling
	debug := app.Group("/debug")
	debug.Get("/pprof", adaptor.HTTPHandlerFunc(pprof.Index))
	debug.Get("/pprof/profile", adaptor.HTTPHandlerFunc(pprof.Profile))

	// configuration
	app.Post("/config/get", rest.ConfigurationGetHandler)
	app.Post("/config", validations.ValidateConfig, rest.ConfigurationPostHandler)

	// report
	reportGroup := app.Group("/report")
	reportGroup.Get("/daily/revenue", report.DailyRevenue)
	reportGroup.Get("/hourly/revenue", report.HourlyRevenue)
	reportGroup.Get("/demand", rest.DemandReportGetHandler)
	reportGroup.Get("/demand/hourly", rest.DemandHourlyReportGetHandler)
	reportGroup.Get("/publisher", rest.PublisherReportGetHandler)
	reportGroup.Get("/publisher/hourly", rest.PublisherHourlyReportGetHandler)
	reportGroup.Get("/iiq/hourly", rest.IiqTestingGetHandler)

	app.Post("/metadata/update", rest.MetadataPostHandler)
	app.Get("/price/floor/set", rest.PriceFloorSetHandler)
	app.Get("/price/floor/get", rest.PriceFloorGetHandler)
	app.Get("/price/floor/get/all", rest.PriceFloorGetAllHandler)
	app.Post("/price/fixed", rest.FixedPricePostHandler)
	app.Get("/price/fixed", rest.FixedPriceGetAllHandler)
	app.Get("/had/price/set", rest.HouseAdPriceSetHandler)
	app.Get("/had/price/get", rest.HouseAdPriceGetHandler)
	app.Get("/had/price/get/all", rest.HouseAdPriceGetAllHandler)
	app.Post("/demand/factor", rest.DemandFactorPostHandler)
	app.Get("/demand/factor/set", rest.DemandFactorSetHandler)
	app.Get("/demand/factor/get", rest.DemandFactorGetHandler)
	app.Get("/demand/factor/get/all", rest.DemandFactorGetAllHandler)
	app.Post("/publisher/demand/get", rest.PublisherDemandGetHandler)
	app.Post("/publisher/demand/udpate", validations.ValidateBulkPublisherDemands, rest.PublisherDemandUpdate)
	app.Post("/price/override", rest.PriceOverrideHandler)

	// global factor
	app.Post("/global/factor", validations.ValidateGlobalFactor, omsNP.GlobalFactorPostHandler)
	app.Post("/global/factor/get", omsNP.GlobalFactorGetHandler)

	// block
	app.Post("/block", validations.ValidateBlocks, omsNP.BlockPostHandler)
	app.Post("/block/get", omsNP.BlockGetAllHandler)

	// confiant
	app.Post("/confiant", validations.ValidateConfiant, omsNP.ConfiantPostHandler)
	app.Post("/confiant/get", omsNP.ConfiantGetAllHandler)

	// pixalate
	app.Post("/pixalate", validations.ValidatePixalate, omsNP.PixalatePostHandler)
	app.Post("/pixalate/get", omsNP.PixalateGetAllHandler)
	app.Delete("/pixalate/delete", omsNP.PixalateDeleteHandler)

	// demand partners
	dp := app.Group("/dp")
	dp.Post("/get", omsNP.DemandPartnerGetHandler)
	dp.Post("/set", omsNP.DemandPartnerSetHandler)       // TODO: add validation - validations.ValidateDemandPartner
	dp.Post("/update", omsNP.DemandPartnerUpdateHandler) // TODO: add validation - validations.ValidateDemandPartner
	dp.Post("/seat_owner/get", omsNP.DemandPartnerGetSeatOwnersHandler)

	// dpo
	dpoGroup := app.Group("/dpo")
	dpoGroup.Post("/set", dpo.ValidateDPO, omsNP.DemandPartnerOptimizationSetHandler)
	dpoGroup.Post("/get", omsNP.DemandPartnerOptimizationGetHandler)
	dpoGroup.Delete("/delete", omsNP.DemandPartnerOptimizationDeleteHandler)
	dpoGroup.Get("/update", dpo.ValidateQueryParams, omsNP.DemandPartnerOptimizationUpdateHandler)

	// ads.txt
	adsTxtGroup := app.Group("/ads_txt")
	adsTxtGroup.Post("/main", omsNP.AdsTxtMainHandler)
	adsTxtGroup.Post("/group_by_dp", omsNP.AdsTxtGroupByDPHandler)
	adsTxtGroup.Post("/am", omsNP.AdsTxtAMHandler)
	adsTxtGroup.Post("/cm", omsNP.AdsTxtCMHandler)
	adsTxtGroup.Post("/mb", omsNP.AdsTxtMBHandler)
	adsTxtGroup.Post("/update", validations.AdsTxtValidation, omsNP.AdsTxtUpdateHandler)

	// publisher
	publisher := app.Group("/publisher")
	publisher.Post("/new", validations.CreatePublisherValidation, omsNP.PublisherNewHandler)
	publisher.Post("/update", validations.UpdatePublisherValidation, omsNP.PublisherUpdateHandler)
	publisher.Post("/get", omsNP.PublisherGetHandler)
	publisher.Post("/count", omsNP.PublisherCountHandler)
	publisher.Post("/details/get", omsNP.PublisherDetailsGetHandler)

	// domain
	publisher.Post("/domain/get", omsNP.PublisherDomainGetHandler)
	publisher.Post("/domain", validations.PublisherDomainValidation, omsNP.PublisherDomainPostHandler)

	// bid caching
	bidCachingGroup := app.Group("/bid_caching")
	bidCachingGroup.Post("/get", omsNP.BidCachingGetAllHandler)
	bidCachingGroup.Post("/set", validations.ValidateBidCaching, omsNP.BidCachingSetHandler)
	bidCachingGroup.Post("/update", validations.ValidateUpdateBidCaching, omsNP.BidCachingUpdateHandler)
	bidCachingGroup.Delete("/delete", omsNP.BidCachingDeleteHandler)

	//refresh cache
	refreshCachingGroup := app.Group("/refresh_cache")
	refreshCachingGroup.Post("/get", omsNP.RefreshCacheGetAllHandler)
	refreshCachingGroup.Post("/set", validations.ValidateRefreshCache, omsNP.RefreshCacheSetHandler)
	refreshCachingGroup.Post("/update", validations.ValidateUpdateRefreshCache, omsNP.RefreshCacheUpdateHandler)
	refreshCachingGroup.Delete("/delete", omsNP.RefreshCacheDeleteHandler)

	// factor
	app.Post("/factor/get", omsNP.FactorGetAllHandler)
	app.Post("/factor", validations.ValidateFactor, omsNP.FactorPostHandler)
	app.Delete("/factor/delete", omsNP.FactorDeleteHandler)

	// floor
	app.Post("/floor/get", omsNP.FloorGetAllHandler)
	app.Post("/floor", validations.ValidateFloors, omsNP.FloorPostHandler)
	app.Delete("/floor/delete", omsNP.FloorDeleteHandler)

	// bulk
	bulkGroup := app.Group("/bulk")
	bulkGroup.Post("/factor", validations.ValidateBulkFactors, omsNP.FactorBulkPostHandler)
	bulkGroup.Post("/floor", validations.ValidateBulkFloor, bulk.FloorBulkPostHandler)
	bulkGroup.Post("/dpo", validations.ValidateDPOInBulk, omsNP.DemandPartnerOptimizationBulkPostHandler)
	bulkGroup.Post("/global/factor", validations.ValidateBulkGlobalFactor, omsNP.GlobalFactorBulkPostHandler)

	// adjuster
	app.Post("/adjust/floor", validations.ValidateAdjusterURL, omsNP.FloorAdjusterHandler)
	app.Post("/adjust/factor", validations.ValidateAdjusterURL, omsNP.FactorAdjusterHandler)

	// competitor
	app.Post("/competitor/get", rest.CompetitorGetAllHandler)
	app.Post("/competitor", validations.ValidateCompetitorURL, rest.CompetitorPostHandler)

	// targeting
	targeting := app.Group("/targeting")
	targeting.Post("/get", omsNP.TargetingGetHandler)
	targeting.Post("/set", validations.ValidateTargeting, omsNP.TargetingSetHandler)
	targeting.Post("/update", validations.ValidateTargeting, omsNP.TargetingUpdateHandler)
	targeting.Post("/tags", omsNP.TargetingExportTagsHandler)

	// search
	app.Post("/search", omsNP.SearchHandler)

	// user management (only for users with 'admin' role)
	users.Use(supertokenClient.AdminRoleRequired)
	users.Post("/get", omsNP.UserGetHandler)
	users.Post("/set", validations.ValidateUser, omsNP.UserSetHandler)
	users.Post("/update", validations.ValidateUser, omsNP.UserUpdateHandler)

	// history
	app.Post("/history/get", omsNP.HistoryGetHandler)
	app.Post("/email", omsNP.SendEmailReport)

	app.Listen(":8000")
}

func init() {
	rootCmd.AddCommand(apiCmd)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/oms/api/")
	viper.AddConfigPath("$HOME/.")
	viper.AddConfigPath(".")

	viper.SetDefault("env", "prod")
	viper.SetDefault("ports.http", "8000")
	viper.SetDefault("api.chunkSize", 2000)

	err := viper.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warn().Msg("config file was not found, default values will be used")
	}
}
