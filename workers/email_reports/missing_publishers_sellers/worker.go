package missing_publishers_sellers

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

type Worker struct {
	Cron          string                `json:"cron"`
	Slack         *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv   string                `json:"dbenv"`
	CompassClient *compass.Compass
	skipInitRun   bool
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return err
	}

	credentialsMap, err := config.FetchConfigValues([]string{"demand_sellers"})
	if err != nil {
		return fmt.Errorf("error fetching config values: %w", err)
	}

	creds := credentialsMap["demand_sellers"]

	var emailConfig EmailCreds
	err = json.Unmarshal([]byte(creds), &emailConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	log.Log().Msg(`Starting demand sellers worker`)
	if worker.skipInitRun {
		worker.skipInitRun = false
		log.Log().Msg("Skip init run")
		return nil
	}

	compassReport, err := getCompassData()
	if err != nil {
		return err
	}

	compassDemandData, err := getDemandData()
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}

	return 0
}
