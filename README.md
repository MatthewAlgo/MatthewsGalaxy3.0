# MatthewsGalaxy 3.0 🌌

A comprehensive, space-themed blog platform showcasing **Matei-Alexandru Dinu's** professional journey and technical insights. Built with **Next.js**, **Golang** (Gin), **PostgreSQL**, an **email notification microservice**, and a **Model Context Protocol (MCP) server** for AI-powered blog administration.

> *Ad astra per aspera* ✨ — Through hardships to the stars

![Matthew's Galaxy](https://matthewsgalaxy.com/og-image.png)

---

## 📑 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [Local Development](#-local-development)
- [Project Structure](#-project-structure)
- [API Reference](#-api-reference)
- [MCP Server (AI Integration)](#-mcp-server-ai-integration)
- [Email Service](#-email-service)
- [Design System](#-design-system)
- [Testing](#-testing)
- [CI/CD Pipelines](#-cicd-pipelines)
- [Security](#-security)
- [Deployment](#-deployment)
- [Environment Variables](#-environment-variables)
- [Contributing](#-contributing)
- [License](#-license)
- [Author](#-author)

---

## ✨ Features

### Frontend (Next.js 14 + TypeScript)
- **Space-themed design** with animated star parallax background and nebula effects
- **Responsive layout** optimized for all screen sizes (mobile, tablet, desktop)
- **Blog posts** with markdown-like content rendering and cover images
- **Comments & Likes** for authenticated users with real-time feedback
- **Email subscription** widget for newsletter notifications
- **Admin dashboard** with statistics, post management, user management, and subscriber overview
- **Authentication** with JWT-based login/register flows
- **About page** with professional CV, experience timeline, and skills showcase
- **SEO optimized** with proper meta tags, Open Graph images, and semantic HTML
- **Glassmorphism cards** with backdrop blur, gradient text, hover glow effects, and fade-in animations

### Backend (Golang + Gin Framework)
- **RESTful API** with versioned endpoints (`/api/v1/...`)
- **JWT authentication** with HMAC-SHA256 signing and 7-day token expiry
- **Role-based authorization** separating user, authenticated, and admin routes
- **PostgreSQL database** with `sqlx` for type-safe queries
- **Pagination** with a single reusable `parsePagination()` helper (DRY)
- **Admin endpoints** for complete CRUD on posts, users, subscribers, and email logs
- **Slug generation** with uniqueness guarantees for blog post URLs
- **Health check endpoint** at `/health` for container orchestration
- **CORS configuration** for local development and production domains
- **Connection pooling** with configurable max connections and connection lifetime
- **Retry logic** on database connection (5 attempts with exponential backoff)
- **Fail-fast validation** — crashes immediately if critical env vars (`DATABASE_URL`, `JWT_SECRET`) are missing

### Email Microservice (Golang)
- **Automatic scanning** for new published posts on a configurable interval
- **Bulk email notifications** to all active subscribers per new post
- **Space-themed HTML email templates** matching the platform aesthetic
- **Rate limiting** with 100ms delay between sends to avoid SMTP throttling
- **Email delivery logging** with success/failure tracking in the `email_logs` table
- **SMTP configuration** supporting Gmail, custom SMTP, and TLS verification options
- **Graceful shutdown** with proper database connection cleanup

### MCP Server (AI Blog Administration)
- **10 AI-accessible tools** for content management, analytics, and monitoring
- **Blog management**: create, read, update, delete posts from any AI assistant
- **Dashboard stats**: aggregate metrics across users, posts, comments, likes, subscribers
- **Post engagement analytics**: rank posts by likes + comments
- **Email delivery monitoring**: inspect notification success/failure logs
- **JWT auto-refresh**: transparent re-authentication on token expiry
- **Stdio transport**: compatible with Claude Desktop, Cursor, Windsurf, and all MCP clients

### CI/CD Pipelines
- **GitHub Actions**: Automated testing, linting, TypeScript checking, and Docker image builds
- **Jenkins**: Enterprise-grade pipeline with security scanning (gosec, trivy, npm audit), staging/production deployments, Slack notifications, and manual production approval gates
- **Docker**: Multi-stage builds for all services, pushed to GitHub Container Registry
- **Environments**: Staging (auto-deploy on `develop`) and Production (manual approval on `main`)

---

## 🏗️ Architecture

```
                                    ┌──────────────────────┐
                                    │   AI Agents          │
                                    │   (Claude, etc.)     │
                                    └────────┬─────────────┘
                                             │ MCP (stdio)
                                    ┌────────▼─────────────┐
                                    │   MCP Server         │
                                    │   (TypeScript)       │
                                    └────────┬─────────────┘
                                             │ REST/JSON
┌──────────────┐   HTTP    ┌─────────────────▼─────────────┐
│              │ ────────▶ │                               │
│   Frontend   │           │        Backend API            │
│   (Next.js)  │ ◀──────── │        (Go / Gin)             │
│   :3000      │           │        :8080                  │
└──────────────┘           └───────────────┬───────────────┘
                                           │ SQL
                           ┌───────────────▼───────────────┐
                           │                               │
                           │       PostgreSQL 15            │
                           │       :5432                   │
                           │                               │
                           └───────────────▲───────────────┘
                                           │ SQL
                           ┌───────────────┴───────────────┐
                           │                               │
                           │     Email Service (Go)        │
                           │     (background worker)       │
                           │                               │
                           └───────────────────────────────┘
```

**Key design principles:**
- **Single database** shared across backend and email service — no data duplication
- **Backend is the authority** — MCP server calls REST API, never touches DB directly
- **Email service is fire-and-forget** — polls for new posts, sends notifications, logs results
- **Stateless backend** — JWT tokens carry all auth state; any instance can handle any request
- **Container-ready** — all services have Dockerfiles and are orchestrated via Docker Compose

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 20+ (for local development)

### Running with Docker

```bash
# Clone the repository
git clone https://github.com/MatthewAlgo/MatthewsGalaxy3.0.git
cd MatthewsGalaxy3.0

# Start all services
docker compose up -d

# View logs
docker compose logs -f
```

**Services will be available at:**
| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Backend API | http://localhost:8080 |
| Health Check | http://localhost:8080/health |

### Default Admin Account
- **Email**: `admin@example.com`
- **Password**: `admin123`

> ⚠️ **Change these credentials** before deploying to production.

---

## 🛠️ Local Development

### Backend

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/matthewsgalaxy?sslmode=disable"
export JWT_SECRET="your-secret-key-at-least-16-chars"
export ADMIN_EMAIL="admin@example.com"
export FRONTEND_URL="http://localhost:3000"

# Run the server
go run cmd/server/main.go
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Set environment variables
export NEXT_PUBLIC_API_URL="http://localhost:8080"

# Run development server
npm run dev
```

### Email Service

```bash
cd email-service

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/matthewsgalaxy?sslmode=disable"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USER="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export SCAN_INTERVAL="3600"

# Run the service
go run cmd/service/main.go
```

### MCP Server

```bash
cd mcp-server

# Install dependencies
npm install

# Build
npm run build

# Set environment variables
export MG_ADMIN_EMAIL="admin@example.com"
export MG_ADMIN_PASSWORD="admin123"
export MG_API_URL="http://localhost:8080"

# Run (production)
npm start

# Run (development)
npm run dev
```

---

## 📁 Project Structure

```
MatthewsGalaxy/
├── frontend/                     # Next.js 14 application
│   ├── src/
│   │   ├── app/                 # App Router pages
│   │   │   ├── page.tsx         # Homepage with post listing
│   │   │   ├── layout.tsx       # Root layout with star background
│   │   │   ├── globals.css      # Design tokens and global styles
│   │   │   ├── about/           # About page with CV
│   │   │   ├── admin/           # Admin dashboard
│   │   │   ├── auth/            # Login/register pages
│   │   │   └── blog/            # Blog listing and post detail
│   │   ├── components/          # Reusable React components
│   │   │   ├── Header.tsx       # Navigation header
│   │   │   ├── Footer.tsx       # Site footer
│   │   │   ├── PostCard.tsx     # Blog post preview card
│   │   │   ├── StarBackground.tsx # Animated star parallax
│   │   │   ├── SubscribeForm.tsx  # Newsletter subscription widget
│   │   │   └── __tests__/       # Component unit tests (Jest)
│   │   ├── lib/
│   │   │   ├── api.ts           # Typed API client (ApiClient class)
│   │   │   ├── auth.tsx         # AuthProvider context + useAuth hook
│   │   │   └── __tests__/       # API client tests
│   │   └── styles/              # Additional CSS modules
│   ├── jest.config.ts           # Jest configuration
│   ├── jest.setup.ts            # Test setup with mocks
│   └── Dockerfile               # Multi-stage production build
│
├── backend/                      # Golang API server
│   ├── cmd/server/main.go       # Entry point: routes, CORS, middleware
│   ├── internal/
│   │   ├── handlers/            # HTTP handlers (one file per domain)
│   │   │   ├── auth.go          # Register, Login, GetCurrentUser, UpdateProfile
│   │   │   ├── posts.go         # CRUD + pagination + slug generation
│   │   │   ├── comments.go      # GetComments, CreateComment, DeleteComment
│   │   │   ├── likes.go         # ToggleLike, GetLikeStatus
│   │   │   ├── subscriptions.go # Subscribe, Unsubscribe, GetSubscribers
│   │   │   ├── admin.go         # Dashboard stats, user/role mgmt, email logs
│   │   │   └── *_test.go        # Tests for each handler
│   │   ├── middleware/
│   │   │   ├── auth.go          # JWT validation, admin guard, token generation
│   │   │   └── auth_test.go     # Middleware tests
│   │   ├── models/
│   │   │   ├── models.go        # All domain types + DTOs
│   │   │   └── models_test.go   # Model tests
│   │   └── database/
│   │       ├── db.go            # PostgreSQL connection with retry + pooling
│   │       └── init.sql         # Schema: 6 tables, indexes, seed data
│   └── Dockerfile               # Multi-stage production build
│
├── email-service/                # Email notification microservice
│   ├── cmd/service/main.go      # Entry point
│   └── internal/
│       ├── config/              # Environment helpers
│       ├── scheduler/           # Post scanner + notification dispatcher
│       ├── sender/              # SMTP email sender (gomail)
│       │   ├── sender.go
│       │   └── sender_test.go
│       └── templates/           # Space-themed HTML email templates
│
├── mcp-server/                   # MCP server for AI blog administration
│   ├── src/
│   │   ├── index.ts             # Entry point: env config, stdio transport
│   │   ├── api-client.ts        # HTTP client with JWT auth + auto-refresh
│   │   └── server.ts            # 10 MCP tool definitions
│   ├── tests/
│   │   ├── api-client.test.ts   # 8 API client unit tests
│   │   └── server.test.ts       # 14 MCP tool integration tests
│   ├── package.json
│   ├── tsconfig.json
│   └── vitest.config.ts
│
├── .github/workflows/            # GitHub Actions CI/CD
│   ├── ci.yml                   # Test + lint on PR
│   ├── deploy.yml               # Deploy to environments
│   └── docker.yml               # Build + push Docker images
├── Jenkinsfile                   # Enterprise Jenkins pipeline
├── docker-compose.yml            # Full stack orchestration
├── CLAUDE.md                     # Claude Code project context
├── GEMINI.md                     # Gemini project context
└── README.md                     # This file
```

---

## 🔌 API Reference

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Service health check |
| `POST` | `/api/v1/auth/register` | Register new user account |
| `POST` | `/api/v1/auth/login` | Authenticate and get JWT token |
| `GET` | `/api/v1/posts` | List published posts (paginated) |
| `GET` | `/api/v1/posts/:slug` | Get post by URL slug |
| `GET` | `/api/v1/posts/:slug/comments` | Get post comments |
| `GET` | `/api/v1/posts/:slug/likes` | Get like count (+ user's like status if authenticated) |
| `POST` | `/api/v1/subscribe` | Subscribe email to newsletter |
| `GET` | `/api/v1/unsubscribe?email=` | Unsubscribe from newsletter |

### Authenticated Endpoints (require `Authorization: Bearer <token>`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/me` | Get current user profile |
| `PATCH` | `/api/v1/me` | Update profile (name, bio, avatar) |
| `POST` | `/api/v1/posts/:slug/comments` | Create a comment |
| `DELETE` | `/api/v1/posts/:slug/comments/:commentId` | Delete own comment |
| `POST` | `/api/v1/posts/:slug/like` | Toggle like on a post |

### Admin Endpoints (require admin role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/admin/stats` | Dashboard statistics |
| `GET` | `/api/v1/admin/posts` | List all posts (including drafts) |
| `GET` | `/api/v1/admin/posts/:id` | Get post by UUID |
| `POST` | `/api/v1/admin/posts` | Create new post |
| `PATCH` | `/api/v1/admin/posts/:id` | Update post |
| `DELETE` | `/api/v1/admin/posts/:id` | Delete post |
| `GET` | `/api/v1/admin/users` | List all users |
| `DELETE` | `/api/v1/admin/users/:id` | Delete user (cannot delete self/admins) |
| `PATCH` | `/api/v1/admin/users/:id/role` | Update user role |
| `GET` | `/api/v1/admin/subscribers` | List active subscribers |
| `GET` | `/api/v1/admin/email-logs` | Email delivery logs (paginated) |

### Pagination

All paginated endpoints accept `?page=1&limit=20` query parameters:

```json
{
  "data": [...],
  "page": 1,
  "limit": 20,
  "total": 42,
  "total_pages": 3
}
```

---

## 🤖 MCP Server (AI Integration)

The MCP (Model Context Protocol) server enables AI agents to manage the blog through natural language. It wraps the backend REST API with 10 tools:

| Tool | Description |
|------|-------------|
| `list_posts` | List all posts with pagination |
| `get_post` | Get post by slug |
| `create_post` | Create a new post (draft or published) |
| `update_post` | Partial update by ID |
| `delete_post` | Delete by ID |
| `get_dashboard_stats` | Aggregate blog statistics |
| `list_users` | Paginated user list |
| `list_subscribers` | Active email subscribers |
| `get_post_engagement` | Posts ranked by engagement |
| `get_recent_email_logs` | Email delivery logs |

### Setup for Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "matthewsgalaxy": {
      "command": "node",
      "args": ["/path/to/MatthewsGalaxy/mcp-server/dist/index.js"],
      "env": {
        "MG_API_URL": "http://localhost:8080",
        "MG_ADMIN_EMAIL": "admin@example.com",
        "MG_ADMIN_PASSWORD": "admin123"
      }
    }
  }
}
```

See [mcp-server/](mcp-server/) for full documentation.

---

## 📧 Email Service

The email microservice runs as a background worker that:

1. **Scans** for new published posts at a configurable interval (default: 1 hour)
2. **Queries** all active subscribers from the database
3. **Sends** space-themed HTML notification emails via SMTP
4. **Logs** every send attempt (success/failure) to the `email_logs` table

### Gmail Setup

1. Enable 2-Factor Authentication on your Google account
2. Generate an App Password at https://myaccount.google.com/apppasswords
3. Use the app password as `SMTP_PASSWORD`

---

## 🎨 Design System

### Color Palette
| Name | Hex | Usage |
|------|-----|-------|
| Deep Space | `#0a0a1a` | Primary background |
| Nebula Purple | `#1a0a2e` | Secondary background / gradients |
| Star White | `#ffffff` | Primary text |
| Cosmic Blue | `#4a9eff` | Links, interactive elements, accents |
| Galaxy Pink | `#ff6b9d` | Highlights, badges, emphasis |
| Aurora Green | `#00ff88` | Success states, positive feedback |

### Visual Effects
- Animated star parallax background with randomized star sizes
- Glassmorphism cards with `backdrop-filter: blur()` and translucent borders
- Fade-in animations on scroll with CSS `@keyframes`
- Gradient text headings and buttons
- Hover glow effects on interactive elements
- Smooth page transitions

### Typography
- **Font**: Inter (Google Fonts)
- **Scale**: Semantic sizing from body text to display headings

---

## 🧪 Testing

### Backend Tests (Go)
```bash
cd backend && go test -v -race ./...
```

Covers: auth handlers, admin handlers, posts CRUD, comments, likes, subscriptions, pagination logic, JWT middleware, model validation.

### Email Service Tests (Go)
```bash
cd email-service && go test -v ./...
```

Covers: email sender, scheduler logic, template rendering.

### Frontend Tests (Jest + Testing Library)
```bash
cd frontend && npm test
```

Covers: components (Header, Footer, PostCard, StarBackground, SubscribeForm), API client, auth context.

### MCP Server Tests (Vitest)
```bash
cd mcp-server && npm test
```

Covers: API client (login, auth, auto-refresh, error handling — 8 tests), all 10 MCP tools (parameter forwarding, engagement sorting, error formatting — 14 tests). **Total: 22 tests.**

---

## 🔄 CI/CD Pipelines

### GitHub Actions

| Workflow | Trigger | Steps |
|----------|---------|-------|
| `ci.yml` | PR / push | Lint, test, type-check all services |
| `deploy.yml` | Push to `main`/`develop` | Deploy to staging or production |
| `docker.yml` | Push to `main` | Build & push multi-platform Docker images |

### Jenkins Pipeline

Enterprise-grade pipeline (`Jenkinsfile`) with:
- **Parallel builds** for all 4 services
- **Security scanning**: `gosec` (Go), `npm audit` (Node), `trivy` (containers)
- **Test coverage reporting** with archived HTML reports
- **Docker image builds** and pushes to GitHub Container Registry
- **Staging auto-deploy** on `develop`, **production manual approval** on `main`
- **Post-deploy health checks** with retry
- **Slack notifications** on success/failure/unstable
- **Database backups** before production deploys

---

## 🔒 Security

- **JWT tokens** with HMAC-SHA256 signing and mandatory 16+ character secrets
- **Fail-fast startup** — server crashes if `JWT_SECRET` or `DATABASE_URL` is missing
- **bcrypt password hashing** with default cost factor
- **Role-based access control** — `user` and `admin` roles enforced in middleware
- **Admin self-deletion protection** — admins cannot delete themselves or other admins
- **CORS restricted** to configured origins only
- **SQL injection prevention** — parameterized queries via `sqlx` throughout
- **XSS protection** — no raw HTML rendering from user input
- **SMTP TLS verification** — enabled by default, opt-out requires explicit env var
- **MCP server credentials** — admin email/password from env vars only, never logged

---

## 🚢 Deployment

### Docker Compose (Recommended)

```bash
docker compose up -d
```

This starts all 4 services + PostgreSQL with:
- Persistent database volume
- Health check dependencies (backend waits for Postgres)
- Automatic restarts (`unless-stopped`)
- Internal networking (`matthewsgalaxy-network`)

### Production Checklist

- [ ] Change admin credentials from defaults
- [ ] Set a strong `JWT_SECRET` (min 32 chars, use `openssl rand -hex 32`)
- [ ] Configure real SMTP credentials for email notifications
- [ ] Set `FRONTEND_URL` to your production domain
- [ ] Update CORS `AllowOrigins` to your production domain
- [ ] Enable HTTPS/TLS termination (reverse proxy recommended)
- [ ] Set up database backups
- [ ] Configure monitoring and alerting

---

## 🔧 Environment Variables

### Backend
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | — | JWT signing secret (min 16 chars) |
| `ADMIN_EMAIL` | ❌ | `admin@example.com` | Admin account email |
| `FRONTEND_URL` | ❌ | `http://localhost:3000` | Frontend URL for CORS |
| `PORT` | ❌ | `8080` | HTTP server port |
| `GIN_MODE` | ❌ | `release` | Gin framework mode |

### Frontend
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | ❌ | `http://localhost:8080` | Backend API URL |

### Email Service
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | ✅ | — | PostgreSQL connection string |
| `SMTP_HOST` | ❌ | `smtp.gmail.com` | SMTP server hostname |
| `SMTP_PORT` | ❌ | `587` | SMTP server port |
| `SMTP_USER` | ❌ | — | SMTP username |
| `SMTP_PASSWORD` | ❌ | — | SMTP password / app password |
| `FROM_EMAIL` | ❌ | `noreply@matthewsgalaxy.com` | Sender email address |
| `SCAN_INTERVAL` | ❌ | `3600` | Seconds between post scans |
| `FRONTEND_URL` | ❌ | `http://localhost:3000` | For blog post links in emails |
| `BACKEND_URL` | ❌ | `http://localhost:8080` | For unsubscribe links |

### MCP Server
| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `MG_ADMIN_EMAIL` | ✅ | — | Admin email for authentication |
| `MG_ADMIN_PASSWORD` | ✅ | — | Admin password |
| `MG_API_URL` | ❌ | `http://localhost:8080` | Backend API URL |

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests for all affected services
4. Commit with conventional commits (`feat:`, `fix:`, `docs:`, etc.)
5. Push and open a Pull Request

---

## 📝 License

MIT License — feel free to use this project as a template for your own blog!

---

## 👨‍💻 Author

**Matei-Alexandru Dinu**
- 🌐 Website: [matthewsgalaxy.com](https://matthewsgalaxy.com)
- 💼 LinkedIn: [linkedin.com/in/matei-alexandru-dinu](https://www.linkedin.com/in/matei-alexandru-dinu)
- 🐙 GitHub: [github.com/MatthewAlgo](https://github.com/MatthewAlgo)

Ex-SWE Intern @ Google | Ex-SWE Intern @ Snowflake | Ex-SWE & Intern @ NXP Semiconductors | UoB CS '25

---

*Built with ❤️ and too much coffee, somewhere among the stars.*
