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

	"github.com/prometheus/client_golang/prometheus"
)

var (
	client = &http.Client{}
)

type ScreepCollector struct {
	token string
}

func NewScreepsCollector(token string) *ScreepCollector {
	return &ScreepCollector{token}
}

func (s *ScreepCollector) Describe(ch chan<- *prometheus.Desc) {}

func (s *ScreepCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Printf("Collecting metrics from Screeps\n")

	req, err := http.NewRequest("GET", "https://screeps.com/api/user/memory?shard=shard2&path=metrics", nil)
	req.Header.Add("X-Token", s.token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	memory := memoryResponse{}
	err = json.Unmarshal(body, &memory)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(memory.Data[3:])))
	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	reader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	metricsRaw, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	metrics := []metric{}
	err = json.Unmarshal(metricsRaw, &metrics)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
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
