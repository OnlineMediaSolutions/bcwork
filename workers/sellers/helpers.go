package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	SellerId          string
	DeletedPublishers []string
	DeletedDomains    []string
}

type AdsTxt struct {
	Domain        string
	SellerId      string
	PublisherName string
	AdsTxtStatus  string `json:"ads_txt_status"`
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
		log.Error().Err(err).Msg("Error making request")
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if sellers, ok := data["sellers"]; ok {
		if _, err := CheckSellersArray(sellers); err != nil {
			log.Error().Err(err).Msg("invalid sellers format")
			return nil, fmt.Errorf("invalid sellers format: %w", err)
		}
	} else {
		log.Error().Err(err).Msg("sellers array not found in the response")
		return nil, err
	}

	return data, nil
}

func InsertCompetitor(ctx context.Context, db boil.ContextExecutor, name string, comparisonResult ComparisonResult, backupToday, backupYesterday, backupBeforeYesterday interface{}) error {

	addedDomainsStr := strings.Join(comparisonResult.ExtraDomains, ",")
	addedPublishersStr := strings.Join(comparisonResult.ExtraPublishers, ",")
	deletedPublishersStr := strings.Join(comparisonResult.DeletedPublishers, ",")
	deletedDomainsStr := strings.Join(comparisonResult.DeletedDomains, ",")

	backupTodayJSON, err := json.Marshal(backupToday)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal backupToday")
		return err
	}
	backupYesterdayJSON, err := json.Marshal(backupYesterday)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal backupYesterday")
		return err
	}
	backupBeforeYesterdayJSON, err := json.Marshal(backupBeforeYesterday)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal backupBeforeYesterday")
		return err
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
		log.Error().Err(err).Msg("failed to insert or update competitor")
		return err
	}

	return nil
}

func (worker *Worker) Request(jobs <-chan Competitor, results chan<- map[string]interface{}, failedCompetitors chan<- Competitor, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		data, err := FetchDataFromWebsite(job.URL)
		if err != nil {
			log.Error().Err(err).Msg("Error fetching data for competitor")
			failedCompetitors <- job
			continue
		}

		results <- map[string]interface{}{job.Name: data}
	}
}

func (worker *Worker) GetHistoryData(ctx context.Context, db *sqlx.DB) ([]SellersJSONHistory, error) {
	histories, err := models.SellersJSONHistories().All(ctx, db)
	if err != nil {
		log.Error().Err(err).Msg("failed to query sellers_json_history")
		return nil, err
	}

	var results []SellersJSONHistory
	for _, history := range histories {
		competitor, err := history.CompetitorNameCompetitor().One(ctx, db)
		if err != nil {
			log.Error().Err(err).Msg("failed to query competitor")
			return nil, err
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
	var sellerId string
	var deletedPublishers []string
	var deletedDomains []string

	for _, seller := range todayData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name)

		if _, exists := sellerMapHistory[key]; !exists {
			//TODO
			extraPublishers = append(extraPublishers, seller.Name)
			extraDomains = append(extraDomains, seller.Domain)
			sellerType = seller.SellerType
			sellerId = seller.SellerID
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
		SellerId:          sellerId,
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
			log.Error().Err(err).Msg("failed to send email")
			return err
		}
	}
	return nil
}

func (worker *Worker) prepareAndInsertCompetitors(ctx context.Context, results chan map[string]interface{}, history []SellersJSONHistory, db *sqlx.DB, competitorsData []CompetitorData, positionMap map[string]string) ([]CompetitorData, error) {
	historyMap := make(map[string]SellersJSONHistory)
	var competitorsSlice []string
	var backupTodayMap map[string]interface{}
	var competitorsResult []CompetitorData

	for _, h := range history {
		historyMap[h.CompetitorName] = h
		if err := json.Unmarshal(*h.BackupToday, &backupTodayMap); err != nil {
			log.Error().Err(err).Msg("failed to parse BackupToday")
			return nil, err
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
				log.Error().Err(err).Msg("Error processing backup data")
				return nil, err
			}

			comparisonResult := compareSellers(todayData, historyBackupToday)

			addedPublishers := comparisonResult.ExtraPublishers
			addedDomains := comparisonResult.ExtraDomains
			sellerType := comparisonResult.SellerType
			sellerId := comparisonResult.SellerId

			deletedPublishers := comparisonResult.DeletedPublishers
			deletedDomains := comparisonResult.DeletedDomains

			addedPublisherDomains := worker.prepareAddedData(addedPublishers, addedDomains, sellerType, sellerId)
			deletedPublisherDomains := worker.prepareDeletedData(deletedPublishers, deletedDomains, sellerType)

			adsTxtData := worker.prepareAdsTxtData(addedPublisherDomains)

			fmt.Println("adsTxtData", adsTxtData)

			competitorsResult = worker.prepareCompetitorsData(comparisonResult, competitorsData, name, historyMap, addedPublisherDomains, deletedPublisherDomains, positionMap)
			backupBeforeYesterday := historyRecord.BackupYesterday
			if err := InsertCompetitor(ctx, db, name, comparisonResult, todayData, historyBackupToday, backupBeforeYesterday); err != nil {
				log.Error().Err(err).Msg("failed to insert competitor data")
				return nil, err
			}
		}
	}

	return competitorsResult, nil
}

func (worker *Worker) prepareCompetitorsData(comparisonResult ComparisonResult, competitorData []CompetitorData, name string, historyMap map[string]SellersJSONHistory, addedPublisherDomains []PublisherDomain, deletedPublisherDomains []PublisherDomain, positionMap map[string]string) []CompetitorData {
	deletedPublishers := comparisonResult.DeletedPublishers
	addedPublishers := comparisonResult.ExtraPublishers

	if len(addedPublishers) > 0 || len(deletedPublishers) > 0 {
		data := CompetitorData{
			Name:                   name,
			URL:                    historyMap[name].URL,
			AddedPublisherDomain:   addedPublisherDomains,
			DeletedPublisherDomain: deletedPublisherDomains,
			Position:               positionMap[name],
		}

		competitorData = append(competitorData, data)

	}
	return competitorData
}

func (worker *Worker) prepareDeletedData(deletedPublishers []string, deletedDomains []string, sellerType string) []PublisherDomain {
	deletedPublisherDomains := make([]PublisherDomain, 0)
	if deletedPublishers != nil {
		for i, publisher := range deletedPublishers {
			deletedPublisherDomains = append(deletedPublisherDomains, PublisherDomain{
				Publisher:  publisher,
				Domain:     deletedDomains[i],
				SellerType: sellerType,
			})
		}
	}
	return deletedPublisherDomains
}

func (worker *Worker) prepareAddedData(addedPublishers []string, addedDomains []string, sellerType string, sellerId string) []PublisherDomain {
	addedPublisherDomains := make([]PublisherDomain, 0)
	if addedPublishers != nil {
		for i, publisher := range addedPublishers {
			addedPublisherDomains = append(addedPublisherDomains, PublisherDomain{
				Publisher:  publisher,
				Domain:     addedDomains[i],
				SellerType: sellerType,
				SellerId:   sellerId,
			})
		}
	}
	return addedPublisherDomains
}

func MapBackupTodayData(backupToday interface{}, historyRecord SellersJSONHistory) (SellersJSON, SellersJSON, error) {
	backupTodayMap, ok := backupToday.(map[string]interface{})
	if !ok {
		log.Error().Msg("invalid backupToday format")
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("invalid backupToday format")
	}

	jsonData, err := json.Marshal(backupTodayMap)
	if err != nil {
		log.Error().Msg("failed to marshal map to JSON")
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	var backupTodayData SellersJSON
	if err := json.Unmarshal(jsonData, &backupTodayData); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal map data to SellersJSON")
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to unmarshal map data to SellersJSON: %w", err)
	}

	var historyBackupToday SellersJSON
	if historyRecord.BackupToday != nil {
		if err := json.Unmarshal(*historyRecord.BackupToday, &historyBackupToday); err != nil {
			log.Error().Err(err).Msg("failed to unmarshal BackupToday from history")
			return SellersJSON{}, SellersJSON{}, fmt.Errorf("failed to unmarshal BackupToday from history: %w", err)
		}
	}

	return backupTodayData, historyBackupToday, nil
}

func CheckSellersArray(sellers interface{}) ([]interface{}, error) {
	if sellersArray, ok := sellers.([]interface{}); ok {
		return sellersArray, nil
	}
	log.Error().Msg("sellers should be an array")
	return nil, fmt.Errorf("sellers should be an array, but got %T", sellers)
}
