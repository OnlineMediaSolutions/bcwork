package missing_publishers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/workers/email_reports"
	"github.com/m6yf/bcwork/workers/sellers"
	"golang.org/x/net/context"
	"strings"
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

type DemandData struct {
	DemandPartnerName string
	PublisherId       string
	DPRequest         int64 `json:"BCMBidRequests,omitempty"`
}

type Report struct {
	Data struct {
		Result []DemandData `json:"result"`
	} `json:"data"`
}

type Response struct {
	Data []Data `json:"data"`
}

type Data struct {
	Meta MetaData `json:"meta"`
}

type MetaData struct {
	DirectGroups map[string][]string `json:"directGroups"`
}

type Partner struct {
	Name string
	URL  string
}

type Seller struct {
	SellerId string
}

func prepareRequestData(start time.Time, end time.Time) email_reports.RequestData {

	startDt := start.Format(time.DateTime)
	endDt := end.Format(time.DateTime)

	requestData := email_reports.RequestData{
		Data: email_reports.RequestDetails{
			Date: &email_reports.Date{
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

func getCompassData() ([]DemandData, error) {
	currentTime := time.Now().In(email_reports.Location)
	yesterday := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -1)
	today := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).AddDate(0, 0, 0)

	requestData := prepareRequestData(yesterday, today)
	report, err := email_reports.GetCompassReport("/report-dashboard/report-new-bidder", requestData, true)

	if err != nil {
		return nil, fmt.Errorf("error in getting compass data, %w", err)
	}

	var reportData Report
	err = json.Unmarshal(report, &reportData)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshalling compass data: %w", err)
	}

	return reportData.Data.Result, nil
}

func getDemandData() (map[string]string, error) {

	data := email_reports.RequestData{Data: email_reports.RequestDetails{Group: "Ads.txt Lines"}}

	report, err := email_reports.GetCompassReport("/settings/query", data, false)

	if err != nil {
		return nil, fmt.Errorf("error in getting compass data, %w", err)
	}

	var response Response

	err = json.Unmarshal(report, &response)

	demandSeat := make(map[string]string)

	for _, item := range response.Data {
		directGroups := item.Meta.DirectGroups
		for key, value := range directGroups {
			demandSeat[key] = value[0]
		}
	}

	return demandSeat, nil
}

func dataMappingByDemandPartner(report []byte) map[string]DemandData {
	var reportData Report
	err := json.Unmarshal(report, &reportData)
	if err != nil {
		return nil
	}

	demandPartnerMap := make(map[string]DemandData)

	for _, partner := range reportData.Data.Result {
		demandPartnerMap[partner.DemandPartnerName] = partner
	}

	return demandPartnerMap

}

//func getSellersJsonFiles(ctx context.Context, db *sqlx.DB) (map[string]interface{}, map[string]interface{}, error) {
//	demandData, err := models.MissingPublishersSellers().All(ctx, db)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	demandSellersData := make(map[string]interface{})
//	yesterdaySellersData := make(map[string]interface{})
//
//	//TODO - this part need to be done in async
//	for _, partner := range demandData {
//		sellersData, err := sellers.FetchDataFromWebsite(partner.URL)
//		if err != nil {
//			return nil, nil, err
//		}
//
//		demandSellersData[partner.Name] = sellersData["sellers"]
//		yesterdaySellersData[partner.Name] = partner
//
//	}
//	return demandSellersData, yesterdaySellersData, nil
//}

func getSellersJsonFiles(ctx context.Context, db *sqlx.DB) (map[string][]string, map[string][]string, error) {
	demandData, err := models.MissingPublishersSellers().All(ctx, db)
	if err != nil {
		return nil, nil, err
	}

	demandSellersData := make(map[string][]string)
	yesterdaySellersData := make(map[string][]string)

	// TODO - this part needs to be done in async
	for _, partner := range demandData {
		sellersData, err := sellers.FetchDataFromWebsite(partner.URL)
		if err != nil {
			return nil, nil, err
		}

		if sel, ok := sellersData["sellers"].([]Seller); ok {
			var sellerIds []string
			for _, seller := range sel {
				sellerIds = append(sellerIds, seller.SellerId)
			}
			demandSellersData[partner.Name] = sellerIds
		}

		if partner.Sellers.Valid {
			yesterdaySellersData[partner.Name] = strings.Split(partner.Sellers.String, ",")
		}

	}

	return demandSellersData, yesterdaySellersData, nil
}

func createCompassDataSet(data []DemandData) map[string][]string {
	dataSet := make(map[string][]string)
	for _, entry := range data {
		dataSet[entry.DemandPartnerName] = append(dataSet[entry.DemandPartnerName], entry.PublisherId)
	}
	return dataSet
}

func findMissingIds(compassDataSet, todayDataSet, yesterdayDataSet map[string][]string) {
	for partner, todayIds := range todayDataSet {
		compassIds, exists := compassDataSet[partner]
		if !exists {
			fmt.Printf("No data for partner: %s\n", partner)
			continue
		}

		// Create a map for quick lookup of compassIds
		compassIdSet := make(map[string]struct{})
		for _, id := range compassIds {
			compassIdSet[id] = struct{}{}
		}

		// Check for missing IDs
		for _, todayId := range todayIds {
			if _, found := compassIdSet[todayId]; !found {
				fmt.Printf("Missing ID for partner %s: %s\n", partner, todayId)
			}
		}
	}
}
