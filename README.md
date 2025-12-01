# Team Stack Template

Mono-repo template for a professional team stack:

- **Frontend**: Next.js 14 + Tailwind CSS + TypeScript
- **Backend**: Go 1.22 + Fiber + Hexagonal-ish Architecture
- **Infra**: Docker Compose, Kubernetes (Kustomize), Argo CD
- **Extras**: Redis cache, PostgreSQL, Kafka-ready event bus

## Layout

- `backend/` – Go API
- `frontend/` – Next.js app
- `infra/` – docker-compose + k8s manifests
- `docs/` – architecture and workflow docs

See `docs/DEV_GUIDE.md` for how to run locally and team workflow.
