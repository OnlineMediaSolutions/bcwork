package rest

import (
	"github.com/gofiber/fiber/v2"
)

// DemandReportGetRequest contains filter parameters for retrieving events
type DemandReportGetRequest struct {
}

// DemandReportGetResponse DemandReport array
type DemandReportGetResponse struct {
	// in: body

}

// DemandReportGetHandler fetching challenges based on filter with pagination,order and selected fields.
// @Description fetching challenges based on filter with pagination,order and selected fields
// @Summary Get DemandReports.
// @Tags DemandReport
// @Accept json
// @Produce json
// @Param options body DemandReportGetRequest true "DemandReport Get Options"
// @Success 200 {object} DemandReportGetResponse
// @Security ApiKeyAuth
// @Router /challenge/get [post]
func DemandReportGetHandler(c *fiber.Ctx) error {

	//var err error
	//var publishers []string
	//if pub := c.Query("publisher"); pub != "" {
	//	publishers = strings.Split(pub, "|")
	//}
	//
	//var demands []string
	//if d := c.Query("demand"); d != "" {
	//	demands = strings.Split(d, "|")
	//}
	//
	//var domains []string
	//if domain := c.Query("domain"); domain != "" {
	//	domains = strings.Split(domain, "|")
	//}
	//
	//var full bool
	//if d := c.Query("full"); d != "" && (d == "1" || d == "true" || d == "on") {
	//	full = true
	//}
	//now := time.Now()
	//currentYear, currentMonth, currentDay := now.Date()
	//currentLocation := now.Location()
	//
	//firstOfMonth := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	//lastOfMonth := firstOfMonth.AddDate(0, 1, 0)
	//
	//if fromstr := c.Query("from"); fromstr != "" {
	//	firstOfMonth, err = time.Parse("20060102", fromstr)
	//	if err != nil {
	//		c.SendString("illegal 'from' param")
	//		return c.SendStatus(http.StatusBadRequest)
	//	}
	//}
	//
	//if tostr := c.Query("to"); tostr != "" {
	//	lastOfMonth, err = time.Parse("20060102", tostr)
	//	if err != nil {
	//		c.SendString("illegal 'to' param")
	//		return c.SendStatus(http.StatusBadRequest)
	//	}
	//}

	//htmlReport, err := demand.DemandReportDaily(c.Context(), firstOfMonth, lastOfMonth, publishers, demands, full, domains)
	//if err != nil {
	//	return err
	//}

	c.Set("Content-Type", "text/html")
	return c.SendString("")

	//return c.SendString(htmlReport)
}

func DemandHourlyReportGetHandler(c *fiber.Ctx) error {

	//var err error
	//var publishers []string
	//if pub := c.Query("publisher"); pub != "" {
	//	publishers = strings.Split(pub, "|")
	//}
	//
	//var demands []string
	//if d := c.Query("demand"); d != "" {
	//	demands = strings.Split(d, "|")
	//}
	//
	//var domains []string
	//if domain := c.Query("domain"); domain != "" {
	//	domains = strings.Split(domain, "|")
	//}
	//
	//var full bool
	//if d := c.Query("full"); d != "" && (d == "1" || d == "true" || d == "on") {
	//	full = true
	//}
	//
	//now := time.Now()
	//currentYear, currentMonth, currentDay := now.Date()
	//currentLocation := now.Location()
	//
	//firstOfMonth := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, currentLocation)
	//lastOfMonth := firstOfMonth.AddDate(0, 1, 0)
	//
	//if fromstr := c.Query("from"); fromstr != "" {
	//	firstOfMonth, err = time.Parse("2006010215", fromstr)
	//	if err != nil {
	//		c.SendString("illegal 'from' param")
	//		return c.SendStatus(http.StatusBadRequest)
	//	}
	//}
	//
	//if tostr := c.Query("to"); tostr != "" {
	//	lastOfMonth, err = time.Parse("2006010215", tostr)
	//	if err != nil {
	//		c.SendString("illegal 'to' param")
	//		return c.SendStatus(http.StatusBadRequest)
	//	}
	//}
	//
	//htmlReport, err := demand.DemandReportHourly(c.Context(), firstOfMonth, lastOfMonth, publishers, demands, full, domains)
	//if err != nil {
	//	return err
	//}
	//
	//c.Set("Content-Type", "text/html")

	//return c.SendString(htmlReport)
	return c.SendString("")

}
