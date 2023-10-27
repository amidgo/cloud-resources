package pricestorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/pricemodel"
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

// получает цены из апишки
func (s *Storage) PriceList(ctx context.Context) ([]*pricemodel.Price, error) {
	var priceList []*pricemodel.Price
	resp, err := s.client.MakeRequest(ctx, http.MethodGet, "/price", "", &bytes.Buffer{})
	if err != nil {
		return nil, fmt.Errorf("failed make request, %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed make request, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	err = json.NewDecoder(resp.Body).Decode(&priceList)
	if err != nil {
		return nil, fmt.Errorf("failed decode body, %w", err)
	}
	return priceList, nil
}
