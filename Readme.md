# GoHex Project Template

This is a Go project template based on the principles of **Clean (Hexagonal) Architecture**. It is built using Dependency Injection (Uber FX) and provides a clear separation of layers (domain, application, infrastructure).

**Core Technologies:**
* **Web Framework:** [Fiber v3](https://gofiber.io/)
* **Dependency Injection:** [Uber FX](https://uber.go.org/fx)
* **Database:** [PGX v5](https://github.com/jackc/pgx) (PostgreSQL)
* **Migrations:** [Goose](https://github.com/pressly/goose)
* **Logging:** [Zerolog](https://github.com/rs/zerolog)
* **Configuration:** [caarlos0/env](https://github.com/caarlos0/env)
* **API Documentation:** [Swaggo (Swagger)](https://github.com/swaggo/swag)

---

### Project Structure
```
.
├── cmd
│   └── api - main application entry point
│   └── migration - database migration tool
├── database
│   └── migrations - database migration files
└── internal
    ├── application - business logic composition root
    ├── domain
    │   ├── dto - data transfer objects
    │   ├── entity - business entities
    │   └── error - domain-specific errors
    │   └── port - primary and secondary ports
    │   └── service - service implementations (implementations of primary ports)
    ├── infra - infrastructure layer (secondary adapters)
    │   ├── repository - repository implementations
    ├── presentation - presentation layer (primary adapters)
    │   ├── httpfx - HTTP server using Fiber and Uber FX
    │   │   └── postgres - PostgreSQL repositories
```

## Getting Started

### 1. Prerequisites
* [Go](https://go.dev/) (version 1.25.4 or higher) installed.
* A running instance of [PostgreSQL](https://www.postgresql.org/).

### 2. Configuration
The project uses environment variables for configuration.
1.  Copy the `.env.example` file to a new file named `.env`.
    ```bash
    cp .env.example .env
    ```
2.  Edit `.env` to match your local PostgreSQL connection settings (HOST, PORT, USER, PASSWORD, DB).

### 3. Install Dependencies
```bash
go mod tidy
```
### Run
#### Run inside docker
```bash
cp .env.example .env.docker
docker compose --env-file .env.docker up -d
docker compose logs -f api
```

#### Run locally
> you can run postgres inside docker if you don't have it installed locally
> ```bash
> docker compose --env-file .env.docker up postgres -d
> ```
Run API
```bash
go run cmd/api/main.go
```

### 4. Database Migrations
Create a new migration
```bash
go run cmd/migration/main.go create <migration_name>
```
Migration will be created in `./database/migrations` folder.
> Migrations embedded using `embed` package.

API automatically runs migrations on startup. But you can also run them manually:
Migrate up
```bash
go run cmd/migration/main.go up
```
Rollback all migrations
```bash
go run cmd/migration/main.go reset
```
Rollback last migration and migrate up again
```bash
go run cmd/migration/main.go redo
```
> All commands you can find in official [Goose documentation](https://github.com/pressly/goose)

### Generate Swagger Documentation
```bash
make docs
```

### Generate Mocks
```bash
make mocks
```

