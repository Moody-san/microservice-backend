package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os" // Import the os package
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the hostname from the environment variables
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Error getting hostname: %v", err)
		}

		res := Response{
			Message: "app 1 from pod -> " + hostname,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
}
