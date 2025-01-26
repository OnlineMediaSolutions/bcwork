package real_time_report

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

type EmailCreds struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}

type Worker struct {
	Cron           string                `json:"cron"`
	Quest          []string              `json:"quest"`
	Start          time.Time             `json:"start"`
	End            time.Time             `json:"end"`
	Slack          *messager.SlackModule `json:"slack_instances"`
	DatabaseEnv    string                `json:"dbenv"`
	EmailCreds     map[string]string     `json:"email_creads"`
	Fees           map[string]float64    `json:"fees"`
	ConsultantFees map[string]float64    `json:"consultant_fees"`
	LogSeverity    int                   `json:"logsev"`
	HttpClient     httpclient.Doer
	Publishers     map[string]string
	skipInitRun    bool
	ReportName     string
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	var questExist bool

	worker.skipInitRun, _ = conf.GetBoolValue("skip_init_run")
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	worker.HttpClient = httpclient.New(true)
	worker.ReportName = "real_time_report"
	worker.LogSeverity, _ = conf.GetIntValueWithDefault(config.LogSeverityKey, int(zerolog.InfoLevel))
	zerolog.SetGlobalLevel(zerolog.Level(worker.LogSeverity))

	emailCredsMap, err := config.FetchConfigValues([]string{worker.ReportName})
	worker.EmailCreds = emailCredsMap

	if err != nil {
		return fmt.Errorf("failed to get email credentials %w", err)
	}

	if err = bcdb.InitDB(worker.DatabaseEnv); err != nil {
		return fmt.Errorf("failed initialize DB for real time report in environment %s,%w", worker.DatabaseEnv, err)
	}

	worker.Cron, _ = conf.GetStringValue("cron")
	worker.Quest, questExist = conf.GetStringSlice("quest", ",")
	if !questExist {
		worker.Quest = []string{"amsquest2", "nycquest2"}
	}

	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting real time reports worker task")

	if worker.skipInitRun {
		log.Info().Msg("Skipping work as per the skip_init_run flag real time report.")
		worker.skipInitRun = false
		return nil
	}

	worker.End = time.Now().UTC()
	worker.Start = worker.End.Add(-1 * 24 * time.Hour)
	oneDayReport, err := worker.FetchAndMergeQuestReports(ctx)
	if err != nil {
		return err
	}

	chunks := makeChunks(oneDayReport)
	err = saveReportDBByChunks(ctx, chunks)
	if err != nil {
		return err
	}

	worker.Start = worker.End.Add(-7 * 24 * time.Hour)
	sevenDayReport, err := worker.FetchRealTimeData(ctx)
	if err != nil {
		return err
	}

	err = worker.RemoveOldDataFromDB(ctx)
	if err != nil {
		return err
	}

	var emailCreds EmailCreds
	credsRaw := worker.EmailCreds[worker.ReportName]
	if err := json.Unmarshal([]byte(credsRaw), &emailCreds); err != nil {
		return err
	}

	worker.PrepareEmail(sevenDayReport, err, emailCreds)
	return nil
}

func makeChunks(report map[string]*RealTimeReport) [][]*RealTimeReport {
	var chunks [][]*RealTimeReport
	chunkSize := viper.GetInt(config.APIChunkSizeKey)

	realTime := make([]*RealTimeReport, 0, len(report))
	for _, data := range report {
		realTime = append(realTime, data)
	}

	for i := 0; i < len(realTime); i += chunkSize {
		end := i + chunkSize
		if end > len(realTime) {
			end = len(realTime)
		}
		chunks = append(chunks, realTime[i:end])
	}
	return chunks
}

func saveReportDBByChunks(ctx context.Context, chunks [][]*RealTimeReport) error {

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, chunk := range chunks {

		realTimeReport := buildChunkRealTimeReport(chunk)
		if err := bulk.BulkInsertRealTimeReport(ctx, tx, realTimeReport); err != nil {
			log.Error().Err(err).Msgf("failed to insert real time report chunk %d", i)
			return fmt.Errorf("failed to insert real time report chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction in dpos bulk update")
		return fmt.Errorf("failed to commit transaction in dpos bulk update: %w", err)
	}
	return nil
}

func buildChunkRealTimeReport(chunk []*RealTimeReport) []models.RealTimeReport {
	realTimeReports := make([]models.RealTimeReport, 0, len(chunk))
	for _, data := range chunk {
		realTimeReports = append(realTimeReports, models.RealTimeReport{
			Time:                 data.Time,
			Publisher:            data.Publisher,
			PublisherID:          data.PublisherID,
			Domain:               data.Domain,
			BidRequests:          data.BidRequests,
			Device:               data.Device,
			Country:              data.Country,
			Revenue:              data.Revenue,
			Cost:                 data.Cost,
			SoldImpressions:      data.SoldImpressions,
			PublisherImpressions: data.PublisherImpressions,
			PubFillRate:          data.PubFillRate,
			CPM:                  data.CPM,
			RPM:                  data.RPM,
			DPRPM:                data.DpRPM,
			GP:                   data.GP,
			GPP:                  data.GPP,
			ConsultantFee:        data.ConsultantFee,
			TamFee:               data.TamFee,
			TechFee:              data.TechFee,
			DemandPartnerFee:     data.DemandPartnerFee,
			DataFee:              data.DataFee,
		})
	}
	return realTimeReports
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
