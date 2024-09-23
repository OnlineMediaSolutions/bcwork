package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rotisserie/eris"
	"sync"
	"time"
)

type Competitor struct {
	Name string
	URL  string
}

type SellersJSONHistory struct {
	CompetitorName  string           `db:"competitor_name"`
	AddedDomains    string           `db:"added_domains"`
	AddedPublishers string           `db:"added_publishers"`
	BackupToday     *json.RawMessage `db:"backup_today"`
	BackupYesterday *json.RawMessage `db:"backup_yesterday"`
	CreatedAt       time.Time        `db:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at"`
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

type Worker struct {
	DatabaseEnv string            `json:"dbenv"`
	Cron        map[string]string `json:"cron"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	if err := bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}
	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	db := bcdb.DB()

	competitors, err := FetchCompetitors(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to fetch competitors: %w", err)
	}

	const numWorkers = 5
	var wg sync.WaitGroup
	jobs := make(chan Competitor, len(competitors))
	results := make(chan map[string]interface{}, len(competitors))

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker.Request(jobs, results, &wg)
	}

	for _, competitor := range competitors {
		jobs <- competitor
	}

	close(jobs)
	wg.Wait()
	close(results)

	history, err := worker.GetHistoryData(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to process competitors: %w", err)
	}
	var competitorsData []CompetitorData

	for result := range results {
		for name, backupToday := range result {
			var historyRecord SellersJSONHistory
			found := false

			for _, h := range history {
				if h.CompetitorName == name {
					historyRecord = h
					found = true
					break
				}
			}

			backupTodayMap := backupToday.(map[string]interface{})

			jsonData, err := json.Marshal(backupTodayMap)
			if err != nil {
				return fmt.Errorf("failed to marshal map to JSON for competitor %s: %w", name, err)
			}

			var backupTodayData SellersJSON
			if err := json.Unmarshal(jsonData, &backupTodayData); err != nil {
				return fmt.Errorf("failed to unmarshal map data to SellersJSON for competitor %s: %w", name, err)
			}

			if !found {
				if err := InsertCompetitor(ctx, db, name, []string{}, []string{}, backupTodayData, nil); err != nil {
					return fmt.Errorf("failed to insert new competitor: %w", err)
				}
				continue
			}

			var historyBackupToday SellersJSON
			if err := json.Unmarshal(*historyRecord.BackupToday, &historyBackupToday); err != nil {
				return fmt.Errorf("failed to unmarshal BackupToday from history for competitor %s: %w", name, err)
			}

			addedPublishers, addedDomains := compareSellers(historyBackupToday, backupTodayData)

			if addedDomains != nil || addedPublishers != nil {
				competitorsData = append(competitorsData, CompetitorData{
					Name:       name,
					Publishers: addedPublishers,
					Domains:    addedDomains,
				})
			}

			if err != nil {
				return fmt.Errorf("failed to compare data for competitor %s: %w", name, err)
			}

			if err := InsertCompetitor(ctx, db, name, addedDomains, addedPublishers, backupTodayData, historyBackupToday); err != nil {
				return fmt.Errorf("failed to insert competitor data for %s: %w", name, err)
			}
		}
	}

	if len(competitorsData) > 0 {
		now := time.Now()
		today := now.Format("2006-01-02")
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

		subject := fmt.Sprintf("Competitors sellers.json daily changes - %s", today)
		emailBody := fmt.Sprintf("Below are the sellers.json changes between - %s and %s", yesterday, today)

		err = SendCustomHTMLEmail("sonai@onlinemediasolutions.com", "sonai@onlinemediasolutions.com", subject, emailBody, competitorsData)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}

func (worker *Worker) GetSleep() int {
	return 0
}
