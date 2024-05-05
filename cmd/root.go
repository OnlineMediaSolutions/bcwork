/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "work",
	Short: "Run Gutsy data ETL job",
	Long:  `Gutsy ETL Job for pull ,transform and store data to be  used in Gutsy App`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command error!")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "aerospike.conf", "", "aerospike.conf file (default is ./aerospike.conf.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output for exec")

	rootCmd.PersistentFlags().StringP("dbenv", "", "", "default is local")
	viper.BindPFlag("database.env", rootCmd.PersistentFlags().Lookup("dbenv"))

	rootCmd.PersistentFlags().StringP("worker", "w", "", "worker class name or alias (mandatory)")
	viper.BindPFlag("worker.name", rootCmd.PersistentFlags().Lookup("worker"))
}

// initConfig reads in aerospike.conf file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use aerospike.conf file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search aerospike.conf in home directory with name ".work" (without extension).
		viper.SetConfigName("config.yaml")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/oms/")
		viper.AddConfigPath("$HOME/.")
		viper.AddConfigPath(".")

	}

	//	viper.SetDefault("database.env", "local")
	viper.AutomaticEnv() // read in environment variables that adq

	// If a aerospike.conf file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("aerospike.conf", viper.ConfigFileUsed()).Msgf("aerospike.conf file: %s", viper.ConfigFileUsed())
	}
}
