---
description: How to properly implement a new feature
---

# Adding Features to GoHex Project

Guide for implementing features following **Clean (Hexagonal) Architecture** principles.

## Architecture Layers

- **Core** (`internal/core`): Entities, DTOs, domain errors, ports (interfaces), services
- **Infrastructure** (`internal/infra`): Repositories, external services
- **Presentation** (`internal/presentation`): HTTP handlers

## Implementation Steps

### 1. Domain Errors (`internal/core/error/<entity>.go`)

```go
package error

var (
    ErrProductNotFound = New("product not found").SetCode(404)
    ErrProductExists   = New("product already exists").SetCode(409)
)
```

### 2. Entity (`internal/core/entity/<entity>.go`)

```go
type Product struct {
    ID        types.ID
    Name      string
    Price     float64
    CreatedAt time.Time
}

func NewProduct(name string, price float64) *Product {
    return &Product{ID: uuid.New(), Name: name, Price: price}
}
```

### 3. DTOs (`internal/core/dto/<entity>.go`)

```go
type CreateProduct struct {
    Name  string
    Price float64
}
```

### 4. Ports (`internal/core/port/<entity>.go`)

```go
// Primary port (for handlers)
type ProductService interface {
    Create(ctx context.Context, input dto.CreateProduct) (*entity.Product, error)
    GetByID(ctx context.Context, id types.ID) (*entity.Product, error)
}

// Secondary port (for repositories)
type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    GetByID(ctx context.Context, id types.ID) (*entity.Product, error)
}
```

### 5. Database Migration

```bash
go run cmd/migration/main.go create <migration_name>
```

**Example:** `database/migrations/YYYYMMDDHHMMSS_products.sql`
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
```

**Commands:**
```bash
go run cmd/migration/main.go up    # Apply all
go run cmd/migration/main.go reset # Rollback all
go run cmd/migration/main.go redo  # Redo last
```

### 6. Repository (`internal/infra/repository/postgres/<entity>.go`)

```go
package postgres

import (
    domainErrors "app/internal/core/error"
    // ... other imports
)

type ProductRepository struct {
    dbGetter pgxTransactor.DBGetter
}

func NewProductRepository(dbGetter pgxTransactor.DBGetter) *ProductRepository {
    return &ProductRepository{dbGetter: dbGetter}
}

func (r *ProductRepository) Create(ctx context.Context, p *entity.Product) error {
    sql, args, _ := psql.Insert("products").
        Columns("id", "name", "price").
        Values(p.ID, p.Name, p.Price).
        Suffix("RETURNING created_at").ToSql()

    err := r.dbGetter(ctx).QueryRow(ctx, sql, args...).Scan(&p.CreatedAt)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == duplicateKeyErrorCode {
            return domainErrors.ErrProductExists
        }
        return fmt.Errorf("execute query: %w", err)
    }
    return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
    sql, args, _ := psql.Select("id", "name", "price", "created_at").
        From("products").Where(sq.Eq{"id": id}).ToSql()

    var p entity.Product
    err := r.dbGetter(ctx).QueryRow(ctx, sql, args...).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainErrors.ErrProductNotFound
        }
        return nil, fmt.Errorf("scan: %w", err)
    }
    return &p, nil
}
```

**Key Points:**
- Import: `domainErrors "app/internal/core/error"`
- Use `psql` (Squirrel) for queries
- Return domain errors
- Handle `pgx.ErrNoRows` → domain error

### 7. Service (`internal/core/service/<entity>/service.go`)

```go
package product

import (
    domainErrors "app/internal/core/error"
    // ... other imports
)

type Service struct {
    productRepo port.ProductRepository
}

func NewService(repo port.ProductRepository) *Service {
    return &Service{productRepo: repo}
}

func (s *Service) Create(ctx context.Context, input dto.CreateProduct) (*entity.Product, error) {
    if input.Price < 0 {
        return nil, domainErrors.ErrInvalidPrice
    }

    product := entity.NewProduct(input.Name, input.Price)
    if err := s.productRepo.Create(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}

func (s *Service) GetByID(ctx context.Context, id types.ID) (*entity.Product, error) {
    return s.productRepo.GetByID(ctx, id)
}
```

### 8. HTTP Handler (`internal/presentation/httpfx/handler/<entity>.go`)

```go
package handler

type createProductReq struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type productResp struct {
    ID        types.ID  `json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"`
}

// CreateProduct
//
//  @Summary     Create product
//  @Tags        products
//  @Accept      json
//  @Produce     json
//  @Param       payload body createProductReq true "Product data"
//  @Success     200 {object} productResp
//  @Failure     400 {object} map[string]string
//  @Router      /products [post]
func (h *Handler) CreateProduct(ctx fiber.Ctx) error {
    var req createProductReq
    if err := ctx.Bind().JSON(&req); err != nil {
        return newBindError(err)
    }

    product, err := h.app.ProductService.Create(ctx.Context(), dto.CreateProduct{
        Name: req.Name, Price: req.Price,
    })
    if err != nil {
        return fmt.Errorf("create product: %w", err)
    }

    return ctx.JSON(productResp{
        ID: product.ID, Name: product.Name,
        Price: product.Price, CreatedAt: product.CreatedAt,
    })
}

// GetProductByID
//
//  @Summary     Get product by ID
//  @Tags        products
//  @Param       id path string true "Product ID"
//  @Success     200 {object} productResp
//  @Failure     404 {object} map[string]string
//  @Router      /products/{id} [get]
func (h *Handler) GetProductByID(ctx fiber.Ctx) error {
    var req struct{ ID types.ID `uri:"id"` }
    if err := ctx.Bind().All(&req); err != nil {
        return newBindError(err)
    }

    product, err := h.app.ProductService.GetByID(ctx.Context(), req.ID)
    if err != nil {
        return fmt.Errorf("get product: %w", err)
    }

    return ctx.JSON(productResp{/*...*/})
}
```

**Handler Rules:**
- Add Swagger docs to all handlers
- Use `ctx.Bind().JSON()` for POST/PUT
- Use `ctx.Bind().All()` for GET with URI params
- Wrap errors: `fmt.Errorf("context: %w", err)`
- Separate request/response types

### 9. Register Routes (`internal/presentation/httpfx/handler/router.go`)

```go
func ApplyRoutes(app *fiber.App, handler *Handler) {
    app.Get("/docs/*", swagger.HandlerDefault)
    app.Post("/products", handler.CreateProduct)
    app.Get("/products/:id", handler.GetProductByID)
}
```

### 10. Update Application (`internal/core/application.go`)

```go
type Application struct {
    UserService    port.UserService
    ProductService port.ProductService  // Add
}

func NewApplication(
    userService port.UserService,
    productService port.ProductService,  // Add
) *Application {
    return &Application{
        UserService: userService,
        ProductService: productService,
    }
}
```

### 11. Register with Uber FX

**Provider:** `internal/presentation/httpfx/provider/product.go`
```go
func ProvideProductService() fx.Option {
    return fx.Options(
        fx.Provide(fx.Annotate(
            postgres.NewProductRepository,
            fx.As(new(port.ProductRepository)),
        )),
        fx.Provide(fx.Annotate(
            product.NewService,
            fx.As(new(port.ProductService)),
        )),
    )
}
```

Add to FX module composition.

### 12. Generate Swagger & Test

```bash
make docs                    # Generate Swagger
go run cmd/api/main.go       # Run API
go test ./...                # Run tests
```

Visit: http://localhost:8080/docs

## Key Patterns

### Error Handling

**Repositories:**
```go
import domainErrors "app/internal/core/error"

if errors.Is(err, pgx.ErrNoRows) {
    return nil, domainErrors.ErrNotFound
}

var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == duplicateKeyErrorCode {
    return domainErrors.ErrAlreadyExists
}
```

**Services:**
```go
if price < 0 {
    return nil, domainErrors.ErrInvalidPrice
}
// Don't wrap domain errors from repo
```

**Handlers:**
```go
// Always wrap for context
if err != nil {
    return fmt.Errorf("operation name: %w", err)
}
```

## Architecture Principles

1. **Dependency Rule**: Presentation → Core ← Infrastructure
2. **Domain Errors**: Import as `domainErrors "app/internal/core/error"`
3. **Interfaces**: Define in `core/port`, implement in infra/services
4. **No Direct Deps**: Handlers use services, not repositories
5. **Business Logic**: In services, not repositories/handlers

## Checklist

- [ ] Domain errors (`internal/core/error/`)
- [ ] Entity (`internal/core/entity/`)
- [ ] DTOs (`internal/core/dto/`)
- [ ] Ports (`internal/core/port/`)
- [ ] Database migration
- [ ] Repository (`internal/infra/repository/postgres/`)
- [ ] Service (`internal/core/service/<entity>/`)
- [ ] HTTP handlers with Swagger docs
- [ ] Register routes
- [ ] Update Application struct
- [ ] Register with Uber FX
- [ ] Generate docs: `make docs`
- [ ] Write tests
