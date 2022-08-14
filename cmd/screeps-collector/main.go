package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":8080"
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is required")
	}

	// Create a new registry.
	reg := prometheus.NewRegistry()
	reg.MustRegister(NewScreepsCollector(token))

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Printf("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
