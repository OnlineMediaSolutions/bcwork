package missing_publishers_sellers

import (
	"encoding/json"
	"fmt"
	"github.com/m6yf/bcwork/workers/email_reports"
	"strconv"
	"time"
)

/**
1.get Data from compass
2.Mapping data By map[DemandPartner]struct {pubId, pubName}
3.get demand partner seat name in http://10.166.10.36:8080/demand-management - add  to [demand partner name]-> {pubId, pubName} seat name
4.get sellers json files and map by [demand partner name] -> {pubId, pubName} map[DemandPartner]struct
5.Comparing function
6.SaveSellersToDB
7.SendEmail
*/

type DemandPartner struct {
	DemandPartnerName string `json:"demand_partner_name"`
	PublisherName     string `json:"publisher_name"`
	PublisherId       string `json:"publisher_id"`
	DPRequest         string `json:"dp_request"`
}

type Report struct {
	Data struct {
		Result []DemandPartner `json:"result"`
	} `json:"data"`
}

func getRequestData(start time.Time, end time.Time) email_reports.RequestData {

	startDt := start.Format(time.DateTime)
	endDt := end.Format(time.DateTime)

	requestData := email_reports.RequestData{
		Data: email_reports.RequestDetails{
			Date: email_reports.Date{
				Range: []string{
					startDt,
					endDt,
				},
				Interval: "day",
			},
			Dimensions: []string{
				"PublisherId",
				"Publisher",
				"DemandPartner",
			},
			Metrics: []string{
				"BCMBidRequests",
			},
		},
	}

	return requestData
}

func getCompassData() ([]byte, error) {
	currentTime := time.Now().In(email_reports.Location)
	yesterday := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -1)
	today := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).AddDate(0, 0, 0)

	requestData := getRequestData(yesterday, today)
	report, err := email_reports.GetCompassReport("/report-dashboard/report-new-bidder", requestData)

	if err != nil {
		return nil, fmt.Errorf("error in getting compass data, %w", err)
	}

	return report, nil
}

func getDemandData() (map[string]string, error) {

	data := email_reports.RequestData{Data: email_reports.RequestDetails{Group: "Ads.txt Lines"}}

	report, err := email_reports.GetCompassReport("/settings/query", data)

	if err != nil {
		return nil, fmt.Errorf("error in getting compass data, %w", err)
	}

	demandData := make(map[string]string)

	for index, value := range report {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}

	return report, nil
}

func dataMappingByDemandPartner(report []byte) map[string]DemandPartner {
	var reportData Report
	err := json.Unmarshal(report, &reportData)
	if err != nil {
		return nil
	}

	demandPartnerMap := make(map[string]DemandPartner)

	for _, partner := range reportData.Data.Result {
		demandPartnerMap[partner.DemandPartnerName] = partner
	}

	return demandPartnerMap

}
