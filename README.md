# Screeps Prometheus Grafana

Prometheus collector for Screeps and Docker Compose services for Prometheus and Grafana.

## Notes

curl -X POST -H "Content-Type: application/json" -d '{"email":"...","password":"..."}' http://localhost:21025/api/auth/signin
curl -H "X-Token: ..." http://localhost:21025/api/user/memory
curl -H "X-Token: ..." -H "X-Username: ..." http://localhost:21025/api/user/memory