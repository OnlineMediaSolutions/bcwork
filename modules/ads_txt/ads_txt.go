package adstxt

import (
	"fmt"

	"github.com/m6yf/bcwork/dto"
)

type AdsTxtLinesCreater interface {
}

type AdsTxtModule struct {
	adsTxtCreateTask chan dto.AdsTxt
}

func NewAdsTxtModule(adsTxtCreateTask chan dto.AdsTxt) *AdsTxtModule {
	go func() {
		// TODO:
		for {
			select {
			case task := <-adsTxtCreateTask:
				fmt.Printf("task - %#v\n", task)
				return
			}
		}
	}()

	return &AdsTxtModule{
		adsTxtCreateTask: adsTxtCreateTask,
	}
}

// TODO:
func (a *AdsTxtModule) CreateNewAdsTxtLines() error {
	return nil
}
