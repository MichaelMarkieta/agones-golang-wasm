#agones-golang-wasm

```bash
export COMPUTE_ZONE=northamerica-northeast1-a
export CLUSTER_NAME=agones

gcloud config set compute/zone $COMPUTE_ZONE

gcloud container clusters create $CLUSTER_NAME \
  --cluster-version=1.15 \
  --tags=game-server \
  --scopes=gke-default \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=3 \
  --no-enable-autoupgrade \
  --machine-type=e2-standard-2 \
  --quiet

gcloud config set container/cluster $CLUSTER_NAME
gcloud container clusters get-credentials $CLUSTER_NAME

gcloud compute firewall-rules create game-server-firewall \
  --allow udp:7000-8000,tcp:7000-8000 \
  --target-tags game-server \
  --source-ranges 65.94.196.175/32 \
  --description "Firewall to allow game server udp traffic" \
  --quiet

helm repo add agones https://agones.dev/chart/stable
helm install agones --namespace agones-system --create-namespace agones/agones
```
---
Optional
```bash
gcloud container node-pools create agones-system \
  --cluster=$CLUSTER_NAME \
  --no-enable-autoupgrade \
  --node-taints agones.dev/agones-system=true:NoExecute \
  --node-labels agones.dev/agones-system=true \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=3 \
  --quiet

gcloud container node-pools create agones-metrics \
  --cluster=$CLUSTER_NAME \
  --no-enable-autoupgrade \
  --node-taints agones.dev/agones-metrics=true:NoExecute \
  --node-labels agones.dev/agones-metrics=true \
  --enable-autoscaling \
  --min-nodes=1 \
  --max-nodes=3 \
  --quiet
```
---
Clean Up
```bash
gcloud container cluster delete --cluster=$CLUSTER_NAME
gcloud compute firewall-rules delete game-server-firewall
```