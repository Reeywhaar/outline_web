package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/reeywhaar/outline_web/backend/services"
)

type apiController struct {
	Servers []string
}

func NewApiController(servers []string, ctx context.Context) *apiController {
	return &apiController{
		Servers: servers,
	}
}

func (cnt *apiController) HandleServers(w http.ResponseWriter, r *http.Request) {
	jsd, err := json.Marshal(makeRange(0, len(cnt.Servers)-1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsd)
}

func (cnt *apiController) HandleServersID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unknown server", http.StatusInternalServerError)
		return
	}

	if index < 0 || index > (len(cnt.Servers)-1) {
		http.Error(w, "Unknown server", http.StatusInternalServerError)
		return
	}

	url := cnt.Servers[index]
	controller := services.NewApiService(url, r.Context())

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

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
