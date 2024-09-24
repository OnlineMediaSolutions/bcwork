package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rotisserie/eris"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func FetchCompetitors(ctx context.Context, db *sqlx.DB) ([]Competitor, error) {
	rows, err := db.QueryContext(ctx, "SELECT name, url FROM competitors")
	if err != nil {
		return nil, eris.Wrap(err, "failed to fetch competitors")
	}
	defer rows.Close()

	var competitors []Competitor
	for rows.Next() {
		var competitor Competitor
		if err := rows.Scan(&competitor.Name, &competitor.URL); err != nil {
			return nil, eris.Wrap(err, "failed to scan row")
		}
		competitors = append(competitors, competitor)
	}
	return competitors, nil
}

func FetchDataFromWebsite(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if sellers, ok := data["sellers"]; ok {
		if _, err := CheckSellersArray(sellers); err != nil {
			return nil, fmt.Errorf("invalid sellers format: %w", err)
		}
	} else {
		return nil, fmt.Errorf("sellers array not found in the response")
	}

	return data, nil
}

func InsertCompetitor(ctx context.Context, db *sqlx.DB, name string, addedDomains, addedPublishers []string, backupToday, backupYesterday, backupBeforeYesterday interface{}) error {
	addedDomainsStr := strings.Join(addedDomains, ",")
	addedPublishersStr := strings.Join(addedPublishers, ",")

	backupTodayJSON, err := json.Marshal(backupToday)
	if err != nil {
		return fmt.Errorf("failed to marshal backupToday: %w", err)
	}
	backupYesterdayJSON, err := json.Marshal(backupYesterday)
	if err != nil {
		return fmt.Errorf("failed to marshal backupYesterday: %w", err)
	}
	backupBeforeYesterdayJSON, err := json.Marshal(backupBeforeYesterday)
	if err != nil {
		return fmt.Errorf("failed to marshal backupBeforeYesterday: %w", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO sellers_json_history (competitor_name, added_domains, added_publishers, backup_today, backup_yesterday, backup_before_yesterday, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        ON CONFLICT (competitor_name)
        DO UPDATE SET added_domains = $2, added_publishers = $3, backup_today = $4, backup_yesterday = $5, backup_before_yesterday = $6, updated_at = NOW()`,
		name, addedDomainsStr, addedPublishersStr, backupTodayJSON, backupYesterdayJSON, backupBeforeYesterdayJSON)

	return eris.Wrap(err, "failed to insert competitor")
}

func (worker *Worker) Request(jobs <-chan Competitor, results chan<- map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		data, err := FetchDataFromWebsite(job.URL)
		if err != nil {
			log.Printf("Error fetching data for competitor %s: %v", job.Name, err)
			continue
		}

		results <- map[string]interface{}{job.Name: data}
	}
}

func (worker *Worker) GetHistoryData(ctx context.Context, db *sqlx.DB) ([]SellersJSONHistory, error) {
	query := `
        SELECT 
            h.competitor_name, 
            h.added_domains, 
            h.added_publishers, 
            h.backup_today, 
            h.backup_yesterday, 
            h.backup_before_yesterday,
            h.created_at, 
            h.updated_at, 
            c.url 
        FROM 
            sellers_json_history h
        JOIN 
            competitors c ON h.competitor_name = c.name
    `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sellers_json_history: %w", err)
	}
	defer rows.Close()

	var histories []SellersJSONHistory
	for rows.Next() {
		var history SellersJSONHistory
		var url string
		err := rows.Scan(
			&history.CompetitorName,
			&history.AddedDomains,
			&history.AddedPublishers,
			&history.BackupToday,
			&history.BackupYesterday,
			&history.BackupBeforeYesterday,
			&history.CreatedAt,
			&history.UpdatedAt,
			&url,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}

		history.URL = url
		histories = append(histories, history)
	}

	return histories, nil
}

func compareSellers(historyBackupToday, backupTodayData SellersJSON) (extraPublishers []string, extraDomains []string) {
	sellerMapToday := make(map[string]struct{})

	for _, seller := range historyBackupToday.Sellers {
		key := seller.Domain + ":" + seller.Name
		sellerMapToday[key] = struct{}{}
	}

	for _, seller := range backupTodayData.Sellers {
		key := seller.Domain + ":" + seller.Name

		if _, exists := sellerMapToday[key]; !exists {
			extraPublishers = append(extraPublishers, seller.Name)
			extraDomains = append(extraDomains, seller.Domain)

		}
	}

	return extraPublishers, extraDomains
}

func (worker *Worker) PrepareCompetitors(competitors []Competitor) chan map[string]interface{} {
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
	return results
}

func (worker *Worker) prepareEmail(competitorsData []CompetitorData, err error) error {
	if len(competitorsData) > 0 {
		now := time.Now()
		today := now.Format("2006-01-02")
		yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

		subject := fmt.Sprintf("Competitors sellers.json daily changes - %s", today)
		message := fmt.Sprintf("Below are the sellers.json changes between - %s and %s", yesterday, today)

		err = SendCustomHTMLEmail("sonai@onlinemediasolutions.com", "sonai@onlinemediasolutions.com", subject, message, competitorsData)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}
	return nil
}

func (worker *Worker) prepareAndInsertCompetitors(ctx context.Context, results chan map[string]interface{}, history []SellersJSONHistory, db *sqlx.DB, competitorsData []CompetitorData) ([]CompetitorData, error) {
	historyMap := make(map[string]SellersJSONHistory)
	for _, h := range history {
		historyMap[h.CompetitorName] = h
	}

	for result := range results {
		for name, backupToday := range result {
			var historyRecord SellersJSONHistory
			if record, found := historyMap[name]; found {
				historyRecord = record
			}

			backupTodayData, historyBackupToday, err := MapBackupTodayData(backupToday, historyRecord)
			if err != nil {
				return nil, fmt.Errorf("Error processing backup data for competitor %s: %w", name, err)
			}

			addedPublishers, addedDomains := compareSellers(historyBackupToday, backupTodayData)

			if addedDomains != nil || addedPublishers != nil {
				competitorsData = append(competitorsData, CompetitorData{
					Name:       name,
					URL:        historyMap[name].URL,
					Publishers: addedPublishers,
					Domains:    addedDomains,
				})
			}

			backupBeforeYesterday := historyRecord.BackupYesterday
			if err := InsertCompetitor(ctx, db, name, addedDomains, addedPublishers, backupTodayData, historyBackupToday, backupBeforeYesterday); err != nil {
				return nil, fmt.Errorf("failed to insert competitor data for %s: %w", name, err)
			}
		}
	}
	return competitorsData, nil
}

func MapBackupTodayData(backupToday interface{}, historyRecord SellersJSONHistory) (SellersJSON, SellersJSON, error) {
	backupTodayMap, ok := backupToday.(map[string]interface{})
	if !ok {
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("invalid backupToday format")
	}

	jsonData, err := json.Marshal(backupTodayMap)
	if err != nil {
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	var backupTodayData SellersJSON
	if err := json.Unmarshal(jsonData, &backupTodayData); err != nil {
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to unmarshal map data to SellersJSON: %w", err)
	}

	var historyBackupToday SellersJSON
	if historyRecord.BackupToday != nil {
		if err := json.Unmarshal(*historyRecord.BackupToday, &historyBackupToday); err != nil {
			return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to unmarshal BackupToday from history: %w", err)
		}
	}

	return backupTodayData, historyBackupToday, nil
}

func CheckSellersArray(sellers interface{}) ([]interface{}, error) {
	if sellersArray, ok := sellers.([]interface{}); ok {
		return sellersArray, nil
	}
	return nil, fmt.Errorf("sellers should be an array, but got %T", sellers)
}
