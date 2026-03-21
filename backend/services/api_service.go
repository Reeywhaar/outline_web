package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type apiService struct {
	Base string
	ctx  context.Context
}

func NewApiService(base string, ctx context.Context) *apiService {
	return &apiService{
		Base: base,
		ctx:  ctx,
	}
}

func (cnt *apiService) GetData() (ApiServiceData, error) {
	var serverInfo ServerInfoResponse
	var accessKeys AccessKeysResponse
	var metrics MetricsResponse

	var wg errgroup.Group
	wg.Go(func() error {
		var lerr error
		serverInfo, lerr = cnt.getServerInfo()
		return lerr
	})

	wg.Go(func() error {
		var lerr error
		accessKeys, lerr = cnt.getAccessKeys()
		return lerr
	})

	wg.Go(func() error {
		var lerr error
		metrics, lerr = cnt.getMetrics()
		return lerr
	})

	if err := wg.Wait(); err != nil {
		return ApiServiceData{}, err
	}

	var users = make([]ApiServiceDataUser, 0)

	for k := range accessKeys.AccessKeys {
		key := accessKeys.AccessKeys[k]
		item := ApiServiceDataUser{
			Name:  key.Name,
			Usage: metrics.BytesTransferredByUserId[key.ID],
		}
		users = append(users, item)
	}

	data := ApiServiceData{
		Name:     serverInfo.Name,
		ServerID: serverInfo.ServerID,
		Users:    users,
	}

	return data, nil
}

func (cnt *apiService) getServerInfo() (ServerInfoResponse, error) {
	var data ServerInfoResponse
	err := cnt.callEndpoint("/server", &data)

	if err != nil {
		return ServerInfoResponse{}, err
	}

	return data, nil
}

func (cnt *apiService) getMetrics() (MetricsResponse, error) {
	var data MetricsResponse
	err := cnt.callEndpoint("/metrics/transfer", &data)

	if err != nil {
		return MetricsResponse{}, err
	}

	return data, nil
}

func (cnt *apiService) getAccessKeys() (AccessKeysResponse, error) {
	var data AccessKeysResponse
	err := cnt.callEndpoint("/access-keys", &data)

	if err != nil {
		return AccessKeysResponse{}, err
	}

	return data, nil
}

func (cnt *apiService) callEndpoint(endpoint string, data any) error {
	url := cnt.Base + endpoint
	req, err := http.NewRequestWithContext(cnt.ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed with status %d", resp.StatusCode)
	}

	respString, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respString, data)
	if err != nil {
		return err
	}

	return nil
}

type ServerInfoResponse struct {
	Name     string `json:"name"`
	ServerID string `json:"serverId"`
}

type MetricsResponse struct {
	BytesTransferredByUserId map[string]int `json:"bytesTransferredByUserId"`
}

type AccessKeyData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	AccessURL string `json:"accessUrl"`
}

type AccessKeysResponse struct {
	AccessKeys []AccessKeyData `json:"accessKeys"`
}

type ApiServiceData struct {
	Name     string               `json:"name"`
	ServerID string               `json:"server_id"`
	Users    []ApiServiceDataUser `json:"users"`
}

type ApiServiceDataUser struct {
	Name  string `json:"name"`
	Usage int    `json:"usage"`
}
