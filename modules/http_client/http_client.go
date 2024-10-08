package httpclient

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Doer interface {
	Do(ctx context.Context, method, url, payload string) ([]byte, error)
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
