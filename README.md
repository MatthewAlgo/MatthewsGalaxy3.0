# MatthewsGalaxy 3.0 🌌

A comprehensive, space-themed blog platform showcasing **Matei-Alexandru Dinu's** professional journey and technical insights. Built with **Next.js**, **Golang** (Gin), **PostgreSQL**, and an **email notification microservice**.

![Matthew's Galaxy](https://matthewsgalaxy.com/og-image.png)

## ✨ Features

### Frontend (Next.js)
- **Space-themed design** with animated star background and nebula effects
- **Responsive layout** optimized for all devices
- **Blog posts** with markdown-like content rendering
- **Comments & Likes** for authenticated users
- **Email subscription** for new post notifications
- **Admin dashboard** for content management
- **About page** with CV, experience timeline, and skills

### Backend (Golang)
- **RESTful API** built with Gin framework
- **JWT authentication** with role-based authorization
- **PostgreSQL database** with efficient queries
- **Admin endpoints** for posts, users, and subscribers management

### Email Microservice
- **Automatic scanning** for new posts
- **Bulk email notifications** to subscribers
- **Themed HTML emails** matching the space aesthetic
- **Configurable scan intervals**

### CI/CD Pipelines
- **GitHub Actions**: Automated testing, linting, and Docker builds
- **Jenkins**: Enterprise-grade deployment pipeline with security scanning
- **Docker**: Multi-platform images pushed to GitHub Container Registry
- **Environments**: Staging (auto-deploy) and Production (manual approval)

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

The services will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Email Service**: http://localhost:8081

### Default Admin Account
- **Email**: admin@example.com
- **Password**: admin123

## 🛠️ Local Development

### Backend

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/matthewsgalaxy?sslmode=disable"
export JWT_SECRET="your-secret-key"

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

## 📁 Project Structure

```
MatthewsGalaxy/
├── frontend/                 # Next.js application
│   ├── src/
│   │   ├── app/             # App router pages
│   │   ├── components/      # React components
│   │   └── lib/             # API client & auth
│   └── Dockerfile
├── backend/                  # Golang API server
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── handlers/        # HTTP handlers
│   │   ├── models/          # Database models
│   │   ├── middleware/      # Auth middleware
│   │   └── database/        # DB connection
│   └── Dockerfile
├── email-service/            # Email microservice
│   ├── cmd/service/         # Entry point
│   └── internal/
│       ├── templates/       # HTML email templates
│       ├── scheduler/       # Post scanner
│       └── sender/          # SMTP handling
├── .github/workflows/        # GitHub Actions CI/CD
├── Jenkinsfile               # Jenkins Pipeline
├── docker-compose.yml
└── README.md
```

## 🔌 API Endpoints

### Public
- `GET /api/v1/posts` - List published posts
- `GET /api/v1/posts/:slug` - Get post by slug
- `GET /api/v1/posts/:slug/comments` - Get post comments
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/subscribe` - Subscribe to newsletter

### Authenticated
- `GET /api/v1/me` - Get current user
- `POST /api/v1/posts/:slug/comments` - Create comment
- `POST /api/v1/posts/:slug/like` - Toggle like

### Admin
- `GET /api/v1/admin/stats` - Dashboard statistics
- `POST /api/v1/admin/posts` - Create post
- `PATCH /api/v1/admin/posts/:id` - Update post
- `DELETE /api/v1/admin/posts/:id` - Delete post
- `GET /api/v1/admin/users` - List users
- `DELETE /api/v1/admin/users/:id` - Delete user

## 🎨 Design System

### Colors
| Name | Hex | Usage |
|------|-----|-------|
| Deep Space | `#0a0a1a` | Primary background |
| Nebula Purple | `#1a0a2e` | Secondary background |
| Star White | `#ffffff` | Primary text |
| Cosmic Blue | `#4a9eff` | Links, accents |
| Galaxy Pink | `#ff6b9d` | Highlights |
| Aurora Green | `#00ff88` | Success states |

### Visual Effects
- Animated star parallax background
- Glassmorphism cards with backdrop blur
- Fade-in animations on scroll
- Gradient text and buttons
- Hover glow effects

## 📧 Email Configuration

For the email service, you'll need SMTP credentials. For Gmail:

1. Enable 2-Factor Authentication on your Google account
2. Generate an App Password at https://myaccount.google.com/apppasswords
3. Use the app password as `SMTP_PASSWORD`

## 🧪 Testing

```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd frontend && npm test
```

## 📝 License

MIT License - feel free to use this project as a template for your own blog!

## 👨‍💻 Author

**Matei-Alexandru Dinu**
- LinkedIn: [linkedin.com/in/username](https://www.linkedin.com/in/username)
- GitHub: [github.com/username](https://github.com/username)
- Website: [matthewsgalaxy.com](https://matthewsgalaxy.com)

---

*Ad astra per aspera* ✨ - Through hardships to the stars
