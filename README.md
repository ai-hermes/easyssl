# EasySSL

Certimate-style SSL automation system rebuilt with:

- Backend: Go + Gin + PostgreSQL
- Frontend: React + Vite + TailwindCSS + shadcn-style UI primitives
- Deployment: single backend binary can serve frontend static assets

## Current implementation (V1 base)

- JWT admin auth (`/api/auth/*`)
- Access CRUD (`/api/accesses`)
- Workflow CRUD + runs (`/api/workflows`, `/api/workflows/:id/runs`)
- Workflow dispatcher queue with cancel/stats
- Certificate list/download/revoke (`/api/certificates/*`)
- Statistics + notification test endpoints
- PostgreSQL migration runner (`go run ./cmd/migrate`)

Default bootstrap admin:

- email: `admin@easyssl.local`
- password: `1234567890`

## Quick start

### 1) Start DB + API

```bash
docker compose up -d db
cd server
GOPROXY=https://goproxy.cn,direct go run ./cmd/migrate
GOPROXY=https://goproxy.cn,direct go run ./cmd/api
```

### 2) Start web

```bash
cd web
npm install
npm run dev
```

Visit: `http://127.0.0.1:5173`

## Project structure

- `server/`: Gin API, domain services, workflow dispatcher, migrations
- `web/`: Vite React app, dashboard/access/workflow/certificate/settings pages
- `docker-compose.yml`: local PostgreSQL

## Notes

- Workflow execution currently includes a mock certificate issuance step for end-to-end verification.
- Provider adapters are scaffold-ready via `provider` and workflow node config but not yet fully expanded to all target providers.
