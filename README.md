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

For deployment steps (Kubernetes and Argo CD), see `docs/DEPLOY.md`.

## CI/CD

GitHub Actions runs automated checks for every PR and push to `main`:

- **Backend tests**: `go test ./...` in `backend/` on Go 1.25.
- **Frontend lint**: `npm run lint` in `frontend/` on Node.js 20.
- **Kustomize validation**: renders `infra/k8s/overlays/prod` to ensure manifests stay valid.
- **Docker publish**: builds and pushes `ghcr.io/<org>/team-api:<ref>` from `backend/` on pushes to `main` or version tags.
