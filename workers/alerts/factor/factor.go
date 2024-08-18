package factor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"strings"
)

type Factor struct {
	DatabaseEnv string               `json:"dbenv"`
	Cron        string               `json:"cron"`
	Slack       *modules.SlackModule `json:"slack_instances"`
}

type Report struct {
	Time           string               `json:"time"`
	EvalTime       string               `json:"eval_time"`
	PubImps        int                  `json:"pubImps"`
	SoldImps       int                  `json:"soldImps"`
	Cost           float64              `json:"cost"`
	Revenue        float64              `json:"revenue"`
	GP             float64              `json:"gp"`
	GPP            float64              `json:"gpp"`
	Publisher      string               `json:"publisher"`
	Domain         string               `json:"domain"`
	Country        string               `json:"country"`
	Device         string               `json:"device"`
	OldFactor      float64              `json:"oldFactor"`
	NewFactor      float64              `json:"newFactor"`
	ResponseStatus int                  `json:"responseStatus"`
	Increase       float64              `json:"increase"`
	Slack          *modules.SlackModule `json:"slack_instances"`
}

func (f *Factor) Init(conf config.StringMap) error {

	f.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(f.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	f.Cron, _ = conf.GetStringValue("cron")
	f.Slack, err = modules.NewSlackModule()

	return nil
}

func (f *Factor) Do(ctx context.Context) error {
	var reports []Report
	db := bcdb.DB()
	slackMod, err := modules.NewSlackModule()
	jsonData, err := getDataFromDB(ctx, db)

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(jsonData), &reports)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON for factor logs: %w", err)
	}

	csvData, err := ConvertReportsToCSV(reports)

	reader := csv.NewReader(strings.NewReader(csvData))

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Error reading CSV data for factor logs:", err)
	}

	if len(rows) > 1 {
		err = slackMod.SendMessage("Factor Logs:\n" + "```" + csvData + "```")
	} else {
		fmt.Println("Received empty body for csv file factor logs")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error sending CSV file to Slack for factor logs: %v\n", err)
	}

	return fmt.Errorf("CSV file for factor logs sent successfully via slack")

}

func (f *Factor) GetSleep() int {
	if f.Cron != "" {
		return bccron.Next(f.Cron)
	}
	return 0
}
