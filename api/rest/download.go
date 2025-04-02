package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
	"github.com/m6yf/bcwork/utils/constant"
)

const (
	downloadRequestTypeAdsTxtMainKey      = "ads_txt/main"
	downloadRequestTypeAdsTxtGroupByDPKey = "ads_txt/group_by_dp"
)

// DownloadHandler Download body data as file according to format in request
// @Description Download body data as file according to format in request. Data should be passed as array of json objects which have same structure
// @Tags Download
// @Accept json
// @Produce json
// @Param options body dto.DownloadRequest true "request"
// @Success 200 {object} utils.BaseResponse
// @Router /download [post]
func (o *OMSNewPlatform) DownloadHandler(c *fiber.Ctx) error {
	var req *dto.DownloadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error parsing download request", err)
	}

	if req.Request.Type != "" {
		data, err := o.getDataForFile(c.Context(), req.Request.Type, req.Request.Body)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting data for file by request", err)
		}

		req.Data = data
	}

	fileData, err := o.downloadService.CreateFile(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error creating %v file", req.FileFormat), err)
	}

	return sendFile(c, req.FilenamePrefix, fileData, req.FileFormat)
}

func (o *OMSNewPlatform) getDataForFile(ctx context.Context, requestType string, requestBody []byte) ([]json.RawMessage, error) {
	switch requestType {
	case downloadRequestTypeAdsTxtMainKey:
		var ops *core.AdsTxtGetMainOptions
		err := json.Unmarshal(requestBody, &ops)
		if err != nil {
			return nil, err
		}

		data, err := o.adsTxtService.GetMainAdsTxtTable(ctx, ops)
		if err != nil {
			return nil, err
		}

		result := make([]json.RawMessage, 0, len(data.Data))
		for _, row := range data.Data {
			byteRow, err := json.Marshal(row)
			if err != nil {
				return nil, err
			}
			result = append(result, byteRow)
		}

		return result, nil
	case downloadRequestTypeAdsTxtGroupByDPKey:
		var ops *core.AdsTxtGetGroupByDPOptions
		err := json.Unmarshal(requestBody, &ops)
		if err != nil {
			return nil, err
		}

		data, err := o.adsTxtService.GetGroupByDPAdsTxtTable(ctx, ops)
		if err != nil {
			return nil, err
		}

		result := make([]json.RawMessage, 0, len(data.Data))
		for _, row := range data.Data {
			byteRow, err := json.Marshal(row)
			if err != nil {
				return nil, err
			}
			result = append(result, byteRow)
		}

		return result, nil
	}

	return nil, fmt.Errorf("unknown request type [%v]", requestType)
}

func sendFile(c *fiber.Ctx, filenamePrefix string, data []byte, format dto.DownloadFormat) error {
	filename := fmt.Sprintf("%v.%v.%v", filenamePrefix, time.Now().Format("2006_01_02_15_04_05"), format)
	c.Set(constant.HeaderContentDescription, "File Transfer")
	c.Set(fiber.HeaderContentDisposition, "attachment; filename="+filename)
	switch format {
	case dto.CSV:
		c.Set(fiber.HeaderContentType, constant.MIMETextCSV)
	case dto.XLSX:
		c.Set(fiber.HeaderContentType, fiber.MIMEOctetStream)
	}

	return c.Send(data)
}
