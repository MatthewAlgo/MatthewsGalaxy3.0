# GEMINI.md — Project Context for Gemini / Antigravity

> This file provides Gemini and Antigravity with comprehensive context about the MatthewsGalaxy project.
> Use this as the primary reference when working on any service in this monorepo.

---

## What Is This Project?

**MatthewsGalaxy 3.0** is a space-themed personal blog platform by Matei-Alexandru Dinu. It consists of four services deployed together:

1. **Backend** (`backend/`) — Go 1.21 / Gin REST API with PostgreSQL, JWT auth, and admin endpoints
2. **Frontend** (`frontend/`) — Next.js 14 App Router with space-themed UI, comments, likes, admin dashboard
3. **Email Service** (`email-service/`) — Go background worker that emails subscribers when new posts are published
4. **MCP Server** (`mcp-server/`) — TypeScript MCP server exposing 10 AI-accessible tools for blog administration

All services share a single PostgreSQL 15 database. Docker Compose orchestrates the full stack.

---

## Architecture Constraints

- **Backend owns the data layer.** MCP server uses the REST API — no direct DB access.
- **Email service is read-heavy.** It queries for new posts and subscribers, sends emails, writes to `email_logs`. It does NOT create/update posts or users.
- **JWT tokens are stateless.** 7-day expiry, HMAC-SHA256, claims include `user_id`, `email`, `role`.
- **Admin routes require two middleware layers:** `AuthMiddleware()` validates the JWT, then `AdminMiddleware()` checks `role == "admin"`.
- **All list endpoints use shared pagination.** `parsePagination()` and `buildPaginatedResponse()` in the backend — always use these, never roll your own.

---

## Technology Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend | Go + Gin | 1.21 |
| Database | PostgreSQL | 15 |
| DB Driver | sqlx | latest |
| Auth | JWT (golang-jwt/jwt/v5) | — |
| Frontend | Next.js (App Router) | 14 |
| Frontend State | React Context (AuthProvider) | — |
| Styling | Vanilla CSS with custom properties | — |
| Email | gomail (SMTP) | — |
| MCP Server | @modelcontextprotocol/sdk + zod | 1.12+ |
| Testing (Go) | stdlib testing | — |
| Testing (TS frontend) | Jest + Testing Library | — |
| Testing (TS MCP) | Vitest | 3.x |
| CI/CD | GitHub Actions + Jenkins | — |
| Containerization | Docker + Docker Compose | — |

---

## Database Schema Reference

**6 tables**, all with `UUID` primary keys (via `pgcrypto`):

### `users`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | `uuid_generate_v4()` |
| email | VARCHAR(255) | UNIQUE, NOT NULL |
| password_hash | VARCHAR(255) | bcrypt |
| name | VARCHAR(100) | NOT NULL |
| role | VARCHAR(20) | `'user'` default, can be `'admin'` |
| avatar_url | TEXT | nullable |
| bio | TEXT | nullable |
| created_at | TIMESTAMPTZ | default NOW() |
| updated_at | TIMESTAMPTZ | default NOW() |

### `posts`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| user_id | UUID FK→users | author |
| title | VARCHAR(200) | NOT NULL |
| slug | VARCHAR(250) | UNIQUE, auto-generated from title |
| content | TEXT | NOT NULL |
| excerpt | TEXT | nullable summary |
| cover_image | TEXT | nullable URL |
| published | BOOLEAN | default FALSE |
| created_at, updated_at | TIMESTAMPTZ | |

### `comments`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| post_id | UUID FK→posts | CASCADE delete |
| user_id | UUID FK→users | CASCADE delete |
| content | TEXT | NOT NULL |
| created_at | TIMESTAMPTZ | |

### `likes`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| post_id | UUID FK→posts | CASCADE delete |
| user_id | UUID FK→users | CASCADE delete |
| created_at | TIMESTAMPTZ | |
| | | UNIQUE(post_id, user_id) |

### `subscriptions`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| email | VARCHAR(255) | UNIQUE |
| active | BOOLEAN | default TRUE |
| subscribed_at | TIMESTAMPTZ | |

### `email_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| post_id | UUID FK→posts | nullable |
| subscriber_email | VARCHAR(255) | |
| sent_at | TIMESTAMPTZ | default NOW() |
| status | VARCHAR(20) | 'sent' or 'failed' |

---

## Complete API Route Map

### Public Routes
```
POST   /api/v1/auth/register          → handlers.Register
POST   /api/v1/auth/login             → handlers.Login
GET    /api/v1/posts?page=&limit=     → handlers.GetPosts (published only)
GET    /api/v1/posts/:slug            → handlers.GetPostBySlug
GET    /api/v1/posts/:slug/comments   → handlers.GetComments
GET    /api/v1/posts/:slug/likes      → handlers.GetLikeStatus (OptionalAuth)
POST   /api/v1/subscribe              → handlers.Subscribe
GET    /api/v1/unsubscribe?email=     → handlers.Unsubscribe
GET    /health                        → health check
```

### Authenticated Routes (AuthMiddleware)
```
GET    /api/v1/me                     → handlers.GetCurrentUser
PATCH  /api/v1/me                     → handlers.UpdateProfile
POST   /api/v1/posts/:slug/comments   → handlers.CreateComment
DELETE /api/v1/posts/:slug/comments/:id → handlers.DeleteComment
POST   /api/v1/posts/:slug/like       → handlers.ToggleLike
```

### Admin Routes (AuthMiddleware + AdminMiddleware)
```
GET    /api/v1/admin/stats            → handlers.GetDashboardStats
GET    /api/v1/admin/posts?page=&limit= → handlers.GetAllPosts
GET    /api/v1/admin/posts/:id        → handlers.GetPostByID
POST   /api/v1/admin/posts            → handlers.CreatePost
PATCH  /api/v1/admin/posts/:id        → handlers.UpdatePost
DELETE /api/v1/admin/posts/:id        → handlers.DeletePost
GET    /api/v1/admin/users?page=&limit= → handlers.GetAllUsers
DELETE /api/v1/admin/users/:id        → handlers.DeleteUser
PATCH  /api/v1/admin/users/:id/role   → handlers.UpdateUserRole
GET    /api/v1/admin/subscribers      → handlers.GetSubscribers
GET    /api/v1/admin/email-logs?page=&limit= → handlers.GetEmailLogs
```

---

## MCP Server (10 Tools)

The MCP server in `mcp-server/` wraps the admin REST API for AI agents:

| Tool | HTTP Call | Description |
|------|-----------|-------------|
| `list_posts` | `GET /admin/posts` | Paginated post listing |
| `get_post` | `GET /posts/:slug` | Post by slug with content |
| `create_post` | `POST /admin/posts` | Create draft or published post |
| `update_post` | `PATCH /admin/posts/:id` | Partial update |
| `delete_post` | `DELETE /admin/posts/:id` | Permanent delete |
| `get_dashboard_stats` | `GET /admin/stats` | Aggregate counts |
| `list_users` | `GET /admin/users` | Paginated user list |
| `list_subscribers` | `GET /admin/subscribers` | Active subscribers |
| `get_post_engagement` | `GET /admin/posts` + sort | Posts by likes+comments |
| `get_recent_email_logs` | `GET /admin/email-logs` | Email delivery status |

**Auth**: Logs in as admin via `MG_ADMIN_EMAIL`/`MG_ADMIN_PASSWORD`, auto-refreshes JWT on 401.
**Transport**: stdio (standard MCP protocol).
**Tests**: 22 vitest tests (api-client: 8, server tools: 14).

---

## File Location Quick Reference

| What you want to do | Where to look |
|---------------------|---------------|
| Add/modify API handler | `backend/internal/handlers/<domain>.go` |
| Register a new route | `backend/cmd/server/main.go` |
| Add/modify a database model | `backend/internal/models/models.go` |
| Change database schema | `backend/internal/database/init.sql` |
| Modify JWT auth/middleware | `backend/internal/middleware/auth.go` |
| Add a frontend page | `frontend/src/app/<path>/page.tsx` |
| Add a React component | `frontend/src/components/` |
| Change frontend API calls | `frontend/src/lib/api.ts` |
| Change auth context/hook | `frontend/src/lib/auth.tsx` |
| Modify global CSS/design tokens | `frontend/src/app/globals.css` |
| Add an MCP tool | `mcp-server/src/server.ts` |
| Change MCP auth/HTTP | `mcp-server/src/api-client.ts` |
| Modify email templates | `email-service/internal/templates/` |
| Change email scheduling | `email-service/internal/scheduler/scheduler.go` |
| CI/CD workflows | `.github/workflows/` + `Jenkinsfile` |
| Docker config | `docker-compose.yml` + per-service `Dockerfile` |

---

## Testing

```bash
# All backend tests (with race detection)
cd backend && go test -v -race ./...

# All email service tests
cd email-service && go test -v ./...

# Frontend Jest tests
cd frontend && npm test

# MCP server Vitest tests (22 tests)
cd mcp-server && npm test

# Full stack via Docker
docker compose up -d
```

---

## Design System

| Token | Value | Usage |
|-------|-------|-------|
| `--bg-primary` | `#0a0a1a` | Deep space background |
| `--bg-secondary` | `#1a0a2e` | Card backgrounds, gradients |
| `--text-primary` | `#ffffff` | Body text |
| `--accent-blue` | `#4a9eff` | Links, CTAs, interactive elements |
| `--accent-pink` | `#ff6b9d` | Highlights, badges |
| `--accent-green` | `#00ff88` | Success states |
| Font | Inter | Google Fonts via `next/font` |

Visual effects: animated star parallax, glassmorphism cards (backdrop blur), fade-in-on-scroll, gradient text, hover glow.

---

## Common Gotchas

1. **Post slugs are auto-generated** from titles with uniqueness suffix. Don't assume `slug == urlize(title)`.
2. **Admin cannot delete self** or other admins — `DeleteUser` handler enforces this.
3. **`OptionalAuthMiddleware()`** on likes endpoint — works without auth but enriches response if token is present.
4. **MCP server stdout is reserved** for MCP protocol. All logs go to stderr via `console.error()`.
5. **Email service writes to backend's database** — `email_logs` table only. Backend reads these logs via `/admin/email-logs`.
6. **Docker internal networking** — services refer to each other as `backend:8080`, `postgres:5432`, not `localhost`.
7. **Frontend uses `'use client'`** only where needed (auth context, interactive components). Pages are Server Components by default.
8. **JWT_SECRET minimum** is 16 chars. Backend crashes at startup if shorter. Use `openssl rand -hex 32` for production.
