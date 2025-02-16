package nodpresponse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/core/bulk"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/modules/export"
	"github.com/m6yf/bcwork/quest"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type Worker struct {
	cron        string
	quest       []string
	databaseEnv string
	chunkSize   int
	emailList   struct {
		to  string
		bcc string
	}
	logSeverity  int
	skipInitRun  bool
	exportModule export.Exporter
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	const (
		questDefault        = "amsquest2,nycquest2"
		dbEnvDefault        = "local_prod"
		cronDefault         = "0 8 * * *" // every day at 08:00 am
		logSeverityDefault  = 1           // info
		chunkSizeDefault    = 200
		sqlDebugDefault     = false
		emailListToDefault  = "inbarl@onlinemediasolutions.com"
		emailListBCCDefault = "devsupport.onlinemediasolutions.com"
	)

	w.skipInitRun, _ = conf.GetBoolValue(config.SkipInitRunKey)
	boil.DebugMode = conf.GetBoolValueWithDefault(config.SQLDebugKey, sqlDebugDefault)
	w.databaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, dbEnvDefault)
	w.chunkSize, _ = conf.GetIntValueWithDefault(config.ChunkSizeKey, chunkSizeDefault)
	w.logSeverity, _ = conf.GetIntValueWithDefault(config.LogSeverityKey, logSeverityDefault)
	zerolog.SetGlobalLevel(zerolog.Level(w.logSeverity))

	w.emailList.to = conf.GetStringValueWithDefault(config.EmailToKey, emailListToDefault)
	w.emailList.bcc = conf.GetStringValueWithDefault(config.EmailBCCKey, emailListBCCDefault)

	if err := bcdb.InitDB(w.databaseEnv); err != nil {
		return fmt.Errorf("failed initialize db for no dp responses report in environment [%s]: %w", w.databaseEnv, err)
	}

	w.cron = conf.GetStringValueWithDefault(config.CronExpressionKey, cronDefault)
	questString := conf.GetStringValueWithDefault(config.QuestKey, questDefault)
	w.quest = strings.Split(questString, ",")

	w.exportModule = export.NewExportModule()

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("starting no dp responses report worker task")

	if w.skipInitRun {
		log.Info().Msg("skipping work as per the skip_init_run flag")
		w.skipInitRun = false

		return nil
	}

	now := time.Now().UTC()
	end := now.Format(time.DateOnly)
	start := now.AddDate(0, 0, -1).Format(time.DateOnly)

	questReport, err := w.fetchAndMergeQuestReports(ctx, start, end)
	if err != nil {
		return err
	}

	log.Info().Msg("saving quest report to postgres")
	err = w.saveQuestReport(ctx, questReport)
	if err != nil {
		return err
	}

	log.Info().Msg("getting data from postgres to build report")
	reportStart := now.AddDate(0, 0, -3).Format(time.DateOnly)
	mods, err := fetchPostgresReport(ctx, reportStart, end)
	if err != nil {
		return fmt.Errorf("failed to fetch data for report from postgres: %w", err)
	}

	log.Info().Msg("deleting old data")
	limit := now.AddDate(0, 0, -7).Format(time.DateOnly)
	_, err = models.NoDPResponseReports(models.NoDPResponseReportWhere.Time.LT(limit)).DeleteAll(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("failed to remove old data from no_dp_response table: %w", err)
	}

	data := make([]json.RawMessage, 0, len(mods))
	for _, mod := range mods {
		b, err := json.Marshal(mod)
		if err != nil {
			return err
		}
		data = append(data, b)
	}

	log.Info().Msg("creating xlsx file")
	filedata, err := w.exportModule.ExportXLSX(ctx, getDownloadXLSXRequest(data))
	if err != nil {
		return err
	}

	log.Info().Msg("sending email")
	err = modules.SendEmail(modules.EmailRequest{
		To:       strings.Split(w.emailList.to, ","),
		Bcc:      strings.Split(w.emailList.bcc, ","),
		Subject:  fmt.Sprintf("No DP Responses - Domain Level Report %s", end),
		Body:     fmt.Sprintf("No DP Responses - Domain Level Report between %s - %s", reportStart, end),
		IsHTML:   false,
		Attach:   bytes.NewBuffer(filedata),
		Filename: fmt.Sprintf("%v.%v.%v", "no_dp_response-domain_level", now.Format("2006_01_02_15_04_05"), dto.XLSX),
	})
	if err != nil {
		return err
	}

	log.Info().Msg("no dp responses report worker task finished successfully")

	return nil
}

func (w *Worker) GetSleep() int {
	next := bccron.Next(w.cron)

	log.Info().Msg(fmt.Sprintf("next run in: %v", time.Duration(next)*time.Second))
	if w.cron != "" {
		return next
	}

	return 0
}

func (w *Worker) fetchAndMergeQuestReports(ctx context.Context, start, end string) ([]*dto.NoDPResponseReport, error) {
	baseQuery := `
		SELECT
			DATE_TRUNC('day', to_timezone(timestamp, 'America/New_York')) as time,
			pubid,
			domain,
			dpid,
			sum(count) as %s
		FROM %s
		WHERE
			to_timezone(timestamp, 'America/New_York') >= '%s'
			AND to_timezone(timestamp, 'America/New_York') < '%s'
		GROUP BY 1,2,3,4;
	`
	requestsQuery := fmt.Sprintf(baseQuery, "bid_requests", "demand_request_placement", start, end)
	responsesQuery := fmt.Sprintf(baseQuery, "bid_responses", "demand_response_placement", start, end)

	reportMap := make(map[string]*dto.NoDPResponseReport)
	for _, instance := range w.quest {
		if err := quest.InitDB(instance); err != nil {
			return nil, fmt.Errorf("failed to initialize quest instance [%s]: %w", instance, err)
		}

		var requests []*dto.NoDPResponseReport
		log.Info().Msgf("instance [%v]: getting requests", instance)
		if err := queries.Raw(requestsQuery).Bind(ctx, quest.DB(), &requests); err != nil {
			return nil, fmt.Errorf("failed to query dp requests from quest instance [%s]: %w", instance, err)
		}
		fillReportMap(reportMap, requests)

		var responses []*dto.NoDPResponseReport
		log.Info().Msgf("instance [%v]: getting responses", instance)
		if err := queries.Raw(responsesQuery).Bind(ctx, quest.DB(), &responses); err != nil {
			return nil, fmt.Errorf("failed to query dp responses from quest instance [%s]: %w", instance, err)
		}
		fillReportMap(reportMap, responses)
	}

	log.Info().Msg("processing results")
	result := make([]*dto.NoDPResponseReport, 0, len(reportMap))
	for _, value := range reportMap {
		if value.BidResponses == 0 {
			result = append(result, value)
		}
	}

	return result, nil
}

func (w *Worker) saveQuestReport(ctx context.Context, report []*dto.NoDPResponseReport) error {
	chunks := makeChunks(report, w.chunkSize)

	tx, err := bcdb.DB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("no dp response: failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	log.Info().Int("len(chunks)", len(chunks)).Msg("total amount of chunks")
	for i, chunk := range chunks {
		if err := bulk.BulkInsertNoDPResponseReport(ctx, tx, chunk); err != nil {
			return fmt.Errorf("no dp response: failed to insert chunk %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("no dp response: failed to commit transaction: %w", err)
	}

	return nil
}

func fetchPostgresReport(ctx context.Context, start, end string) ([]*dto.NoDPResponseReport, error) {
	query := `
		WITH res AS (
			SELECT
				demand_partner_id AS dpid,
				publisher_id AS pubid,
				"domain",
				sum(bid_requests) AS bid_requests
			FROM no_dp_response_report AS x 
			WHERE
				"time" >= $1 AND "time" < $2
			GROUP BY demand_partner_id, publisher_id, "domain"
			HAVING count(demand_partner_id) = $3
		)
		SELECT
			r.dpid,
			p."name" AS publisher_name,
			r.pubid,
			r."domain",
			r.bid_requests
		FROM res AS r
		JOIN publisher AS p ON p.publisher_id = r.pubid
		ORDER BY r.bid_requests DESC;
	`

	var mods []*dto.NoDPResponseReport
	err := queries.Raw(query, start, end, 3).
		Bind(ctx, bcdb.DB(), &mods)
	if err != nil {
		return nil, err
	}

	return mods, nil
}

func fillReportMap(reportMap map[string]*dto.NoDPResponseReport, table []*dto.NoDPResponseReport) {
	for _, row := range table {
		key := row.BuildKey()
		data, ok := reportMap[key]
		if !ok {
			name, ok := constant.DemandPartnerMap[row.DPID]
			if !ok {
				name = row.DPID
			}

			reportMap[key] = &dto.NoDPResponseReport{
				Time:         row.Time,
				DPID:         name,
				PubID:        row.PubID,
				Domain:       row.Domain,
				BidRequests:  row.BidRequests,
				BidResponses: row.BidResponses,
			}
		} else {
			data.BidRequests += row.BidRequests
			data.BidResponses += row.BidResponses
		}
	}
}

func makeChunks(report []*dto.NoDPResponseReport, chunkSize int) [][]*dto.NoDPResponseReport {
	var chunks [][]*dto.NoDPResponseReport
	for i := 0; i < len(report); i += chunkSize {
		end := i + chunkSize
		if end > len(report) {
			end = len(report)
		}
		chunks = append(chunks, report[i:end])
	}

	return chunks
}

func getDownloadXLSXRequest(data []json.RawMessage) *dto.DownloadRequest {
	return &dto.DownloadRequest{
		Columns: []*dto.Column{
			{Name: "dpid", DisplayName: "Demand Partner"},
			{Name: "publisher_name", DisplayName: "Publisher"},
			{Name: "pubid", DisplayName: "Publisher ID"},
			{Name: "domain", DisplayName: "Domain"},
			{Name: "bid_requests", DisplayName: "DP Bid Requests", Style: export.IntColumnStyle},
		},
		Data: data,
	}
}
