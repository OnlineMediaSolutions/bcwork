package adstxt

import (
	"github.com/m6yf/bcwork/dto"
)

type AdsTxtLinesCreater interface {
}

type AdsTxtModule struct {
	adsTxtCreateTask chan dto.AdsTxt
}

func NewAdsTxtModule() *AdsTxtModule {
	return &AdsTxtModule{}
}

// TODO:
func (a *AdsTxtModule) CreateNewAdsTxtLines() error {
	return nil
}
