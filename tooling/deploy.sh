#!/usr/bin/env bash
set -euo pipefail

# Automated deployment script for the team API.
# Builds and pushes the image, renders the production kustomize overlay with the
# provided image tag, applies it to the cluster, and optionally bootstraps the
# Argo CD Application.

usage() {
  cat <<'USAGE'
Usage: REGISTRY=ghcr.io/your-org/team-api TAG=v1.2.3 ./tooling/deploy.sh [options]

Options:
  --registry <registry>  Image repository (default: ${REGISTRY:-ghcr.io/your-org/team-api})
  --tag <tag>            Image tag to build/push and deploy (default: ${TAG:-latest})
  --overlay <path>       Kustomize overlay path (default: infra/k8s/overlays/prod)
  --skip-build           Skip docker build/push (assumes the image already exists)
  --argocd               Also apply infra/k8s/argocd-app.yaml after kubectl deploy
  --no-git-tag           Skip creating a git tag matching --tag after a successful deploy
  -h, --help             Show this help message
USAGE
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

require_cmd() {
  local cmd="$1"
  if ! command_exists "$cmd"; then
    echo "error: missing required command '$cmd'" >&2
    exit 1
  fi
}

REGISTRY="${REGISTRY:-ghcr.io/your-org/team-api}"
TAG="${TAG:-latest}"
OVERLAY="infra/k8s/overlays/prod"
SKIP_BUILD=false
APPLY_ARGO=false
CREATE_GIT_TAG=true

while [[ $# -gt 0 ]]; do
  case "$1" in
    --registry)
      REGISTRY="$2"; shift 2 ;;
    --tag)
      TAG="$2"; shift 2 ;;
    --overlay)
      OVERLAY="$2"; shift 2 ;;
    --skip-build)
      SKIP_BUILD=true; shift 1 ;;
    --argocd)
      APPLY_ARGO=true; shift 1 ;;
    --no-git-tag)
      CREATE_GIT_TAG=false; shift 1 ;;
    -h|--help)
      usage; exit 0 ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1 ;;
  esac
done

IMAGE="$REGISTRY:$TAG"

require_cmd kubectl
require_cmd docker

KUSTOMIZE_CMD=()
if command_exists kustomize; then
  KUSTOMIZE_CMD=(kustomize)
elif command_exists kubectl && kubectl kustomize --help >/dev/null 2>&1; then
  KUSTOMIZE_CMD=(kubectl kustomize)
else
  echo "error: kustomize (or 'kubectl kustomize') is required" >&2
  exit 1
fi

echo "Using image: $IMAGE"
echo "Using overlay: $OVERLAY"

die() {
  echo "error: $1" >&2
  exit 1
}

create_git_tag() {
  if [[ "$CREATE_GIT_TAG" != true ]]; then
    return 0
  fi

  if ! command_exists git; then
    echo "git not available; skipping tag creation" >&2
    return 0
  fi

  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "not inside a git repository; skipping tag creation" >&2
    return 0
  fi

  if [[ -n "$(git status --porcelain)" ]]; then
    echo "working tree is not clean; skipping git tag" >&2
    return 0
  fi

  if git tag --list "$TAG" | grep -q "^$TAG$"; then
    echo "tag '$TAG' already exists; skipping" >&2
    return 0
  fi

  git tag -a "$TAG" -m "Deploy $TAG"
  echo "Created git tag '$TAG'."
  echo "Push it with: git push --tags"
}

if [[ ! -d "$OVERLAY" ]]; then
  die "overlay path '$OVERLAY' not found"
fi

K8S_ROOT="infra/k8s"
if [[ ! -d "$K8S_ROOT" ]]; then
  die "expected Kubernetes manifests under '$K8S_ROOT'"
fi

OVERLAY_ABS=$(cd "$OVERLAY" && pwd)
K8S_ROOT_ABS=$(cd "$K8S_ROOT" && pwd)

if [[ "$OVERLAY_ABS" != "$K8S_ROOT_ABS"* ]]; then
  die "overlay must be within $K8S_ROOT"
fi

if [[ "$SKIP_BUILD" != true ]]; then
  echo "Building image..."
  docker build -t "$IMAGE" backend
  echo "Pushing image..."
  docker push "$IMAGE"
else
  echo "Skipping build/push; assuming $IMAGE is available"
fi

TMP_DIR=$(mktemp -d)
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

cp -a "$K8S_ROOT" "$TMP_DIR/k8s"

OVERLAY_REL=${OVERLAY_ABS#${K8S_ROOT_ABS}/}
OVERLAY_TMP="$TMP_DIR/k8s/${OVERLAY_REL}"

pushd "$OVERLAY_TMP" >/dev/null
"${KUSTOMIZE_CMD[@]}" edit set image ghcr.io/your-org/team-api="$IMAGE"
popd >/dev/null

echo "Applying manifests..."
"${KUSTOMIZE_CMD[@]}" build "$OVERLAY_TMP" | kubectl apply -f -

echo "Waiting for rollout..."
kubectl rollout status deploy/team-api -n team --timeout=120s

if [[ "$APPLY_ARGO" == true ]]; then
  echo "Applying Argo CD Application..."
  kubectl apply -f infra/k8s/argocd-app.yaml
fi

create_git_tag

echo "Deployment completed successfully."

