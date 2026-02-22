ArgoCD manifests for Lab 13

Usage:
1. Install ArgoCD via Helm (if cluster reachable):

```bash
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update
helm install argo argo/argo-cd --namespace argocd --create-namespace
```

2. Apply namespace and application manifests:

```bash
kubectl apply -f k8s/ArgoCD/namespaces.yaml
kubectl apply -f k8s/ArgoCD/argocd-python-app.yaml
kubectl apply -f k8s/ArgoCD/argocd-python-dev.yaml
kubectl apply -f k8s/ArgoCD/argocd-python-prod.yaml
```

3. Verify in ArgoCD UI or CLI:

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:443 &
# retrieve admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode
# login
argocd login localhost:8080 --insecure
argocd app list
```

Notes:
- Replace `repoURL` in application manifests with your repository URL.
- Ensure `targetRevision` matches your branch/tag (lab13).
