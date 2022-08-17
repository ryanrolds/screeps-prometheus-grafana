package main

import (
	"errors"
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
	Name              string `yaml:"name"`
	Host              string `yaml:"host"`
	Path              string `yaml:"path"`
	Shard             string `yaml:"shard"`
	OverrideShardName string `yaml:"overrideShardName"`
	Token             string `yaml:"token"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
}

func (s *Server) Validate() error {
	if s.Token == "" && s.Username == "" && s.Password == "" {
		return errors.New("No auth credentials provided")
	}

	return nil
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

	for _, server := range servers.Servers {
		err = server.Validate()
		if err != nil {
			log.Fatal(err)
		}

		scraper := NewScreepsScraper().WithShard(server.Shard)

		if server.Host != "" {
			scraper.WithHost(server.Host)
		}

		if server.Path != "" {
			scraper.WithPath(server.Path)
		}

		if server.Token != "" {
			scraper.WithToken(server.Token)
		} else if server.Username != "" && server.Password != "" {
			scraper.WithUserPass(server.Username, server.Password)
		}

		if server.OverrideShardName != "" {
			scraper.WithOverrideShardName(server.OverrideShardName)
		}

		log.Infof("Registering metrics from %s %s", server.Host, server.Shard)
		reg.MustRegister(scraper)
	}

	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":8080"
	}

	log.Infof("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
