package real_time_report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/m6yf/bcwork/modules"
	"github.com/m6yf/bcwork/utils/helpers"
	"reflect"
	"strings"
)

type EmailReport struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"pubid" json:"pubid" toml:"pubid" yaml:"pubid"`
	Publisher            string  `boil:"publisher" json:"publisher" toml:"publisher" yaml:"publisher"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	BidRequests          float64 `boil:"bid_requests" json:"bid_requests" toml:"bid_requests" yaml:"bid_requests"`
	Device               string  `boil:"dtype" json:"dtype" toml:"dtype" yaml:"dtype"`
	Country              string  `boil:"country" json:"country" toml:"country" yaml:"country"`
	Revenue              float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64 `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	SoldImpressions      float64 `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions float64 `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	PubFillRate          float64 `boil:"fill_rate" json:"fill_rate" toml:"fill_rate" yaml:"fill_rate"`
	CPM                  float64 `boil:"cpm" json:"cpm" toml:"cpm" yaml:"cpm"`
	RPM                  float64 `boil:"rpm" json:"rpm" toml:"rpm" yaml:"rpm"`
	DpRPM                float64 `boil:"dp_rpm" json:"dp_rpm" toml:"dp_rpm" yaml:"dp_rpm"`
	GP                   float64 `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	GPP                  float64 `boil:"gpp" json:"gpp" toml:"gpp" yaml:"gpp"`
}

func SendCustomHTMLEmail(to, bcc, subject string, htmlBody string, report []RealTimeReport, reportName string) error {
	toRecipients := strings.Split(to, ",")

	csvData, err := createCSVData(report)
	if err != nil {
		return err
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

	reportType := reflect.TypeOf(EmailReport{})

	var columnNames []string
	for i := 0; i < reportType.NumField(); i++ {
		field := reportType.Field(i)
		columnNames = append(columnNames, field.Name)
	}

	if err := writer.Write(columnNames); err != nil {
		return nil, err
	}
	//formatter := &helpers.FormatValues{}

	for _, record := range report {
		formattedTime := helpers.FormatDate(record.Time)
		row := []string{
			formattedTime,
			record.PublisherID,
			record.Publisher,
			record.Domain,
			fmt.Sprintf("%.2f", record.BidRequests),
			record.Device,
			record.Country,
			fmt.Sprintf("%.2f", record.Revenue),
			fmt.Sprintf("%.2f", record.Cost),
			fmt.Sprintf("%.2f", record.SoldImpressions),
			fmt.Sprintf("%.2f", record.PublisherImpressions),
			fmt.Sprintf("%.2f", record.PubFillRate),
			fmt.Sprintf("%.2f", record.CPM),
			fmt.Sprintf("%.2f", record.RPM),
			fmt.Sprintf("%.2f", record.DpRPM),
			fmt.Sprintf("%.2f", record.GP),
			fmt.Sprintf("%.2f", record.GPP),
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
