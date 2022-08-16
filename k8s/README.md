# Screeps Grafana K8s

kubectl create secret generic scrapper-envvars \
  --from-literal=TOKEN=...

```
export DOCKER_REPO=example.domain.com

docker build .
<get final tag (75519489812) from output >
export TAG_NAME=75519489812
docker tag $TAG_NAME $DOCKER_REPO/screeps-scraper:$TAG_NAME
docker push $DOCKER_REPO/screeps-scraper:$TAG_NAME
```

K8s

```
export DOCKER_REPO=example.domain.com
kubectl create secret docker-registry regcred --docker-server=$DOCKER_REPO --docker-username=<username> --docker-password=<password> --docker-email=<email> --namespace screeps
```

```
export DOCKER_REPO=example.domain.com
export TAG_NAME=75519489812
envsubst < k8s/scraper-deployment.yaml | kubectl apply --namespace screeps -f - 
envsubst < k8s/scraper-service.yaml | kubectl apply --namespace screeps -f - 
```
