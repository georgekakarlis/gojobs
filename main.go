package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gojobs.sh/config"
	"gojobs.sh/database"
	"gojobs.sh/handlers"
)


func init () {
	config.LoadEnvVariables()
	// db connect
	database.ConnectDB()
}

func main() {
	

	r := mux.NewRouter()

	// handle static
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Handle routes
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/postajob", handlers.Postajob)
	r.HandleFunc("/api/auth/login", handlers.LoginHandler)
	r.HandleFunc("/callback", handlers.HandleGoogleCallback)
	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "You've requested the healthcheck url")
    })
	


	fmt.Println("Server is starting on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error occurred: %s", err)
	}
}