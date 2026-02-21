# Kubernetes Deployment Lab 9

This directory contains Kubernetes manifests for deploying the Moscow Time App (Python) and Wordle Game App (Go) on a local Minikube cluster.

## Contents

- `deployment-python.yml` - Deployment manifest for Moscow Time App (3 replicas)
- `service-python.yml` - NodePort Service for Moscow Time App
- `deployment-go.yml` - Deployment manifest for Wordle Game App (3 replicas)
- `service-go.yml` - NodePort Service for Wordle Game App
- `ingress.yml` - Ingress manifest for routing to both applications

## Setup Instructions

### Prerequisites

- Minikube installed and running
- kubectl configured to access the Minikube cluster
- Docker images pushed to Docker Hub:
  - `escalopax/moscow-time-app:latest`
  - `escalopax/wordle-game:latest`

### Deployment Steps

1. Start Minikube:

```bash
minikube start
```

2. Deploy applications using kubectl:

```bash
kubectl apply -f deployment-python.yml
kubectl apply -f service-python.yml
kubectl apply -f deployment-go.yml
kubectl apply -f service-go.yml
```

3. Enable Kubernetes Dashboard (bonus visualization):

```bash
minikube addons enable dashboard
# Dashboard will be accessible via minikube dashboard command
```

4. For Ingress (bonus), enable the ingress addon:

```bash
minikube addons enable ingress
kubectl apply -f ingress.yml
```

## Verification

### Kubernetes Dashboard (Visual Monitoring)

Access the dashboard UI:

```bash
minikube dashboard
```

This will automatically open your default browser to the Kubernetes Dashboard where you can visualize:
- Pods and their status
- Deployments and replicas
- Services and endpoints
- Resource usage and health
- Logs from containers

**Dashboard Features:**
- Real-time cluster monitoring
- Pod logs and resource metrics
- Service endpoints and port forwarding
- Namespace management

**Screenshot:**
[Add screenshot showing dashboard with pods and services]

### Check Pods and Services

```bash
kubectl get pods,svc
```

**Output:**

```txt
[Add output from kubectl get pods,svc here]
```

### Access Applications

Get the Minikube IP and service endpoints:

```bash
minikube service --all
```

**Output:**

```txt
[Add output from minikube service --all here]
```

### Test Connectivity

#### Python App (Moscow Time)

Via NodePort:

```bash
curl http://$(minikube ip):30001/health
curl http://$(minikube ip):30001/metrics
```

**Output:**
```
[Add curl output here]
```

#### Go App (Wordle Game)

Via NodePort:

```bash
curl http://$(minikube ip):30002/health
curl http://$(minikube ip):30002/metrics
```

**Output:**

```
[Add curl output here]
```

### Ingress Testing (Bonus)

Add hosts to /etc/hosts (or equivalent):

```bash
# Get Minikube IP
minikube ip

# Add to /etc/hosts:
# <minikube-ip> moscow-time.local
# <minikube-ip> wordle-game.local
```

Test via Ingress:
```bash
curl http://moscow-time.local/
curl http://wordle-game.local/
```

**Screenshot:**
[Add screenshot showing browser access]

## Cleanup

Remove all Kubernetes resources:

```bash
kubectl delete -f deployment-python.yml
kubectl delete -f service-python.yml
kubectl delete -f deployment-go.yml
kubectl delete -f service-go.yml
kubectl delete -f ingress.yml
```

Or delete the entire namespace:
```bash
kubectl delete namespace default  # Be careful with this!
```

## Key Features

- **3 Replicas**: Each app runs with 3 replicas for redundancy
- **Health Checks**: Liveness and readiness probes configured
- **Resource Limits**: CPU and memory requests/limits set
- **NodePort Services**: Direct access from outside the cluster
- **Ingress Routing**: DNS-based routing for domain names (bonus)
- **Kubernetes Dashboard**: Web UI for cluster visualization and management (bonus)

## Notes

- Minikube provides a single-node Kubernetes cluster for local development
- The NodePort service type exposes apps on fixed ports (30001, 30002)
- For production, use ClusterIP with a proper Ingress controller
- Resource limits prevent apps from consuming excessive cluster resources

## Dashboard Access Troubleshooting

### If `minikube dashboard` doesn't open a browser:

```bash
# Manually get dashboard URL
minikube dashboard --url

# Then open in your browser:
# http://<url-from-above>
```

### Dashboard Access Token (if needed):

```bash
kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret -o name | grep proxy-token)
```

### Alternative Dashboard Access (manual port-forward):

```bash
kubectl proxy
# Then access: http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/
```

## Advanced kubectl Commands

Monitor pod logs:
```bash
kubectl logs -f deployment/moscow-time-app
kubectl logs -f deployment/wordle-game-app
```

Port-forward to test services locally:
```bash
# Python app
kubectl port-forward svc/moscow-time-app 8000:80

# Go app  
kubectl port-forward svc/wordle-game-app 8001:80
```

Scale deployments:
```bash
kubectl scale deployment moscow-time-app --replicas=5
kubectl scale deployment wordle-game-app --replicas=5
```

Describe resources for troubleshooting:
```bash
kubectl describe pod <pod-name>
kubectl describe deployment moscow-time-app
kubectl describe svc moscow-time-app
```
