package cmd

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/m6yf/bcwork/api/rest"
	"github.com/m6yf/bcwork/api/rest/report"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils/pointer"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"strings"

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

	// Should be one of "NONE" or "VIA_CUSTOM_HEADER" or "VIA_TOKEN"
	antiCsrf := "NONE"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: viper.GetString("supertokens.uri"),
			APIKey:        viper.GetString("supertokens.key"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         viper.GetString("supertokens.appInfo.appName"),
			APIDomain:       viper.GetString("supertokens.appInfo.apiDomain"),
			APIBasePath:     pointer.String(viper.GetString("supertokens.appInfo.apiBasePath")),
			WebsiteDomain:   viper.GetString("supertokens.appInfo.websiteDomain"),
			WebsiteBasePath: pointer.String(viper.GetString("supertokens.appInfo.websiteBasePath")),
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				/*
				   We use different credentials for different platforms when required. For example the redirect URI for Github
				   is different for Web and mobile. In such a case we can provide multiple providers with different client Ids.

				   When the frontend makes a request and wants to use a specific clientId, it needs to send the clientId to use in the
				   request. In the absence of a clientId in the request the SDK uses the default provider, indicated by `isDefault: true`.
				   When adding multiple providers for the same type (Google, Github etc), make sure to set `isDefault: true`.
				*/
				Providers: []tpmodels.ProviderInput{
					// We have provided you with development keys which you can use for testing.
					// IMPORTANT: Please replace them with your own OAuth keys for production use.
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
									ClientSecret: "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",
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID:     "467101b197249757c71f",
									ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "apple",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientID: "4398792-io.supertokens.example.service",
									AdditionalConfig: map[string]interface{}{
										"keyId":      "7M48Y4RYDL",
										"privateKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
										"teamId":     "YWQCXGJRJL",
									},
								},
							},
						},
					},
				},
			}),
			dashboard.Init(&dashboardmodels.TypeInput{
				ApiKey: viper.GetString("supertokens.dashboardApiKey"),
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &antiCsrf,
			}),
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to supertokens")
	}

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

	//adding the supertokens middleware
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))

	app.Get("/sessioninfo", verifySession(nil), sessioninfo)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("UP")
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/report/daily/revenue", report.DailyRevenue)
	app.Get("/report/hourly/revenue", report.HourlyRevenue)

	app.Get("/report/demand", rest.DemandReportGetHandler)
	app.Get("/report/demand/hourly", rest.DemandHourlyReportGetHandler)
	app.Get("/report/publisher", rest.PublisherReportGetHandler)
	app.Get("/report/publisher/hourly", rest.PublisherHourlyReportGetHandler)
	app.Get("/report/iiq/hourly", rest.IiqTestingGetHandler)
	app.Post("/metadata/update", rest.MetadataPostHandler)
	//app.Post("/price/factor", rest.PriceFactorPostHandler)

	app.Get("/price/factor/set", rest.PriceFactorSetHandler)
	app.Get("/price/factor/get", rest.PriceFactorGetHandler)
	app.Get("/price/factor/get/all", rest.PriceFactorGetAllHandler)
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
	app.Post("/confiant", rest.ConfiantPostHandler)
	app.Post("/confiant/get", rest.ConfiantGetAllHandler)
	app.Post("/block", rest.BlockPostHandler)
	app.Post("/block/get", rest.BlockGetAllHandler)

	app.Post("/dpo/set", rest.DemandPartnerOptimizationSetHandler)
	app.Get("/dpo/get", rest.DemandPartnerOptimizationGetHandler)
	app.Delete("/dpo/delete", rest.DemandPartnerOptimizationDeleteHandler)
	app.Get("/dpo/update", rest.DemandPartnerOptimizationUpdateHandler)

	app.Post("/publisher/new", rest.PublisherNewHandler)
	app.Post("/publisher/update", rest.PublisherUpdateHandler)
	app.Post("/publisher/get", rest.PublisherGetHandler)
	app.Post("/publisher/count", rest.PublisherCountHandler)

	app.Post("/price/factor/get", rest.FactorGetAllHandler)
	app.Post("/price/factor", rest.FactorPostHandler)

	app.Get("/ping", rest.PingPong)

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
	viper.SetDefault("supertokens.appInfo.appName", "OMS-API")
	viper.SetDefault("supertokens.appInfo.apiDomain", "http://localhost:8000")
	viper.SetDefault("supertokens.appInfo.apiBasePath", "/auth")
	viper.SetDefault("supertokens.appInfo.websiteDomain", "http://localhost:8001")
	viper.SetDefault("supertokens.appInfo.websiteBasePath", "/auth")

	err := viper.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warn().Msg("config file was not found, default values will be used")
	}

	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("/etc/bcwork/")
	//viper.AddConfigPath("$HOME/.")
	//viper.AddConfigPath(".")
	//viper.SetDefault("dbenv", "prod")
	//viper.SetDefault("ports.http", "8000")
	//viper.SetDefault("assets", "./api")
	//
	//err := viper.ReadInConfig() // Find and read the config file
	//if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	//	log.Warn().Msg("config file was not found, default values will be used")
	//} else if err != nil {
	//	log.Fatal().Msg(errors.WithStack(err).Error())
	//}

	// Init firebase auth
	//err = firebase.SetupFirebase()
	//if err != nil {
	//	log.Fatal().Err(errors.WithStack(err)).Msg("failed to init firebase auth")
	//}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// wrapper of the original implementation of verify session to match the required function signature
func verifySession(options *sessmodels.VerifySessionOptions) fiber.Handler {
	shouldCallNext := false
	return func(c *fiber.Ctx) error {
		err := adaptor.HTTPHandlerFunc(session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.SetUserContext(r.Context())
			authUserID := session.GetSessionFromRequestContext(r.Context()).GetUserID()
			c.Locals("auth_user_id", authUserID)
			c.Context().SetUserValue("auth_user_id", authUserID)
			userID, _ := core.AuthToUserID(c.Context(), authUserID, r.URL.Path != "/impersonate")
			if userID != "" {
				c.Locals("user_id", userID)
				c.Context().SetUserValue("user_id", userID)
			}
			shouldCallNext = true
		}))(c)
		if err != nil {
			return eris.Cause(err)
		}
		if shouldCallNext {
			return c.Next()
		}
		return nil
	}
}

func sessioninfo(c *fiber.Ctx) error {
	sessionContainer := session.GetSessionFromRequestContext(c.UserContext())
	if sessionContainer == nil {
		return c.Status(500).JSON("no session found")
	}
	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	currAccessTokenPayload := sessionContainer.GetAccessTokenPayload()
	counter, ok := currAccessTokenPayload["counter"]
	if !ok {
		counter = 1
	} else {
		counter = int(counter.(float64) + 1)
	}
	err = sessionContainer.MergeIntoAccessTokenPayload(map[string]interface{}{
		"counter": counter.(int),
	})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
}
