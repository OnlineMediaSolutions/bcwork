package rest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/report/iiq"
	"net/http"
	"strings"
	"time"
)

// IiqTestingGetRequest contains filter parameters for retrieving events
type IiqTestingGetRequest struct {
}

// IiqTestingGetResponse IiqTesting array
type IiqTestingGetResponse struct {
	// in: body

}

// IiqTestingGetHandler fetching challenges based on filter with pagination,order and selected fields.
// @Description fetching challenges based on filter with pagination,order and selected fields
// @Summary Get IiqTestings.
// @Tags IiqTesting
// @Accept json
// @Produce json
// @Param options body IiqTestingGetRequest true "IiqTesting Get Options"
// @Success 200 {object} IiqTestingGetResponse
// @Security ApiKeyAuth
// @Router /challenge/get [post]
func IiqTestingGetHandler(c *fiber.Ctx) error {

	var err error
	var demands []string
	if d := c.Query("demand"); d != "" {
		demands = strings.Split(d, "|")
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

	data, err := iiq.IiqReportHourlyCSV(c.Context(), iiq.IiqReportOptions{
		FromTime:     firstOfMonth,
		ToTime:       lastOfMonth,
		WithTotals:   totals,
		WithHeaders:  header,
		DemandFilter: demands,
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

	filename := fmt.Sprintf("iiq.report.%s.csv", time.Now().Format("200601021504"))
	c.Set("Content-Description", "File Transfer")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Type", "text/csv")

	return c.Send(b.Bytes())
}
