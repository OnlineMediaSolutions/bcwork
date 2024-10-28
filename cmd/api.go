package cmd

import (
	"strings"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/m6yf/bcwork/api/rest"
	"github.com/m6yf/bcwork/api/rest/bulk"
	"github.com/m6yf/bcwork/api/rest/report"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
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

// @title Swagger Brightcom API
// @version 1.0
// @description API for Brightcom game.

// @contact.name Brightcom Support
// @contact.url http://www.nanoook.com/support
// @contact.email support@gutsy.me

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// apiCmd represents api server command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "OMS API Server",
	Long:  ``,

	Run: ApiCmd,
}

func ApiCmd(cmd *cobra.Command, args []string) {
	dbEnv := viper.GetString("database.env")
	if dbEnv != "" {
		err := bcdb.InitDB(dbEnv)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect init DB")
		}
	}

	boil.DebugMode = true

	//err = bcdwh.InitDB("prod-do")
	//if err != nil {
	//	log.Fatal().Err(err).Msg("failed to connect DWH")
	//}

	apiURL, webURL, initFunc := supertokens_module.GetSuperTokensConfig()
	supertokenClient, err := supertokens_module.NewSuperTokensClient(apiURL, webURL, initFunc)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to supertokens")
	}
	userService := core.NewUserService(supertokenClient, true)
	userManagementSystem := rest.NewUserManagementSystem(userService)

	// Log sql
	if viper.GetBool("sqlboiler.debug") {
		boil.DebugMode = true
	}

	app := fiber.New(fiber.Config{ErrorHandler: rest.ErrorHandler})
	allowedHeaders := append([]string{"Content-Type", "x-amz-acl"}, supertokens.GetAllCORSHeaders()...)
	allowedHeadersInCommaSeparetedStringFormat := strings.Join(allowedHeaders, ", ")

	app.Use(cors.New(cors.Config{
		Next:         nil,
		AllowOrigins: "http://localhost:3000,https://app-dev.nanoook.com,https://app.nanoook.com,https://login.nanoook.com,https://admin.nanoook.com,https://api.nanoook.com",
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
	// app.Use(loggingMiddleware)

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("UP") })
	app.Get("/ping", rest.PingPong)
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	users := app.Group("/user")
	users.Get("/info", userManagementSystem.UserGetInfoHandler) // unsecured endpoint for inner usage

	// adding the supertokens middleware + session verification
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))
	// app.Use(adaptor.HTTPMiddleware(supertokenClient.VerifySession))

	// Configuration
	app.Post("/config/get", rest.ConfigurationGetHandler)
	app.Post("/config", validations.ValidateConfig, rest.ConfigurationPostHandler)

	app.Get("/report/daily/revenue", report.DailyRevenue)
	app.Get("/report/hourly/revenue", report.HourlyRevenue)

	app.Get("/report/demand", rest.DemandReportGetHandler)
	app.Get("/report/demand/hourly", rest.DemandHourlyReportGetHandler)
	app.Get("/report/publisher", rest.PublisherReportGetHandler)
	app.Get("/report/publisher/hourly", rest.PublisherHourlyReportGetHandler)
	app.Get("/report/iiq/hourly", rest.IiqTestingGetHandler)
	app.Post("/metadata/update", rest.MetadataPostHandler)
	app.Get("/price/floor/set", rest.PriceFloorSetHandler)
	app.Get("/price/floor/get", rest.PriceFloorGetHandler)
	app.Get("/price/floor/get/all", rest.PriceFloorGetAllHandler)
	app.Get("/had/price/set", rest.HouseAdPriceSetHandler)
	app.Get("/had/price/get", rest.HouseAdPriceGetHandler)
	app.Get("/had/price/get/all", rest.HouseAdPriceGetAllHandler)
	app.Post("/demand/factor", rest.DemandFactorPostHandler)
	app.Get("/demand/factor/set", rest.DemandFactorSetHandler)
	app.Get("/demand/factor/get", rest.DemandFactorGetHandler)
	app.Get("/demand/factor/get/all", rest.DemandFactorGetAllHandler)
	app.Post("/price/fixed", rest.FixedPricePostHandler)
	app.Get("/price/fixed", rest.FixedPriceGetAllHandler)

	app.Post("/confiant", validations.ValidateConfiant, rest.ConfiantPostHandler)
	app.Post("/confiant/get", rest.ConfiantGetAllHandler)

	app.Post("/publisher/demand/get", rest.PublisherDemandGetHandler)
	app.Post("/publisher/demand/udpate", validations.ValidateBulkPublisherDemands, rest.PublisherDemandUpdate)

	app.Post("/global/factor", validations.ValidateGlobalFactor, rest.GlobalFactorPostHandler)
	app.Post("/global/factor/get", rest.GlobalFactorGetHandler)

	app.Post("/pixalate", validations.ValidatePixalate, rest.PixalatePostHandler)
	app.Post("/pixalate/get", rest.PixalateGetAllHandler)
	app.Delete("/pixalate/delete", rest.PixalateDeleteHandler)

	app.Post("/block", rest.BlockPostHandler)
	app.Post("/block/get", rest.BlockGetAllHandler)
	app.Post("/dp/get", rest.DemandPartnerGetHandler)

	app.Post("/dpo/set", dpo.ValidateDPO, rest.DemandPartnerOptimizationSetHandler)
	app.Post("/dpo/get", rest.DemandPartnerOptimizationGetHandler)
	app.Delete("/dpo/delete", rest.DemandPartnerOptimizationDeleteHandler)
	app.Get("/dpo/update", dpo.ValidateQueryParams, rest.DemandPartnerOptimizationUpdateHandler)

	app.Post("/publisher/new", validations.PublisherValidation, rest.PublisherNewHandler)
	app.Post("/publisher/update", rest.PublisherUpdateHandler)
	app.Post("/publisher/get", rest.PublisherGetHandler)
	app.Post("/publisher/count", rest.PublisherCountHandler)
	app.Post("/publisher/details/get", rest.PublisherDetailsGetHandler)

	app.Post("/publisher/domain/get", rest.PublisherDomainGetHandler)
	app.Post("/publisher/domain", validations.PublisherDomainValidation, rest.PublisherDomainPostHandler)

	app.Post("/factor/get", rest.FactorGetAllHandler)
	app.Post("/factor", validations.ValidateFactor, rest.FactorPostHandler)

	app.Post("/floor/get", rest.FloorGetAllHandler)
	app.Post("/floor", validations.ValidateFloors, rest.FloorPostHandler)

	app.Post("/bulk/factor", validations.ValidateBulkFactors, bulk.FactorBulkPostHandler)
	app.Post("/bulk/floor", validations.ValidateBulkFloor, bulk.FloorBulkPostHandler)
	app.Post("/bulk/dpo", validations.ValidateDPOInBulk, bulk.DemandPartnerOptimizationBulkPostHandler)
	app.Post("/bulk/global/factor", validations.ValidateBulkGlobalFactor, bulk.GlobalFactorBulkPostHandler)

	app.Post("/competitor/get", rest.CompetitorGetAllHandler)
	app.Post("/competitor", validations.ValidateCompetitorURL, rest.CompetitorPostHandler)

	app.Post("/download", validations.ValidateDownload, rest.DownloadPostHandler)
	// Targeting
	targeting := app.Group("/targeting")
	targeting.Post("/get", rest.TargetingGetHandler)
	targeting.Post("/set", validations.ValidateTargeting, rest.TargetingSetHandler)
	targeting.Post("/update", validations.ValidateTargeting, rest.TargetingUpdateHandler)
	targeting.Post("/tags", rest.TargetingExportTagsHandler)
	// User management (only for users with 'admin' role)
	// users.Use(supertokenClient.AdminRoleRequired)
	users.Post("/get", userManagementSystem.UserGetHandler)
	users.Post("/set", validations.ValidateUser, userManagementSystem.UserSetHandler)
	users.Post("/update", validations.ValidateUser, userManagementSystem.UserUpdateHandler)

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

func loggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	c.Next()

	log.Info().
		Str("method", string(c.Request().Header.Method())).
		Str("url", c.Request().URI().String()).
		// Str("request", string(c.Request().Body())).
		// Str("response", string(c.Response().Body())).
		Str("duration", time.Since(start).String()).
		Msg("logging middleware")

	return nil
}
