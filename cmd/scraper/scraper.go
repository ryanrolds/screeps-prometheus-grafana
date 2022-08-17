package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	authTypeToken    = "token"
	authTypeUserPass = "userpass"
)

type memoryResponse struct {
	Ok   int    `json:"ok"`
	Data string `json:"data"`
}

type metric struct {
	Key    string            `json:"key"`
	Value  float64           `json:"value"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
	Time   int               `json:"time"`
}

type ScreepScraper struct {
	host   string
	shards []string
	path   string
	client *http.Client

	authType string
	token    string
	username string
	password string
}

func NewScreepsScraper() *ScreepScraper {
	return &ScreepScraper{
		host:   "https://screeps.com",
		shards: []string{},
		client: &http.Client{},
	}
}

func (s *ScreepScraper) WithHost(host string) *ScreepScraper {
	s.host = host
	return s
}

func (s *ScreepScraper) WithPath(path string) *ScreepScraper {
	s.path = path
	return s
}

func (s *ScreepScraper) WithToken(token string) *ScreepScraper {
	s.authType = authTypeToken
	s.token = token
	return s
}

func (s *ScreepScraper) WithUserPass(username, password string) *ScreepScraper {
	s.authType = authTypeUserPass
	s.username = username
	s.password = password
	return s
}

func (s *ScreepScraper) WithShards(shards []string) *ScreepScraper {
	s.shards = shards
	return s
}

func (s *ScreepScraper) WithClient(client *http.Client) *ScreepScraper {
	s.client = client
	return s
}

func (s *ScreepScraper) Describe(ch chan<- *prometheus.Desc) {}

func (s *ScreepScraper) Collect(ch chan<- prometheus.Metric) {
	log.Info("Collecting metrics from Screeps")

	shard := ""
	if len(s.shards) > 0 {
		shard = s.shards[0]
	}

	// Fetch the metrics data from Shard2
	memoryUrl := fmt.Sprintf("%s/api/user/memory?shard=%s&path=%s", s.host, shard, s.path)
	log := log.WithField("url", memoryUrl)

	req, err := http.NewRequest("GET", memoryUrl, nil)
	req.Header.Add("X-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	// Report API rate limit
	limitRaw := resp.Header.Get("X-Ratelimit-Remaining")
	if limitRaw != "" {
		limit, err := strconv.ParseFloat(limitRaw, 64)
		if err != nil {
			log.Errorf("Error: %s", err)
			return
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("screeps_api_rate_limit_remaining", "Screeps API rate limit", nil, nil),
			prometheus.GaugeValue,
			limit,
		)
	}

	// Report API rate limit reset
	resetRaw := resp.Header.Get("Retry-After")
	if resetRaw != "" {
		reset, err := strconv.ParseFloat(resetRaw, 64)
		if err != nil {
			log.Errorf("Error: %s", err)
			return
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("screeps_api_rate_limit_reset", "Screeps API rate limit reset", nil, nil),
			prometheus.GaugeValue,
			reset,
		)
	}

	// Check if failed (rate limit?)
	if resp.StatusCode != 200 {
		log.Errorf("Error: %d %s", resp.StatusCode, string(body))
		log.Errorf("%v", resp.Header)
		return
	}

	memory := memoryResponse{}
	err = json.Unmarshal(body, &memory)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(memory.Data[3:])))
	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	reader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	metricsRaw, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	metrics := []metric{}
	err = json.Unmarshal(metricsRaw, &metrics)
	if err != nil {
		log.Errorf("Error: %s", err)
		return
	}

	for _, metric := range metrics {
		valueType := prometheus.GaugeValue
		if metric.Type == "counter" {
			valueType = prometheus.CounterValue
		}

		desc := prometheus.NewDesc(metric.Key, "", nil, metric.Labels)
		ch <- prometheus.MustNewConstMetric(desc, valueType, metric.Value)
	}
}
