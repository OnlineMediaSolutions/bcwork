package httpclient

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
)

type Doer interface {
	Do(ctx context.Context, method, url string, body io.Reader) ([]byte, int, error)
	DoWithRequest(ctx context.Context, req *http.Request) ([]byte, int, error)
}

var _ Doer = (*HttpClient)(nil)

type HttpClient struct {
	isForInternalUse bool
	httpClient       *http.Client
}

func New(isForInternalUse bool) *HttpClient {
	return &HttpClient{
		isForInternalUse: isForInternalUse,
		httpClient:       &http.Client{},
	}
}

func (h *HttpClient) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	if h.isForInternalUse {
		req.Header.Add(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))
	}

	log.Println(viper.GetString(config.CronWorkerAPIKeyKey))

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return data, res.StatusCode, nil
}

func (h *HttpClient) DoWithRequest(ctx context.Context, req *http.Request) ([]byte, int, error) {
	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return data, res.StatusCode, nil
}
