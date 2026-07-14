# istream

Helping people find work that fits them and earn money: a personality-driven job-fit
test, job search with customized ATS-safe resumes, and guided career pivoting
(switch out, reduce hours, change within field, move to consultancy).

## Architecture

Two codebases:

- `backend/` — Go (Gin) REST API, PostgreSQL, JWT auth
- `frontend/` — React + TypeScript (Vite), talks to the backend via `/api` proxy

## Requirements

- Go 1.26+
- Node 20+
- PostgreSQL running locally (user `mohitrawat`, no password)
- Databases: `istream` (dev) and `istream_test` (tests) — `createdb istream istream_test`

## Run

```bash
# backend (port 8080; runs migrations on startup)
cd backend && go run .

# frontend (port 5173; proxies /api to :8080)
cd frontend && npm install && npm run dev
```

Configuration is via env vars with dev defaults: `PORT`, `DATABASE_URL`,
`JWT_SECRET`, `FRONTEND_ORIGIN`. Set a real `JWT_SECRET` outside local dev.

## Test

```bash
cd backend && go test ./...    # uses istream_test (TEST_DATABASE_URL to override)
cd frontend && npm test        # Vitest + Testing Library
```

## API (Phase 1)

| Method | Path               | Auth   | Description                          |
|--------|--------------------|--------|--------------------------------------|
| GET    | /api/health        | none   | Liveness check                       |
| POST   | /api/auth/register | none   | Create account → `{token, user}`     |
| POST   | /api/auth/login    | none   | Log in → `{token, user}`             |
| GET    | /api/auth/me       | Bearer | Current user                         |

## Roadmap

- **Phase 1 — done**: scaffold, register/login, JWT-protected user data, tests
- **Phase 2**: personality test engine → scored job-fit results (weighted with
  2026 in-demand jobs research)
- **Phase 3**: job search (LinkedIn URL builder) + guided Q&A → customized
  ATS-safe resume
- **Phase 4**: career pivot module — question tree with forkable threads
  (switch out / reduce hours / change within field / consultancy) + resources
- **Phase 5**: trends layer — 2026-2030 demand data baked into recommendations,
  end-to-end tests