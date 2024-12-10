package real_time_report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/helpers"
	"strings"
)

type EmailReport struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"pubid" json:"pubid" toml:"pubid" yaml:"pubid"`
	Publisher            string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Device               string  `boil:"dtype" json:"dtype" toml:"dtype" yaml:"dtype"`
	Country              string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	PubFillRate          float64 `boil:"fill_rate" json:"fill_rate" toml:"fill_rate" yaml:"fill_rate"`
	BidRequests          float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	PublisherImpressions float64 `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	SoldImpressions      float64 `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	Cost                 float64 `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	Revenue              float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	CPM                  float64 `boil:"cpm" json:"cpm" toml:"cpm" yaml:"cpm"`
	RPM                  float64 `boil:"rpm" json:"rpm" toml:"rpm" yaml:"rpm"`
	DPRPM                float64 `boil:"dp_rpm" json:"dp_rpm" toml:"dp_rpm" yaml:"dp_rpm"`
	GP                   float64 `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	GPP                  float64 `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
}

func SendCustomHTMLEmail(to, bcc, subject string, htmlBody string, report []RealTimeReport, reportName string) error {
	toRecipients := strings.Split(to, ",")

	csvData, err := createCSVData(report)
	if err != nil {
		return fmt.Errorf("failed to create csv data, %w", err)
	}

	emailReq := modules.EmailRequest{
		To:       toRecipients,
		Bcc:      bcc,
		Subject:  subject,
		Body:     htmlBody,
		IsHTML:   false,
		Attach:   csvData,
		Filename: reportName,
	}

	return modules.SendEmail(emailReq)
}

func createCSVData(report []RealTimeReport) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(ColumnNames); err != nil {
		return nil, fmt.Errorf("failed to create csv data, %w", err)
	}
	formatter := &helpers.FormatValues{}

	for _, record := range report {
		formattedTime := helpers.FormatDate(record.Time)
		row := []string{
			formattedTime,
			record.PublisherID,
			record.Publisher,
			record.Domain,
			record.Device,
			record.Country,
			formatter.FillRate(record.PubFillRate),
			formatter.BidRequests(record.BidRequests),
			formatter.PubImps(int(record.PublisherImpressions)),
			formatter.SoldImps(int(record.SoldImpressions)),
			formatter.Cost(record.Cost),
			formatter.Revenue(record.Revenue),
			formatter.CPM(record.CPM),
			formatter.RPM(record.RPM),
			formatter.DPRPM(record.DpRPM),
			formatter.GP(record.GP),
			formatter.GPP(record.GPP),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return &buf, nil
}
