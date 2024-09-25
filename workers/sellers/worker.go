package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"time"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
}

type Competitor struct {
	Name string
	URL  string
}

type Config struct {
	EmailTo string
}

type SellersJSONHistory struct {
	CompetitorName        string           `db:"competitor_name"`
	AddedDomains          string           `db:"added_domains"`
	AddedPublishers       string           `db:"added_publishers"`
	BackupToday           *json.RawMessage `db:"backup_today"`
	BackupYesterday       *json.RawMessage `db:"backup_yesterday"`
	BackupBeforeYesterday *json.RawMessage `db:"backup_before_yesterday"`
	CreatedAt             time.Time        `db:"created_at"`
	UpdatedAt             time.Time        `db:"updated_at"`
	URL                   string
}

type Seller struct {
	SellerID   string `json:"seller_id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	SellerType string `json:"seller_type"`
}

type SellersJSON struct {
	Sellers []Seller `json:"sellers"`
}

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	if err := bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return eris.Wrapf(err, "Failed to initialize DB for sellers")
	}
	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	db := bcdb.DB()

	emailCredsMap, err := config.FetchConfigValues([]string{"sellers_json_crawler"})
	if err != nil {
		fmt.Println("Error fetching email credentials:", err)
		return nil
	}

	credsRaw, found := emailCredsMap["sellers_json_crawler"]
	if !found {
		fmt.Println("'sellers_json_crawler' not found in the config")
		return nil
	}

	var emailCreds EmailCreds
	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		fmt.Println("Error unmarshalling email credentials:", err)
		return nil
	}

	competitors, err := FetchCompetitors(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to fetch competitors: %w", err)
	}

	results := worker.PrepareCompetitors(competitors)

	history, err := worker.GetHistoryData(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to process competitors: %w", err)
	}

	var competitorsData []CompetitorData

	competitorsData, err = worker.prepareAndInsertCompetitors(ctx, results, history, db, competitorsData)
	if err != nil {
		return err
	}

	err = worker.prepareEmail(competitorsData, err, emailCreds)
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	fmt.Println("worker.Cron", worker.Cron)
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
