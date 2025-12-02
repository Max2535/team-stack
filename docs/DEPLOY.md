# Deployment Guide

This guide covers deploying the API to Kubernetes either directly with `kubectl` or via Argo CD.

## Prerequisites
- Access to a Kubernetes cluster and `kubectl` configured for it
- Docker installed for building images
- For local development: Docker Desktop with Kubernetes enabled
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
- Create a local git tag matching `TAG` after successful deployment by default (disable with `--no-git-tag`; requires a clean working tree).

## Build Docker Image
Build the backend API image locally:
```bash
docker build -t team-api:latest ./backend
```

## Deploy Dependencies
The application requires PostgreSQL and Redis. Deploy them to your cluster:

```bash
# Create namespace
kubectl create namespace team

# Deploy PostgreSQL and Redis
kubectl apply -f infra/k8s/base/postgres-deployment.yaml -n team
kubectl apply -f infra/k8s/base/redis-deployment.yaml -n team

# Verify they're running
kubectl get pods -n team
```

## Configure Application
1. Review and set runtime values in `infra/k8s/base/api-configmap.yaml`:
   - App settings (port, host, environment)
   - Redis connection (`REDIS_ADDR`)
   - Kafka settings (if using)
   - JWT access token TTL
   - Database connection pool settings

2. Set secret values in `infra/k8s/base/api-secret.yaml`:
   - `DB_DSN`: PostgreSQL connection string (default: `postgres://team:team@postgres:5432/teamdb?sslmode=disable`)
   - `JWT_SECRET`: Strong secret for JWT signing (⚠️ change in production!)

## Deploy with kubectl

### Option 1: Local Development (Docker Desktop)
For local development with Docker Desktop Kubernetes:

1. Ensure the image uses local tag and pull policy in `infra/k8s/base/api-deployment.yaml`:
   ```yaml
   image: team-api:latest
   imagePullPolicy: Never
   ```

2. Deploy the application:
   ```bash
   kubectl apply -k infra/k8s/base -n team
   ```

3. Verify deployment:
   ```bash
   kubectl get pods -n team
   kubectl get svc -n team
   ```

4. Access the API:
   ```bash
   kubectl port-forward -n team svc/team-api 8080:80
   curl http://localhost:8080/api/health
   ```

### Option 2: Production Deployment
For production deployment to a real cluster:

1. Build and push to a container registry:
   ```bash
   docker build -t ghcr.io/your-org/team-api:v1.0.0 ./backend
   docker push ghcr.io/your-org/team-api:v1.0.0
   ```

2. Update image in `infra/k8s/overlays/prod/kustomization.yaml`:
   ```yaml
   images:
     - name: ghcr.io/your-org/team-api
       newTag: v1.0.0
   ```

3. Apply production overlay:
   ```bash
   kubectl apply -k infra/k8s/overlays/prod
   ```

4. Verify deployment:
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

## Post-deploy Verification

### Check Pod Status
```bash
# View all resources
kubectl get all -n team

# Check pod status
kubectl get pods -n team

# View pod details
kubectl describe pod <pod-name> -n team
```

### Check Logs
```bash
# View API logs
kubectl logs -n team deployment/team-api

# Follow logs in real-time
kubectl logs -n team deployment/team-api -f

# View logs from specific pod
kubectl logs -n team <pod-name>
```

### Test API Endpoints
```bash
# Port forward to access locally
kubectl port-forward -n team svc/team-api 8080:80

# Health check
curl http://localhost:8080/api/health

# Expected response:
# {"success":true,"data":{"status":"ok"}}
```

### Verify Configuration
```bash
# Check ConfigMap
kubectl get configmap team-api-config -n team -o yaml

# Check Secrets (values are base64 encoded)
kubectl get secret team-api-secrets -n team -o yaml
```

## Common Operations

### Restart Deployment
```bash
kubectl rollout restart deployment/team-api -n team
```

### Scale Replicas
```bash
# Scale to 3 replicas
kubectl scale deployment/team-api -n team --replicas=3

# Check scaling status
kubectl get deployment team-api -n team
```

### Update Image
```bash
# Update to new image version
kubectl set image deployment/team-api -n team api=team-api:v2.0.0

# Check rollout status
kubectl rollout status deployment/team-api -n team
```

### Rollback Deployment
```bash
# Rollback to previous version
kubectl rollout undo deployment/team-api -n team

# Rollback to specific revision
kubectl rollout undo deployment/team-api -n team --to-revision=2
```

## Troubleshooting

### Pods Not Starting
```bash
# Check pod events
kubectl describe pod <pod-name> -n team

# Common issues:
# - Image pull errors: Check image name and imagePullPolicy
# - CrashLoopBackOff: Check logs for application errors
# - Pending: Check resource requests and node capacity
```

### Database Connection Issues
```bash
# Verify PostgreSQL is running
kubectl get pods -n team | grep postgres

# Check PostgreSQL logs
kubectl logs -n team deployment/postgres

# Test connection from API pod
kubectl exec -n team <api-pod-name> -- wget -O- postgres:5432
```

### Clean Up
```bash
# Delete all resources in namespace
kubectl delete namespace team

# Delete specific resources
kubectl delete -k infra/k8s/base -n team
```
