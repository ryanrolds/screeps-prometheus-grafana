# Screeps Grafana K8s

> This is a work in progress.

1. Fill out the values in `grafana-config.cm.yaml`.
2. Build and upload the Scraper to your Docker repo
```
export DOCKER_REPO=example.domain.com

docker build .
<get final tag (75519489812) from output >
export TAG_NAME=75519489812
docker tag $TAG_NAME $DOCKER_REPO/screeps-scraper:$TAG_NAME
docker push $DOCKER_REPO/screeps-scraper:$TAG_NAME
```
3. If you're Docker repo requires creds, make sure they are setup in the namespace and the deployment file references `regcred` for Docker repo credentials
K8s

```
export DOCKER_REPO=example.domain.com
kubectl create secret docker-registry regcred --docker-server=$DOCKER_REPO --docker-username=<username> --docker-password=<password> --docker-email=<email> --namespace screeps
```
4. When applying the scraper deployment, make sure to substitute the env vars
```
export DOCKER_REPO=example.domain.com
export TAG_NAME=75519489812
envsubst < k8s/scraper-deployment.yaml | kubectl apply -f - 
```
5. Apply the remaining yaml files

