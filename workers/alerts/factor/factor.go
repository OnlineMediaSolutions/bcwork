package factor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
)

type Factor struct {
	DatabaseEnv string                `json:"dbenv"`
	Cron        string                `json:"cron"`
	Slack       *messager.SlackModule `json:"slack_instances"`
}

type Report struct {
	Time           string                `json:"time"`
	EvalTime       string                `json:"eval_time"`
	PubImps        int                   `json:"pubImps"`
	SoldImps       int                   `json:"soldImps"`
	Cost           float64               `json:"cost"`
	Revenue        float64               `json:"revenue"`
	GP             float64               `json:"gp"`
	GPP            float64               `json:"gpp"`
	Publisher      string                `json:"publisher"`
	Domain         string                `json:"domain"`
	Country        string                `json:"country"`
	Device         string                `json:"device"`
	OldFactor      float64               `json:"oldFactor"`
	NewFactor      float64               `json:"newFactor"`
	ResponseStatus int                   `json:"responseStatus"`
	Increase       float64               `json:"increase"`
	Slack          *messager.SlackModule `json:"slack_instances"`
}

func (factor *Factor) Init(conf config.StringMap) error {
	factor.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(factor.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}

	factor.Cron, _ = conf.GetStringValue("cron")
	factor.Slack, err = messager.NewSlackModule()
	if err != nil {
		return eris.Wrapf(err, "failed to initialize slack module")
	}

	return nil
}

func (factor *Factor) Do(ctx context.Context) error {
	var reports []Report
	db := bcdb.DB()
	slackMod, err := messager.NewSlackModule()
	if err != nil {
		return eris.Wrapf(err, "failed to initialize slack module")
	}

	jsonData, err := getDataFromDB(ctx, db)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(jsonData), &reports)
	if err != nil {
		return fmt.Errorf("Error unmarshalling JSON for factor logs: %w", err)
	}

	csvData, err := ConvertReportsToCSV(reports)
	if err != nil {
		return fmt.Errorf("error converting reports to csv: %w", err)
	}

	reader := csv.NewReader(strings.NewReader(csvData))

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Error reading CSV data for factor logs: %w", err)
	}

	if len(rows) > 1 {
		err = slackMod.SendMessage("Factor Logs:\n" + "```" + csvData + "```")
	} else {
		fmt.Println("Received empty body for csv file factor logs")

		return nil
	}

	if err != nil {
		return fmt.Errorf("Error sending CSV file to Slack for factor logs: %w", err)
	}

	return fmt.Errorf("CSV file for factor logs sent successfully via slack")
}

func (factor *Factor) GetSleep() int {
	if factor.Cron != "" {
		return bccron.Next(factor.Cron)
	}

	return 0
}
