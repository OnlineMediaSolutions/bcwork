package report

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/bcdwh"
	"github.com/m6yf/bcwork/core/report/revenue"
	"net/http"
	"strings"
	"time"
)

func DailyRevenue(c *fiber.Ctx) error {

	db := bcdb.DB()
	if c.Query("src") == "dwh" {
		db = bcdwh.DB()
	}

	if strings.ToLower(c.Query("auth")) != "omsrev2024" {
		return c.SendStatus(http.StatusForbidden)
	}

	var month string
	if month = c.Query("month"); month != "" {
		if len(month) != 6 {
			c.SendString("illegal 'month' format (should be YYYYMM)")
			return c.SendStatus(http.StatusBadRequest)
		}
	} else {
		month = time.Now().Format("200601")
	}

	htmlReport, err := revenue.DailyHtmlReport(c.Context(), month, db)
	if err != nil {
		return err
	}
	if c.Query("src") == "dwh" {
		htmlReport = strings.Replace(htmlReport, "Revenue Report", "Revenue Report (DWH)", -1)
	}

	c.Set("Content-Type", "text/html")

	return c.SendString(htmlReport)
}

func HourlyRevenue(c *fiber.Ctx) error {

	db := bcdb.DB()
	if c.Query("src") == "dwh" {
		db = bcdwh.DB()
	}

	if strings.ToLower(c.Query("auth")) != "omsrev2024" {
		return c.SendStatus(http.StatusForbidden)
	}

	var date string
	if date = c.Query("date"); date != "" {
		if len(date) != 8 {
			c.SendString("illegal 'date' format (should be YYYYMMDD)")
			return c.SendStatus(http.StatusBadRequest)
		}
	} else {
		date = time.Now().Format("20060102")
	}

	htmlReport, err := revenue.HourlyHtmlReport(c.Context(), date, db)
	if err != nil {
		return err
	}

	if c.Query("src") == "dwh" {
		htmlReport = strings.Replace(htmlReport, "Revenue Report", "Revenue Report (DWH)", -1)
	}
	c.Set("Content-Type", "text/html")

	return c.SendString(htmlReport)
}
