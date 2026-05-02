# Repository Guidelines

## Project Structure & Module Organization
This repository is split into two apps:
- `server/`: Go backend (Gin API, services, repositories, workflow dispatcher, migrations).
  - Entrypoints: `server/cmd/api`, `server/cmd/migrate`
  - Core code: `server/internal/{config,db,handler,middleware,model,repository,service,workflow}`
  - SQL migrations: `server/migrations/`
- `web/`: React + Vite + TypeScript frontend.
  - App code: `web/src/{api,components,pages,lib}`
  - Static build output: `web/dist/` (generated; do not edit manually)
- Root infra: `docker-compose.yml`, `.env.example`

## Build, Test, and Development Commands
- `docker compose up -d db`: start local PostgreSQL.
- `cp .env.example .env`: initialize root runtime config.
- `cd server && go run ./cmd/migrate`: apply DB migrations.
- `cd server && go run ./cmd/api`: run backend at `127.0.0.1:8090` (from `.env`).
- `cd web && npm install && npm run dev`: run frontend locally (Vite dev server).
- `cd web && npm run build`: type-check and create production build in `web/dist`.
- `cd web && npm run preview`: preview the built frontend locally.
- `cd server && go test ./...`: run backend tests (add/maintain as features evolve).

## Coding Style & Naming Conventions
- Go: format with `gofmt` (or `go fmt ./...`), keep packages lowercase, exported names in `PascalCase`.
- TypeScript/React: 2-space indentation, `PascalCase` for page/component files (for example `DashboardPage.tsx`), camelCase for functions/variables.
- Frontend implementation must follow `DESIGN.md` (in repo root) as the primary design reference.
- Prefer `@/` imports in web code (`@` maps to `web/src`).
- Keep handlers/services/repositories separated by responsibility; avoid cross-layer shortcuts.

## Testing Guidelines
- Current coverage is minimal; add tests with each non-trivial change.
- Go tests live as `*_test.go` beside implementation files.
- Frontend tests are not scaffolded yet; if introduced, place under `web/src` and use `*.test.ts(x)`.
- Before opening a PR, run `go test ./...` and `npm run build`.

## Commit & Pull Request Guidelines
- Follow the existing Conventional Commit style seen in history:
  - `feat(server): ...`
  - `feat(web): ...`
  - `docs(readme): ...`
  - `chore(web): ...`
- Keep commits scoped to one concern and module (`server` vs `web`).
- When a requirement is complete, create commits using the correct `scope` (for example `feat(server): ...`), then **must** push the branch to remote and create a PR.
- PRs should include:
  - concise summary of behavior changes,
  - linked issue/task,
  - validation steps run locally,
  - UI screenshots for frontend changes.

## Security & Configuration Tips
- Never commit real secrets in `.env`; use `.env.example` templates only.
- Use only repository-root `.env` / `.env.example` for configuration; do not add per-module `.env.example` duplicates.
- Rotate `JWT_SECRET` and database credentials for non-local environments.
