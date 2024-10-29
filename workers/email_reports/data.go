package email_reports

import (
	"context"
	"fmt"
	"github.com/m6yf/bcwork/quest"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"sync"
	"time"
)

type RealTimeReport struct {
	PublisherID string  `boil:"pubid" json:"pubid" toml:"pubid" yaml:"pubid"`
	Domain      string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	BidRequests float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	Device      string  `boil:"dtype" json:"dtype" toml:"dtype" yaml:"dtype"`
	Country     string  `boil:"country" json:"country" toml:"country" yaml:"country"`
}

var QuestImpressions = `SELECT 
                     pubid, domain, dtype,country ,sum(count) bid_requests
                     FROM request_placement
                    WHERE timestamp >='%s' AND timestamp <'%s'
                    GROUP BY pubid,domain,dtype,country
                    ORDER BY pubid desc, domain desc`

func (worker *Worker) FetchFromQuest(ctx context.Context, start time.Time, end time.Time) (map[string]*RealTimeReport, error) {
	realTimeRecordsMap := make(map[string]*RealTimeReport)
	var wg sync.WaitGroup

	startString := start.Format("2006-01-02T15:04:05Z")
	endString := end.Format("2006-01-02T15:04:05Z")

	impressionsQuery := fmt.Sprintf(QuestImpressions, startString, endString)
	log.Info().Str("query", impressionsQuery).Msg("processImpressionsCounters")

	results := make(chan *RealTimeReport, len(worker.Quest))
	errChan := make(chan error, len(worker.Quest))

	for _, instance := range worker.Quest {
		wg.Add(1)

		go func(instance string) {
			defer wg.Done()

			err := quest.InitDB(instance)
			if err != nil {
				errChan <- errors.Wrapf(err, "Failed to initialize Quest instance: %s", instance)
				return
			}

			defer func() {
				err := quest.CloseDB()
				if err != nil {

				}
			}()

			var realTimeRecords []*RealTimeReport
			err = queries.Raw(impressionsQuery).Bind(ctx, quest.DB(), &realTimeRecords)
			if err != nil {
				errChan <- errors.Wrapf(err, "Failed to query Quest instance: %s", instance)
				return
			}

			for _, record := range realTimeRecords {
				key := record.Key()
				record.PublisherID = key
				results <- record
			}
		}(instance)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case record, ok := <-results:
			if ok {
				key := record.PublisherID

				if existingRecord, found := realTimeRecordsMap[key]; found {
					existingRecord.BidRequests += record.BidRequests
				} else {
					realTimeRecordsMap[key] = record
				}
			}
		case err, ok := <-errChan:
			if ok {
				return nil, err
			}
		case <-done:
			return realTimeRecordsMap, nil
		}
	}
}

func (record *RealTimeReport) Key() string {
	return fmt.Sprintf("%s - %s - %s - %s", record.PublisherID, record.Domain, record.Device, record.Country)
}
