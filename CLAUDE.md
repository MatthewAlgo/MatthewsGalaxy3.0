# CLAUDE.md — Project Context for Claude Code

> This file provides Claude Code with comprehensive context about the MatthewsGalaxy project.
> It is the single source of truth for AI-assisted development on this codebase.

---

## Project Overview

**MatthewsGalaxy** is a full-stack blog platform with 4 services in a monorepo:

| Service | Language | Framework | Port | Purpose |
|---------|----------|-----------|------|---------|
| `backend/` | Go 1.21 | Gin | 8080 | REST API, auth, database |
| `frontend/` | TypeScript | Next.js 14 | 3000 | Web UI, SSR, admin dashboard |
| `email-service/` | Go 1.21 | stdlib | — | Background email notifications |
| `mcp-server/` | TypeScript | MCP SDK | stdio | AI tool integration (10 tools) |

**Database**: PostgreSQL 15 — shared by backend and email-service.

---

## Architecture Rules

1. **Backend is the single authority for data access.** The MCP server calls the REST API — it never touches the database directly.
2. **One database, two readers.** Backend and email-service both connect to Postgres, but only backend does CUD operations. Email-service only reads (`email_logs` writes are its only exception).
3. **JWT auth with fail-fast.** `JWT_SECRET` is validated at startup (≥16 chars). Tokens expire after 7 days. HMAC-SHA256 only.
4. **Middleware chain**: `AuthMiddleware()` → `AdminMiddleware()` → handler. Auth sets `userID`, `userEmail`, `userRole` in Gin context.
5. **Pagination is centralized.** `parsePagination(c, defaultLimit)` + `buildPaginatedResponse(data, page, limit, total)` — used by all list endpoints. Do not duplicate this pattern.

---

## Code Conventions

### Go (backend + email-service)

- **File-per-domain**: `auth.go`, `posts.go`, `comments.go`, `likes.go`, `subscriptions.go`, `admin.go`
- **Test files**: `*_test.go` alongside source files
- **Models in one file**: `internal/models/models.go` contains ALL domain types and DTOs
- **Request/Response DTOs**: `CreatePostRequest`, `UpdatePostRequest`, `UserResponse`, etc. — always use structs, never raw `map[string]interface{}`
- **Error responses**: Always `c.JSON(http.StatusXxx, gin.H{"error": "descriptive message"})` then `return`
- **UUID primary keys**: All tables use `uuid_generate_v4()` via `pgcrypto`
- **DB queries**: Use `sqlx` with `$1`, `$2` positional params — never string interpolation
- **Password hashing**: `bcrypt.GenerateFromPassword` / `bcrypt.CompareHashAndPassword`
- **Imports**: Group as stdlib → external → internal

### TypeScript (frontend + mcp-server)

- **Frontend**: Next.js 14 App Router (not Pages Router), `'use client'` only where needed
- **API client**: `src/lib/api.ts` — `ApiClient` class with typed methods, token in `localStorage`
- **Auth**: `AuthProvider` React context + `useAuth()` hook. Never access auth state outside this.
- **MCP server**: `@modelcontextprotocol/sdk` with zod schemas. All tools return `{ content: [{ type: "text", text: string }] }`.
- **Error handling in MCP tools**: All wrapped in `handleToolCall()` which catches `ApiClientError` and returns structured MCP errors.

---

## Database Schema

6 tables in `backend/internal/database/init.sql`:

```
users           UUID PK, email (unique), password_hash, name, role, avatar_url, bio, timestamps
posts           UUID PK, user_id FK→users, title, slug (unique), content, excerpt, cover_image, published, timestamps
comments        UUID PK, post_id FK→posts, user_id FK→users, content, timestamps
likes           UUID PK, post_id FK→posts, user_id FK→users, unique(post_id, user_id)
subscriptions   UUID PK, email (unique), active, subscribed_at
email_logs      UUID PK, post_id FK→posts, subscriber_email, sent_at, status
```

Key indexes: `idx_posts_slug`, `idx_posts_user_id`, `idx_comments_post_id`, `idx_likes_post_user`, `idx_subscriptions_email`.

---

## Key API Routes

```
# Public
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/posts                    (paginated, published only)
GET    /api/v1/posts/:slug
GET    /api/v1/posts/:slug/comments
GET    /api/v1/posts/:slug/likes
POST   /api/v1/subscribe
GET    /api/v1/unsubscribe?email=

# Authenticated
GET    /api/v1/me
PATCH  /api/v1/me
POST   /api/v1/posts/:slug/comments
DELETE /api/v1/posts/:slug/comments/:id
POST   /api/v1/posts/:slug/like

# Admin (auth + admin middleware)
GET    /api/v1/admin/stats
GET    /api/v1/admin/posts              (paginated, all posts)
POST   /api/v1/admin/posts
PATCH  /api/v1/admin/posts/:id
DELETE /api/v1/admin/posts/:id
GET    /api/v1/admin/users              (paginated)
DELETE /api/v1/admin/users/:id
PATCH  /api/v1/admin/users/:id/role
GET    /api/v1/admin/subscribers
GET    /api/v1/admin/email-logs         (paginated)
```

---

## MCP Server Tools

The MCP server (`mcp-server/`) provides 10 tools that wrap the backend API:

```
Blog:       list_posts, get_post, create_post, update_post, delete_post
Dashboard:  get_dashboard_stats, list_users, list_subscribers
Analytics:  get_post_engagement, get_recent_email_logs
```

- Auth: logs in as admin via `MG_ADMIN_EMAIL` / `MG_ADMIN_PASSWORD` env vars
- Auto-refresh: retries once on 401 with fresh token
- Transport: stdio (MCP standard)
- Tests: 22 vitest tests (8 api-client + 14 server tools)

---

## Testing Commands

```bash
# Backend
cd backend && go test -v -race ./...

# Email service
cd email-service && go test -v ./...

# Frontend
cd frontend && npm test

# MCP server
cd mcp-server && npm test
```

---

## Running Locally

```bash
# Full stack via Docker
docker compose up -d

# Or individually:
cd backend    && go run cmd/server/main.go       # needs DATABASE_URL, JWT_SECRET
cd frontend   && npm run dev                      # needs NEXT_PUBLIC_API_URL
cd email-service && go run cmd/service/main.go    # needs DATABASE_URL, SMTP_*
cd mcp-server && npm run dev                      # needs MG_ADMIN_EMAIL, MG_ADMIN_PASSWORD
```

Default admin: `admin@example.com` / `admin123`

---

## File Navigation Guide

When working on specific features, here's where to look:

| Feature | Files |
|---------|-------|
| Add a new API endpoint | `backend/internal/handlers/` + `backend/cmd/server/main.go` (route registration) |
| Add a new data model | `backend/internal/models/models.go` + `backend/internal/database/init.sql` |
| Change auth logic | `backend/internal/middleware/auth.go` + `frontend/src/lib/auth.tsx` |
| Add a frontend page | `frontend/src/app/<route>/page.tsx` |
| Add a React component | `frontend/src/components/` |
| Change API client | `frontend/src/lib/api.ts` |
| Add an MCP tool | `mcp-server/src/server.ts` (tool def) + `mcp-server/tests/server.test.ts` (test) |
| Change email template | `email-service/internal/templates/` |
| Update DB schema | `backend/internal/database/init.sql` |
| Modify CI/CD | `.github/workflows/` or `Jenkinsfile` |

---

## Pitfalls and Gotchas

1. **Slug generation**: Posts auto-generate slugs from titles. If two posts have the same title, a numeric suffix is appended. Don't assume slug = title.
2. **Admin deletion guard**: The `DeleteUser` handler prevents deleting yourself or any admin. Don't bypass this in new admin endpoints.
3. **Optional auth**: `/api/v1/posts/:slug/likes` uses `OptionalAuthMiddleware()` — it works without auth (returns just count) but shows user's like status if authenticated.
4. **Email logs vs email sending**: The `email_logs` table is written by the email-service, but read by the backend's `/admin/email-logs` endpoint. This is the only cross-service data shared through the DB.
5. **MCP server logs to stderr**: All `console.error()` — never `console.log()`. The stdout channel is reserved for MCP protocol messages.
6. **Frontend font**: Inter from Google Fonts, loaded via `next/font/google`. Don't add other font loading mechanisms.
7. **Docker networking**: Services communicate through the `matthewsgalaxy-network` Docker bridge. The backend is `backend:8080` inside Docker, not `localhost:8080`.
