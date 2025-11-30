=== FILE: README.md ===
# Go Clean Architecture - Feature-Based

Backend aplikasi dengan Clean Architecture menggunakan pendekatan feature-based yang simple, scalable, dan mudah di-maintain.

## 🚀 Features

- ✅ Clean Architecture dengan feature-based structure
- ✅ PostgreSQL untuk data persistence
- ✅ Redis untuk caching
- ✅ MongoDB untuk logging (optional)
- ✅ JWT Authentication
- ✅ Request validation
- ✅ Structured logging
- ✅ Database Migration
- ✅ Unit Testing dengan coverage
- ✅ CORS middleware
- ✅ Dependency Injection
- ✅ Error handling
- ✅ Best practices & SOLID principles

## 📁 Project Structure

```
my-app/
├── cmd/
│   ├── api/              # Application entry point
│   └── migrate/          # Database migration CLI
├── internal/
│   ├── features/         # Feature modules (product, order)
│   │   ├── product/
│   │   │   ├── domain.go
│   │   │   ├── repository.go
│   │   │   ├── repository_test.go
│   │   │   ├── usecase.go
│   │   │   ├── usecase_test.go
│   │   │   ├── handler.go
│   │   │   ├── handler_test.go
│   │   │   └── dto.go
│   │   └── order/
│   │       └── ... (same structure)
│   ├── middleware/       # HTTP middlewares
│   ├── database/         # Database connections
│   └── router/           # Route definitions
├── pkg/                  # Shared packages
│   ├── logger/
│   ├── response/
│   ├── jwt/
│   ├── validator/
│   └── errors/
├── migrations/           # SQL migration files
│   ├── 000001_create_products_table.up.sql
│   ├── 000001_create_products_table.down.sql
│   └── ...
├── tests/
│   ├── mocks/           # Mock implementations
│   └── testutil/        # Test utilities
├── config/              # Configuration
└── docker-compose.yaml
```

## 🛠️ Tech Stack

- **Go 1.23**
- **Gin** - HTTP framework
- **GORM** - ORM
- **PostgreSQL** - Primary database
- **Redis** - Caching
- **MongoDB** - Logging/Analytics
- **JWT** - Authentication
- **Validator** - Request validation
- **Logrus** - Structured logging
- **golang-migrate** - Database migrations
- **Testify** - Testing framework
- **SQLMock** - Database mocking

## 🏃 Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional, tapi recommended)

### Installation

1. Clone repository
```bash
git clone <repo-url>
cd my-app
```

2. Install dependencies
```bash
make install
```

3. Copy environment file
```bash
cp .env.example .env
```

4. Start development environment (All-in-one)
```bash
make dev
```

Atau manual:

```bash
# Start databases
make docker-up

# Wait for databases ready, then run migrations
make migrate-up

# Run application
make run
```

Server akan berjalan di `http://localhost:8080`

## 📝 API Endpoints

### Health Check
```
GET /health
```

### Products
```
GET    /api/v1/products          # Get all products
GET    /api/v1/products/:id      # Get product by ID
POST   /api/v1/products          # Create product (auth required)
PUT    /api/v1/products/:id      # Update product (auth required)
DELETE /api/v1/products/:id      # Delete product (auth required)
```

### Orders
```
POST   /api/v1/orders            # Create order (auth required)
GET    /api/v1/orders            # Get user orders (auth required)
GET    /api/v1/orders/:id        # Get order by ID (auth required)
PATCH  /api/v1/orders/:id/status # Update order status (auth required)
```

## 🗄️ Database Migration

### Create New Migration

```bash
# Manual creation
touch migrations/000003_create_users_table.up.sql
touch migrations/000003_create_users_table.down.sql
```

### Run Migration

```bash
# Migrate up
make migrate-up

# Migrate down
make migrate-down

# Check version
make migrate-version

# Force to specific version
make migrate-force VERSION=2
```

## 🧪 Testing

### Run All Tests

```bash
make test
```

### Run Unit Tests Only

```bash
make test-unit
```

### Run Specific Test

```bash
make test-specific TEST=TestCreateProduct_Success
```

### Generate Coverage Report

```bash
make test-coverage
# Opens coverage.html in browser
```

### Test Coverage Example

```bash
$ make test-coverage

=== RUN   TestCreateProduct_Success
--- PASS: TestCreateProduct_Success (0.00s)
=== RUN   TestGetProduct_Success
--- PASS: TestGetProduct_Success (0.00s)
=== RUN   TestGetProduct_NotFound
--- PASS: TestGetProduct_NotFound (0.00s)

PASS
coverage: 85.4% of statements
Coverage report generated: coverage.html
```

### Testing Best Practices

#### 1. **Unit Test Structure (AAA Pattern)**
```go
func TestCreateProduct_Success(t *testing.T) {
    // Arrange - Setup test data & mocks
    mockRepo := new(mocks.ProductRepositoryMock)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Act - Execute the function
    result, err := usecase.CreateProduct(ctx, req)
    
    // Assert - Verify results
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

#### 2. **Test Coverage Goals**
- Repository Layer: 80%+
- Usecase Layer: 90%+
- Handler Layer: 75%+

#### 3. **Mock Best Practices**
- Mock eksternal dependencies (database, redis, APIs)
- Jangan mock domain objects
- Use interface untuk dependency injection

## 🔧 Development Commands

```bash
# Run application
make run

# Build binary
make build

# Run tests
make test

# Test with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Database migrations
make migrate-up
make migrate-down

# Docker
make docker-up
make docker-down
make docker-logs

# Development (all-in-one)
make dev
```

## 🧪 Testing Example

### Product Usecase Test
```go
// tests successful product creation
func TestCreateProduct_Success(t *testing.T) {
    mockRepo := new(mocks.ProductRepositoryMock)
    logger := testutil.GetTestLogger()
    uc := NewUsecase(mockRepo, logger)

    req := &CreateProductRequest{
        Name:     "Laptop",
        Price:    15000000,
        Stock:    10,
        Category: "Electronics",
    }

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*product.Product")).
        Return(nil)

    result, err := uc.CreateProduct(context.Background(), req)

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Laptop", result.Name)
}
```

### Order Integration Test
```go
// tests order creation with stock validation
func TestCreateOrder_InsufficientStock(t *testing.T) {
    mockOrderRepo := new(mocks.OrderRepositoryMock)
    mockProductRepo := new(mocks.ProductRepositoryMock)
    logger := testutil.GetTestLogger()
    uc := NewUsecase(mockOrderRepo, mockProductRepo, logger)

    product := &product.Product{
        ID: 1, Stock: 5, // Only 5 in stock
    }
    
    mockProductRepo.On("FindByID", mock.Anything, int64(1)).
        Return(product, nil)

    req := &CreateOrderRequest{
        Items: []OrderItemRequest{{ProductID: 1, Quantity: 10}},
    }

    result, err := uc.CreateOrder(context.Background(), 1, req)

    assert.Error(t, err)
    assert.Equal(t, errors.ErrInsufficientStock, err)
    assert.Nil(t, result)
}
```

## 🔑 Key Testing Principles

### 1. **Separation of Concerns**
Setiap layer di-test terpisah dengan mock dependencies-nya:
- **Repository Test**: Mock database (sqlmock)
- **Usecase Test**: Mock repository
- **Handler Test**: Mock usecase

### 2. **Test Independence**
Setiap test harus independent dan bisa run dalam order apapun.

### 3. **Table-Driven Tests**
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateProductRequest
        wantErr bool
    }{
        {
            name: "valid product",
            input: CreateProductRequest{Name: "Product", Price: 100},
            wantErr: false,
        },
        {
            name: "invalid price",
            input: CreateProductRequest{Name: "Product", Price: -100},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.Validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 🚀 Migration Best Practices

### 1. **Naming Convention**
```
000001_create_products_table.up.sql
000001_create_products_table.down.sql
```
- Sequential numbering (6 digits)
- Descriptive name
- Separate up/down files

### 2. **Migration File Structure**

**Up Migration:**
```sql
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
```

**Down Migration:**
```sql
DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
```

### 3. **Migration Rules**
- ✅ Always create indexes for foreign keys
- ✅ Add constraints (NOT NULL, CHECK, etc.)
- ✅ Use IF EXISTS / IF NOT EXISTS
- ✅ Always provide down migration
- ❌ Never modify existing migrations (create new ones)
- ❌ Avoid data migrations in schema migrations

### 4. **Running Migrations in Production**

```bash
# Check current version
make migrate-version

# Test migration in staging first
make migrate-up

# If issues, rollback
make migrate-down

# Force to specific version if needed
make migrate-force VERSION=1
```

## 🏗️ Architecture Principles

### 1. **Dependency Rule**
```
External → Framework → Interface Adapters → Use Cases → Entities
```
Inner layers tidak depend ke outer layers.

### 2. **Single Responsibility**
Setiap file punya satu tanggung jawab:
- `domain.go` - Business entities
- `repository.go` - Data access
- `usecase.go` - Business logic
- `handler.go` - HTTP handling

### 3. **Interface Segregation**
```go
// Good - Small, focused interfaces
type ProductRepository interface {
    Create(ctx context.Context, p *Product) error
    FindByID(ctx context.Context, id int64) (*Product, error)
}

// Bad - Large, monolithic interface
type Repository interface {
    CreateProduct(...)
    CreateOrder(...)
    CreateUser(...)
    // ... 50 more methods
}
```

### 4. **Dependency Injection**
```go
// Constructor injection
func NewUsecase(repo Repository, logger *logger.Logger) Usecase {
    return &usecase{
        repo:   repo,
        logger: logger,
    }
}

// Easy to test with mocks
mockRepo := new(mocks.ProductRepositoryMock)
uc := NewUsecase(mockRepo, testLogger)
```

## 📦 Migration to Microservices

Feature-based structure membuat migration ke microservices mudah:

### Option 1: Extract by Feature
```bash
# Extract product service
mkdir ../product-service
cp -r internal/features/product ../product-service/internal/
cp -r pkg/ ../product-service/pkg/
cp -r migrations/000001* ../product-service/migrations/

# Extract order service
mkdir ../order-service
cp -r internal/features/order ../order-service/internal/
cp -r pkg/ ../order-service/pkg/
cp -r migrations/000002* ../order-service/migrations/
```

### Option 2: Shared Libraries
```bash
# Create shared library
mkdir ../shared-lib
cp -r pkg/ ../shared-lib/

# Update go.mod
go get github.com/yourorg/shared-lib
```

### Communication Between Services
```go
// Add gRPC or HTTP client in infrastructure
type ProductClient interface {
    GetProduct(ctx context.Context, id int64) (*Product, error)
}

// Use in order service
type orderUsecase struct {
    orderRepo     Repository
    productClient ProductClient // External service
    logger        *logger.Logger
}
```

## 🔒 Security Best Practices

1. **Environment Variables**: Never commit `.env` file
2. **JWT Secret**: Use strong random string in production
3. **Database Credentials**: Rotate regularly
4. **Input Validation**: Always validate user input
5. **SQL Injection**: Use parameterized queries (GORM handles this)
6. **Rate Limiting**: Add middleware untuk production

## 📊 Monitoring & Logging

### Structured Logging Example
```go
logger.WithFields(map[string]interface{}{
    "user_id":    userID,
    "product_id": productID,
    "action":     "create_order",
}).Info("Order created successfully")
```

### MongoDB for Logs (Optional)
```go
// Store logs to MongoDB for analytics
type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
    UserID    int64
    Action    string
}

mongoCollection.InsertOne(ctx, logEntry)
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure tests pass (`make test`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push to branch (`git push origin feature/amazing-feature`)
7. Create Pull Request

### Code Review Checklist
- ✅ Tests written and passing
- ✅ Code formatted (`make fmt`)
- ✅ No linter errors (`make lint`)
- ✅ Documentation updated
- ✅ Migration files if DB changes

## 📚 Additional Resources

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://dave.cheney.net/2016/08/20/solid-go-design)
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)

## 📄 License

MIT License

---

**Happy Coding! 🚀**

Jika ada pertanyaan atau issue, silakan buat issue di repository atau hubungi maintainer.