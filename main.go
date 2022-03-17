package main

import (
	"crypto/tls"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/", handleMain)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/api/data", handleApiData)

	addr := "127.0.0.1:" + os.Getenv("PORT")
	log.Printf("Starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

//

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

type ApiController struct {
	endpoint string
}

func (cnt *ApiController) GetData() ([]ViewDataUser, error) {
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

func (cnt *ApiController) getMetrics() (MetricsResponse, error) {
	url := cnt.endpoint + "/metrics/transfer"
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

func (cnt *ApiController) getAccessKeys() (AccessKeysResponse, error) {
	url := cnt.endpoint + "/access-keys"
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

type ViewData struct {
	Data []ViewDataUser `json:"data"`
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	metrics := ViewData{
		Data: []ViewDataUser{
			{
				Name:  "Dima",
				Usage: 132123,
			},
		},
	}
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, metrics)
}

func handleApiData(w http.ResponseWriter, r *http.Request) {
	controller := ApiController{endpoint: os.Getenv("OUTLINE_API_URL")}

	data, err := controller.GetData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsd, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsd)
}
