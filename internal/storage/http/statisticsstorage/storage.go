package statisticsstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/statiscticsmodel"
	"github.com/amidgo/cloud-resources/pkg/httpclient"
	"io"
	"net/http"
)

type Storage struct {
	client *httpclient.HttpClient
}

func New(client *httpclient.HttpClient) *Storage {
	return &Storage{client: client}
}

// получение статистика из апишки
func (s *Storage) Statisctics(ctx context.Context) (statiscticsmodel.Statistics, error) {
	var statistics statiscticsmodel.Statistics
	resp, err := s.client.MakeRequest(ctx, http.MethodGet, "/statistic", "", &bytes.Buffer{})
	if err != nil {
		return statistics, fmt.Errorf("failed make request, %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return statistics, fmt.Errorf("failed make request, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	err = json.NewDecoder(resp.Body).Decode(&statistics)
	if err != nil {
		return statistics, fmt.Errorf("failed decode request body, %w", err)
	}
	return statistics, nil
}
