package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := Response{
			Message: "Running2",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
	})

	log.Println("app2 listening on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}
