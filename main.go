package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/christodenny/lastwin/teamfetcher"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, teamfetcher.GetCfbTeams())
}

func main() {
	fmt.Println(teamfetcher.GetCfbTeams())
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not provided, using default port 8080")
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	// log.Fatal(http.ListenAndServe(":"+port, r))
}
