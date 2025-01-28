package adstxt

import (
	"fmt"
)

type AdsTxtLinesCreater interface {
}

type AdsTxtModule struct {
	adsTxtTaskChan chan AdsTxtTask
}

func NewAdsTxtModule(adsTxtTaskChan chan AdsTxtTask) *AdsTxtModule {
	go func() {
		// TODO:
		for {
			select {
			case task := <-adsTxtTaskChan:
				fmt.Printf("task - %#v\n", task)
				return
			}
		}
	}()

	return &AdsTxtModule{
		adsTxtTaskChan: adsTxtTaskChan,
	}
}

// TODO:
func (a *AdsTxtModule) CreateNewAdsTxtLines() error {
	return nil
}
