package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/m6yf/bcwork/models"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Updater interface {
	Update(context.Context, *models.MetadataQueue) error
}

type HttpUpdater struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func (updater HttpUpdater) Update(ctx context.Context, record *models.MetadataQueue) error {

	body, err := json.Marshal(map[string]interface{}{
		"key":   record.Key,
		"value": string(record.Value),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to marshal metadata payload")
	}
	resp, err := http.Post(updater.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrapf(err, "failed to send http post metadata payload")
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("error while sending http post metadata(code:%d)", resp.StatusCode)
	}

	return nil
}

type RedisUpdater struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (updater RedisUpdater) Update(ctx context.Context, record *models.MetadataQueue) error {

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", updater.Host, updater.Port),
	})
	defer rdb.Close()

	err := rdb.Set(ctx, record.Key, record.Value, 0).Err()
	if err != nil {
		return errors.Wrapf(err, "failed to send redis metadata set")
	}

	return nil
}
