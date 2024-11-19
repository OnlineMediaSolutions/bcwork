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
	skipInitRun bool
}

type Competitor struct {
	Name     string
	URL      string
	Type     string
	Position string
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
	SellerID   interface{} `json:"seller_id"`
	Name       string      `json:"name"`
	Domain     string      `json:"domain"`
	SellerType string      `json:"seller_type"`
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
	worker.Cron, _ = conf.GetStringValue("cron")

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")

	if err := bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return eris.Wrapf(err, "Failed to initialize DB for sellers")
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {

	if worker.skipInitRun {
		fmt.Println("Skipping work as per the skip_init_run flag.")
		worker.skipInitRun = false
		return nil
	}

	db := bcdb.DB()

	competitors, err := FetchCompetitors(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to fetch competitors: %w", err)
	}

	emailCredsMap, err := config.FetchConfigValues([]string{"sellers_json_crawler_web", "sellers_json_crawler_inapp"})
	if err != nil {
		fmt.Println("Error fetching email credentials:", err)
		return nil
	}

	competitorsByType := make(map[string][]Competitor)
	for _, competitor := range competitors {
		competitorsByType[competitor.Type] = append(competitorsByType[competitor.Type], competitor)
	}

	for competitorType, competitorsGroup := range competitorsByType {
		var emailCreds EmailCreds
		var credsRaw string
		var found bool

		switch competitorType {
		case "inapp":
			credsRaw, found = emailCredsMap["sellers_json_crawler_inapp"]
		default:
			credsRaw, found = emailCredsMap["sellers_json_crawler_web"]
		}

		if !found {
			fmt.Printf("Email credentials not found for type %s\n", competitorType)
			continue
		}

		if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
			fmt.Printf("Error unmarshalling email credentials for type %s: %v\n", competitorType, err)
			continue
		}

		var competitorsData []CompetitorData
		results := worker.PrepareCompetitors(competitorsGroup)
		history, err := worker.GetHistoryData(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to process competitors: %w", err)
		}

		positionMap := make(map[string]string)

		for _, competitor := range competitors {
			positionMap[competitor.Name] = competitor.Position
		}

		competitorsData, err = worker.prepareAndInsertCompetitors(ctx, results, history, db, competitorsData, positionMap)
		if err != nil {
			return err
		}

		var competitorsEmailData []CompetitorData

		for _, competitor := range competitorsData {
			if len(competitor.AddedPublisherDomain) > 0 || len(competitor.DeletedPublisherDomain) > 0 {
				competitorsEmailData = append(competitorsEmailData, competitor)
			}
		}

		if len(competitorsEmailData) > 0 {
			err = worker.prepareEmail(competitorsData, nil, emailCreds, competitorType)
			if err != nil {
				fmt.Printf("Error sending email for type %s: %v\n", competitorType, err)
				continue
			}
			fmt.Printf("Email sent successfully for type %s\n", competitorType)
		} else {
			fmt.Printf("No competitors data to send for type %s\n", competitorType)
		}
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
