package cmd

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/kvdb"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// @title Swagger Brightcom API
// @version 1.0
// @description API for Brightcom game.

// @contact.name Brightcom Support
// @contact.url http://www.nanoook.com/support
// @contact.email support@gutsy.me

// @securityDefinitions.kvdbkey KeyValueDBKeyAuth
// @in header
// @name Authorization

// kvdbCmd represents kvdb server command
var kvdbCmd = &cobra.Command{
	Use:   "kvdb",
	Short: "Brightcom Key Value DB Server",
	Long:  ``,
	Run:   KeyValueDBCmd,
}

func KeyValueDBCmd(cmd *cobra.Command, args []string) {
	go kvdb.KvLoop()

	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Get("/get", func(c *fiber.Ctx) error {
		return c.SendString(kvdb.Get(c.Query("k")))
	})

	app.Get("/save", func(c *fiber.Ctx) error {
		if kvdb.Saving > 0 {
			return c.SendString(strconv.Itoa(kvdb.Saving))
		}
		go func() {
			_, err := kvdb.Save("/root/kvdb.db")
			if err != nil {
				log.Error().Err(err).Msg("error saving data to key-value db")
			}
		}()

		return c.SendString("1")
	})

	app.Get("/load", func(c *fiber.Ctx) error {
		go func() {
			err := kvdb.Load("/root/kvdb.db")
			if err != nil {
				log.Error().Err(err).Msg("error loading data for key-value db")
			}
		}()

		return c.SendString("1")
	})

	app.Get("/scan", func(c *fiber.Ctx) error {
		go kvdb.Scan()
		return c.SendStatus(http.StatusOK)
	})

	app.Get("/count", func(c *fiber.Ctx) error {
		return c.SendString(strconv.Itoa(kvdb.Count()))
	})

	app.Post("/set", func(c *fiber.Ctx) error {
		body := c.Body()
		k := c.Query("k")
		if k != "" && len(body) > 0 {
			kvdb.Q <- kvdb.Pair{K: strings.Clone(string(k)), V: strings.Clone(string(body))}
		}

		return c.SendStatus(http.StatusOK)
	})

	err := app.Listen(":8090")
	if err != nil {
		log.Error().Err(err).Msg("failed to bind")
	}
}

func init() {
	rootCmd.AddCommand(kvdbCmd)

	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("/etc/bcwork/")
	//viper.AddConfigPath("$HOME/.")
	//viper.AddConfigPath(".")
	//viper.SetDefault("dbenv", "prod")
	//viper.SetDefault("ports.http", "8000")
	//viper.SetDefault("assets", "./kvdb")
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
