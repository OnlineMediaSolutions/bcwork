package rest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/report/publisher"
	"net/http"
	"strings"
	"time"
)

// PublisherReportGetRequest contains filter parameters for retrieving events
type PublisherReportGetRequest struct {
}

// PublisherReportGetResponse PublisherReport array
type PublisherReportGetResponse struct {
	// in: body

}

// PublisherReportGetHandler fetching challenges based on filter with pagination,order and selected fields.
// @Description fetching challenges based on filter with pagination,order and selected fields
// @Summary Get PublisherReports.
// @Tags PublisherReport
// @Accept json
// @Produce json
// @Param options body PublisherReportGetRequest true "PublisherReport Get Options"
// @Success 200 {object} PublisherReportGetResponse
// @Security ApiKeyAuth
// @Router /report/publisher [get]
func PublisherReportGetHandler(c *fiber.Ctx) error {

	var err error
	var publishers []string
	if pub := c.Query("publisher"); pub != "" {
		publishers = strings.Split(pub, "|")
	}

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0)

	if fromstr := c.Query("from"); fromstr != "" {
		firstOfMonth, err = time.Parse("20060102", fromstr)
		if err != nil {
			c.SendString("illegal 'from' param")
			return c.SendStatus(http.StatusBadRequest)
		}
	}

	if tostr := c.Query("to"); tostr != "" {
		lastOfMonth, err = time.Parse("20060102", tostr)
		if err != nil {
			c.SendString("illegal 'to' param")
			return c.SendStatus(http.StatusBadRequest)
		}
	}

	var header bool
	if str := strings.ToLower(c.Query("header")); str == "on" || str == "true" || str == "1" {
		header = true
	}

	var totals bool
	if str := strings.ToLower(c.Query("header")); str == "on" || str == "true" || str == "1" {
		totals = true
	}

	data, err := publisher.PublisherReportDailyCSV(c.Context(), publisher.PublisherReportOptions{
		FromTime:        firstOfMonth,
		ToTime:          lastOfMonth,
		WithTotals:      totals,
		WithHeaders:     header,
		PublisherFilter: publishers,
	})
	if err != nil {
		return err
	}

	var b bytes.Buffer
	csvwriter := csv.NewWriter(&b)
	for _, d := range data {
		s := make([]string, len(d))
		for i, v := range d {
			s[i] = fmt.Sprint(v)
		}
		_ = csvwriter.Write(s)
	}
	csvwriter.Flush()

	filename := fmt.Sprintf("publisher.report.%s.csv", time.Now().Format("200601021504"))
	c.Set("Content-Description", "File Transfer")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Type", "text/csv")

	return c.Send(b.Bytes())
}

func PublisherHourlyReportGetHandler(c *fiber.Ctx) error {

	var err error
	var publishers []string
	if pub := c.Query("publisher"); pub != "" {
		publishers = strings.Split(pub, "|")
	}

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 0, 1)

	if fromstr := c.Query("from"); fromstr != "" {
		firstOfMonth, err = time.Parse("20060102", fromstr)
		if err != nil {
			c.SendString("illegal 'from' param")
			return c.SendStatus(http.StatusBadRequest)
		}
	}

	if tostr := c.Query("to"); tostr != "" {
		lastOfMonth, err = time.Parse("20060102", tostr)
		if err != nil {
			c.SendString("illegal 'to' param")
			return c.SendStatus(http.StatusBadRequest)
		}
	}

	var header bool
	if str := strings.ToLower(c.Query("header")); str == "on" || str == "true" || str == "1" {
		header = true
	}

	var totals bool
	if str := strings.ToLower(c.Query("header")); str == "on" || str == "true" || str == "1" {
		totals = true
	}

	data, err := publisher.PublisherReportHourlyCSV(c.Context(), publisher.PublisherReportOptions{
		FromTime:        firstOfMonth,
		ToTime:          lastOfMonth,
		WithTotals:      totals,
		WithHeaders:     header,
		PublisherFilter: publishers,
	})
	if err != nil {
		return err
	}

	var b bytes.Buffer
	csvwriter := csv.NewWriter(&b)
	for _, d := range data {
		s := make([]string, len(d))
		for i, v := range d {
			s[i] = fmt.Sprint(v)
		}
		_ = csvwriter.Write(s)
	}
	csvwriter.Flush()

	filename := fmt.Sprintf("publisher.report.%s.csv", time.Now().Format("200601021504"))
	c.Set("Content-Description", "File Transfer")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Type", "text/csv")

	return c.Send(b.Bytes())
}
