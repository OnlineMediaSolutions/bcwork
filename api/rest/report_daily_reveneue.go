package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/report/monthly"
	"net/http"
	"time"
)

func DailyRevenueReport(c *fiber.Ctx) error {

	var date string
	if date = c.Query("date"); date != "" {
		if len(date) != 8 {
			c.SendString("illegal 'date' format (should be YYYYMMDD)")
			return c.SendStatus(http.StatusBadRequest)
		}
	} else {
		date = time.Now().Format("20060102")
	}

	htmlReport, err := monthly.DailyHtmlReport(c.Context(), date)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")

	return c.SendString(htmlReport)
}
