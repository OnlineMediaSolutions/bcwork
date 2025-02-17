package dto

import (
	"encoding/json"
)

type DownloadFormat string

const (
	CSV  = "csv"
	XLSX = "xlsx"
)

type DownloadRequest struct {
	FilenamePrefix string            `json:"filename_prefix"`
	Columns        []*Column         `json:"columns"`
	FileFormat     DownloadFormat    `json:"file_format"`
	Data           []json.RawMessage `json:"data"`
}

type Column struct {
	Name               string              `json:"name"`
	DisplayName        string              `json:"display_name"`
	Style              string              `json:"style"`
	Multiply           *float64            `json:"multiply"`
	BooleanReplacement *BooleanReplacement `json:"boolean_replacement"`
}

func (c Column) GetMultiply() float64 {
	if c.Multiply != nil {
		return *c.Multiply
	}

	return 1
}

func (c Column) GetBooleanReplacementValue(isTrue bool) string {
	if c.BooleanReplacement != nil {
		if isTrue {
			return c.BooleanReplacement.True
		}

		return c.BooleanReplacement.False
	}

	return ""
}

type BooleanReplacement struct {
	True  string `json:"true"`
	False string `json:"false"`
}
