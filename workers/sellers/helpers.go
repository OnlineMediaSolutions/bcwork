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

func fetchDataFromWebsite(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func InsertCompetitor(ctx context.Context, db *sqlx.DB, name string, addedDomains, addedPublishers []string, backupToday, backupYesterday interface{}) error {
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

	_, err = db.ExecContext(ctx,
		`INSERT INTO sellers_json_history (competitor_name, added_domains, added_publishers, backup_today, backup_yesterday, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (competitor_name)
		DO UPDATE SET added_domains = $2, added_publishers = $3, backup_today = $4, backup_yesterday = $5, updated_at = NOW()`,
		name, addedDomainsStr, addedPublishersStr, backupTodayJSON, backupYesterdayJSON)

	return eris.Wrap(err, "failed to insert competitor")
}

func (worker *Worker) Request(jobs <-chan Competitor, results chan<- map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		data, err := fetchDataFromWebsite(job.URL)
		if err != nil {
			log.Printf("Error fetching data for competitor %s: %v", job.Name, err)
			continue
		}

		results <- map[string]interface{}{job.Name: data}
	}
}

func (worker *Worker) GetHistoryData(ctx context.Context, db *sqlx.DB) ([]SellersJSONHistory, error) {
	rows, err := db.QueryContext(ctx, "SELECT competitor_name, added_domains, added_publishers, backup_today, backup_yesterday, created_at, updated_at FROM sellers_json_history")
	if err != nil {
		return nil, fmt.Errorf("failed to query sellers_json_history: %w", err)
	}
	defer rows.Close()

	var histories []SellersJSONHistory
	for rows.Next() {
		var history SellersJSONHistory
		err := rows.Scan(
			&history.CompetitorName,
			&history.AddedDomains,
			&history.AddedPublishers,
			&history.BackupToday,
			&history.BackupYesterday,
			&history.CreatedAt,
			&history.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}
		histories = append(histories, history)
	}

	return histories, nil
}

func compareSellers(json1, json2 SellersJSON) (extraPublishers []string, extraDomains []string) {
	sellerMap1 := make(map[string]struct{})

	for _, seller := range json1.Sellers {
		key := seller.Domain + ":" + seller.Name
		sellerMap1[key] = struct{}{}
	}

	for _, seller := range json2.Sellers {
		key := seller.Domain + ":" + seller.Name

		if _, exists := sellerMap1[key]; !exists {
			extraPublishers = append(extraPublishers, seller.Name)
			extraDomains = append(extraDomains, seller.Domain)
		}
	}

	return extraPublishers, extraDomains
}
