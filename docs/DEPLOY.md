# Deployment Guide

This guide covers deploying the API to Kubernetes either directly with `kubectl` or via Argo CD.

## Prerequisites
- Access to a Kubernetes cluster and `kubectl` configured for it
- `kustomize` (or `kubectl kustomize`) available locally
- A container registry with push access for `ghcr.io/your-org/team-api` (or your chosen image)
- Argo CD installed in the cluster (for GitOps deployment)
- `docker` CLI for building and pushing the backend image

## One-command deploy
Run the helper script to build, push, and deploy a specific image tag in one go:

```bash
REGISTRY=ghcr.io/your-org/team-api TAG=v1.2.3 ./tooling/deploy.sh --argocd
```

The script will:
- Build and push `REGISTRY:TAG` from the `backend` directory (skip with `--skip-build` if already pushed).
- Render the production overlay with the provided image tag and apply it with `kubectl`.
- Wait for the `team-api` Deployment rollout in the `team` namespace.
- Optionally apply the Argo CD Application when `--argocd` is supplied.

## Prepare runtime configuration
1. Review and set non-secret runtime values in `infra/k8s/base/api-configmap.yaml` (e.g., log level, Kafka topic, Redis and database tuning).
2. Set secret values in `infra/k8s/base/api-secret.yaml`, including `DB_DSN` and a strong `JWT_SECRET`.
3. Update the image tag in `infra/k8s/overlays/prod/kustomization.yaml` (`images[0].newTag`) to the version you intend to deploy.
4. Build and push the API image:
   ```bash
   docker build -t ghcr.io/your-org/team-api:<tag> backend
   docker push ghcr.io/your-org/team-api:<tag>
   ```

## Deploy with kubectl
1. Render the manifests for the production overlay:
   ```bash
   kustomize build infra/k8s/overlays/prod
   ```
2. Apply them to the cluster:
   ```bash
   kustomize build infra/k8s/overlays/prod | kubectl apply -f -
   ```
3. Confirm resources are running in the `team` namespace:
   ```bash
   kubectl get pods -n team
   kubectl get svc -n team
   ```

## Deploy with Argo CD
1. Ensure the `repoURL` in `infra/k8s/argocd-app.yaml` points to your Git repository (including `.git`).
2. Apply the Argo CD Application to the control-plane namespace:
   ```bash
   kubectl apply -f infra/k8s/argocd-app.yaml
   ```
3. Sync the application (via the Argo CD UI or CLI):
   ```bash
   argocd app sync team-api
   ```
   The Application is configured to prune and self-heal, and will create the `team` namespace automatically.

## Post-deploy verification
- Check pod readiness and logs:
  ```bash
  kubectl get pods -n team
  kubectl logs deploy/team-api -n team
  ```
- Verify health endpoint responds via the Service (port-forward example):
  ```bash
  kubectl port-forward -n team svc/team-api 8080:80
  curl http://localhost:8080/api/health
  ```
- Confirm JWT TTL and other runtime settings match your environment in `team-api-config` and `team-api-secrets`.
