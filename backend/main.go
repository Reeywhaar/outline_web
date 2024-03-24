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
	"github.com/reeywhaar/outline_web/backend/middlewares"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r := mux.NewRouter()

	authMiddleware := middlewares.AuthMiddleware{}
	authMiddleware.Init("/api/auth", os.Getenv("ADMIN_PASSWORD"))

	r.HandleFunc("/", handleMain)

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware.Middleware)
	apiController := controllers.ApiController{
		Servers: strings.Split(os.Getenv("OUTLINE_API_URL"), ","),
	}
	apiRouter.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {})
	apiRouter.HandleFunc("/servers", apiController.HandleServers)
	apiRouter.HandleFunc("/servers/{id}", apiController.HandleServersID)

	http.Handle("/", r)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	addr := "127.0.0.1:" + os.Getenv("PORT")
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
