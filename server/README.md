# Backend (Go API)

Run all commands from this directory (`server/`).

## Setup

- **Create `server/.env`** — copy from `.env.example` in this directory:
  ```bash
  cp .env.example .env
  ```
  Then set at least: `BLUEPRINT_DB_DATABASE`, `BLUEPRINT_DB_USERNAME`, `BLUEPRINT_DB_PASSWORD`, `JWT_SECRET_KEY`. Optional: `BLUEPRINT_DB_HOST` (default localhost), `BLUEPRINT_DB_PORT` (default 5432), `BLUEPRINT_DB_SCHEMA` (default public).
- **Run from `server/`** — the app loads `.env` from the current working directory, so always `cd server` before `make run`.
- Start PostgreSQL: `make docker-run`

## Commands

| Command | Description |
|--------|-------------|
| `make run` | Start API (default port 8080) |
| `make build` | Build binary |
| `make watch` | Live reload (Air) |
| `make test` | Run tests |
| `make itest` | Integration tests (DB) |
| `make docker-run` | Start PostgreSQL container |
| `make docker-down` | Stop PostgreSQL |
| `make clean` | Remove built binary |

## API

- Base URL: `http://localhost:8080`
- Public: `POST /signup`, `POST /signin`, `POST /refresh`, `GET /`, `GET /health`
- Protected: `GET/POST/PUT/DELETE /accounts`, `GET/POST /accounts/:id/transactions`, `GET/POST/PUT/DELETE /users`, `GET /transactions`

Use `Authorization: Bearer <token>` for protected routes (token from `POST /signin`).
