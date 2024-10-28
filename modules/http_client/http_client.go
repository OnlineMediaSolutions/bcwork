package httpclient

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/utils/constant"
	"github.com/spf13/viper"
)

type Doer interface {
	Do(ctx context.Context, method, url, payload string) ([]byte, error)
	DoWithRequest(ctx context.Context, req *http.Request) ([]byte, error)
}

var _ Doer = (*HttpClient)(nil)

type HttpClient struct {
	httpClient *http.Client
}

func New() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{},
	}
}

func (h *HttpClient) Do(ctx context.Context, method, url, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(constant.HeaderOMSWorkerAPIKey, viper.GetString(config.CronWorkerAPIKeyKey))

	log.Println(viper.GetString(config.CronWorkerAPIKeyKey))

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *HttpClient) DoWithRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
