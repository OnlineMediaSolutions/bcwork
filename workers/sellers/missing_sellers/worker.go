package missing_sellers

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/compass"
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
	Cron          string `json:"cron"`
	DatabaseEnv   string `json:"dbenv"`
	CompassClient *compass.Compass
	skipInitRun   bool
	emailConfig   EmailCreds
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return err
	}

	credentialsMap, err := config.FetchConfigValues([]string{"missing_sellers"})
	if err != nil {
		return fmt.Errorf("error fetching config values: %w", err)
	}

	creds := credentialsMap["missing_sellers"]

	var emailConfig EmailCreds
	err = json.Unmarshal([]byte(creds), &emailConfig)

	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.emailConfig = emailConfig

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	log.Log().Msg(`Starting missing sellers worker`)

	if worker.skipInitRun {
		worker.skipInitRun = false
		log.Log().Msg("Skip init run")

		return nil
	}

	compassData, err := fetchCompassData()
	if err != nil {
		return fmt.Errorf("error getting compass data: %w", err)
	}

	compassDemandData, err := fetchDemandData()
	if err != nil {
		return fmt.Errorf("error getting demand data: %w", err)
	}

	todaySellersData, yesterdaySellersData, err := getSellersJsonFiles(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("error getting  today, yestarday json files: %w", err)
	}

	err = insert(ctx, todaySellersData, err)
	if err != nil {
		return fmt.Errorf("error inserting today sellers: %w", err)
	}

	compassDataSet := createCompassDataSet(compassData, compassDemandData)

	statusMap := findMissingIds(compassDataSet, todaySellersData, yesterdaySellersData)

	err = prepareEmailAndSend(statusMap, worker.emailConfig)
	if err != nil {
		return fmt.Errorf("error preparing email: %w", err)
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}

	return 0
}
