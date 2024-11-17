package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
)

type Publisher struct {
	PublisherId string `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Name        string `boil:"name" json:"name" toml:"name" yaml:"name"`
}

func FetchPublishers(ctx context.Context) (map[string]string, error) {
	HttpClient := httpclient.New(true)

	requestBody := map[string]interface{}{
		"filter": map[string]interface{}{},
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, fmt.Errorf("error in marshaling body: %d", err)
	}

	publisherData, statusCode, err := HttpClient.Do(ctx, http.MethodPost, constant.ProductionApiUrl+"/publisher/get", bytes.NewBuffer(body))

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", statusCode)
	}

	var publishers []Publisher
	if err := json.Unmarshal(publisherData, &publishers); err != nil {
		return nil, errors.Wrapf(err, "Error parsing publisher data  from API")
	}

	publisherMap := make(map[string]string)
	for _, publisher := range publishers {
		publisherMap[publisher.PublisherId] = publisher.Name
	}

	return publisherMap, nil
}
