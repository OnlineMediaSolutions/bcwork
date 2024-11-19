package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("failed to fetch competitors")
		return err
	}

	emailCredsMap, err := config.FetchConfigValues([]string{"sellers_json_crawler_web", "sellers_json_crawler_inapp"})
	if err != nil {
		log.Error().Err(err).Msg("Error fetching email credentials")
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
			log.Error().Msg("email credentials not found")
			continue
		}

		if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
			log.Error().Err(err).Msg("Error unmarshalling email credentials")
			continue
		}

		var competitorsData []CompetitorData
		results := worker.PrepareCompetitors(competitorsGroup)
		history, err := worker.GetHistoryData(ctx, db)
		if err != nil {
			log.Error().Err(err).Msg("failed to process competitors")
			return err
		}

		positionMap := make(map[string]string)

		for _, competitor := range competitors {
			positionMap[competitor.Name] = competitor.Position
		}

		competitorsData, err = worker.prepareAndInsertCompetitors(ctx, results, history, db, competitorsData, positionMap)
		if err != nil {
			return err
		}

		if len(competitorsData) > 0 {
			err = worker.prepareEmail(competitorsData, nil, emailCreds, competitorType)
			if err != nil {
				message := fmt.Sprintf("Error sending email for type %s: %v", competitorType, err)
				log.Error().Err(err).Msg(message)
				continue
			}
			log.Info().Msg(fmt.Sprintf("Email sent successfully for type %s", competitorType))
		} else {
			log.Info().Msg(fmt.Sprintf("No competitors data to send for type %s", competitorType))
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
