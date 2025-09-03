package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/reeywhaar/outline_web/backend/controllers"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r := mux.NewRouter()

	r.HandleFunc("/", handleMain)

	servers := strings.Split(os.Getenv("OUTLINE_API_URL"), ",")

	apiController := controllers.ApiController{
		Servers: servers,
	}
	r.HandleFunc("/api/servers", apiController.HandleServers)
	r.HandleFunc("/api/servers/{id}", apiController.HandleServersID)

	http.Handle("/", r)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	bind_addr := os.Getenv("ADDR")
	if bind_addr == "" {
		bind_addr = "127.0.0.1"
	}

	addr := bind_addr + ":" + os.Getenv("PORT")
	log.Printf("Starting server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, logRequest(http.DefaultServeMux)))
}

//

func handleMain(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, nil)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
