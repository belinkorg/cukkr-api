# Barber Booking API

[![CI](https://github.com/username/barber-api/workflows/CI/badge.svg)](https://github.com/username/barber-api/actions)
[![Docker](https://github.com/username/barber-api/workflows/Docker/badge.svg)](https://github.com/username/barber-api/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Backend API untuk aplikasi booking barbershop dengan Clean Architecture, CI/CD automation, dan production-ready deployment.

## 🚀 Features

- ✅ Clean Architecture dengan feature-based structure
- ✅ PostgreSQL untuk data persistence
- ✅ Redis untuk caching dan async job queue (Asynq)
- ✅ JWT Authentication & Authorization
- ✅ Email OTP verification
- ✅ Async email processing dengan Asynq
- ✅ Request validation
- ✅ Structured logging
- ✅ Database Migration
- ✅ Unit Testing dengan coverage
- ✅ CI/CD dengan GitHub Actions
- ✅ Docker containerization
- ✅ CORS middleware
- ✅ Dependency Injection
- ✅ Error handling
- ✅ Best practices & SOLID principles

## 📁 Project Structure

```
barber-api/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml              # CI workflow (test & build)
│   │   └── docker.yml          # Docker build & push
│   └── PULL_REQUEST_TEMPLATE.md
├── cmd/
│   ├── api/                     # Application entry point
│   │   └── main.go
│   └── migrate/                 # Database migration CLI
│       └── main.go
├── internal/
│   ├── features/                # Feature modules
│   │   ├── user/
│   │   │   ├── domain.go
│   │   │   ├── repository.go
│   │   │   ├── repository_test.go
│   │   │   ├── usecase.go
│   │   │   ├── usecase_test.go
│   │   │   ├── handler.go
│   │   │   ├── handler_test.go
│   │   │   └── dto.go
│   │   ├── barbershop/
│   │   ├── booking/
│   │   └── ...
│   ├── middleware/              # HTTP middlewares
│   ├── database/                # Database connections
│   └── router/                  # Route definitions
├── pkg/                         # Shared packages
│   ├── logger/
│   ├── response/
│   ├── jwt/
│   ├── validator/
│   ├── email/
│   └── errors/
├── migrations/                  # SQL migration files
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   └── ...
├── docs/                        # Swagger documentation
├── .env.example                 # Environment variables template
├── .gitignore
├── .dockerignore
├── Dockerfile                   # Production Docker image
├── docker-compose.yml           # Local development
├── Taskfile.yml                 # Task runner
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Tech Stack

- **Go 1.21** - Programming language
- **Gin** - HTTP framework
- **GORM** - ORM
- **PostgreSQL** - Primary database
- **Redis** - Caching & job queue
- **Asynq** - Async task processing
- **JWT** - Authentication
- **Validator** - Request validation
- **Logrus** - Structured logging
- **golang-migrate** - Database migrations
- **Testify** - Testing framework
- **Docker** - Containerization
- **GitHub Actions** - CI/CD
- **GHCR** - Container registry

## 🏃 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Task (optional, tapi recommended)

### Installation

1. **Clone repository**
```bash
git clone https://github.com/username/barber-api.git
cd barber-api
```

2. **Install Task runner (optional)**
```bash
# macOS
brew install go-task/tap/go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Windows
choco install go-task
```

3. **Setup environment**
```bash
# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi lokal kamu
nano .env
```

4. **Start development environment**
```bash
# All-in-one (start docker, migrate, run)
task dev

# Atau manual:
task docker-up      # Start PostgreSQL & Redis
task migrate-up     # Run migrations
task run            # Run application
```

Server akan berjalan di `http://localhost:8080`

### Quick Commands

```bash
# Development
task run              # Run application
task dev              # Start all (docker + migrate + run)

# Docker
task docker-up        # Start databases
task docker-down      # Stop and remove volumes
task docker-logs      # View logs

# Testing
task test             # Run all tests
task test-unit        # Run unit tests only
task test-coverage    # Generate coverage report

# Database
task migrate-up       # Run migrations
task migrate-down     # Rollback migration
task migrate-version  # Check current version

# Build
task build            # Build binary
task clean            # Clean build artifacts

# Code Quality
task fmt              # Format code
task lint             # Run linter
```

## 📝 API Endpoints

### Health Check
```
GET    /health                      # Health check endpoint
```

## 🔄 Git Workflow & Branch Strategy

### Branch Structure

```
main (production)
  ↑
  └── dev (staging)
       ↑
       ├── feature/booking-system
       ├── feature/payment-integration
       ├── bugfix/email-validation
       └── hotfix/critical-bug
```

### Branch Naming Convention

```bash
# Features
feature/booking-system
feature/user-profile
feature/payment-integration

# Bug fixes
bugfix/login-error
bugfix/email-validation

# Hotfixes (urgent production fixes)
hotfix/security-patch
hotfix/payment-crash

# Chores
chore/update-dependencies
chore/refactor-code
```

### Commit Message Convention

```bash
# Format: <type>: <description>

feat: add booking cancellation feature
fix: resolve email validation bug
docs: update API documentation
style: format code with gofmt
refactor: simplify booking logic
test: add unit tests for user service
chore: update dependencies
```

### Development Workflow

#### 1. **Start New Feature**

```bash
# Update dev branch
git checkout dev
git pull origin dev

# Create feature branch
git checkout -b feature/booking-system

# Work on your feature
# ... make changes ...

# Run tests locally
task test

# Commit changes
git add .
git commit -m "feat: add booking system with validation"

# Push to GitHub
git push origin feature/booking-system
```

#### 2. **Create Pull Request**

```
1. Go to GitHub repository
2. Click "Pull requests" → "New pull request"
3. Base: dev ← Compare: feature/booking-system
4. Fill PR template:
   - Description
   - Type of change
   - Checklist
5. Create pull request
```

**CI akan otomatis run:**
- ✅ Run all tests
- ✅ Build application
- ✅ Comment hasil di PR

#### 3. **Review & Merge**

```
1. Wait for CI to pass ✅
2. Request review (jika ada team)
3. Address review comments
4. Merge pull request ke dev
5. Delete feature branch
```

**Setelah merge ke dev:**
- ✅ CI run lagi
- ✅ Docker image di-build → `ghcr.io/username/barber-api:dev`

#### 4. **Deploy to Production**

```bash
# Buat PR: dev → main di GitHub
# Setelah merged:

git checkout main
git pull origin main

# Create release tag
git tag v1.0.0
git push origin v1.0.0
```

**Setelah push tag:**
- ✅ Docker images di-build:
    - `ghcr.io/username/barber-api:v1.0.0`
    - `ghcr.io/username/barber-api:latest`

### Hotfix Workflow (Emergency)

```bash
# Create hotfix dari main
git checkout main
git pull origin main
git checkout -b hotfix/critical-payment-bug

# Fix the bug
# ... make changes ...

# Test
task test

# Commit & push
git add .
git commit -m "fix: resolve critical payment bug"
git push origin hotfix/critical-payment-bug

# Buat 2 PR:
# 1. hotfix/xxx → main (urgent)
# 2. hotfix/xxx → dev (keep in sync)

# Setelah merged ke main, create tag
git checkout main
git pull origin main
git tag v1.0.1 -m "Hotfix: critical payment bug"
git push origin v1.0.1
```

## 🔒 Branch Protection Rules

### `dev` Branch

**Required:**
- ✅ Require pull request before merging
- ✅ Require approvals: 1 (atau 0 untuk solo dev)
- ✅ Require status checks to pass
    - ✅ test (CI workflow)
- ✅ Require conversation resolution before merging

### `main` Branch

**Required:**
- ✅ Require pull request before merging
- ✅ Require approvals: 1
- ✅ Require status checks to pass
    - ✅ test (CI workflow)
- ✅ Require conversation resolution before merging
- ✅ Do not allow bypassing the above settings

**Setup di GitHub:**
```
Repository → Settings → Branches → Add branch protection rule
```

## 🐳 Docker & Deployment

### Local Development

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api

# Stop all
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Build Docker Image

```bash
# Build
docker build -t barber-api:latest .

# Run
docker run -d \
  --name barber-api \
  -p 8080:8080 \
  --env-file .env \
  barber-api:latest

# Check logs
docker logs -f barber-api
```

### Pull from GitHub Container Registry

```bash
# Login (untuk private images)
echo $GITHUB_TOKEN | docker login ghcr.io -u username --password-stdin

# Pull
docker pull ghcr.io/username/barber-api:latest
docker pull ghcr.io/username/barber-api:dev
docker pull ghcr.io/username/barber-api:v1.0.0

# Run
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  ghcr.io/username/barber-api:latest
```

### Available Docker Tags

```
ghcr.io/username/barber-api:latest    # Latest stable (from main)
ghcr.io/username/barber-api:dev       # Development (from dev)
ghcr.io/username/barber-api:v1.0.0    # Specific version (from tag)
```

## 🧪 Testing

### Test Structure

```
internal/features/user/
├── repository_test.go    # Database layer tests (with sqlmock)
├── usecase_test.go       # Business logic tests (with mocks)
└── handler_test.go       # HTTP handler tests (with mocks)
```

### Writing Tests

#### Repository Test (with sqlmock)
```
func TestCreateUser_Success(t *testing.T) {
    // Setup
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    require.NoError(t, err)

    repo := NewRepository(gormDB)

    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }

    // Mock expectations
    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "users"`).
        WithArgs(user.Name, user.Email, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()

    // Execute
    err = repo.Create(context.Background(), user)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

#### Usecase Test (with mocks)
```
func TestCreateBooking_Success(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.BookingRepositoryMock)
    mockBarberRepo := new(mocks.BarbershopRepositoryMock)
    logger := testutil.GetTestLogger()
    uc := NewUsecase(mockRepo, mockBarberRepo, logger)

    req := &CreateBookingRequest{
        BarbershopID: 1,
        DateTime:     time.Now().Add(24 * time.Hour),
        ServiceType:  "Haircut",
    }

    barbershop := &Barbershop{ID: 1, Name: "Great Barber"}
    
    mockBarberRepo.On("FindByID", mock.Anything, int64(1)).
        Return(barbershop, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*booking.Booking")).
        Return(nil)

    // Act
    result, err := uc.CreateBooking(context.Background(), 1, req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(1), result.BarbershopID)
    mockRepo.AssertExpectations(t)
    mockBarberRepo.AssertExpectations(t)
}
```

### Run Tests

```bash
# All tests
task test

# Unit tests only
task test-unit

# With coverage
task test-coverage

# Specific test
task test-specific TEST=TestCreateBooking_Success

# Verbose
go test -v ./...
```

### Coverage Goals

- Repository Layer: **80%+**
- Usecase Layer: **90%+**
- Handler Layer: **75%+**

## 🗄️ Database Migrations

### Create New Migration

```bash
task migrate-create NAME=create_bookings_table

# Creates:
# migrations/000XXX_create_bookings_table.up.sql
# migrations/000XXX_create_bookings_table.down.sql
```

### Migration File Example

**Up Migration:**
```
-- migrations/000002_create_bookings_table.up.sql
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    barbershop_id BIGINT NOT NULL REFERENCES barbershops(id) ON DELETE CASCADE,
    booking_date TIMESTAMP NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_price DECIMAL(10,2) NOT NULL CHECK (total_price > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_barbershop_id ON bookings(barbershop_id);
CREATE INDEX idx_bookings_status ON bookings(status);
```

**Down Migration:**
```
-- migrations/000002_create_bookings_table.down.sql
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_barbershop_id;
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP TABLE IF EXISTS bookings;
```

### Migration Commands

```bash
# Run migrations
task migrate-up

# Rollback last migration
task migrate-down

# Check version
task migrate-version

# Force to specific version
task migrate-force VERSION=1
```

### Migration Best Practices

- ✅ Use sequential numbering (000001, 000002, ...)
- ✅ Descriptive names (create_bookings_table, add_user_phone)
- ✅ Add indexes for foreign keys
- ✅ Use constraints (NOT NULL, CHECK, etc.)
- ✅ Always provide down migration
- ✅ Test in development first
- ❌ Never modify existing migrations (create new ones)
- ❌ Don't mix schema and data migrations

## 🚀 CI/CD Pipeline

### Workflows

#### 1. **CI Workflow** (`.github/workflows/ci.yml`)

**Triggers:**
- Pull request ke `main` atau `dev`
- Push ke `main` atau `dev`

**Steps:**
1. Checkout code
2. Setup Go 1.21
3. Download dependencies
4. Run tests dengan PostgreSQL & Redis services
5. Build binary
6. Comment hasil di PR (jika PR)

#### 2. **Docker Workflow** (`.github/workflows/docker.yml`)

**Triggers:**
- Push ke `main` atau `dev`
- Push tag `v*`

**Steps:**
1. Checkout code
2. Login ke GitHub Container Registry
3. Build Docker image
4. Push dengan tags:
    - `dev` branch → `:dev`
    - `main` branch → `:latest`
    - Tag `v1.0.0` → `:v1.0.0` dan `:latest`

### Workflow Diagram

```
┌─────────────────────────────────────────────────────┐
│ Developer: Push code                                │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Create PR: feature/xxx → dev                        │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ CI Workflow                                         │
│ ✅ Run tests (with PostgreSQL & Redis)             │
│ ✅ Build binary                                     │
│ ✅ Comment PR                                       │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Review & Merge PR                                   │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Push to dev branch                                  │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ CI Workflow (again)                                 │
│ ✅ Run tests                                        │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Docker Workflow                                     │
│ ✅ Build image                                      │
│ ✅ Push to ghcr.io/username/barber-api:dev         │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Staging Server: Deploy dev image                    │
└─────────────────────────────────────────────────────┘

                 (When ready for production)
                 
┌─────────────────────────────────────────────────────┐
│ Create PR: dev → main                               │
│ Review & Merge                                      │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Create Release Tag: v1.0.0                          │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Docker Workflow                                     │
│ ✅ Build image                                      │
│ ✅ Push to ghcr.io/username/barber-api:v1.0.0     │
│ ✅ Push to ghcr.io/username/barber-api:latest     │
└────────────────┬────────────────────────────────────┘
                 ↓
┌─────────────────────────────────────────────────────┐
│ Production Server: Deploy v1.0.0 image              │
└─────────────────────────────────────────────────────┘
```

## 🏗️ Architecture Principles

### 1. Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│ External (Framework & Drivers)          │
│ - Gin, GORM, PostgreSQL, Redis          │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ Interface Adapters                      │
│ - Handlers, Repositories                │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ Use Cases (Business Logic)              │
│ - Usecases                              │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ Entities (Domain)                       │
│ - Domain models, Business rules         │
└─────────────────────────────────────────┘
```

**Dependency Rule:** Inner layers tidak depend ke outer layers.

### 2. Feature-Based Structure

Setiap feature adalah module yang self-contained:

```
internal/features/booking/
├── domain.go       # Business entities
├── repository.go   # Data access interface & implementation
├── usecase.go      # Business logic
├── handler.go      # HTTP handlers
├── dto.go          # Data transfer objects
└── *_test.go       # Tests
```

**Benefits:**
- ✅ Easy to understand
- ✅ Easy to test
- ✅ Easy to scale (extract to microservice)
- ✅ Clear boundaries

### 3. Dependency Injection

```
// Constructor injection
type bookingUsecase struct {
    bookingRepo     Repository
    barbershopRepo  barbershop.Repository
    emailService    email.Service
    logger          *logger.Logger
}

func NewUsecase(
    bookingRepo Repository,
    barbershopRepo barbershop.Repository,
    emailService email.Service,
    logger *logger.Logger,
) Usecase {
    return &bookingUsecase{
        bookingRepo:    bookingRepo,
        barbershopRepo: barbershopRepo,
        emailService:   emailService,
        logger:         logger,
    }
}
```

**Benefits:**
- ✅ Easy to test (inject mocks)
- ✅ Loose coupling
- ✅ Flexible configuration

## 🔐 Security Best Practices

### Environment Variables

**Never commit:**
- ❌ `.env` file
- ❌ Secrets or credentials
- ❌ API keys

**Always use:**
- ✅ `.env.example` (template without values)
- ✅ Environment variables in production
- ✅ Secret managers (for production)

### Authentication

```
// JWT middleware
router.Use(middleware.AuthMiddleware(jwtService))

// Protected routes
authorized := router.Group("/api/v1")
authorized.Use(middleware.AuthMiddleware(jwtService))
{
    authorized.POST("/bookings", bookingHandler.Create)
    authorized.GET("/bookings", bookingHandler.GetUserBookings)
}
```

### Input Validation

```
type CreateBookingRequest struct {
    BarbershopID int64     `json:"barbershop_id" validate:"required,gt=0"`
    DateTime     time.Time `json:"date_time" validate:"required"`
    ServiceType  string    `json:"service_type" validate:"required,min=3,max=100"`
}

// Validate
if err := validator.Validate(req); err != nil {
    return response.BadRequest(c, err.Error())
}
```

### Security Headers

```
// CORS
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))

// Security headers
router.Use(func(c *gin.Context) {
    c.Header("X-Content-Type-Options", "nosniff")
    c.Header("X-Frame-Options", "DENY")
    c.Header("X-XSS-Protection", "1; mode=block")
    c.Next()
})
```

## 📊 Monitoring & Logging

### Structured Logging

```
logger.WithFields(map[string]interface{}{
    "user_id":       userID,
    "barbershop_id": barbershopID,
    "booking_id":    bookingID,
    "action":        "create_booking",
    "status":        "success",
}).Info("Booking created successfully")
```

### Error Tracking

```
logger.WithFields(map[string]interface{}{
    "user_id":    userID,
    "error":      err.Error(),
    "stack_trace": debug.Stack(),
}).Error("Failed to create booking")
```

### Health Check Endpoint

```
// GET /health
{
    "status": "healthy",
    "database": "connected",
    "redis": "connected",
    "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🤝 Contributing

### Pull Request Process

1. **Fork** repository
2. **Create** feature branch (`feature/amazing-feature`)
3. **Write** tests for changes
4. **Ensure** tests pass (`task test`)
5. **Format** code (`task fmt`)
6. **Commit** changes (`git commit -m 'feat: add amazing feature'`)
7. **Push** to branch (`git push origin feature/amazing-feature`)
8. **Create** Pull Request

### Code Review Checklist

- ✅ Tests written and passing
- ✅ Code follows project structure
- ✅ No linter errors
- ✅ Documentation updated
- ✅ Migration files (if DB changes)
- ✅ No breaking changes (or documented)

### Coding Standards

- Follow Go best practices
- Use meaningful variable names
- Add comments for complex logic
- Keep functions small and focused
- Follow DRY (Don't Repeat Yourself)
- Write tests for new features

## 📚 Additional Resources

### Learning Materials
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles in Go](https://dave.cheney.net/2016/08/20/solid-go-design)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Documentation
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Docker Documentation](https://docs.docker.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

### Tools
- [Task](https://taskfile.dev/) - Task runner
- [golangci-lint](https://golangci-lint.run/) - Go linters
- [Air](https://github.com/cosmtrek/air) - Live reload
- [Swagger](https://swagger.io/) - API documentation

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🎯 Quick Reference

### Most Used Commands

```bash
# Development
task dev              # Start everything
task run              # Run app only
task test             # Run tests

# Docker
task docker-up        # Start services
task docker-down      # Stop services
task docker-logs      # View logs

# Database
task migrate-up       # Run migrations
task migrate-down     # Rollback
```

### Environment Variables Template

```
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=barber_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### Git Flow Summary

```bash
# Feature development
feature/xxx -> PR -> dev -> merge -> Docker :dev

# Production release
dev -> PR -> main -> merge
git tag v1.0.0 → Docker :v1.0.0 :latest

# Hotfix
main -> hotfix/xxx -> PR -> main -> merge
git tag v1.0.1 -> Docker :v1.0.1 :latest
```

---

**Happy Coding! 🚀**

For questions or issues, please create an issue or contact the maintainer.

**Maintainer:** [Haryanda Alfitroh](https://github.com/Haryandaal)