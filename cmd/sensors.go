package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/sensors/core"
	"github.com/m6yf/bcwork/sensors/rest"

	//	"github.com/m6yf/bcwork/sensors/rest"
	//	_ "github.com/m6yf/bcwork/cmd/sensors/docs"
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

// @securityDefinitions.sensorskey CountKeyAuth
// @in header
// @name Authorization

// sensorsCmd represents sensors server command
var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Brightcom Count Server",
	Long:  ``,
	Run:   CountCmd,
}

func CountCmd(cmd *cobra.Command, args []string) {
	err := bcdb.InitDB("prod-do")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect init DB")
	}

	boil.DebugMode = false
	app := fiber.New(fiber.Config{
		AppName:   "Brightcom Sensors Digest",
		BodyLimit: 100 * 1024 * 1024,
	})
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("UP")
	})

	go core.SensorWorker()

	app.Post("/sensors/digest", rest.DigestSensorsHandler)
	app.Get("/select", rest.SelectHandler)
	app.Get("/sumcount", rest.SumCountHandler)

	//rest.Routes(app)

	//app.Static("/swagger", viper.GetString("assets") + "/swagger")
	//app.Get("/swagger/*", fiberSwagger.Handler)

	app.Listen(":8001")
}
