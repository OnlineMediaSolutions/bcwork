package core

import (
	"context"
	"fmt"

	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/modules/export"
)

type DownloadService struct {
	exporter export.Exporter
}

func NewDownloadService(exporter export.Exporter) *DownloadService {
	return &DownloadService{
		exporter: exporter,
	}
}

func (d *DownloadService) CreateFile(ctx context.Context, req *dto.DownloadRequest) ([]byte, error) {
	switch req.FileFormat {
	case dto.CSV:
		return d.exporter.ExportCSV(ctx, req.Data)
	case dto.XLSX:
		return d.exporter.ExportXLSX(ctx, req)
	}

	return nil, fmt.Errorf("unknown format [%v]", req.FileFormat)
}
