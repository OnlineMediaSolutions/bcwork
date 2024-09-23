package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/utils"
)

// DownloadPostHandler Download body data as csv
// @Description Download body data as csv. Data should be passed as array of json objects which have same structure
// @Tags Download
// @Accept json
// @Produce json
// @Param options body []core.DownloadDataExample true "options"
// @Success 200 {object} utils.BaseResponse
// @Router /download [post]
func DownloadPostHandler(c *fiber.Ctx) error {
	var data []json.RawMessage
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error parsing download request", err)
	}

	b, err := core.CreateCSVFile(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating CSV file", err)
	}

	// TODO: receive the place where endpoint was invoked and add this information to filename
	filename := fmt.Sprintf("file.%s.csv", time.Now().Format("20060102150405"))
	c.Set("Content-Description", "File Transfer")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	c.Set("Content-Type", "text/csv")

	return c.Send(b)
}
