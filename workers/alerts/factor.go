package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	slackmodule "github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"log"
	"time"
)

type Factor struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
}

func (f *Factor) Init(ctx context.Context, conf config.StringMap) error {

	f.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(f.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	f.Cron, _ = conf.GetStringValue("cron")
	fmt.Println()

	return nil
}

func (f *Factor) Do(ctx context.Context) error {

	//	db := bcdb.DB()

	//report, _ := getDataFromDB(ctx, db)

	slackMod := slackmodule.NewSlackModule()

	err := slackMod.SendMessage("Hello from slack module!")

	if err != nil {
		log.Fatalf("Error sending message to Slack: %v", err)
	}

	return nil
}

func getDataFromDB(ctx context.Context, db *sqlx.DB) (string, error) {
	records := make(models.PriceFactorLogSlice, 0)
	eval_time := time.Now().UTC().Truncate(time.Duration(30) * time.Minute)
	time_string := eval_time.Format("2006-01-02T15:04:05Z")

	sql := ` SELECT * FROM public.price_factor_log
            where eval_time >= TO_TIMESTAMP(time_string, 'YYYY-MM-DD HH24:MI:SS')
            and response_status != 400`

	fmt.Println(sql)
	query := fmt.Sprintf(sql, time_string)

	err := queries.Raw(query).Bind(ctx, db, &records)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f *Factor) GetSleep() int {
	if f.Cron != "" {
		return bccron.Next(f.Cron)
	}
	return 0
}
