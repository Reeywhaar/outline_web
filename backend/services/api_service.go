package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type ApiService struct {
	Endpoint string
}

func (cnt *ApiService) GetData() (ApiServiceData, error) {
	var serverInfo ServerInfoResponse
	var accessKeys AccessKeysResponse
	var metrics MetricsResponse
	var err error
	var errMutex sync.Mutex

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		var lerr error
		serverInfo, lerr = cnt.getServerInfo()
		if lerr != nil {
			errMutex.Lock()
			err = lerr
			errMutex.Unlock()
		}
		wg.Done()
	}()

	go func() {
		var lerr error
		accessKeys, lerr = cnt.getAccessKeys()
		if lerr != nil {
			errMutex.Lock()
			err = lerr
			errMutex.Unlock()
		}
		wg.Done()
	}()

	go func() {
		var lerr error
		metrics, lerr = cnt.getMetrics()
		if lerr != nil {
			errMutex.Lock()
			err = lerr
			errMutex.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()

	if err != nil {
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

func (cnt *ApiService) getServerInfo() (ServerInfoResponse, error) {
	var data ServerInfoResponse
	err := cnt.callEndpoint("/server", &data)

	if err != nil {
		return ServerInfoResponse{}, err
	}

	return data, nil
}

func (cnt *ApiService) getMetrics() (MetricsResponse, error) {
	var data MetricsResponse
	err := cnt.callEndpoint("/metrics/transfer", &data)

	if err != nil {
		return MetricsResponse{}, err
	}

	return data, nil
}

func (cnt *ApiService) getAccessKeys() (AccessKeysResponse, error) {
	var data AccessKeysResponse
	err := cnt.callEndpoint("/access-keys", &data)

	if err != nil {
		return AccessKeysResponse{}, err
	}

	return data, nil
}

func (cnt *ApiService) callEndpoint(endpoint string, data interface{}) error {
	url := cnt.Endpoint + endpoint
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed with status %d", resp.StatusCode)
	}

	respString, err := ioutil.ReadAll(resp.Body)
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
