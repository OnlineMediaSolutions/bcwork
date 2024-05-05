package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core/report/monthly"
	"net/http"
	"time"
)

func MonthlyRevenueReport(c *fiber.Ctx) error {

	var month string
	if month = c.Query("month"); month != "" {
		if len(month) != 6 {
			c.SendString("illegal 'month' format (should be YYYYMM)")
			return c.SendStatus(http.StatusBadRequest)
		}
	} else {
		month = time.Now().Format("200601")
	}

	htmlReport, err := monthly.MonthlyHtmlReport(c.Context(), month)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")

	return c.SendString(htmlReport)
}
