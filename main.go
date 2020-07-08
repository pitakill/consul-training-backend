package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const VERSION = "2.0.0"

var count = 0

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		// Better with a middleware
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// App logic (very basic)
		count++

		// Handle response
		j, err := json.Marshal(struct {
			Count   int    `json:"count"`
			Version string `json:"version"`
		}{count, VERSION})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(j)
	})

	ifaces, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("Can't get the ifaces where is listening")
	}

	for _, iface := range ifaces {
		ip := strings.Split(iface.String(), "/")
		log.Printf("Listening on: %s:%s", ip[0], port)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
