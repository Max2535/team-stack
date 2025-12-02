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
npm install --legacy-peer-deps
npm run dev
```

If `npm ci` complains about missing entries in the lockfile, run `npm install --package-lock-only --legacy-peer-deps` first to refresh it, then re-run `npm ci`.

## Branching Strategy

- `main` – always deployable
- feature branches: `feature/<short-desc>`
- bugfix branches: `fix/<short-desc>`

Use PR + code review before merge.
