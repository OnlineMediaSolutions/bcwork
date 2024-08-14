package alerts

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"github.com/m6yf/bcwork/bcdb"
//	"github.com/m6yf/bcwork/config"
//	"github.com/m6yf/bcwork/models"
//	"github.com/m6yf/bcwork/utils/bccron"
//	"github.com/rotisserie/eris"
//	"github.com/volatiletech/sqlboiler/v4/queries"
//	"time"
//)
//
//type Factor struct {
//	DatabaseEnv string `json:"dbenv"`
//	Cron        string `json:"cron"`
//}
//
//func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
//
//	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
//	err := bcdb.InitDB(w.DatabaseEnv)
//	if err != nil {
//		return eris.Wrapf(err, "Failed to initialize DB")
//	}
//
//	w.Cron, _ = conf.GetStringValue("cron")
//
//	return nil
//}
//
//func (w *Worker) Do(ctx context.Context) error {
//	//slackToken := viper.GetString("slack.token")
//
//	db := bcdb.DB()
//
//	report, _ := getDataFromDB(ctx, db)
//
//	fmt.Println("report", report)
//
//	//api := slack.New(slackToken)
//	//
//	//_, _, err := api.PostMessage(
//	//	"C07G8DB7DMX",
//	//	slack.MsgOptionText("Hello, Slack channel!", false),
//	//)
//
//	//if err != nil {
//	//	log.Fatalf("Error sending message to Slack: %v", err)
//	//}
//
//	fmt.Println("Message sent successfully")
//
//	return nil
//}
//
//func getDataFromDB(ctx context.Context, db *sqlx.DB) (string, error) {
//	records := make(models.PriceFactorLogSlice, 0)
//	eval_time := time.Now().UTC().Truncate(time.Duration(30) * time.Minute)
//	time_string := eval_time.Format("2006-01-02T15:04:05Z")
//
//	sql := ` SELECT * FROM public.price_factor_log
//             where eval_time >= TO_TIMESTAMP(time_string, 'YYYY-MM-DD HH24:MI:SS')
//             and response_status != 400`
//	query := fmt.Sprintf(sql, time_string)
//
//	err := queries.Raw(query).Bind(ctx, db, &records)
//	if err != nil {
//		return "", err
//	}
//
//	data, err := json.Marshal(records)
//	if err != nil {
//		return "", err
//	}
//
//	return string(data), nil
//}
//
//func (w *Worker) GetSleep() int {
//	if w.Cron != "" {
//		return bccron.Next(w.Cron)
//	}
//	return 0
//}
