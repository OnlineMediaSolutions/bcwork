package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/constant"
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
	ExtraPublishers    []string
	ExtraDomains       []string
	SellerType         []string
	SellerId           []string
	DeletedPublishers  []string
	DeletedDomains     []string
	DeletedSellerTypes []string
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
		return nil, fmt.Errorf("failed to get competitors from db: %w", err)
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
		return nil, fmt.Errorf("error making request for getting sellers")
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

func InsertCompetitor(ctx context.Context, db boil.ContextExecutor, name string, comparisonResult ComparisonResult, backupToday, backupYesterday, backupBeforeYesterday interface{}) error {

	addedDomainsStr := strings.Join(comparisonResult.ExtraDomains, ",")
	addedPublishersStr := strings.Join(comparisonResult.ExtraPublishers, ",")
	deletedPublishersStr := strings.Join(comparisonResult.DeletedPublishers, ",")
	deletedDomainsStr := strings.Join(comparisonResult.DeletedDomains, ",")

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
		return fmt.Errorf("failed to insert or update competitor: %w", err)
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

func normalizeKey(domain, name, sellerId string) string {
	domain = strings.ReplaceAll(domain, "http://", "")
	domain = strings.ReplaceAll(domain, "https://", "")
	return strings.TrimSpace(strings.ToLower(domain)) + ":" + strings.TrimSpace(strings.ToLower(name)) + ":" + strings.TrimSpace(strings.ToLower(sellerId))
}

func compareSellers(todayData, historyData SellersJSON) ComparisonResult {
	var extraPublishers []string
	var extraDomains []string
	var sellerTypes []string
	var sellerIds []string
	var deletedPublishers []string
	var deletedDomains []string
	var deletedSellerTypes []string

	sellerMapHistory := make(map[string]struct{})
	sellerMapToday := make(map[string]struct{})

	for _, seller := range historyData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name, seller.SellerID)
		sellerMapHistory[key] = struct{}{}
	}

	for _, seller := range todayData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name, seller.SellerID)
		sellerMapToday[key] = struct{}{}
	}

	extraPublishers, extraDomains, sellerTypes, sellerIds = getSellersTodayData(todayData, sellerMapHistory, extraPublishers, extraDomains, sellerTypes, sellerIds)
	deletedPublishers, deletedDomains, deletedSellerTypes = getSellersHistoryData(historyData, sellerMapToday, deletedPublishers, deletedDomains, deletedSellerTypes)

	return ComparisonResult{
		ExtraPublishers:    extraPublishers,
		ExtraDomains:       extraDomains,
		SellerType:         sellerTypes,
		SellerId:           sellerIds,
		DeletedPublishers:  deletedPublishers,
		DeletedDomains:     deletedDomains,
		DeletedSellerTypes: deletedSellerTypes,
	}
}

func getSellersHistoryData(historyData SellersJSON, sellerMapToday map[string]struct{}, deletedPublishers []string, deletedDomains []string, deletedSellersType []string) ([]string, []string, []string) {
	for _, seller := range historyData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name, seller.SellerID)

		if _, exists := sellerMapToday[key]; !exists {
			deletedPublishers = append(deletedPublishers, seller.Name)
			deletedDomains = append(deletedDomains, seller.Domain)
			deletedSellersType = append(deletedSellersType, seller.SellerType)
		}
	}
	return deletedPublishers, deletedDomains, deletedSellersType
}

func getSellersTodayData(todayData SellersJSON, sellerMapHistory map[string]struct{}, extraPublishers []string, extraDomains []string, sellerTypes []string, sellerIds []string) ([]string, []string, []string, []string) {
	for _, seller := range todayData.Sellers {
		key := normalizeKey(seller.Domain, seller.Name, seller.SellerID)

		if _, exists := sellerMapHistory[key]; !exists {
			extraPublishers = append(extraPublishers, seller.Name)
			extraDomains = append(extraDomains, seller.Domain)
			sellerTypes = append(sellerTypes, seller.SellerType)
			sellerIds = append(sellerIds, seller.SellerID)
		}
	}
	return extraPublishers, extraDomains, sellerTypes, sellerIds
}

func (worker *Worker) PrepareCompetitors(competitors []Competitor) chan map[string]interface{} {
	var wg sync.WaitGroup
	jobs := make(chan Competitor, len(competitors))
	results := make(chan map[string]interface{}, len(competitors))
	failedCompetitors := make(chan Competitor, len(competitors))

	for i := 1; i <= constant.SellersJsonWorkerCount; i++ {
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
	if len(competitorsData) > 0 {
		sort.Slice(competitorsData, func(i, j int) bool {
			return competitorsData[i].Position < competitorsData[j].Position
		})

		now := time.Now()
		today := now.Format(constant.PostgresTimestamp)
		yesterday := now.AddDate(0, 0, -1).Format(constant.PostgresTimestamp)

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

	data, err := worker.prepereCompetitorSlice(history, historyMap, backupTodayMap, competitorsSlice)
	if err != nil {
		return data, err
	}

	for result := range results {
		for name, backupToday := range result {
			historyRecord, found := historyMap[name]

			if !found {
				log.Info().Msgf("No history record found for %s, skipping", name)
				continue
			}

			todayData, historyBackupToday, err := MapBackupTodayData(backupToday, historyRecord)
			if err != nil {
				return nil, fmt.Errorf("error processing backup data for competitor %s: %w", name, err)

			}

			comparisonResult := compareSellers(todayData, historyBackupToday)

			addedPublisherDomains := worker.prepareAddedData(
				comparisonResult.ExtraPublishers,
				comparisonResult.ExtraDomains,
				comparisonResult.SellerType,
				comparisonResult.SellerId,
			)

			deletedPublisherDomains := worker.prepareDeletedData(
				comparisonResult.DeletedPublishers,
				comparisonResult.DeletedDomains,
				comparisonResult.DeletedSellerTypes,
			)

			competitorsResult = worker.prepareCompetitorsData(
				comparisonResult,
				competitorsData,
				name,
				historyMap,
				addedPublisherDomains,
				deletedPublisherDomains,
				positionMap,
			)

			if err := InsertCompetitor(ctx, db, name, comparisonResult, todayData, historyBackupToday, historyRecord.BackupYesterday); err != nil {
				return nil, err
			}
		}
	}

	return competitorsResult, nil
}

func (worker *Worker) prepereCompetitorSlice(history []SellersJSONHistory, historyMap map[string]SellersJSONHistory, backupTodayMap map[string]interface{}, competitorsSlice []string) ([]CompetitorData, error) {
	for _, h := range history {
		historyMap[h.CompetitorName] = h
		if err := json.Unmarshal(*h.BackupToday, &backupTodayMap); err != nil {
			return nil, fmt.Errorf("failed to parse BackupToday: %w", err)
		}

		if len(backupTodayMap) == 2 {
			competitorsSlice = append(competitorsSlice, h.CompetitorName)
		}
	}
	return nil, nil
}

func (worker *Worker) prepareCompetitorsData(comparisonResult ComparisonResult, competitorData []CompetitorData, name string, historyMap map[string]SellersJSONHistory, addedPublisherDomains []PublisherDomain, deletedPublisherDomains []PublisherDomain, positionMap map[string]string) []CompetitorData {
	deletedPublishers := comparisonResult.DeletedPublishers
	enhancedAddedDomains := worker.enhancePublisherDomains(addedPublisherDomains)

	if len(enhancedAddedDomains) > 0 || len(deletedPublishers) > 0 {
		data := CompetitorData{
			Name:                   name,
			URL:                    historyMap[name].URL,
			AddedPublisherDomain:   enhancedAddedDomains,
			DeletedPublisherDomain: deletedPublisherDomains,
			Position:               positionMap[name],
		}

		competitorData = append(competitorData, data)

	}
	return competitorData
}

func (worker *Worker) prepareDeletedData(deletedPublishers []string, deletedDomains []string, sellerTypes []string) []PublisherDomain {
	deletedPublisherDomains := make([]PublisherDomain, 0)
	if deletedPublishers != nil {
		for i, publisher := range deletedPublishers {
			deletedPublisherDomains = append(deletedPublisherDomains, PublisherDomain{
				Publisher:  publisher,
				Domain:     deletedDomains[i],
				SellerType: sellerTypes[i],
			})
		}
	}
	return deletedPublisherDomains
}

func (worker *Worker) prepareAddedData(addedPublishers []string, addedDomains []string, sellerTypes []string, sellerIds []string) []PublisherDomain {
	addedPublisherDomains := make([]PublisherDomain, 0)
	if addedPublishers != nil {
		for i, publisher := range addedPublishers {
			addedPublisherDomains = append(addedPublisherDomains, PublisherDomain{
				Publisher:  publisher,
				Domain:     addedDomains[i],
				SellerType: sellerTypes[i],
				SellerId:   sellerIds[i],
			})
		}
	}
	return addedPublisherDomains
}

func MapBackupTodayData(backupToday interface{}, historyRecord SellersJSONHistory) (SellersJSON, SellersJSON, error) {
	backupTodayMap, ok := backupToday.(map[string]interface{})
	if !ok {
		return SellersJSON{}, SellersJSON{}, fmt.Errorf("invalid backupToday format for today map")
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
