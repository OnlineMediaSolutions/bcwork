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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/job"
	"github.com/m6yf/bcwork/structs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ArbitraryArgs,
	Run:  runExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(sensorsCmd)
}

type JobExecutionInfo struct {
	Class  string                 `json:"class"`
	Config map[string]interface{} `json:"aerospike.conf"`
}

//var queue *lane.Queue

func runExec(cmd *cobra.Command, args []string) {
	workerName := viper.GetString("worker.name")
	initLogger(workerName)

	//creat worker class
	worker, err := structs.NewInstance(workerName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create job worker instance ")
	}
	w := worker.(job.Worker)

	//initialize worker with configuration
	conf, err := parseArgs(args)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse args")
	}
	ctx := context.Background()

	err = w.Init(ctx, conf)
	if err != nil {
		log.Fatal().Err(err).Msg("worker initialization error")
	}

	for {
		err = w.Do(ctx)
		if err != nil {
			log.Error().Err(err).Msg("worker error")
		}
		if w.GetSleep() == 0 {
			break
		}
		time.Sleep(time.Duration(w.GetSleep()) * time.Second)
	}

	log.Info().Msg("worker off")
}

func parseArgs(args []string) (map[string]string, error) {
	res := map[string]string{}
	for _, arg := range args {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("failed to parse worker arguments '%s' is not in the format k=v", arg)
		}
		res[kv[0]] = kv[1]
	}

	return res, nil
}

func initLogger(workerName string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Logger.With().Str("worker.name", workerName).Logger()
}

//func OnJobMessage(ctx context.Context, msg *pubsub.Message) {
//	data := make(map[string]interface{})
//	err := json.Unmarshal(msg.Data, &data)
//	if err != nil {
//		log.Error().Err(err).Msg("failed to parse job incoming message")
//	}
//	queue.Enqueue(data)
//}
