package email_reports

import (
	"fmt"
	"github.com/m6yf/bcwork/bcdb/filter"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils/constant"
	"golang.org/x/net/context"
	"time"
)

var Location, _ = time.LoadLocation(constant.AmericaNewYorkTimeZone)

type RequestData struct {
	Data RequestDetails `json:"data"`
}
type RequestDetails struct {
	Date       Date     `json:"date"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}
type Date struct {
	Range    []string `json:"range"`
	Interval string   `json:"interval"`
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
