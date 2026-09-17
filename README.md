# Single Sign-On (SSO) gRPC Service

A lightweight, robust **Single Sign-On (SSO)** authentication microservice built with **Go (Golang)** and **gRPC**, designed to handle user authentication, token generation, and permission verification across multiple applications.

> **Credits & Acknowledgements:**  
> This project was developed as an educational project following the YouTube tutorial series by **Nikolay Tuzov** ([Николай Тузов - Golang](https://www.youtube.com/@NikolayTuzov)).

---

## Tech Stack

- **Language:** Go 1.26+
- **API Protocol:** [gRPC](https://grpc.io/) (Protocol Buffers v3)
- **Protobuf Contracts:** [loundxr/protos](https://github.com/loundxr/protos)
- **Database:** SQLite 3 ([mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
- **Tokens & Security:** JWT (`golang-jwt/jwt/v5`), password hashing via `bcrypt`
- **Configuration:** [cleanenv](https://github.com/ilyakaznacheev/cleanenv) (YAML + Environment variables)
- **Logging:** Structured logging with `log/slog` (with custom `slogpretty` handler)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Testing:** Integration testing suite via [testify](https://github.com/stretchr/testify)
- **Automation:** [Taskfile](https://taskfile.dev/)

---

## Key Features

- **Centralized Authentication (SSO):** Issue JWT tokens for client applications based on their unique `app_id`.
- **Multi-App Architecture:** Each registered application (`app`) has its own unique secret key for JWT signing.
- **gRPC API:**
  - `Register`: Creates a new user with secure password hashing.
  - `Login`: Verifies credentials and generates a signed JWT token with a configurable TTL.
  - `IsAdmin`: Validates whether a given user possesses administrative privileges.
- **Automated Test Suite:** Complete end-to-end integration tests (`tests/suite`) running against an isolated test SQLite database.

---

## Project Structure

```text
.
├── cmd/
│   └── sso/                 # Service entry point (main.go)
├── config/                  # Configuration files (local.yaml)
├── internal/
│   ├── app/                 # Application runner & gRPC server lifecycle
│   ├── config/              # Configuration loader
│   ├── domain/models/       # Business models (User, App)
│   ├── grpc/auth/           # gRPC server controller (handlers for Register, Login, IsAdmin)
│   ├── lib/
│   │   ├── jwt/             # JWT token creation with app-specific secret
│   │   └── logger/          # slog setup & pretty handler
│   ├── services/auth/       # Business logic layer
│   └── storage/sqlite/      # SQLite repository layer
├── migrations/              # Database schema migrations
├── storage/                 # Local SQLite database file (sso.db)
├── tests/                   # Integration test suite with dedicated test migrations
└── Taskfile.yaml            # Automation commands
```

---

## Getting Started

### Prerequisites
- [Go 1.21+](https://go.dev/)
- [Task](https://taskfile.dev/) (optional, but recommended)

### Installation & Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/loundxr/sso.git
   cd sso
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Run database migrations:**
   ```bash
   task migrate-up
   ```
   *(Or run migrations manually against `storage/sso.db`).*

4. **Start the gRPC server:**
   ```bash
   go run cmd/sso/main.go --config=./config/local.yaml
   ```
   The gRPC server will start listening on the configured port (e.g., `:44044`).

---

## Running Tests

The project includes an integration test suite verifying user registration, login, and token validation against a test database:

```bash
task test
# or
go test -v ./tests/...
```