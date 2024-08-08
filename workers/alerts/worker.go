package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/rotisserie/eris"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"log"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	slackToken := viper.GetString("slack.token")

	db := bcdb.DB()

	report, _ := getDataFromDB(ctx, db)

	fmt.Println("report", report)

	api := slack.New(slackToken)

	_, _, err := api.PostMessage(
		"C07G8DB7DMX",
		slack.MsgOptionText("Hello, Slack channel!", false),
	)

	if err != nil {
		log.Fatalf("Error sending message to Slack: %v", err)
	}

	fmt.Println("Message sent successfully")

	return nil
}

func getDataFromDB(ctx context.Context, db *sqlx.DB) (string, error) {
	records := make(models.PriceFactorLogSlice, 0)
	sql := `SELECT * FROM price_factor_log`

	err := queries.Raw(sql).Bind(ctx, db, &records)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(records)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (w *Worker) GetSleep() int {
	return 0
}
