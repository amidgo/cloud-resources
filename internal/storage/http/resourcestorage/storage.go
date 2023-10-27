package resourcestorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amidgo/cloud-resources/internal/model/resourcemodel"
	"github.com/amidgo/cloud-resources/pkg/httpclient"
	"io"
	"log"
	"net/http"
)

type Storage struct {
	client *httpclient.HttpClient
}

func New(client *httpclient.HttpClient) *Storage {
	return &Storage{client: client}
}

// добавление ресурса
func (s *Storage) AddResource(ctx context.Context, addResource resourcemodel.AddResource) (resourcemodel.Resource, error) {
	var resource resourcemodel.Resource
	body, err := json.Marshal(addResource)
	if err != nil {
		return resource, fmt.Errorf("failed marshal add resource, %w", err)
	}
	resp, err := s.client.MakeRequest(ctx, http.MethodPost, "/resource", "application/json", bytes.NewReader(body))
	if err != nil {
		return resource, fmt.Errorf("failed make request, %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return resource, fmt.Errorf("request failed, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	err = json.NewDecoder(resp.Body).Decode(&resource)
	if err != nil {
		return resource, fmt.Errorf("failed decode body, %w", err)
	}
	return resource, nil
}

// удаление ресурса
func (s *Storage) DeleteResource(ctx context.Context, id int) error {
	url := fmt.Sprintf("/resource/%d", id)
	resp, err := s.client.MakeRequest(ctx, http.MethodDelete, url, "", &bytes.Buffer{})
	if err != nil {
		return fmt.Errorf("failed delete resource, %w", err)
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	return nil
}

// обновление ресурса
func (s *Storage) UpdateResource(ctx context.Context, id int, updateResource resourcemodel.UpdateResource) error {
	url := fmt.Sprintf("/resource/%d", id)
	body, err := json.Marshal(updateResource)
	if err != nil {
		return fmt.Errorf("failed marshal body, %w", err)
	}
	resp, err := s.client.MakeRequest(ctx, http.MethodPut, url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed delete resource, %w", err)
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		log.Printf("id: %d, resource: %v\n, rawBody: %s", id, updateResource, string(body))
		return fmt.Errorf("request failed, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	return nil
}

// получение ресурса по идентификатору
func (s *Storage) Resource(ctx context.Context, id int) (resourcemodel.Resource, error) {
	var resource resourcemodel.Resource
	url := fmt.Sprintf("/resource/%d", id)
	resp, err := s.client.MakeRequest(ctx, http.MethodGet, url, "", &bytes.Buffer{})
	if err != nil {
		return resource, fmt.Errorf("failed make request, %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return resource, fmt.Errorf("request failed, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	err = json.NewDecoder(resp.Body).Decode(&resource)
	if err != nil {
		return resource, fmt.Errorf("failed decode body, %w", err)
	}
	return resource, nil
}

// получение всех ресурсов
func (s *Storage) ResourceList(ctx context.Context) ([]*resourcemodel.Resource, error) {
	var resourceList []*resourcemodel.Resource
	resp, err := s.client.MakeRequest(ctx, http.MethodGet, "/resource", "", &bytes.Buffer{})
	if err != nil {
		return nil, fmt.Errorf("failed make request, %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		message, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed, wrong status code %d, message %s", resp.StatusCode, string(message))
	}
	err = json.NewDecoder(resp.Body).Decode(&resourceList)
	if err != nil {
		return nil, fmt.Errorf("failed decode body, %w", err)
	}
	return resourceList, nil
}
