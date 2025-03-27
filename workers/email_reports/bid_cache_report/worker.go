package bid_cache_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

type Worker struct {
	Cron        string                `json:"cron"`
	Quest       []string              `json:"quest"`
	Start       time.Time             `json:"start"`
	End         time.Time             `json:"end"`
	Slack       *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv string                `json:"dbenv"`
	EmailCreds  map[string]string     `json:"email_creads"`
	LogSeverity int                   `json:"logsev"`
	HttpClient  httpclient.Doer
	Publishers  map[string]string
	skipInitRun bool
	ReportName  string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var questExist bool

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	worker.HttpClient = httpclient.New(true)
	worker.ReportName = "bid_cache_report"
	worker.LogSeverity, _ = conf.GetIntValueWithDefault(config.LogSeverityKey, int(zerolog.InfoLevel))
	zerolog.SetGlobalLevel(zerolog.Level(worker.LogSeverity))

	emailCredsMap, err := config.FetchConfigValues([]string{worker.ReportName})
	worker.EmailCreds = emailCredsMap

	if err != nil {
		return fmt.Errorf("failed to get email credentials %w", err)
	}

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return fmt.Errorf("failed initialize DB %s,%w", worker.DatabaseEnv, err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"nycquest2", "amsquest2"}
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting Bid Caching worker")

	if worker.skipInitRun {
		log.Info().Msg("Skipping work as per the skip_init_run flag Bid cache report.")
		worker.skipInitRun = false

		return nil
	}

	activeBidCachings, err := fetchActiveBidCache(ctx)
	if err != nil {
		return err
	}

	pubDomains := buildPublisherDomainString(activeBidCachings)
	bidCacheResponse, pubDom, err := worker.fetchDataFromQuest(ctx, pubDomains)
	if err != nil {
		return err
	}

	responseMap := filterData(bidCacheResponse, pubDom)

	if len(responseMap) > 0 {
		var emailCreds EmailCreds
		credsRaw := worker.EmailCreds[worker.ReportName]
		if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
			return err
		}
		err = worker.SendEmail(responseMap, emailCreds)
		if err != nil {
			return err
		}
	}
	log.Info().Msg("Finished Bid Caching worker")

	return nil
}

func filterData(data map[string]*BidCacheData, pubDom map[string]bool) []*BidCacheData {
	var responseMap []*BidCacheData
	for key, _ := range pubDom {
		recordA := data[key+" - A"]
		recordB := data[key+" - B"]

		recordA.Time = strings.Split(recordA.Time, "T")[0]
		recordB.Time = strings.Split(recordB.Time, "T")[0]
		recordA.GP = (recordA.Revenue - (recordA.Cost + recordA.DataFee))
		recordB.GP = (recordB.Revenue - (recordB.Cost + recordB.DataFee))
		recordA.GPperPubImp = recordA.GP / float64(recordA.PublisherImpressions) * 1000
		recordB.GPperPubImp = recordB.GP / float64(recordB.PublisherImpressions) * 1000

		if (recordA.GPperPubImp > recordB.GPperPubImp*THRESHOLD) && (recordA.PublisherImpressions+recordB.PublisherImpressions > MINIMUM_IMPRESSIONS) {
			responseMap = append(responseMap, recordA)
			responseMap = append(responseMap, recordB)
		}
	}

	return responseMap
}

func buildPublisherDomainString(caching models.BidCachingSlice) string {
	var builder strings.Builder
	builder.WriteString("(")

	for i := 0; i < len(caching); i++ {
		if i > 0 {
			builder.WriteString(" OR ")
		}
		bidCache := caching[i]
		builder.WriteString(fmt.Sprintf("(domain='%s' and publisher='%s')", bidCache.Domain.String, bidCache.Publisher))
	}

	builder.WriteString(")")

	return builder.String()
}

func (worker *Worker) fetchDataFromQuest(ctx context.Context, pubdomains string) (map[string]*BidCacheData, map[string]bool, error) {
	responseMap := make(map[string]*BidCacheData)
	pubDom := make(map[string]bool)

	yesterday := time.Now().UTC().Add(-1 * 24 * time.Hour).Format("2006-01-02")
	query := fmt.Sprintf(questQuery, yesterday, pubdomains)

	for _, instance := range worker.Quest {
		if err := quest.InitDB(instance); err != nil {
			return nil, nil, fmt.Errorf("failed to initialize Quest instance: %s", instance)
		}
		var bidCacheData []*BidCacheData
		if err := queries.Raw(query).Bind(ctx, quest.DB(), &bidCacheData); err != nil {
			return nil, nil, fmt.Errorf("failed to query bid cache from Quest instance: %s", instance)
		}

		responseMap = generateResponseMap(responseMap, bidCacheData, pubDom)
	}

	return responseMap, pubDom, nil
}

func (record *BidCacheData) Key() string {
	return fmt.Sprintf("%s - %s - %s", record.PublisherID, record.Domain, record.Target)
}

func generateResponseMap(responseMap map[string]*BidCacheData, bidRequestRecords []*BidCacheData, pubDom map[string]bool) map[string]*BidCacheData {
	for _, record := range bidRequestRecords {
		key := record.Key()
		item, exists := responseMap[key]
		pubDom[record.PublisherID+" - "+record.Domain] = true
		if exists {
			mergedItem := &BidCacheData{
				Time:                 record.Time,
				PublisherID:          record.PublisherID,
				Domain:               record.Domain,
				Target:               record.Target,
				DataFee:              item.DataFee + record.DataFee,
				SoldImpressions:      item.SoldImpressions + record.SoldImpressions,
				PublisherImpressions: item.PublisherImpressions + record.PublisherImpressions,
				DemandPartnerFee:     item.DemandPartnerFee + record.DemandPartnerFee,
				Revenue:              item.Revenue + record.Revenue,
				Cost:                 item.Cost + record.Cost,
			}
			responseMap[key] = mergedItem
		} else {
			responseMap[key] = record
		}
	}

	return responseMap
}

func fetchActiveBidCache(ctx context.Context) (models.BidCachingSlice, error) {
	bidCaching, err := models.BidCachings(models.BidCachingWhere.Active.EQ(true)).All(ctx, bcdb.DB())
	if err != nil {
		return nil, fmt.Errorf("cannot fetch bid caching data: %w", err)
	}

	return bidCaching, nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}

	return 0
}
