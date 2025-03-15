package email_reports

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/modules/compass"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/utils/helpers"
	"golang.org/x/net/context"
	"time"
)

var Location, _ = time.LoadLocation(constant.AmericaNewYorkTimeZone)

type RequestData struct {
	Data RequestDetails `json:"data"`
}
type RequestDetails struct {
	Date       *Date    `json:"date,omitempty"`
	Dimensions []string `json:"dimensions,omitempty"`
	Metrics    []string `json:"metrics,omitempty"`
	Group      string   `json:"group,omitempty"`
}

type Date struct {
	Range    []string `json:"range,omitempty"`
	Interval string   `json:"interval,omitempty"`
}

type ReportDetails struct {
	Date                 string  `json:"date"`
	DataStamp            int64   `json:"DateStamp"`
	Publisher            string  `json:"Publisher"`
	Domain               string  `json:"Domain"`
	PaymentType          string  `json:"PaymentType"`
	AM                   string  `json:"AM"`
	PubImps              int64   `json:"PubImps"`
	LoopingRatio         float64 `json:"nbLR"`
	Ratio                float64 `json:"nbRatio"`
	CPM                  float64 `json:"nbCpm"`
	Cost                 float64 `json:"Cost"`
	RPM                  float64 `json:"nbRpm"`
	DpRPM                float64 `json:"nbDpRpm"`
	Revenue              float64 `json:"Revenue"`
	GP                   float64 `json:"nbGp"`
	GPP                  float64 `json:"nbGpp"`
	PublisherBidRequests int64   `json:"PublisherBidRequests"`
}

type Report struct {
	Data struct {
		Result []ReportDetails `json:"result"`
	} `json:"data"`
}

type AggregatedReport struct {
	Date                 string  `json:"date"`
	DataStamp            int64   `json:"DateStamp"`
	Publisher            string  `json:"publisher"`
	Domain               string  `json:"domain"`
	PaymentType          string  `json:"PaymentType"`
	AM                   string  `json:"am"`
	PubImps              string  `json:"PubImps"`
	LoopingRatio         float64 `json:"looping_ratio"`
	Ratio                float64 `json:"ratio"`
	CPM                  float64 `json:"cpm"`
	Cost                 float64 `json:"cost"`
	RPM                  float64 `json:"rpm"`
	DpRPM                float64 `json:"dpRpm"`
	Revenue              float64 `json:"Revenue"`
	GP                   float64 `json:"Gp"`
	GPP                  float64 `json:"Gpp"`
	PublisherBidRequests string  `json:"PublisherBidRequests"`
}

var userService = core.UserService{}

func GetUsers(responsiblePerson string) (map[string]string, error) {
	filters := core.UserFilter{
		Types: filter.String2DArrayFilter(filter.StringArrayFilter{responsiblePerson}),
	}

	options := core.UserOptions{
		Filter:     filters,
		Pagination: nil,
		Order:      nil,
		Selector:   "",
	}

	users, err := userService.GetUsers(context.Background(), &options)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]string)

	for _, user := range users {
		key := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		userMap[key] = user.Email
	}

	return userMap, nil
}

func Aggregate(reports []AggregatedReport) map[string][]AggregatedReport {
	aggregated := make(map[string][]AggregatedReport)

	for _, r := range reports {
		key := fmt.Sprintf("%s|%s|%s|%s", r.AM, r.Domain, r.Publisher, r.PaymentType)
		aggregated[key] = append(aggregated[key], AggregatedReport{
			AM:                   r.AM,
			Domain:               r.Domain,
			Publisher:            r.Publisher,
			PaymentType:          r.PaymentType,
			Date:                 r.Date,
			PubImps:              r.PubImps,
			DataStamp:            r.DataStamp,
			RPM:                  r.RPM,
			Ratio:                r.Ratio,
			LoopingRatio:         r.LoopingRatio,
			CPM:                  r.CPM,
			Cost:                 r.Cost,
			DpRPM:                r.DpRPM,
			Revenue:              r.Revenue,
			GP:                   r.GP,
			GPP:                  r.GPP,
			PublisherBidRequests: r.PublisherBidRequests,
		})
	}

	return aggregated
}

// TODO -change this to GetCOmpassReport
func GetReport(pubImpsThreshold int64) ([]AggregatedReport, error) {
	compassClient := compass.NewCompass()

	requestData := getRequestData()

	data, err := json.Marshal(requestData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling request data, %w", err)
	}

	reportData, err := compassClient.Request("/report-dashboard/report-new-bidder", "POST", data, true)

	if err != nil {
		return nil, fmt.Errorf("error getting report data,%w", err)
	}

	var report Report
	err = json.Unmarshal(reportData, &report)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling report data: %w", err)
	}

	aggregatedReports := prepareReport(report, pubImpsThreshold)

	return aggregatedReports, nil
}

func prepareReport(report Report, pubImpsThreshold int64) []AggregatedReport {
	aggregatedReports := make([]AggregatedReport, len(report.Data.Result))
	formatter := &helpers.FormatValues{}
	for i, r := range report.Data.Result {
		if r.PubImps >= pubImpsThreshold {
			aggregatedReports[i] = AggregatedReport{
				Date:                 r.Date,
				DataStamp:            r.DataStamp,
				Publisher:            r.Publisher,
				Domain:               r.Domain,
				PaymentType:          r.PaymentType,
				AM:                   r.AM,
				PubImps:              formatter.PubImps(int(r.PubImps)),
				PublisherBidRequests: formatter.PubBidRequests(int(r.PublisherBidRequests)),
				LoopingRatio:         helpers.RoundFloat(r.LoopingRatio),
				Ratio:                helpers.RoundFloat(r.Ratio),
				CPM:                  helpers.RoundFloat(r.CPM),
				Cost:                 helpers.RoundFloat(r.Cost),
				RPM:                  helpers.RoundFloat(r.RPM),
				DpRPM:                helpers.RoundFloat(r.DpRPM),
				Revenue:              helpers.RoundFloat(r.Revenue),
				GP:                   helpers.RoundFloat(r.GP),
				GPP:                  helpers.RoundFloat(r.GPP),
			}
		}
	}

	return aggregatedReports
}

func getRequestData() RequestData {
	currentTime := time.Now().In(Location)
	startOfLast7Days := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -7)

	endOfLast7Days := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).AddDate(0, 0, 0)

	startOfLast7DaysStr := startOfLast7Days.Format(time.DateTime)
	endOfLast7DaysStr := endOfLast7Days.Format(time.DateTime)

	requestData := RequestData{
		Data: RequestDetails{
			Date: &Date{
				Range: []string{
					startOfLast7DaysStr,
					endOfLast7DaysStr,
				},
				Interval: "day",
			},
			Dimensions: []string{
				"Publisher",
				"Domain",
				"AM",
				"PaymentType",
			},
			Metrics: []string{
				"PublisherBidRequests",
				"nbLR",
				"nbRatio",
				"nbRpm",
				"nbDpRpm",
				"Revenue",
				"nbGp",
				"nbGpp",
				"PubImps",
				"nbCpm",
				"Cost",
			},
		},
	}

	return requestData
}

func GetTimestampsForDateRange() (int64, int64, int64) {
	now := time.Now().In(Location)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := startOfDay.AddDate(0, 0, -1).Unix() / 100
	startOfLastWeek := startOfDay.AddDate(0, 0, -7).Unix() / 100
	today := startOfDay.Unix() / 100

	return yesterday, startOfLastWeek, today
}

func FilterReportsByDate(aggregated map[string][]AggregatedReport) map[string][]AggregatedReport {
	yesterday, startOfLastWeek, today := GetTimestampsForDateRange()

	amDomainData := make(map[string][]AggregatedReport)

	for key, aggs := range aggregated {
		for _, agg := range aggs {
			if agg.DataStamp == yesterday || agg.DataStamp == startOfLastWeek || agg.DataStamp == today {
				amDomainData[key] = append(amDomainData[key], agg)
			}
		}
	}

	return amDomainData
}

func GetCompassReport(url string, requestData RequestData, isReporting bool) ([]byte, error) {
	compassClient := compass.NewCompass()

	data, err := json.Marshal(requestData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling request data, %w", err)
	}

	report, err := compassClient.Request(url, "POST", data, isReporting)

	if err != nil {
		return nil, fmt.Errorf("error getting report data,%w", err)
	}

	return report, nil
}
