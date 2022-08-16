# Screeps Grafana K8s

kubectl create secret generic poller-envvars \
  --from-literal=GRAPHITE_PORT_8125_UDP_ADDR=statsd.screeps.svc.cluster.local \
  --from-literal=SCREEPS_EMAIL=... \
  --from-literal=SCREEPS_PASSWORD=... \
  --from-literal=SCREEPS_SHARD=... \
  --from-literal=SCREEPS_USERNAME=...

```
export DOCKER_REPO=example.domain.com

docker build .
<get final tag (75519489812) from output >
export TAG_NAME=75519489812
docker tag $TAG_NAME $DOCKER_REPO/screeps-poller:$TAG_NAME
docker push $DOCKER_REPO/screeps-poller:$TAG_NAME
```

K8s

```
export DOCKER_REPO=example.domain.com
kubectl create secret docker-registry regcred --docker-server=$DOCKER_REPO --docker-username=<username> --docker-password=<password> --docker-email=<email> --namespace screeps
```

```
export DOCKER_REPO=example.domain.com
export TAG_NAME=75519489812
envsubst < k8s/poller-deployment.yaml | kubectl apply --namespace screeps -f - 
```
