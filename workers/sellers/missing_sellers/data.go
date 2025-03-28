package missing_sellers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/m6yf/bcwork/workers/email_reports"
	"github.com/m6yf/bcwork/workers/sellers"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"time"
)

type DemandData struct {
	DemandPartnerName string `json:"DemandPartner"`
	PublisherName     string `json:"Publisher"`
	PublisherId       int64  `json:"PublisherId"`
	DPRequest         int64  `json:"BCMBidRequests"`
	SeatOwner         string `json:"SeatOwner,omitempty"`
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

type MissingPublisherInfo struct {
	PublisherId   string `json:"publisher_id"`
	PublisherName string `json:"publisher_name"`
	Status        string `json:"status"`
	SeatOwner     string `json:"seatOwner"`
	SeatURL       string `json:"seat_url"`
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

func fetchCompassData() ([]DemandData, error) {
	currentTime := time.Now().In(email_reports.Location)
	yesterday := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).AddDate(0, 0, -1)
	today := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).AddDate(0, 0, 0)

	requestData := prepareRequestData(yesterday, today)
	report, err := email_reports.GetCompassReport(constant.NewBidderReportingURL, requestData, true)

	if err != nil {
		return nil, err
	}

	var reportData Report
	err = json.Unmarshal(report, &reportData)
	if err != nil {
		return nil, err
	}

	return reportData.Data.Result, nil
}

func (worker *Worker) fetchDemandData(ctx context.Context) (map[string]string, error) {
	filters := core.DemandPartnerGetFilter{}
	options := core.DemandPartnerGetOptions{
		Filter:     filters,
		Pagination: nil,
		Order:      nil,
		Selector:   "",
	}

	dpoDemand, err := worker.demandPartnerService.GetDemandPartners(ctx, &options)
	if err != nil {
		return nil, fmt.Errorf("cannot get demand partners from database: %w", err)
	}

	demandSeat := make(map[string]string)

	for _, item := range dpoDemand {
		demandSeat[strings.ToLower(item.DemandPartnerID)] = item.SeatOwnerName
	}

	return demandSeat, nil
}

func getDemandMap(ctx context.Context, db *sqlx.DB) (map[string]string, error) {
	demandData, err := models.MissingSellers().All(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("error in getting sellers data for db, %w", err)
	}

	data := make(map[string]string)
	for _, partner := range demandData {
		data[partner.Name] = partner.URL
	}

	return data, nil
}

func getSellersJsonFiles(ctx context.Context, db *sqlx.DB) (map[string][]string, map[string][]string, error) {
	demandData, err := models.MissingSellers().All(ctx, db)
	if err != nil {
		return nil, nil, err
	}
	yesterdaySellersData := make(map[string][]string)
	todaySellersData := make(map[string][]string)

	for _, partner := range demandData {
		sellersData, err := sellers.FetchDataFromWebsite(partner.URL)
		if err != nil {
			return nil, nil, err
		}

		mapTodaySellersData(sellersData, partner, todaySellersData)

		if partner.Sellers.Valid {
			yesterdaySellersData[partner.Name] = strings.Split(partner.Sellers.String, ",")
		}
	}

	return todaySellersData, yesterdaySellersData, nil
}

func mapTodaySellersData(sellersData map[string]interface{}, partner *models.MissingSeller, todaySellersData map[string][]string) map[string][]string {
	if rawSellers, ok := sellersData["sellers"].([]interface{}); ok {
		sellerIds := make([]string, 0, len(rawSellers))
		dpMap := sellers.GetDPMap()
		mappingValue, exists := dpMap[partner.Name]

		for _, rawSeller := range rawSellers {
			sellerMap := rawSeller.(map[string]interface{})

			sellerId, ok := sellerMap["seller_id"].(string)
			if !ok {
				continue
			}

			if exists {
				sellerId = sellers.TrimSellerIdByDemand(mappingValue, sellerId)
			}

			sellerIds = append(sellerIds, sellerId)
		}

		todaySellersData[partner.Name] = sellerIds
	}

	return todaySellersData
}

func createCompassDataSet(data []DemandData, demandData map[string]string) map[string][]DemandData {
	dataSet := make(map[string][]DemandData)

	for _, entry := range data {
		seatOwner, exists := demandData[strings.ToLower(entry.DemandPartnerName)]
		if !exists || seatOwner == "" {
			continue
		}

		dataSet[seatOwner] = append(dataSet[seatOwner], DemandData{
			DemandPartnerName: entry.DemandPartnerName,
			PublisherId:       entry.PublisherId,
			PublisherName:     entry.PublisherName,
			DPRequest:         entry.DPRequest,
			SeatOwner:         seatOwner,
		})
	}

	return dataSet
}

func findMissingIds(compassData map[string][]DemandData, todayData, yesterdayData map[string][]string) map[string]MissingPublisherInfo {
	statusMap := make(map[string]MissingPublisherInfo)

	for dpName, dataList := range compassData {
		todaySet := make(map[string]struct{})
		yesterdaySet := make(map[string]struct{})

		if todayIds, found := todayData[dpName]; found {
			for _, id := range todayIds {
				todaySet[id] = struct{}{}
			}
		}

		if yesterdayIds, found := yesterdayData[dpName]; found {
			for _, id := range yesterdayIds {
				yesterdaySet[id] = struct{}{}
			}
		}

		status, err := prepareStatuses(dataList, todaySet, yesterdaySet, statusMap)
		if err != nil {
			return nil
		}

		statusMap = status
	}

	return statusMap
}

func prepareStatuses(dataList []DemandData, todaySet map[string]struct{}, yesterdaySet map[string]struct{}, statusMap map[string]MissingPublisherInfo) (map[string]MissingPublisherInfo, error) {
	for _, data := range dataList {
		publisherId := strconv.FormatInt(data.PublisherId, 10)

		_, todayExists := todaySet[publisherId]
		_, yesterdayExists := yesterdaySet[publisherId]

		var status string
		switch {
		case !todayExists && !yesterdayExists:
			status = "missing"
		case yesterdayExists && !todayExists:
			status = "deleted"
		default:
			continue
		}

		demandMap, err := getDemandMap(context.Background(), bcdb.DB())
		if err != nil {
			return nil, err
		}

		statusMap[publisherId] = MissingPublisherInfo{
			PublisherId:   publisherId,
			PublisherName: data.PublisherName,
			Status:        status,
			SeatOwner:     demandMap[data.SeatOwner],
		}
	}

	return statusMap, nil
}

func insertToDB(ctx context.Context, todaySellersData map[string][]string) error {
	for partnerName, todaySellerData := range todaySellersData {
		existingSeller, err := models.MissingSellers(
			models.MissingSellerWhere.Name.EQ(partnerName),
		).One(ctx, bcdb.DB())

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error fetching existing seller data: %w", err)
		}

		var yesterdayBackup string
		if existingSeller != nil {
			yesterdayBackup = existingSeller.Sellers.String
		}

		sellersData := &models.MissingSeller{
			Name:            partnerName,
			Sellers:         null.StringFrom(strings.Join(todaySellerData, ",")),
			Yesterdaybackup: yesterdayBackup,
		}

		err = sellersData.Upsert(ctx, bcdb.DB(), true,
			[]string{"name"},
			boil.Whitelist("sellers", "yesterdaybackup", "updated_at"),
			boil.Infer(),
		)
		if err != nil {
			return fmt.Errorf("error upserting data: %w", err)
		}
	}

	return nil
}

func prepareEmailAndSend(statusMap map[string]MissingPublisherInfo, emailCred EmailCreds) error {
	if len(statusMap) > 0 {
		now := time.Now()
		today := now.Format(time.DateOnly)
		subject := fmt.Sprintf("Missing Publishers in seller.json - %s", today)
		htmlReport, err := GenerateHTMLFromMissingPublishers(statusMap)

		err = sendCustomHTMLEmail(emailCred.TO, emailCred.BCC, subject, htmlReport)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendCustomHTMLEmail(to, bcc, subject, htmlBody string) error {
	toRecipients := strings.Split(to, ",")
	bccString := strings.Split(bcc, ",")

	emailReq := modules.EmailRequest{
		To:      toRecipients,
		Bcc:     bccString,
		Subject: subject,
		Body:    htmlBody,
		IsHTML:  true,
	}

	return modules.SendEmail(emailReq)
}
