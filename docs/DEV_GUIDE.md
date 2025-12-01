# Dev Guide

## Local Dev

1. Start DB, Redis, API, Web via Docker Compose:

```bash
cd infra
docker-compose -f docker-compose.dev.yml up -d
```

2. Or run each service manually:

Backend:

```bash
cd backend
go run ./cmd/api
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

## Branching Strategy

- `main` – always deployable
- feature branches: `feature/<short-desc>`
- bugfix branches: `fix/<short-desc>`

Use PR + code review before merge.
