# Web Application

Monorepo: **Go API** (backend) in `server/` and **Next.js** (frontend) in `frontend/`.

## Structure

```
webapplication/
├── server/          # Go backend (API, DB, auth)
│   ├── cmd/
│   ├── internal/
│   ├── auth/
│   ├── Makefile
│   ├── docker-compose.yml
│   └── go.mod
├── frontend/        # Next.js 15 (TypeScript, Tailwind)
│   ├── src/app/
│   └── package.json
└── README.md
```

## Quick start

### 1. Backend (API)

```bash
cd server
cp .env.example .env   # if you have one; set PORT, DB vars, JWT_SECRET_KEY
make docker-run        # start PostgreSQL
make run               # start API on http://localhost:8080
```

### 2. Frontend

```bash
cd frontend
npm install
npm run dev            # http://localhost:3000
```

The frontend can call the API at `http://localhost:8080` or use the dev proxy at `/api/*` (see `frontend/next.config.ts`).

## Backend (server/)

See **server/README.md** for Makefile targets, DB, and live reload.

## Frontend (frontend/)

See **frontend/README.md** for npm scripts and setup.
