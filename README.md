# Screeps Prometheus Grafana

Prometheus collector for Screeps and Docker Compose services for Prometheus and Grafana.

## Setup

1. Instrument Screeps bot and write metrics to `Memory.metrics`. See [example library](https://github.com/ryanrolds/screeps-bot/blob/605bcfc1fa7176a2f5b1699b72ceda17ca4eef87/src/lib/metrics.ts#L12) and [example usage](https://github.com/ryanrolds/screeps-bot/blob/605bcfc1fa7176a2f5b1699b72ceda17ca4eef87/src/main.ts#L151-L153).
2. Copy `config.example.yaml` to `config.yaml` and edit as needed, token or username+password required

## Running

`docker-compose up`

Visit: http://localhost:3000/, login with admin/admin, then import `grafana/dashboard.json`.

## Troubleshooting

Check that the Scraper is getting metrics, `http://localhost:8080/metrics`.
