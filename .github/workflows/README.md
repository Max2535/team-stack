# GitHub Actions Workflows

This directory contains CI/CD workflows for the team-stack project.

## Workflows

### 🔨 CI - Test and Lint (`ci.yml`)
Runs on every push and pull request to validate code quality.

**Jobs:**
- **lint-backend**: Go code linting (golangci-lint, go vet, go fmt)
- **test-backend**: Go unit tests with coverage
- **lint-frontend**: TypeScript/ESLint checks
- **build-frontend**: Next.js build verification
- **validate-k8s**: Kubernetes manifest validation

**Triggers:**
- Push to `main` or `feature/**` branches
- Pull requests to `main`

---

### 🚀 Build and Deploy (`deploy.yml`)
Automated build and deployment pipeline.

**Jobs:**

#### 1. Build
- Builds Docker image from `backend/`
- Pushes to GitHub Container Registry (ghcr.io)
- Tags: `branch-name`, `branch-name-sha`, `pr-number`
- Uses Docker layer caching for faster builds

#### 2. Deploy
- Updates image tag in `infra/k8s/base/api-deployment.yaml`
- Commits changes back to repository
- Argo CD automatically syncs changes to Kubernetes cluster

#### 3. Notify
- Reports deployment status

**Triggers:**
- Push to `main` or `feature/deploy_k8s`
- Changes in `backend/**`, `infra/k8s/**`, or workflow file

---

## Setup Instructions

### 1. Enable GitHub Container Registry

Your images will be pushed to: `ghcr.io/max2535/team-stack/team-api`

No additional secrets needed - uses `GITHUB_TOKEN` automatically.

### 2. Configure Repository Settings

Go to **Settings → Actions → General**:
- ✅ Enable "Read and write permissions" for workflows
- ✅ Enable "Allow GitHub Actions to create and approve pull requests"

### 3. Update Argo CD Application

Make sure `argocd-app.yaml` points to your repository:
```yaml
source:
  repoURL: https://github.com/Max2535/team-stack.git
  targetRevision: main  # or feature/deploy_k8s
```

### 4. Image Pull from Private Registry (Optional)

If your repository is private, create a Kubernetes secret:

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=Max2535 \
  --docker-password=YOUR_GITHUB_PAT \
  --docker-email=your-email@example.com \
  -n team
```

Then update `api-deployment.yaml`:
```yaml
spec:
  template:
    spec:
      imagePullSecrets:
        - name: ghcr-secret
```

---

## How It Works

### Deployment Flow

```
1. Developer pushes code
   ↓
2. GitHub Actions builds Docker image
   ↓
3. Image pushed to ghcr.io/max2535/team-stack/team-api:branch-sha
   ↓
4. Workflow updates api-deployment.yaml with new image tag
   ↓
5. Changes committed back to Git
   ↓
6. Argo CD detects change (within 3 minutes)
   ↓
7. Argo CD syncs to Kubernetes cluster
   ↓
8. New pods deployed with updated image
```

### Image Tagging Strategy

- `main-abc1234` - Commits to main branch
- `feature-deploy-k8s-abc1234` - Feature branch commits
- `pr-42` - Pull request builds

---

## Usage Examples

### Manual Workflow Trigger

```bash
# Trigger deploy workflow manually
gh workflow run deploy.yml

# View workflow runs
gh run list --workflow=deploy.yml

# View logs
gh run view --log
```

### Check Build Status

```bash
# Get latest run status
gh run list --limit 5

# View specific run
gh run view RUN_ID
```

### Rollback Deployment

```bash
# Find previous image tag
kubectl get deployment team-api -n team -o yaml | grep image:

# Update to previous version
kubectl set image deployment/team-api -n team \
  api=ghcr.io/max2535/team-stack/team-api:main-abc1234

# Or let Argo CD sync from Git
git revert HEAD
git push
```

---

## Troubleshooting

### Build Fails: Permission Denied

Check repository settings → Actions → General → Workflow permissions

### Image Pull Error in Kubernetes

1. Check if repository is private
2. Create `ghcr-secret` (see setup instructions)
3. Verify imagePullSecrets in deployment

### Argo CD Not Syncing

```bash
# Force refresh
kubectl patch application team-api -n argocd \
  --type merge \
  -p '{"metadata":{"annotations":{"argocd.argoproj.io/refresh":"hard"}}}'

# Check Argo CD logs
kubectl logs -n argocd deployment/argocd-application-controller
```

### Workflow Commits Create Loop

Workflows include `[skip ci]` in commit messages to prevent loops.

---

## Advanced Configuration

### Add Slack Notifications

Add to `deploy.yml`:

```yaml
- name: Notify Slack
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

### Run Tests in Kubernetes

Add integration test job that deploys to test namespace first.

### Multi-Environment Deployment

Create separate workflows for staging/production with different branches and overlays.

---

## Security Best Practices

- ✅ GITHUB_TOKEN has minimal required permissions
- ✅ Images scanned with Trivy (can be added)
- ✅ Secrets stored in GitHub Secrets
- ✅ Workflows use pinned action versions with SHA
- ⚠️ Review workflow permissions regularly
