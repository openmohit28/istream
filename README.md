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

## API

| Method | Path                   | Auth   | Description                          |
|--------|------------------------|--------|--------------------------------------|
| GET    | /api/health            | none   | Liveness check                       |
| POST   | /api/auth/register     | none   | Create account → `{token, user}`     |
| POST   | /api/auth/login        | none   | Log in → `{token, user}`             |
| GET    | /api/auth/me           | Bearer | Current user                         |
| GET    | /api/quiz/questions    | Bearer | 24-question RIASEC bank + scale      |
| POST   | /api/quiz/submit       | Bearer | Score answers → `{id, scores, matches}` |
| GET    | /api/quiz/results      | Bearer | Past results (summaries)             |
| GET    | /api/quiz/results/:id  | Bearer | One result (owner-scoped)            |
| GET    | /api/jobs/search-url   | Bearer | Filters → pre-filtered LinkedIn jobs URL |
| POST   | /api/resumes           | Bearer | Create resume (structured document)  |
| GET    | /api/resumes           | Bearer | List resumes (summaries)             |
| GET    | /api/resumes/:id       | Bearer | One resume (owner-scoped)            |
| PUT    | /api/resumes/:id       | Bearer | Update resume                        |
| DELETE | /api/resumes/:id       | Bearer | Delete resume                        |
| POST   | /api/resumes/:id/keyword-check | Bearer | Score resume vs job description (ATS) |

## How matching works (Phase 2)

The test measures the six RIASEC dimensions (Holland Codes) with 24 Likert
questions. Each catalog job carries a 3-letter Holland code plus demand
outlook and AI-risk level sourced from 2026 labor research (WEF Future of
Jobs, PwC AI Jobs Barometer, Guardian AI-safe careers). Fit = cosine
similarity between the user's profile and the job vector; ranking adds a
small boost for growing fields. Results are snapshotted per user in
`test_results`.

## Job search & resumes (Phase 3)

Job search generates pre-filtered LinkedIn deep links (keywords, location,
remote/hybrid, experience, job type, recency) - no scraping, no ToS risk.
The resume builder is a 7-step guided wizard producing a structured document
rendered as a single-column ATS-safe sheet (print to PDF from the browser).
The keyword check tokenizes a pasted job description, filters stopwords and
posting boilerplate, and reports matched/missing terms with a coverage score -
because 64% of ATS setups auto-reject poor keyword matches (Jobscan, 2026).

## Roadmap

- **Phase 1 — done**: scaffold, register/login, JWT-protected user data, tests
- **Phase 2 — done**: RIASEC personality test → scored job-fit results
  (weighted with 2026 in-demand jobs research), per-user history
- **Phase 3 — done**: LinkedIn job-search URL builder + guided resume wizard
  → ATS-safe resume with keyword coverage check
- **Phase 4**: career pivot module — question tree with forkable threads
  (switch out / reduce hours / change within field / consultancy) + resources
- **Phase 5**: trends layer — 2026-2030 demand data baked into recommendations,
  end-to-end tests