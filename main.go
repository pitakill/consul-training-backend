package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var count = 0

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// App logic (very basic)
		count++

		// Handle response
		j, err := json.Marshal(struct {
			Count int `json:"count"`
		}{count})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(j)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
