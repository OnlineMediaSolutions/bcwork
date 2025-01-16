package core

import (
	"context"

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
	default:
		return d.exporter.ExportXLSX(ctx, req)
	}
}
