package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/rotisserie/eris"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type ComparisonResult struct {
	ExtraPublishers   []string
	ExtraDomains      []string
	SellerType        string
	DeletedPublishers []string
	DeletedDomains    []string
}

func FetchCompetitors(ctx context.Context, db *sqlx.DB) ([]Competitor, error) {
	competitorModels, err := models.Competitors(qm.Select("name, url,type,position ")).All(ctx, db)
	if err != nil {
		return nil, err
	}

	competitors := make([]Competitor, len(competitorModels))
	for i, c := range competitorModels {
		competitors[i] = Competitor{
			Name:     c.Name,
			URL:      c.URL,
			Type:     c.Type,
			Position: c.Position,
		}
	}

	return competitors, nil
}

func FetchDataFromWebsite(url string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {

		return nil, fmt.Errorf("failed to decode JSON: %w", err)
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

func InsertCompetitor(ctx context.Context, db boil.ContextExecutor, name string, addedDomains, addedPublishers []string, backupToday, backupYesterday, backupBeforeYesterday interface{}, deletedPublishers []string, deletedDomains []string) error {
	addedDomainsStr := strings.Join(addedDomains, ",")
	addedPublishersStr := strings.Join(addedPublishers, ",")
	deletedPublishersStr := strings.Join(deletedPublishers, ",")
	deletedDomainsStr := strings.Join(deletedDomains, ",")

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

	history := &models.SellersJSONHistory{
		CompetitorName:        name,
		AddedDomains:          addedDomainsStr,
		AddedPublishers:       addedPublishersStr,
		BackupToday:           backupTodayJSON,
		BackupYesterday:       backupYesterdayJSON,
		BackupBeforeYesterday: backupBeforeYesterdayJSON,
		DeletedPublishers:     deletedPublishersStr,
		DeletedDomains:        deletedDomainsStr,
	}

	err = history.Upsert(ctx, db, true, []string{"competitor_name"}, boil.Whitelist("added_domains", "added_publishers", "backup_today", "backup_yesterday", "backup_before_yesterday", "deleted_publishers", "deleted_domains", "updated_at"),
		boil.Infer())
	if err != nil {
		return eris.Wrap(err, "failed to insert or update competitor")
	}

	return nil
}

func (worker *Worker) Request(jobs <-chan Competitor, results chan<- map[string]interface{}, failedCompetitors chan<- Competitor, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		data, err := FetchDataFromWebsite(job.URL)
		if err != nil {
			log.Printf("Error fetching data for competitor %s: %v", job.Name, err)
			failedCompetitors <- job
			continue
		}

		results <- map[string]interface{}{job.Name: data}
	}
}

func (worker *Worker) GetHistoryData(ctx context.Context, db *sqlx.DB) ([]SellersJSONHistory, error) {
	histories, err := models.SellersJSONHistories().All(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("failed to query sellers_json_history: %w", err)
	}

	var results []SellersJSONHistory
	for _, history := range histories {
		competitor, err := history.CompetitorNameCompetitor().One(ctx, db)
		if err != nil {
			return nil, fmt.Errorf("failed to query competitor: %w", err)
		}

		results = append(results, SellersJSONHistory{
			CompetitorName:        history.CompetitorName,
			AddedDomains:          history.AddedDomains,
			AddedPublishers:       history.AddedPublishers,
			BackupToday:           (*json.RawMessage)(&history.BackupToday),
			BackupYesterday:       (*json.RawMessage)(&history.BackupYesterday),
			BackupBeforeYesterday: (*json.RawMessage)(&history.BackupBeforeYesterday),
			CreatedAt:             history.CreatedAt.Time,
			UpdatedAt:             history.UpdatedAt.Time,
			URL:                   competitor.URL,
		})
	}

	return results, nil
}

func normalizeKey(domain, name string) string {
	domain = strings.ReplaceAll(domain, "http://", "")
	domain = strings.ReplaceAll(domain, "https://", "")
	return strings.TrimSpace(strings.ToLower(domain)) + ":" + strings.TrimSpace(strings.ToLower(name))
}

func compareSellers(todayData, historyData SellersJSON) ComparisonResult {
	sellerMapHistory := make(map[string]struct{})
	sellerMapToday := make(map[string]struct{})

	for _, seller := range historyData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name)
		sellerMapHistory[key] = struct{}{}
	}

	for _, seller := range todayData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name)
		sellerMapToday[key] = struct{}{}
	}

	var extraPublishers []string
	var extraDomains []string
	var sellerType string
	var deletedPublishers []string
	var deletedDomains []string

	for _, seller := range todayData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name)

		if _, exists := sellerMapHistory[key]; !exists {
			extraPublishers = append(extraPublishers, seller.Name)
			extraDomains = append(extraDomains, seller.Domain)
			sellerType = seller.SellerType

		}
	}

	for _, seller := range historyData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name)

		if _, exists := sellerMapToday[key]; !exists {
			deletedPublishers = append(deletedPublishers, seller.Name)
			deletedDomains = append(deletedDomains, seller.Domain)
		}
	}

	return ComparisonResult{
		ExtraPublishers:   extraPublishers,
		ExtraDomains:      extraDomains,
		SellerType:        sellerType,
		DeletedPublishers: deletedPublishers,
		DeletedDomains:    deletedDomains,
	}
}

func (worker *Worker) PrepareCompetitors(competitors []Competitor) chan map[string]interface{} {
	const numWorkers = 5
	var wg sync.WaitGroup
	jobs := make(chan Competitor, len(competitors))
	results := make(chan map[string]interface{}, len(competitors))
	failedCompetitors := make(chan Competitor, len(competitors))

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker.Request(jobs, results, failedCompetitors, &wg)
	}

	for _, competitor := range competitors {
		jobs <- competitor
	}

	close(jobs)
	wg.Wait()
	close(results)

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		worker.SendSlackMessageToFailedCompetitors(failedCompetitors)
	}()

	close(failedCompetitors)
	wg2.Wait()

	return results
}

func (worker *Worker) SendSlackMessageToFailedCompetitors(failedCompetitors chan Competitor) {
	var failedCompetitorsList []string
	slackMod, err := messager.NewSlackModule()

	if err != nil {
		return
	}

	for competitor := range failedCompetitors {
		failedCompetitorsList = append(failedCompetitorsList, competitor.Name)
	}

	if len(failedCompetitorsList) > 0 {
		failedCompetitorsString := strings.Join(failedCompetitorsList, ", ")
		err = slackMod.SendMessage("Sellers crawler- Failed to get data for following competitors : " + failedCompetitorsString)

		if err != nil {
			return
		}
	}
}

func (worker *Worker) prepareEmail(competitorsData []CompetitorData, err error, emailCred EmailCreds, competitorType string) error {

	const dateFormat = "2006-01-02"

	if len(competitorsData) > 0 {
		sort.Slice(competitorsData, func(i, j int) bool {
			return competitorsData[i].Position < competitorsData[j].Position
		})

		now := time.Now()
		today := now.Format(dateFormat)
		yesterday := now.AddDate(0, 0, -1).Format(dateFormat)

		subject := fmt.Sprintf("Competitors sellers.json daily changes  for %s - %s", competitorType, today)
		message := fmt.Sprintf("Below are the sellers.json changes for %s between - %s and %s", competitorType, yesterday, today)

		err = SendCustomHTMLEmail(emailCred.TO, emailCred.BCC, subject, message, competitorsData)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}
	return nil
}

func (worker *Worker) prepareAndInsertCompetitors(ctx context.Context, results chan map[string]interface{}, history []SellersJSONHistory, db *sqlx.DB, competitorsData []CompetitorData, positionMap map[string]string) ([]CompetitorData, error) {
	historyMap := make(map[string]SellersJSONHistory)
	var competitorsSlice []string
	var backupTodayMap map[string]interface{}

	for _, h := range history {
		historyMap[h.CompetitorName] = h
		if err := json.Unmarshal(*h.BackupToday, &backupTodayMap); err != nil {
			return nil, fmt.Errorf("failed to parse BackupToday for %s: %w", h.CompetitorName, err)
		}

		if len(backupTodayMap) == 2 {
			competitorsSlice = append(competitorsSlice, h.CompetitorName)
		}
	}

	for result := range results {
		for name, backupToday := range result {
			var historyRecord SellersJSONHistory

			if record, found := historyMap[name]; found {
				historyRecord = record
			}

			todayData, historyBackupToday, err := MapBackupTodayData(backupToday, historyRecord)
			if err != nil {
				return nil, fmt.Errorf("Error processing backup data for competitor %s: %w", name, err)
			}

			result := compareSellers(todayData, historyBackupToday)

			addedPublishers := result.ExtraPublishers
			addedDomains := result.ExtraDomains
			sellerType := result.SellerType
			deletedPublishers := result.DeletedPublishers
			deletedDomains := result.DeletedDomains

			addedPublisherDomains := make([]PublisherDomain, 0)
			deletedPublisherDomains := make([]PublisherDomain, 0)

			if addedPublishers != nil {
				for i, publisher := range addedPublishers {
					addedPublisherDomains = append(addedPublisherDomains, PublisherDomain{
						Publisher:  publisher,
						Domain:     addedDomains[i],
						SellerType: sellerType,
					})
				}
			}

			if deletedPublishers != nil {
				for i, publisher := range deletedPublishers {
					deletedPublisherDomains = append(deletedPublisherDomains, PublisherDomain{
						Publisher:  publisher,
						Domain:     deletedDomains[i],
						SellerType: sellerType,
					})
				}
			}

			competitorsData = append(competitorsData, CompetitorData{
				Name:                   name,
				URL:                    historyMap[name].URL,
				AddedPublisherDomain:   addedPublisherDomains,
				DeletedPublisherDomain: deletedPublisherDomains,
				Position:               positionMap[name],
			})

			backupBeforeYesterday := historyRecord.BackupYesterday
			if err := InsertCompetitor(ctx, db, name, addedDomains, addedPublishers, todayData, historyBackupToday, backupBeforeYesterday, deletedPublishers, deletedDomains); err != nil {
				return nil, fmt.Errorf("failed to insert competitor data for %s: %w", name, err)
			}
		}
	}

	var filteredCompetitorsData []CompetitorData
	for _, competitor := range competitorsData {
		if !isInSlice(competitor.Name, competitorsSlice) {
			filteredCompetitorsData = append(filteredCompetitorsData, competitor)
		}
	}

	return filteredCompetitorsData, nil
}

func isInSlice(competitor string, competitorsSlice []string) bool {
	for _, comp := range competitorsSlice {
		if comp == competitor {
			return true
		}
	}
	return false
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
