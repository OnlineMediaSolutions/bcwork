package missing_publishers

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

	credentialsMap, err := config.FetchConfigValues([]string{"missing_publishers_sellers"})
	if err != nil {
		return fmt.Errorf("error fetching config values: %w", err)
	}

	creds := credentialsMap["missing_publishers_sellers"]

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
	log.Log().Msg(`Starting missing publishers worker`)
	if worker.skipInitRun {
		worker.skipInitRun = false
		log.Log().Msg("Skip init run")
		return nil
	}

	compassData, err := getCompassData()
	if err != nil {
		return err
	}

	compassDemandData, err := getDemandData()
	fmt.Println(compassDemandData)

	demandSellersData, yesterdaySellersData, err := getSellersJsonFiles(ctx, bcdb.DB())

	// TODO - Insert to DB

	if err != nil {
		return err
	}

	compassDataSet := createCompassDataSet(compassData)

	data := findMissingIds(compassDataSet, demandSellersData, yesterdaySellersData)

	fmt.Println(data)
	//	sellersFiles := sellers.FetchDataFromWebsite(demandUrls)
	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}

	return 0
}
