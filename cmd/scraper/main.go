package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Servers struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Path     string `yaml:"path"`
	Token    string `yaml:"token"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func main() {
	// Setup logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Fetch and pass config.yaml to get Servers to fetch metrics from
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var servers Servers
	err = yaml.Unmarshal(configFile, &servers)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Found %d servers", len(servers.Servers))

	// Create a new registry
	reg := prometheus.NewRegistry()

	// Scraper for MMO shards and seasonal events
	mmoScraper := NewScreepsScraper().
		//WithToken(token).
		WithPath("metrics").
		WithShards([]string{"shard2"})
	reg.MustRegister(mmoScraper)

	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":8080"
	}

	log.Infof("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
