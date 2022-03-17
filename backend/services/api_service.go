package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type ApiService struct {
	Endpoint string
}

func (cnt *ApiService) GetData() ([]ViewDataUser, error) {
	var accessKeys AccessKeysResponse
	var metrics MetricsResponse
	var err error
	var errMutex sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

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
		return nil, err
	}

	var items []ViewDataUser

	for k := range accessKeys.AccessKeys {
		key := accessKeys.AccessKeys[k]
		item := ViewDataUser{
			Name:  key.Name,
			Usage: metrics.BytesTransferredByUserId[key.ID],
		}
		items = append(items, item)
	}

	return items, nil
}

func (cnt *ApiService) getMetrics() (MetricsResponse, error) {
	url := cnt.Endpoint + "/metrics/transfer"
	resp, err := http.Get(url)
	if err != nil {
		return MetricsResponse{}, err
	}
	defer resp.Body.Close()

	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MetricsResponse{}, err
	}

	var data MetricsResponse

	err = json.Unmarshal(respString, &data)
	if err != nil {
		return MetricsResponse{}, err
	}

	return data, nil
}

func (cnt *ApiService) getAccessKeys() (AccessKeysResponse, error) {
	url := cnt.Endpoint + "/access-keys"
	resp, err := http.Get(url)
	if err != nil {
		return AccessKeysResponse{}, err
	}
	defer resp.Body.Close()

	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AccessKeysResponse{}, err
	}

	var data AccessKeysResponse

	err = json.Unmarshal(respString, &data)
	if err != nil {
		return AccessKeysResponse{}, err
	}

	return data, nil
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

type ViewDataUser struct {
	Name  string `json:"name"`
	Usage int    `json:"usage"`
}
